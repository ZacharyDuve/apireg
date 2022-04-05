package godistributedapiregistry

type ApiRegistry interface {
	RegisterApi(name string, version string, port int) error
	GetAvailableApis() []Api
	GetApisByApiName(name string) []Api
}
