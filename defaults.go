package defaults

import (
	"sync"
	"time"
)

var (
	defaultFiller *filler
	once          sync.Once
)

func SetDefaults(variable interface{}) {
	GetDefaultFiller().SetDefaults(variable)
}

type Filler interface {
	SetDefaults(variable interface{})
}

func NewFiller(opts ...Option) Filler {
	return newFiller(opts...)
}

func GetDefaultFiller() Filler {
	initDefaultFiller()
	return defaultFiller
}

func SetDefaultTag(tag string) {
	initDefaultFiller()
	UseDefaultTag(tag)(defaultFiller)
}

func SetOmitKey(key string) {
	initDefaultFiller()
	UseOmitKey(key)(defaultFiller)
}

func SetDiveKey(key string) {
	initDefaultFiller()
	UseDiveKey(key)(defaultFiller)
}

func RegisterDefaultType(defVal interface{}) {
	initDefaultFiller()
	UseDefaultType(defVal)(defaultFiller)
}

func RegisterTimeLayout(layout string) {
	initDefaultFiller()
	UseTimeFormat(layout)(defaultFiller)
}

func initDefaultFiller() {
	once.Do(func() {
		defaultFiller = newFiller(UseDefault(), UseTimeFormat(time.RFC3339), ParseDuration())
	})
}