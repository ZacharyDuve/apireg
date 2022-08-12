package apireg

import (
	"git.zmanhobbies.com/software/apireg/api"
	"git.zmanhobbies.com/software/apireg/event"
)

type ApiRegistry interface {
	RegisterApi(name string, version *api.Version, port int) error
	GetAvailableApis() []api.Api
	GetApisByApiName(name string) []api.Api
	AddListener(event.RegistrationListener)
	RemoveListener(event.RegistrationListener)
}
