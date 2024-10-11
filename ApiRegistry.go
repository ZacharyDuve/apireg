package apireg

type ApiRegistry interface {
	RegisterApi(name string, version Version, port int) error
	GetAvailableApis() []Api
	GetApisByApiName(name string) []Api
	AddEventListener(RegistrationListener)
	RemoveEventListener(RegistrationListener)
}
