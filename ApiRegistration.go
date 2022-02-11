package godistributedapiregistry

import "time"

type apiRegistration struct {
	api            Api
	timeRegistered time.Time
}
