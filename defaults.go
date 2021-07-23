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
	once.Do(func() {
		defaultFiller = newFiller(UseDefault(), UseTimeFormat(time.RFC3339), ParseDuration())
	})

	return defaultFiller
}

func SetDefaultTag(tag string) {
	UseDefaultTag(tag)(defaultFiller)
}

func SetOmitKey(key string) {
	UseOmitKey(key)(defaultFiller)
}

func RegisterDefaultType(defVal interface{}) {
	UseDefaultType(defVal)(defaultFiller)
}

func RegisterTimeLayout(layout string) {
	UseTimeFormat(layout)(defaultFiller)
}
