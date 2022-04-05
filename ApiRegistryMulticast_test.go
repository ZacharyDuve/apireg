package godistributedapiregistry

import (
	"log"
	"testing"
	"time"
)

func TestMainForMulticastRegistry(t *testing.T) {
	log.Println("Hello")
	r, err := NewRegistry()

	failOnErr(err, t)

	err = r.RegisterApi("Something something", "v1.0.0")
	failOnErr(err, t)
	err = r.RegisterApi("Somethingelse", "V0.1.3")
	failOnErr(err, t)

	time.Sleep(time.Second * 1)
	availableApis := r.GetAvailableApis()
	t.Log("Available Apis", len(availableApis))
	for _, curApi := range availableApis {
		t.Log(curApi.Name(), curApi.Version(), curApi.HostIP(), curApi.HostPort())
	}
	if len(availableApis) != 2 {
		t.Fail()
	}
}

func TestThatTwoRegistriesRegisterEachOther(t *testing.T) {
	reg0, err := NewRegistry()
	failOnErr(err, t)

	reg1, err := NewRegistry()
	failOnErr(err, t)
	reg0ApiName := "Something"
	reg0ApiVersion := "V0.0.1"

	reg1ApiName := "Other"
	reg1ApiVersion := "V199138831231"

	reg0.RegisterApi(reg0ApiName, reg0ApiVersion)
	reg1.RegisterApi(reg1ApiName, reg1ApiVersion)

	time.Sleep(time.Second * 2)
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
