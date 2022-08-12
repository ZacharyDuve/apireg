package apireg

import (
	"git.zmanhobbies.com/software/apireg/api"
	"git.zmanhobbies.com/software/apireg/environment"
)

type apiRegisterMessageJSON struct {
	ApiName     string                  `json:"api-name"`
	ApiVersion  *api.Version            `json:"api-version"`
	ApiPort     int                     `json:"api-port"`
	SenderId    string                  `json:"sender-id"`
	Environment environment.Environment `json:"env"`
}
