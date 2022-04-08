package apireg

import (
	"github.com/ZacharyDuve/apireg/api"
)

type apiRegisterMessageJSON struct {
	ApiName    string       `json:"api-name"`
	ApiVersion *api.Version `json:"api-version"`
	ApiPort    int          `json:"api-port"`
}
