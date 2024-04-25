package cfg

import (
	"context"
	"github.com/snzysnk/gs/v2/container/sarray"
	"github.com/snzysnk/gs/v2/container/smap"
)

var _ Adapter = (*AdapterFile)(nil)

var (
	customConfigContentMap = smap.NewStrStrMap(true) // Customized configuration content.
)

type AdapterFileOptions struct {
	configFileName string
}

func NewDefaultAdapterFileConfig() AdapterFileOptions {
	return AdapterFileOptions{
		configFileName: "default",
	}
}

type AdapterFile struct {
	defaultName   string
	searchPaths   sarray.IStrArray
	jsonMap       smap.IStrAnyMap
	violenceCheck bool
}

func (a *AdapterFile) SetPath(path string) error {
	//TODO implement me
	panic("implement me")
}

func (a *AdapterFile) Available(ctx context.Context, resource ...string) (ok bool) {
	//TODO implement me
	panic("implement me")
}

func (a *AdapterFile) Get(ctx context.Context, pattern string) (value interface{}, err error) {
	//TODO implement me
	panic("implement me")
}

func (a *AdapterFile) Data(ctx context.Context) (data map[string]interface{}, err error) {
	//TODO implement me
	panic("implement me")
}

func NewAdapterFile(config AdapterFileOptions) (adapterFile *AdapterFile) {
	adapterFile = &AdapterFile{
		defaultName: config.configFileName,
		searchPaths: sarray.NewStrArray(true),
		jsonMap:     smap.NewStrAnyMap(true),
	}
	return adapterFile
}

// to do
func (a *AdapterFile) getJson(fileName string) (err error) {
	a.jsonMap.GetOrSet(fileName, func() {

	})
	return nil
}
