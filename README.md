defaults [![Test](https://github.com/sidai/defaults/actions/workflows/test.yml/badge.svg)](https://github.com/sidai/defaults/actions/workflows/test.yml) [![GitHub release](https://img.shields.io/github/release/sidai/defaults.svg)](https://github.com/sidai/defaults/releases) [![License](https://img.shields.io/github/license/sidai/defaults.svg)](./LICENSE)
=======
Structures default value filling with support in almost all types of data using [struct tags](http://golang.org/pkg/reflect/#StructTag) or [struct type](https://pkg.go.dev/reflect#Type) <br>

Notice
-------
This repo is inspired by [go-defaults](https://github.com/sidai/go-defaults) and applies the same [LICENSE](https://github.com/sidai/defaults/blob/master/LICENSE). 

The aforementioned repo provides basic default value setting for simple data type. However, 
1. It does not support complex structure like `pointer`, `interface`, `map` or `slice of map`. 
2. It always recursively fill the struct with default value but there is case struct filling should be skipped.
3. The default filler provided is not exported which makes it hard for customization.  

I created this repo to provide more data types support, more flexibility in default value filling for struct and export the function for better customization. 


Supported Data Types
-------
- **Primitive Types:** 
    - `bool`
    - `int`, `int8`, `int16`, `int32`, `int64`
    - `uint`, `uint8`, `uint16`, `uint32`, `uint64`
    - `float32`, `float64`
    - `[]byte`, `string`
    
- **Custom Types:**
    - `time.Duration`, `time.Time`
    - Aliased types. e.g `type UserName string`
    - Self-defined type. e.g. `type User struct {Name string, Age int}`
    
- **Complex Types:**
    - `map`,`slice`, `struct`, `interface`
    
- **Nested Types:**
    - e.g. `map[int][]*User`, `[]map[int]*User`, `map[int]map[int]*User` 
    
- **Pointer Types:**
    - e.g. `*int`, `*User`, `**int`, `**User` 
   
Rules
-------
- Filling rules can be defined both by [Kind](https://pkg.go.dev/reflect#Kind) and [Type](https://pkg.go.dev/reflect#Type) <br>
  If both rules found when filling a field, 
  [FuncsByKind](https://github.com/sidai/defaults/blob/master/filler.go#L17) is used first before
  [FuncsByType](https://github.com/sidai/defaults/blob/master/filler.go#L18) 

- Skip default filling for non-zero fields to prevent fields with initial value being reset

- By default struct is recursively filled only when it is *empty* <br> 
  Use `default:"omit"` to skip struct filling <br>
  Use `default:"dive"` to always apply struct filling even when it is not empty

Usage
-------
- **Installation**: ```go get github.com/sidai/defaults```

- **[FuncsByKind](https://github.com/sidai/defaults/blob/master/filler.go#L17) Examples**:
    ```go
    type Role string
    
    type Admin struct {
        Name string
        Role Role   `default:"admin"`
    }
    
    func (a *Admin) GetRole() Role {
        return a.Role
    }
    
    type User interface {
        GetRole() Role
    }
    
    type ExampleFuncsByKind struct {
        Int              int           `default:"1"`                         // Primitive
        IntPtrPtr        **int         `default:"1"`                         // Ptr type
        Role             Role          `default:"DBA"`                       // Alias of Primitive
        Duration         time.Duration `default:"1s"`                        // Duration
        Time             time.Time     `default:"2007-07-07T07:07:07.007Z"`  // Time
        ListOfInt        []int         `default:"[1,2,3,4]"`                 // Slice
        ListOfIntList    [][]int       `default:"[[1,2],[3,4]]"`             // 2D Slice
        ListOfIntMap     []map[int]int `default:"[{1:10,2:20},{3:30,4:40}]"` // Slice of Map
    
        Admin            Admin                                               // Struct 
        AdminPtr         *Admin                                              // Struct Ptr
        AdminOmit        Admin         `default:"omit"`                      // Struct w Omit
        AdminWithVal     Admin                                               // Struct w Initial Value
        AdminWithValDive Admin         `default:"dive"`                      // Struct w Dive
        User             User                                                // Interface
        UserWithVal      User                                                // Interface w Implementation
        UserWithValDive  User          `default:"dive"`                      // Interface w Implementation & Dive
    }
    
    ...
    
    foo := ExampleFuncsByKind{
   		AdminWithVal:     Admin{Name: "admin1"},
   		AdminWithValDive: Admin{Name: "admin2"},
   		UserWithVal:      &Admin{Name: "admin3"},
   		UserWithValDive:  &Admin{Name: "admin4"}, 
    }
  
    SetDefaults(&foo)
    
    foo = {
        "Int": 1,
        "IntPtrPtr": (**int) 1,
        "Role": "DBA",
        "Duration": 1s,
        "Time": "2007-07-07 07:07:07.007 +0000 UTC",
        "ListOfInt": [1, 2, 3, 4],
        "ListOfIntList": [[1, 2], [3, 4]],
        "ListOfIntMap": [{1: 10, 2: 20}, {3: 30, 4: 40}],
        "Admin": {"Name": "", "Role": "admin"},
        "AdminPtr": (*Admin) {"Name": "", "Role": "admin"},
        "AdminOmit": {"Name": "", "Role": ""},
        "AdminWithVal": {"Name": "admin1", "Role": ""},
        "AdminWithValDive": {"Name": "admin2", "Role": "admin"},
        "User": nil,
        "UserWithVal": (*Admin) {"Name": "admin3", "Role": ""},
        "UserWithValDive": (*Admin) {"Name": "admin4", "Role": "admin"}
    }
    ```
  
- **[FuncsByType](https://github.com/sidai/defaults/blob/master/filler.go#L18)  Examples**:
    ```go
    type Enum string
    
    type DefaultData struct {
    	DefaultString string
    	DefaultInt    int
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
    
    foo := ExampleFuncsByType{
        EnumWithValueNTag:      Enum("value"),
        DefaultDataWithVal:     DefaultData{Int: 1},
        DefaultDataWithValDive: DefaultData{Int: 1},
    }

    RegisterDefaultType(Enum("type"))
    RegisterDefaultType(DefaultData{String: "type", Int: 7})
    SetDefaults(&foo)
  
    ...
  
    foo = {
        "Enum": "type",                                // <= Use FuncsByType
        "EnumWithTag": "tag",                          // <= Use tag as FuncsByKind has higher precedence
        "EnumWithValueNTag": "value",                  // <= No filling applied as value is not empty
        "DefaultData": {String: "type", Int: 7},       
        "DefaultDataOmit": {String: "", Int: 0},       // <= Omit tag works for FuncsByType
        "DefaultDataWithVal": {String: "", Int: 1},    // <= FuncsByType skip filling when value is not empty
        "DefaultDataWithValDive": {String: "", Int: 1} // <= FuncsByType ignores dive tag as it works on the extra type only
    }
    ```
    
- More Examples [*HERE*](https://github.com/sidai/defaults/blob/master/filler_test.go)
