package store

import (
	"git.zmanhobbies.com/software/apireg/api"
)

type ApiStore interface {
	Add(api.Api) bool
	Contains(api.Api) bool
	All() []api.Api
	Remove(api.Api) bool
}
