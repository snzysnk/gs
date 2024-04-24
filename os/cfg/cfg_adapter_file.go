package cfg

import (
	"github.com/snzysnk/gs/v2/container/sarray"
	"github.com/snzysnk/gs/v2/container/smap"
)

type AdapterFile struct {
	defaultName   string
	searchPaths   sarray.IStrArray
	jsonMap       smap.IStrAnyMap
	violenceCheck bool
}
