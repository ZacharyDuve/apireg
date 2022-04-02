package godistributedapiregistry

import (
	"log"
	"testing"
	"time"
)

func TestMainForMulticastRegistry(t *testing.T) {
	log.Println("Hello")
	r, err := NewRegistry()

	if err != nil {
		t.Fail()
	}

	err = r.RegisterApi("Something something", "v1.0.0")
	if err != nil {
		t.Fail()
	}
	err = r.RegisterApi("Somethingelse", "V0.1.3")
	if err != nil {
		t.Fail()
	}

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
