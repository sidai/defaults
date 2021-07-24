package defaults

import (
	. "gopkg.in/check.v1"
	"time"
)

type OptionSuite struct{}

var _ = Suite(&OptionSuite{})

type ExampleCustomTag struct {
	DefaultTag string `default:"string"`
	CustomTag  string `custom:"string"`
}

func (s *OptionSuite) TestCustomDefaultTag(c *C) {
	var foo ExampleCustomTag

	SetDefaults(&foo)

	c.Assert(foo.DefaultTag, Equals, "string")
	c.Assert(foo.CustomTag, Equals, "")

	var bar ExampleCustomTag

	NewFiller(UseDefault(), UseDefaultTag("custom")).SetDefaults(&bar)

	c.Assert(bar.DefaultTag, Equals, "")
	c.Assert(bar.CustomTag, Equals, "string")
}

type ExampleCustomOmitKey struct {
	StructOmit Struct `default:"omit"`
	StructSkip Struct `default:"skip"`
}

func (s *OptionSuite) TestCustomOmitKey(c *C) {
	var foo ExampleCustomOmitKey

	SetDefaults(&foo)

	c.Assert(foo.StructOmit, Equals, Struct{})
	c.Assert(foo.StructSkip, Equals, Struct{String: "string", Integer: 1})

	var bar ExampleCustomOmitKey

	NewFiller(UseDefault(), UseOmitKey("skip")).SetDefaults(&bar)

	c.Assert(bar.StructOmit, Equals, Struct{String: "string", Integer: 1})
	c.Assert(bar.StructSkip, Equals, Struct{})
}

type ExampleCustomDiveKey struct {
	StructDive Struct `default:"dive"`
	StructDig  Struct `default:"dig"`
}

func (s *OptionSuite) TestCustomDiveKey(c *C) {
	structWithVal := Struct{Integer: 7}
	foo := ExampleCustomDiveKey{
		StructDive: structWithVal,
		StructDig:  structWithVal,
	}

	SetDefaults(&foo)

	c.Assert(foo.StructDive, Equals, Struct{String: "string", Integer: 7})
	c.Assert(foo.StructDig, Equals, Struct{Integer: 7})

	bar := ExampleCustomDiveKey{
		StructDive: structWithVal,
		StructDig:  structWithVal,
	}

	NewFiller(UseDefault(), UseDiveKey("dig")).SetDefaults(&bar)

	c.Assert(bar.StructDive, Equals, Struct{Integer: 7})
	c.Assert(bar.StructDig, Equals, Struct{String: "string", Integer: 7})
}

type ExampleCustomLayout struct {
	TimeRFC3339 time.Time `default:"2007-07-07T07:07:07.007Z"`
	TimeRFC822Z time.Time `default:"07 Jul 07 07:07 +0700"`
}

func (s *OptionSuite) TestCustomTimeLayout(c *C) {
	var foo ExampleCustomLayout

	SetDefaults(&foo)

	TimeRFC3339, _ := time.Parse(time.RFC3339, "2007-07-07T07:07:07.007Z")
	c.Assert(foo.TimeRFC3339, DeepEquals, TimeRFC3339)
	c.Assert(foo.TimeRFC822Z, DeepEquals, time.Time{})

	var bar ExampleCustomLayout

	NewFiller(UseDefault(), UseTimeFormat(time.RFC822Z)).SetDefaults(&bar)

	TimeRFC822Z, _ := time.Parse(time.RFC822Z, "07 Jul 07 07:07 +0700")
	c.Assert(bar.TimeRFC3339, DeepEquals, time.Time{})
	c.Assert(bar.TimeRFC822Z, DeepEquals, TimeRFC822Z)
}

type Default string
type DefaultStruct struct {
	Integer int
	String  string
}

type ExampleDefaultType struct {
	Default        Default
	DefaultWithTag Default `default:"string"`
	DefaultPtr     *Default

	Struct              DefaultStruct
	StructOmit          DefaultStruct `default:"omit"`
	StructPtr           *DefaultStruct
	StructWithValue     DefaultStruct
	StructWithValueDive DefaultStruct `default:"dive"`
	StructList          []DefaultStruct
}

func (s *OptionSuite) TestUseDefaultType(c *C) {
	structWithVal := DefaultStruct{Integer: 1}
	foo := ExampleDefaultType{
		StructWithValue:     structWithVal,
		StructWithValueDive: structWithVal,
		StructList:          []DefaultStruct{structWithVal},
	}

	NewFiller(UseDefault(), UseDefaultType(Default("7")), UseDefaultType(DefaultStruct{Integer: 7, String: "7"})).SetDefaults(&foo)

	c.Assert(foo.Default, Equals, Default("7"))
	c.Assert(foo.DefaultWithTag, Equals, Default("string")) // value from tag since it has higher precedence
	c.Assert(*foo.DefaultPtr, Equals, Default("7"))

	c.Assert(foo.Struct, Equals, DefaultStruct{Integer: 7, String: "7"})
	c.Assert(foo.StructOmit, Equals, DefaultStruct{})
	c.Assert(*foo.StructPtr, Equals, DefaultStruct{Integer: 7, String: "7"})

	c.Assert(foo.StructWithValue, Equals, DefaultStruct{Integer: 1, String: ""})
	// default type for structs apply on entirely struct only when the struct is of zero value
	// it does not recursively set default for inner filed so dive tag does not make any difference
	c.Assert(foo.StructWithValueDive, Equals, DefaultStruct{Integer: 1, String: ""})
	c.Assert(foo.StructList[0], Equals, DefaultStruct{Integer: 1, String: ""})
}
