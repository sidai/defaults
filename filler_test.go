package defaults

import (
	. "gopkg.in/check.v1"
	"reflect"
	"testing"
	"time"
)

func Test(t *testing.T) { TestingT(t) }

type FillerSuite struct{}

var _ = Suite(&FillerSuite{})

type IntegerVal int
type StringVal string

type Interface interface {
	StringVal() string
	IntegerVal() int
}

type Struct struct {
	String  string `default:"string"`
	Integer int    `default:"1"`
}

func (s *Struct) IntegerVal() int {
	return s.Integer
}

func (s *Struct) StringVal() string {
	return s.String
}

type ExampleBasic struct {
	BoolNoTag   bool
	BoolInvalid bool  `default:"invalid"`
	Bool        bool  `default:"true"`
	BoolPtr     *bool `default:"true"`

	IntNoTag   int
	IntInvalid int        `default:"invalid"`
	Int        int        `default:"1"`
	Int8       int8       `default:"8"`
	Int16      int16      `default:"16"`
	Int32      int32      `default:"32"`
	Int64      int64      `default:"64"`
	IntegerVal IntegerVal `default:"1"`

	UintNoTag   uint
	UintInvalid uint   `default:"invalid"`
	Uint        uint   `default:"1"`
	Uint8       uint8  `default:"8"`
	Uint16      uint16 `default:"16"`
	Uint32      uint32 `default:"32"`
	Uint64      uint64 `default:"64"`

	Float32NoTag   float32
	Float32Invalid float32 `default:"invalid"`
	Float32        float32 `default:"0.32"`
	Float64        float64 `default:"0.64"`

	StringNoTag string
	String      string    `default:"string"`
	StringVal   StringVal `default:"string"`

	BytesNoTag []byte
	Bytes      []byte `default:"bytes"`

	DurationNoTag   time.Duration
	DurationInvalid time.Duration `default:"invalid"`
	Duration        time.Duration `default:"1s"`
	TimeNoTag       time.Time
	TimeInvalid     time.Time `default:"invalid"`
	Time            time.Time `default:"2000-00-00T00:00:00.000Z"`

	IntListNoTag   []int
	IntListInvalid []int      `default:"invalid"`
	IntListEmpty   []int      `default:"[]"`
	IntList        []int      `default:"[1,2,3,4]"`
	Int2DList      [][]int    `default:"[[1,2],[3,4]]"`
	StringList     []string   `default:"[a,b,c,d]"`
	String2DList   [][]string `default:"[[a,b],[c,d]]"`

	MapNoTag        map[int]string
	MapInvalid      map[int]string         `default:"invalid"`
	MapEmpty        map[int]string         `default:"{}"`
	Map             map[int]string         `default:"{1:a,2:b}"`
	MapDuplicateKey map[int]string         `default:"{1:a,1:b}"`
	MapList         []map[int]string       `default:"[{1:a,2:b},{3:c,4:d}]"`
	MapListValue    map[int][]string       `default:"{1:[a,b],2:[c,d]}"`
	MapMapValue     map[int]map[int]string `default:"{1:{1:a,2:b},2:{3:c,4:d}}"`

	Interface         Interface
	InterfaceOmit     Interface `default:"omit"`
	InterfaceList     []Interface
	InterfaceListOmit []Interface `default:"omit"`

	Struct            Struct
	StructOmit        Struct `default:"omit"`
	StructPtr         *Struct
	StructList        []Struct
	StructListOmit    []Struct `default:"omit"`
	StructMap         map[int]Struct
	StructMapOmit     map[int]Struct `default:"omit"`
	StructListMap     map[int][]Struct
	StructListMapOmit map[int][]Struct `default:"omit"`
	StructMapList     []map[int]Struct
	StructMapListOmit []map[int]Struct `default:"omit"`
}

func (s *FillerSuite) TestSetDefaults(c *C) {
	var foo ExampleBasic

	SetDefaults(&foo)

	c.Assert(foo.BoolNoTag, Equals, false)
	c.Assert(foo.BoolInvalid, Equals, false)
	c.Assert(foo.Bool, Equals, true)
	c.Assert(*foo.BoolPtr, Equals, true)

	c.Assert(foo.IntNoTag, Equals, 0)
	c.Assert(foo.IntInvalid, Equals, 0)
	c.Assert(foo.Int, Equals, 1)
	c.Assert(foo.Int8, Equals, int8(8))
	c.Assert(foo.Int16, Equals, int16(16))
	c.Assert(foo.Int32, Equals, int32(32))
	c.Assert(foo.Int64, Equals, int64(64))
	c.Assert(foo.IntegerVal, Equals, IntegerVal(1))

	c.Assert(foo.UintNoTag, Equals, uint(0))
	c.Assert(foo.UintInvalid, Equals, uint(0))
	c.Assert(foo.Uint, Equals, uint(1))
	c.Assert(foo.Uint8, Equals, uint8(8))
	c.Assert(foo.Uint16, Equals, uint16(16))
	c.Assert(foo.Uint32, Equals, uint32(32))
	c.Assert(foo.Uint64, Equals, uint64(64))

	c.Assert(foo.Float32NoTag, Equals, float32(0))
	c.Assert(foo.Float32Invalid, Equals, float32(0))
	c.Assert(foo.Float32, Equals, float32(0.32))
	c.Assert(foo.Float64, Equals, 0.64)

	c.Assert(foo.StringNoTag, Equals, "")
	c.Assert(foo.String, Equals, "string")
	c.Assert(foo.StringVal, Equals, StringVal("string"))

	c.Assert(string(foo.BytesNoTag), Equals, "")
	c.Assert(string(foo.Bytes), Equals, "bytes")

	expectedTime, _ := time.Parse(time.RFC3339, "2000-00-00T00:00:00.000Z")
	c.Assert(foo.DurationNoTag, DeepEquals, time.Duration(0))
	c.Assert(foo.DurationInvalid, DeepEquals, time.Duration(0))
	c.Assert(foo.Duration, DeepEquals, time.Second)
	c.Assert(foo.TimeNoTag, DeepEquals, time.Time{})
	c.Assert(foo.TimeInvalid, DeepEquals, time.Time{})
	c.Assert(foo.Time, DeepEquals, expectedTime)

	c.Assert(foo.IntListNoTag, IsNil)
	c.Assert(foo.IntListInvalid, IsNil)
	c.Assert(foo.IntListEmpty, DeepEquals, []int{})
	c.Assert(foo.IntList, DeepEquals, []int{1, 2, 3, 4})
	c.Assert(foo.Int2DList, DeepEquals, [][]int{{1, 2}, {3, 4}})
	c.Assert(foo.StringList, DeepEquals, []string{"a", "b", "c", "d"})
	c.Assert(foo.String2DList, DeepEquals, [][]string{{"a", "b"}, {"c", "d"}})

	c.Assert(foo.MapNoTag, IsNil)
	c.Assert(foo.MapInvalid, IsNil)
	c.Assert(foo.MapEmpty, DeepEquals, map[int]string{})
	c.Assert(foo.Map, DeepEquals, map[int]string{1: "a", 2: "b"})
	c.Assert(foo.MapDuplicateKey, DeepEquals, map[int]string{1: "b"})
	c.Assert(foo.MapList, DeepEquals, []map[int]string{{1: "a", 2: "b"}, {3: "c", 4: "d"}})
	c.Assert(foo.MapListValue, DeepEquals, map[int][]string{1: {"a", "b"}, 2: {"c", "d"}})
	c.Assert(foo.MapMapValue, DeepEquals, map[int]map[int]string{1: {1: "a", 2: "b"}, 2: {3: "c", 4: "d"}})

	c.Assert(foo.Interface, IsNil)
	c.Assert(foo.InterfaceOmit, IsNil)
	c.Assert(foo.InterfaceList, IsNil)
	c.Assert(foo.InterfaceListOmit, IsNil)

	defaultStruct := Struct{String: "string", Integer: 1}
	c.Assert(foo.Struct, Equals, defaultStruct)
	c.Assert(foo.StructOmit, Equals, Struct{})
	c.Assert(*foo.StructPtr, Equals, defaultStruct)
	c.Assert(foo.StructList, IsNil)
	c.Assert(foo.StructListOmit, IsNil)
	c.Assert(foo.StructMap, IsNil)
	c.Assert(foo.StructMapOmit, IsNil)
	c.Assert(foo.StructListMap, IsNil)
	c.Assert(foo.StructListMapOmit, IsNil)
	c.Assert(foo.StructMapList, IsNil)
	c.Assert(foo.StructMapListOmit, IsNil)
}

func (s *FillerSuite) TestSetDefaultsWithValue(c *C) {
	expectedTime, _ := time.Parse(time.RFC3339, "2007-07-07T07:07:07.007Z")
	foo := &ExampleBasic{
		Int:               7,
		Uint:              7,
		Float64:           0.7,
		String:            "7",
		Bytes:             []byte("7"),
		Duration:          7 * time.Second,
		Time:              expectedTime,
		IntList:           []int{7, 7, 7, 7},
		String2DList:      [][]string{{"7", "7"}, {"7", "7"}},
		Map:               map[int]string{7: "7"},
		MapList:           []map[int]string{{7: "7"}},
		MapListValue:      map[int][]string{7: {"7"}},
		MapMapValue:       map[int]map[int]string{7: {7: "7"}},
		Interface:         &Struct{Integer: 7},
		InterfaceOmit:     &Struct{Integer: 7},
		InterfaceList:     []Interface{&Struct{Integer: 7}},
		InterfaceListOmit: []Interface{&Struct{Integer: 7}},
		Struct:            Struct{Integer: 7},
		StructOmit:        Struct{Integer: 7},
		StructPtr:         &Struct{Integer: 7},
		StructList:        []Struct{{Integer: 7}},
		StructListOmit:    []Struct{{Integer: 7}},
		StructListMap:     map[int][]Struct{7: {{Integer: 7}}},
		StructListMapOmit: map[int][]Struct{7: {{Integer: 7}}},
	}

	SetDefaults(foo)

	expectedStruct := Struct{String: "string", Integer: 7}
	expectedOmitStruct := Struct{String: "", Integer: 7}

	c.Assert(foo.Int, Equals, 7)
	c.Assert(foo.Uint, Equals, uint(7))
	c.Assert(foo.Float64, Equals, 0.7)
	c.Assert(foo.String, Equals, "7")
	c.Assert(string(foo.Bytes), Equals, "7")
	c.Assert(foo.Duration, Equals, 7*time.Second)
	c.Assert(foo.Time, Equals, expectedTime)
	c.Assert(foo.IntList, DeepEquals, []int{7, 7, 7, 7})
	c.Assert(foo.String2DList, DeepEquals, [][]string{{"7", "7"}, {"7", "7"}})
	c.Assert(foo.Map, DeepEquals, map[int]string{7: "7"})
	c.Assert(foo.MapList, DeepEquals, []map[int]string{{7: "7"}})
	c.Assert(foo.MapListValue, DeepEquals, map[int][]string{7: {"7"}})
	c.Assert(foo.MapMapValue, DeepEquals, map[int]map[int]string{7: {7: "7"}})
	c.Assert(foo.Interface.IntegerVal(), Equals, 7)
	c.Assert(foo.Interface.StringVal(), Equals, "string")
	c.Assert(foo.InterfaceOmit.IntegerVal(), Equals, 7)
	c.Assert(foo.InterfaceOmit.StringVal(), Equals, "")
	c.Assert(foo.InterfaceList[0].IntegerVal(), Equals, 7)
	c.Assert(foo.InterfaceList[0].StringVal(), Equals, "string")
	c.Assert(foo.InterfaceListOmit[0].IntegerVal(), Equals, 7)
	c.Assert(foo.InterfaceListOmit[0].StringVal(), Equals, "")
	c.Assert(foo.Struct, Equals, expectedStruct)
	c.Assert(foo.StructOmit, Equals, expectedOmitStruct)
	c.Assert(*foo.StructPtr, Equals, expectedStruct)
	c.Assert(foo.StructList[0], Equals, expectedStruct)
	c.Assert(foo.StructListOmit[0], Equals, expectedOmitStruct)
	c.Assert(foo.StructListMap[7][0], Equals, expectedStruct)
	c.Assert(foo.StructListMapOmit[7][0], Equals, expectedOmitStruct)
}

func (s *FillerSuite) TestGetValueInternalKind(c *C) {
	fn := func(field interface{}) reflect.Kind {
		return GetValueInternalKind(reflect.ValueOf(field))
	}

	type String string

	c.Assert(fn(true), Equals, reflect.Bool)
	c.Assert(fn(new(bool)), Equals, reflect.Bool)
	c.Assert(fn([1]bool{false}), Equals, reflect.Bool)
	c.Assert(fn([]bool{false}), Equals, reflect.Bool)
	c.Assert(fn([]*bool{new(bool)}), Equals, reflect.Bool)
	c.Assert(fn(1), Equals, reflect.Int)
	c.Assert(fn(new(int)), Equals, reflect.Int)
	c.Assert(fn([1]int{1}), Equals, reflect.Int)
	c.Assert(fn([]int{1}), Equals, reflect.Int)
	c.Assert(fn([]*int{new(int)}), Equals, reflect.Int)
	c.Assert(fn(int8(1)), Equals, reflect.Int8)
	c.Assert(fn(int16(1)), Equals, reflect.Int16)
	c.Assert(fn(int32(1)), Equals, reflect.Int32)
	c.Assert(fn(int64(1)), Equals, reflect.Int64)
	c.Assert(fn(uint(1)), Equals, reflect.Uint)
	c.Assert(fn(new(uint)), Equals, reflect.Uint)
	c.Assert(fn([1]uint{1}), Equals, reflect.Uint)
	c.Assert(fn([]uint{1}), Equals, reflect.Uint)
	c.Assert(fn([]*uint{new(uint)}), Equals, reflect.Uint)
	c.Assert(fn(uint8(1)), Equals, reflect.Uint8)
	c.Assert(fn(uint16(1)), Equals, reflect.Uint16)
	c.Assert(fn(uint32(1)), Equals, reflect.Uint32)
	c.Assert(fn(uint64(1)), Equals, reflect.Uint64)
	c.Assert(fn(uintptr(1)), Equals, reflect.Uintptr)
	c.Assert(fn(float32(1)), Equals, reflect.Float32)
	c.Assert(fn(new(float32)), Equals, reflect.Float32)
	c.Assert(fn([1]float32{1}), Equals, reflect.Float32)
	c.Assert(fn([]float32{1}), Equals, reflect.Float32)
	c.Assert(fn([]*float32{new(float32)}), Equals, reflect.Float32)
	c.Assert(fn(float64(1)), Equals, reflect.Float64)
	c.Assert(fn("A"), Equals, reflect.String)
	c.Assert(fn(new(string)), Equals, reflect.String)
	c.Assert(fn([1]string{"A"}), Equals, reflect.String)
	c.Assert(fn([]string{"A"}), Equals, reflect.String)
	c.Assert(fn([]*string{new(string)}), Equals, reflect.String)
	c.Assert(fn(String("A")), Equals, reflect.String)
	c.Assert(fn(new(String)), Equals, reflect.String)
	c.Assert(fn([1]String{"A"}), Equals, reflect.String)
	c.Assert(fn([]String{"A"}), Equals, reflect.String)
	c.Assert(fn([]*String{new(String)}), Equals, reflect.String)
	c.Assert(fn(map[int]string{1: "A"}), Equals, reflect.String)
	c.Assert(fn(map[string]int{"A": 1}), Equals, reflect.Int)
	c.Assert(fn(Struct{"A", 1}), Equals, reflect.Struct)
	c.Assert(fn(&Struct{"A", 1}), Equals, reflect.Struct)
	c.Assert(fn(Interface(&Struct{"A", 1})), Equals, reflect.Struct)
	c.Assert(fn(Interface(nil)), Equals, reflect.Invalid)
	c.Assert(fn([]Interface{}), Equals, reflect.Interface)
	c.Assert(fn([]*Interface{}), Equals, reflect.Interface)
	c.Assert(fn([]Interface{&Struct{"A", 1}}), Equals, reflect.Interface)
	c.Assert(fn([1]int{1}), Equals, reflect.Int)
	c.Assert(fn([1]*int{new(int)}), Equals, reflect.Int)
	c.Assert(fn([1][1]*int{{new(int)}}), Equals, reflect.Int)
	c.Assert(fn([1][]*int{{new(int)}}), Equals, reflect.Int)
	c.Assert(fn([1]Struct{{"A", 1}}), Equals, reflect.Struct)
	c.Assert(fn([1]*Struct{{"A", 1}}), Equals, reflect.Struct)
	c.Assert(fn([1][1]*Struct{{{"A", 1}}}), Equals, reflect.Struct)
	c.Assert(fn([1][]*Struct{{{"A", 1}}}), Equals, reflect.Struct)
	c.Assert(fn([1]Interface{&Struct{"A", 1}}), Equals, reflect.Interface)
	c.Assert(fn([1][1]Interface{{&Struct{"A", 1}}}), Equals, reflect.Interface)
	c.Assert(fn([1][]Interface{{&Struct{"A", 1}}}), Equals, reflect.Interface)
	c.Assert(fn([1]Interface{nil}), Equals, reflect.Interface)
	c.Assert(fn([1][1]Interface{{nil}}), Equals, reflect.Interface)
	c.Assert(fn([1][]Interface{{nil}}), Equals, reflect.Interface)
}
