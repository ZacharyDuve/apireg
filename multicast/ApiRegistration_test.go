package multicast

import (
	"net"
	"testing"
	"time"

	"github.com/ZacharyDuve/apireg"
)

func TestThatNewApiRegistrationReturnsErrorIfApiIsNil(t *testing.T) {
	_, err := newApiRegistration(nil, time.Time{}, 0)

	if err == nil {
		t.Fail()
	}
}

func TestThatNewApiRegistrationReturnsNilForRegIfApiIsNil(t *testing.T) {
	reg, _ := newApiRegistration(nil, time.Time{}, 0)

	if reg != nil {
		t.Fail()
	}
}

func TestThatNewApiRegistrationReturnsNoErrorIfApiIsNotNil(t *testing.T) {
	_, err := newApiRegistration(getValidApi(), time.Time{}, 0)

	if err != nil {
		t.Fail()
	}
}

func TestThatNewApiRegistrationReturnsApiRegIfApiIsNotNil(t *testing.T) {
	reg, _ := newApiRegistration(getValidApi(), time.Time{}, 0)

	if reg == nil {
		t.Fail()
	}
}

func TestThatRegistrationIsExpiredIfTimePassedInLessThanRegTimePlusLife(t *testing.T) {
	now := time.Now()
	life := time.Second * 30
	reg, _ := newApiRegistration(getValidApi(), now, life)

	if !reg.Expired(now.Add(life).Add(time.Second * 1)) {
		t.Fail()
	}
}

func TestThatRegistrationIsNotExpiredIfTimePassedInMoreThanRegTimePlusLife(t *testing.T) {
	now := time.Now()
	life := time.Second * 30
	reg, _ := newApiRegistration(getValidApi(), now, life)

	if reg.Expired(now.Add(life).Add(time.Second * -1)) {
		t.Fail()
	}
}

func getValidApi() apireg.Api {
	api, _ := apireg.NewApi("someApi", apireg.NewVersion(0, 0, 0), apireg.All, net.IPv4(192, 168, 0, 3), 8080)
	return api
}
