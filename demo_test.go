package defaults

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

type Role string

type Admin struct {
	Name string
	Role Role `default:"admin"`
}

func (a *Admin) GetRole() Role {
	return a.Role
}

type User interface {
	GetRole() Role
}

type ExampleFuncsByKind struct {
	Int           int           `default:"1"`                         // Primitive
	IntPtrPtr     **int         `default:"1"`                         // Ptr type
	Role          Role          `default:"DBA"`                       // Alias of Primitive
	Duration      time.Duration `default:"1s"`                        // Duration
	Time          time.Time     `default:"2007-07-07T07:07:07.007Z"`  // Time
	ListOfInt     []int         `default:"[1,2,3,4]"`                 // Slice
	ListOfIntList [][]int       `default:"[[1,2],[3,4]]"`             // 2D Slice
	ListOfIntMap  []map[int]int `default:"[{1:10,2:20},{3:30,4:40}]"` // Slice of Map

	Admin            Admin  // Struct
	AdminPtr         *Admin // Struct Ptr
	AdminOmit        Admin  `default:"omit"` // Struct w Omit
	AdminWithVal     Admin  // Struct w Initial Value
	AdminWithValDive Admin  `default:"dive"` // Struct w Dive
	User             User   // Interface
	UserWithVal      User   // Interface Implementation
	UserWithValDive  User   `default:"dive"` // Interface Implementation w Dive
}

func TestExampleFuncsByKind(t *testing.T) {
	foo := ExampleFuncsByKind{
		AdminWithVal:     Admin{Name: "admin1"},
		AdminWithValDive: Admin{Name: "admin2"},
		UserWithVal:      &Admin{Name: "admin3"},
		UserWithValDive:  &Admin{Name: "admin4"},
	}

	SetDefaults(&foo)

	jsonPrint(foo)
}

type Enum string

type DefaultData struct {
	String string
	Int    int
}

type ExampleFuncsByType struct {
	Enum                   Enum
	EnumWithTag            Enum `default:"tag"`
	EnumWithValueNTag      Enum `default:"tag"`
	DefaultData            DefaultData
	DefaultDataOmit        DefaultData `default:"omit"`
	DefaultDataWithVal     DefaultData
	DefaultDataWithValDive DefaultData `default:"dive"`
}

func TestExampleFuncsByType(t *testing.T) {
	foo := ExampleFuncsByType{
		EnumWithValueNTag:      Enum("value"),
		DefaultDataWithVal:     DefaultData{Int: 1},
		DefaultDataWithValDive: DefaultData{Int: 1},
	}

	RegisterDefaultType(Enum("type"))
	RegisterDefaultType(DefaultData{String: "type", Int: 7})
	SetDefaults(&foo)

	jsonPrint(foo)
}

func jsonPrint(obj interface{}) {
	b, _ := json.MarshalIndent(obj, "", "    ")
	fmt.Println(string(b))
}