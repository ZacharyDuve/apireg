package store

import (
	"github.com/ZacharyDuve/apireg/api"
)

type ApiRegistrationStore interface {
	//Attempts to add ApiRegistration. If it is already there then nothing occurs
	AddReg(ApiRegistration)
	//Get all of the apis for a given name
	GetAllRegsForName(name string) []ApiRegistration
	//Get every registration that we have stored
	GetAllRegs() []ApiRegistration
	//Remove the registration that we have for a given Api
	RemoveRegForApi(api.Api) error
}
