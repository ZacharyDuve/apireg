package godistributedapiregistry

import "time"

type apiRegisterMessageJSON struct {
	ApiName    string        `json:"api-name"`
	ApiVersion string        `json:"api-version"`
	ApiPort    int           `json:"api-port"`
	LifeSpan   time.Duration `json:"reg-lifespan"`
}
