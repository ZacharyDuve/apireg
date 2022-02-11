package godistributedapiregistry

type ApiRegistry interface {
	RegisterApi(name string, version string) error
	GetAvailableApis() []Api
	GetApisByApiName(name string) []Api
}
