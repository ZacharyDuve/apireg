package apireg

import "time"

type apiRegistration struct {
	api            Api
	timeRegistered time.Time
	lifeSpan       time.Duration
}
