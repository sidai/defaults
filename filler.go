package defaults

import (
	"reflect"
)

const (
	defaultTag = "default"
	omitKey    = "omit"
)

type FillFn func(field *Field)

type filler struct {
	FuncByKind map[reflect.Kind]FillFn
	FuncByType map[reflect.Type]FillFn
	OmitKey    string
	DefaultTag string
}

type Field struct {
	Value  reflect.Value
	Tag    string
	Parent *Field
}

func newFiller(opts ...Option) *filler {
	f := &filler{
		FuncByKind: make(map[reflect.Kind]FillFn),
		FuncByType: make(map[reflect.Type]FillFn),
		DefaultTag: defaultTag,
		OmitKey:    omitKey,
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
		Value:  Indirect(value),
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
	// Fill the field when field is empty in precedence of Kind (via Tag) > Type (via Type Default)
	if fn, ok := f.FuncByKind[field.Value.Kind()]; ok && f.isFieldEmpty(field) {
		fn(field)
	}

	if fn, ok := f.FuncByType[Indirect(field.Value).Type()]; ok && f.isFieldEmpty(field) {
		fn(field)
	}
}

func (f *filler) isFieldEmpty(field *Field) bool {
	switch GetValueInternalKind(field.Value) {
	case reflect.Struct, reflect.Interface:
		// always assume the structs/interface is empty and can be filled
		// the actually struct filling logic should take care of the rest
		return true
	default:
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
		fallthrough
	default:
		return typ.Kind()
	}
}

// Indirect keeps de-referencing the value until a non-ptr value found
func Indirect(field reflect.Value) reflect.Value {
	if field.Kind() != reflect.Ptr {
		return field
	}

	elem := field.Elem()
	for elem.Kind() == reflect.Ptr {
		elem = elem.Elem()
	}

	return elem
}
