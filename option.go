package defaults

import (
	"bytes"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Option func(filler *filler)

func UseDefaultTag(tag string) Option {
	return func(f *filler) {
		f.DefaultTag = tag
	}
}

func UseOmitKey(key string) Option {
	return func(f *filler) {
		f.OmitKey = key
	}
}

func ParseDuration() Option {
	return func(f *filler) {
		f.FuncByKind[reflect.Int64] = f.skipIfTagEmpty(func(field *Field) {
			if field.Value.Type() == reflect.TypeOf(time.Second) {
				value, _ := time.ParseDuration(field.Tag)
				field.Value.Set(reflect.ValueOf(value))
				return
			}
			f.FuncByKind[reflect.Int](field)
		})
	}
}

func UseTimeFormat(layout string) Option {
	return func(f *filler) {
		f.FuncByType[reflect.TypeOf(time.Time{})] = f.skipIfTagEmpty(func(field *Field) {
			if field.Value.IsZero() {
				value, _ := time.Parse(layout, field.Tag)
				field.Value.Set(reflect.ValueOf(value))
			}
		})
	}
}

func UseDefaultType(defVal interface{}) Option {
	return func(f *filler) {
		value := Indirect(reflect.ValueOf(defVal))
		f.FuncByType[value.Type()] = func(field *Field) {
			if field.Value.IsZero() && field.Tag != f.OmitKey {
				reflect.Indirect(field.Value).Set(value)
			}
		}
	}
}

func UseDefault() Option {
	return func(f *filler) {
		f.useDefaultKindFuncs()
	}
}

func (f *filler) useDefaultKindFuncs() {
	fns := f.FuncByKind

	fns[reflect.Bool] = f.skipIfTagEmpty(func(field *Field) {
		value, _ := strconv.ParseBool(field.Tag)
		field.Value.SetBool(value)
	})

	fns[reflect.Int] = f.skipIfTagEmpty(func(field *Field) {
		value, _ := strconv.ParseInt(field.Tag, 10, 64)
		field.Value.SetInt(value)
	})
	fns[reflect.Int8] = fns[reflect.Int]
	fns[reflect.Int16] = fns[reflect.Int]
	fns[reflect.Int32] = fns[reflect.Int]
	fns[reflect.Int64] = fns[reflect.Int]

	fns[reflect.Uint] = f.skipIfTagEmpty(func(field *Field) {
		value, _ := strconv.ParseUint(field.Tag, 10, 64)
		field.Value.SetUint(value)
	})
	fns[reflect.Uint8] = fns[reflect.Uint]
	fns[reflect.Uint16] = fns[reflect.Uint]
	fns[reflect.Uint32] = fns[reflect.Uint]
	fns[reflect.Uint64] = fns[reflect.Uint]

	fns[reflect.Float32] = f.skipIfTagEmpty(func(field *Field) {
		value, _ := strconv.ParseFloat(field.Tag, 64)
		field.Value.SetFloat(value)
	})
	fns[reflect.Float64] = fns[reflect.Float32]

	fns[reflect.String] = f.skipIfTagEmpty(func(field *Field) {
		field.Value.SetString(field.Tag)
	})

	fns[reflect.Struct] = func(field *Field) {
		// Skip set default value for struct if `omit` tag found
		if field.Tag != f.OmitKey {
			f.fillStruct(field)
		}
	}

	fns[reflect.Ptr] = func(field *Field) {
		ptr := field.Value

		// Nil Ptr, re-point to address containing a zero value
		if field.Value.IsNil() {
			ptr = reflect.New(field.Value.Type().Elem())
		}

		f.fillField(&Field{
			Value:  ptr.Elem(),
			Tag:    field.Tag,
			Parent: field.Parent,
		})

		if field.Value.IsNil() && !ptr.Elem().IsZero() {
			field.Value.Set(ptr)
		}
	}

	fns[reflect.Interface] = func(field *Field) {
		// Only set default value for the interface if underlying implementation is of kind struct and not nil
		if !field.Value.IsNil() && GetValueInternalKind(field.Value) == reflect.Struct {
			f.fillField(&Field{
				Value:  field.Value.Elem(),
				Tag:    field.Tag,
				Parent: field.Parent,
			})
		}
	}

	fns[reflect.Slice] = func(field *Field) {
		switch GetValueInternalKind(field.Value) {
		case reflect.Uint8:
			if field.Value.Bytes() == nil {
				field.Value.SetBytes([]byte(field.Tag))
			}
		case reflect.Struct, reflect.Interface:
			for i := 0; i < field.Value.Len(); i++ {
				f.fillField(&Field{
					Value:  field.Value.Index(i),
					Tag:    field.Tag,
					Parent: field,
				})
			}
		default:
			// handle slice of data with eventually primitive type like [1,2,3,4], [[1,2], [3,4]], [{1:2},{3:4}]
			tag, slice := field.Tag, field.Value
			if !regexp.MustCompile(`^\[.*]$`).MatchString(tag) {
				return // invalid default value to set slice
			}

			if tag == "[]" {
				slice.Set(reflect.MakeSlice(slice.Type(), 0, 0)) // empty slice value to set
				return
			}

			values := tokenizeValues(tag[1 : len(tag)-1])
			result := reflect.MakeSlice(slice.Type(), len(values), len(values))
			for i := 0; i < len(values); i++ {
				f.fillField(&Field{
					Value:  result.Index(i),
					Tag:    values[i],
					Parent: field,
				})
			}
			slice.Set(result)
		}
	}

	fns[reflect.Map] = func(field *Field) {
		switch GetValueInternalKind(field.Value) {
		case reflect.Struct, reflect.Interface:
			// Map kind has both unaddressable key and value, we can only directly set the map using SetMapIndex
			for _, mapKay := range field.Value.MapKeys() {
				mapVal := field.Value.MapIndex(mapKay)
				item := &Field{
					Value:  reflect.New(mapVal.Type()).Elem(),
					Tag:    field.Tag,
					Parent: field,
				}
				item.Value.Set(mapVal) // copy the original value before set the rest
				f.fillField(item)
				field.Value.SetMapIndex(mapKay, item.Value)
			}
		default:
			// handle slice of data with eventually primitive type like {1:2, 3:4}, {"arr":[1,2,3]} {1: {1:2}, 2: {3:4}}
			tag, mapField := field.Tag, field.Value
			if !regexp.MustCompile(`^{.*}$`).MatchString(tag) {
				return // invalid default value to set map
			}

			if tag == "{}" {
				mapField.Set(reflect.MakeMapWithSize(mapField.Type(), 0)) // empty map value to set
				return
			}

			keyValues := tokenizeValues(tag[1 : len(tag)-1])
			mapField.Set(reflect.MakeMapWithSize(mapField.Type(), len(keyValues)))
			keyType := mapField.Type().Key()
			valType := mapField.Type().Elem()

			for _, entry := range keyValues {
				keyValue := strings.SplitN(entry, ":", 2)
				if len(keyValue) < 2 {
					continue
				}

				keyField := &Field{
					Value:  reflect.New(keyType).Elem(),
					Tag:    keyValue[0],
					Parent: field,
				}
				f.fillField(keyField)

				valField := &Field{
					Value:  reflect.New(valType).Elem(),
					Tag:    keyValue[1],
					Parent: field,
				}
				f.fillField(valField)

				mapField.SetMapIndex(keyField.Value, valField.Value)
			}
		}
	}
}

func (f *filler) skipIfTagEmpty(fn FillFn) FillFn {
	return func(field *Field) {
		if field.Tag != "" {
			fn(field)
		}
	}
}

func tokenizeValues(expr string) []string {
	var tokens []string

	var count int
	buf := bytes.NewBufferString("")

	for _, b := range []rune(expr) {
		switch b {
		case ',':
			if count == 0 { // end of value
				tokens = append(tokens, buf.String())
				buf.Reset()
				continue
			}
			buf.WriteRune(b)
		case '{', '[':
			count++
			buf.WriteRune(b)
		case '}', ']':
			count--
			buf.WriteRune(b)
		default:
			buf.WriteRune(b)
		}
	}

	if buf.Len() > 0 {
		tokens = append(tokens, buf.String())
	}

	return tokens
}
