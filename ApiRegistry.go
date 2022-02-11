package godistributedapiregistry

type Api interface {
	Name() string
	Version() ApiVersion
}

type ApiVersion interface {
	Major() int
	Minor() int
	Patch() int
}

type ApiRegistry interface {
	RegisterApi(Api) error
	UnRegisterApi(Api) error
	GetAvailableApis() []Api
}
