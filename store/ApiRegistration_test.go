package store

import (
	"net"
	"testing"
	"time"

	"github.com/ZacharyDuve/apireg/api"
	"github.com/ZacharyDuve/apireg/environment"
)

func TestThatNewApiRegistrationReturnsErrorIfApiIsNil(t *testing.T) {
	_, err := NewApiRegistration(nil, time.Time{}, 0)

	if err == nil {
		t.Fail()
	}
}

func TestThatNewApiRegistrationReturnsNilForRegIfApiIsNil(t *testing.T) {
	reg, _ := NewApiRegistration(nil, time.Time{}, 0)

	if reg != nil {
		t.Fail()
	}
}

func TestThatNewApiRegistrationReturnsNoErrorIfApiIsNotNil(t *testing.T) {
	_, err := NewApiRegistration(getValidApi(), time.Time{}, 0)

	if err != nil {
		t.Fail()
	}
}

func TestThatNewApiRegistrationReturnsApiRegIfApiIsNotNil(t *testing.T) {
	reg, _ := NewApiRegistration(getValidApi(), time.Time{}, 0)

	if reg == nil {
		t.Fail()
	}
}

func TestThatRegistrationIsExpiredIfTimePassedInLessThanRegTimePlusLife(t *testing.T) {
	now := time.Now()
	life := time.Second * 30
	reg, _ := NewApiRegistration(getValidApi(), now, life)

	if !reg.Expired(now.Add(life).Add(time.Second * 1)) {
		t.Fail()
	}
}

func TestThatRegistrationIsNotExpiredIfTimePassedInMoreThanRegTimePlusLife(t *testing.T) {
	now := time.Now()
	life := time.Second * 30
	reg, _ := NewApiRegistration(getValidApi(), now, life)

	if reg.Expired(now.Add(life).Add(time.Second * -1)) {
		t.Fail()
	}
}

func getValidApi() api.Api {
	api, _ := api.NewApi("someApi", &api.Version{}, environment.All, net.IPv4(192, 168, 0, 3), 8080)
	return api
}
