package store

import (
	"github.com/ZacharyDuve/apireg/api"
)

type ApiStore interface {
	Add(api.Api) bool
	Contains(api.Api) bool
	All() []api.Api
	Remove(api.Api) bool
}
