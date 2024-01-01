package apireg

import (
	"log"
	"testing"
	"time"

	"github.com/ZacharyDuve/apireg/api"
	"github.com/ZacharyDuve/apireg/environment"
)

func TestMainForMulticastRegistry(t *testing.T) {
	log.Println("Hello")
	r, err := NewRegistry(environment.All)

	failOnErr(err, t)

	err = r.RegisterApi("Something something", &api.Version{Major: 0, Minor: 1, BugFix: 3}, 80)
	failOnErr(err, t)
	err = r.RegisterApi("Somethingelse", &api.Version{Major: 0, Minor: 1, BugFix: 3}, 433)
	failOnErr(err, t)
}

func TestThatTwoRegistriesRegisterEachOther(t *testing.T) {
	reg0, err := NewRegistry(environment.All)
	failOnErr(err, t)

	reg1, err := NewRegistry(environment.All)
	failOnErr(err, t)
	reg0ApiName := "Something"
	reg0ApiVersion := &api.Version{Major: 0, Minor: 1, BugFix: 3}

	reg1ApiName := "Other"
	reg1ApiVersion := &api.Version{Major: 0, Minor: 1, BugFix: 3}

	reg0.RegisterApi(reg0ApiName, reg0ApiVersion, 8080)
	reg1.RegisterApi(reg1ApiName, reg1ApiVersion, 8080)

	time.Sleep(time.Second * 2)
	t.Log("reg0.GetApisByApiName(reg1ApiName)", reg0.GetApisByApiName(reg1ApiName))
	t.Log("reg0.GetApisByApiName(reg1ApiName)", reg0.GetApisByApiName(reg1ApiName)[0].Name())
	t.Log("reg1.GetApisByApiName(reg0ApiName)", reg1.GetApisByApiName(reg0ApiName)[0].Name())
	if reg0.GetApisByApiName(reg1ApiName) == nil || len(reg0.GetApisByApiName(reg1ApiName)) == 0 {
		t.Fail()
	}

	if reg1.GetApisByApiName(reg0ApiName) == nil || len(reg1.GetApisByApiName(reg0ApiName)) == 0 {
		t.Fail()
	}

}

func failOnErr(err error, t *testing.T) {
	if err != nil {
		t.Fail()
	}
}
