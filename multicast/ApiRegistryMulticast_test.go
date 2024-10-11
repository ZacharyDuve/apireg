package multicast

import (
	"log"
	"testing"
	"time"

	"github.com/ZacharyDuve/apireg"
	"github.com/google/uuid"
)

func TestMainForMulticastRegistry(t *testing.T) {
	log.Println("Hello")
	r, err := NewMulticastRegistry(nil, apireg.All, uuid.New())

	failOnErr(err, t)
	err = r.RegisterApi("Something something", apireg.NewVersion(0, 1, 3), 80)
	failOnErr(err, t)
	err = r.RegisterApi("Somethingelse", apireg.NewVersion(0, 1, 3), 433)
	failOnErr(err, t)
}

func TestThatTwoRegistriesRegisterEachOther(t *testing.T) {
	reg0, err := NewMulticastRegistry(nil, apireg.All, uuid.New())
	failOnErr(err, t)

	reg1, err := NewMulticastRegistry(nil, apireg.All, uuid.New())
	failOnErr(err, t)
	reg0ApiName := "Something"
	reg0ApiVersion := apireg.NewVersion(0, 1, 3)

	reg1ApiName := "Other"
	reg1ApiVersion := apireg.NewVersion(0, 1, 3)

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
