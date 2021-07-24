package defaults

import (
	"reflect"
)

const (
	defaultTag = "default"

	diveKey = "dive"
	omitKey = "omit"
)

type FillFn func(field *Field)

type filler struct {
	FuncsByKind map[reflect.Kind]FillFn
	FuncsByType map[reflect.Type]FillFn
	DefaultTag  string
	DiveKey     string
	OmitKey     string
}

type Field struct {
	Value  reflect.Value
	Tag    string
	Parent *Field
}

func newFiller(opts ...Option) *filler {
	f := &filler{
		FuncsByKind: make(map[reflect.Kind]FillFn),
		FuncsByType: make(map[reflect.Type]FillFn),
		DefaultTag:  defaultTag,
		DiveKey:     diveKey,
		OmitKey:     omitKey,
	}

	for _, opt := range opts {
		opt(f)
	}

	return f
}

func (f *filler) SetDefaults(variable interface{}) {
	value := reflect.ValueOf(variable)

	// skip if variable is not a ptr to a struct
	if value.Kind() != reflect.Ptr || GetValueInternalKind(value) != reflect.Struct {
		return
	}

	f.fillStruct(&Field{
		Value:  value.Elem(),
		Tag:    "",
		Parent: nil,
	})
}

func (f *filler) fillStruct(field *Field) {
	structVal, structType := field.Value, field.Value.Type()

	for i := 0; i < structVal.NumField(); i++ {
		fieldVal, fieldType := structVal.Field(i), structType.Field(i)
		// Only fill the filed if the field can be set, i.e. exported
		if fieldVal.CanSet() {
			f.fillField(&Field{
				Value:  fieldVal,
				Tag:    fieldType.Tag.Get(f.DefaultTag),
				Parent: field,
			})
		}
	}
}

func (f *filler) fillField(field *Field) {
	// Fill the field when field should be filled in precedence of Kind (via Tag) > Type (via Type Default)
	if fn, ok := f.FuncsByKind[field.Value.Kind()]; ok && f.shouldFill(field) {
		fn(field)
	}

	if fn, ok := f.FuncsByType[field.Value.Type()]; ok && f.shouldFill(field) {
		fn(field)
	}
}

func (f *filler) shouldFill(field *Field) bool {
	switch GetValueInternalKind(field.Value) {
	case reflect.Struct:
		// always fill struct if dive key found
		if field.Tag == f.DiveKey {
			return true
		}
		// always skip struct filling if omit key found
		if field.Tag == f.OmitKey {
			return false
		}
		// otherwise only fill struct if all its field is of zero value
		return field.Value.IsZero()
	case reflect.Interface:
		// always assume interface should be filled util we find the actual implementation
		return true
	default:
		// only fill zero field for primitive type (or its alias)
		return field.Value.IsZero()
	}
}

// GetValueInternalKind returns the actual underlying kind of field value.
// It repeatedly dives into slice, array, interface and pointer kind data until find a struct,
// unimplemented interface or a primitive data kind
// note that for Map underlying kind of its value is used.
func GetValueInternalKind(field reflect.Value) reflect.Kind {
	if field.Kind() == reflect.Invalid {
		return reflect.Invalid
	}

	typ := field.Type()
	if field.Kind() == reflect.Interface && !field.IsNil() {
		typ = field.Elem().Type()
	}

BEGIN:
	switch typ.Kind() {
	case reflect.Slice, reflect.Ptr, reflect.Array, reflect.Map:
		typ = typ.Elem()
		goto BEGIN
	case reflect.Struct, reflect.Interface:
		return typ.Kind()
	default:
		return typ.Kind()
	}
}

// IndirectType keeps de-referencing the type of the value until a non-ptr type found then returns the type
func IndirectType(field reflect.Value) reflect.Type {
	typ := field.Type()
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	return typ
}
