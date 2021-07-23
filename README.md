defaults
=======
Set default values to structs using [struct tags](http://golang.org/pkg/reflect/#StructTag) or [struct type](https://pkg.go.dev/reflect#Type)

Notice
-------
This repo is inspired by [go-defaults](https://github.com/sidai/go-defaults) and applies the same LICENSE. 
The aforementioned provides basic default value setting for simple data type. 
However, it does not support complex structure like `pointer`, `interface`, `map` or `slice of slice`. 
I created this repo to provide more data type support. 


Supported Data Type 
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
    - e.g. `*int`, `*User*` 
   
Rules
-------
- Unless `default:"omit"` tag found in a struct, always recursively set default to all its zero fields

- Skip setting default for non-zero fields to prevent fields with initial value being reset

- If both FuncByKind or FuncByType defined on an aliased type, FuncByKind is used with higher precedence

Usage
-------
- **Installation**: ```go get github.com/sidai/defaults```
- **Examples**:
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

type Example struct {
	User          User
	Admin         *Admin
	Int           int           `default:"1"`
	Role          Role          `default:"DBA"`
	Duration      time.Duration `default:"1s"`
	Time          time.Time     `default:"2007-07-07T07:07:07.007Z"`
	ListOfInt     []int         `default:"[1,2,3,4]"`
	ListOfIntList [][]int       `default:"[[1,2],[3,4]]"`
	ListOfIntMap  []map[int]int `default:"[{1:10,2:20},{3:30,4:40}]"`
}

...

foo := Example{
    User: &Admin{Name: "john doe"},
}

GetDefaultFiller().SetDefaults(&foo)

foo = {
    "User": {"Name": "john doe", "Role": "admin"},
    "Admin": {"Name": "", "Role": "admin"},
    "Int": 1,
    "Role": "DBA",
    "Duration": 1s,
    "Time": "2007-07-07 07:07:07.007 +0000 UTC",
    "ListOfInt": [1, 2, 3, 4],
    "ListOfIntList": [[1, 2], [3, 4]],
    "ListOfIntMap": [{1: 10, 2: 20}, {3: 30, 4: 40}]
}
```
  
License
------
MIT, see LICENSE