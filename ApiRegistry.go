package apireg

import (
	"github.com/ZacharyDuve/apireg/api"
	"github.com/ZacharyDuve/apireg/event"
)

type ApiRegistry interface {
	RegisterApi(name string, version *api.Version, port int) error
	GetAvailableApis() []api.Api
	GetApisByApiName(name string) []api.Api
	AddListener(event.RegistrationListener)
	RemoveListener(event.RegistrationListener)
}
