package godistributedapiregistry

import "time"

type apiRegisterMessageJSON struct {
	ApiName    string        `json:"api-name"`
	ApiVersion string        `json:"api-version"`
	LifeSpan   time.Duration `json:"reg-lifespan"`
}
