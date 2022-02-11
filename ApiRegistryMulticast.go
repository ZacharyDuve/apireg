package godistributedapiregistry

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const (
	MulticastGroupIP             string = "224.0.0.78"
	MulticastGroupPort           int    = 5324
	RegistrationMessageSizeBytes int    = 1400
)

type multicastApiRegistry struct {
	mConn      *net.UDPConn
	mapRWMutex *sync.RWMutex
	apiRegs    map[string][]*apiRegistration
}

func NewRegistry() (ApiRegistry, error) {
	r := &multicastApiRegistry{}
	r.mapRWMutex = &sync.RWMutex{}
	r.apiRegs = make(map[string][]*apiRegistration)

	mC, err := net.ListenMulticastUDP("udp", nil, &net.UDPAddr{IP: net.ParseIP(MulticastGroupIP), Port: MulticastGroupPort})

	if err != nil {
		return nil, err
	}
	r.mConn = mC

	go r.listenMutlicast()
	return r, nil
}

func (this *multicastApiRegistry) RegisterApi(name string, version string) error {
	if name == "" {
		return errors.New("name was empty and name is a required parameter")
	}
	if version == "" {
		return errors.New("version was empty and version is a required parameter")
	}

	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: net.ParseIP(MulticastGroupIP), Port: MulticastGroupPort})

	if err != nil {
		return err
	}

	message := &apiRegisterMessageJSON{ApiName: name, ApiVersion: version}

	dataOut := bytes.NewBuffer(make([]byte, RegistrationMessageSizeBytes))

	err = json.NewEncoder(dataOut).Encode(message)

	if err != nil {
		return err
	}

	if dataOut.Len() > RegistrationMessageSizeBytes {
		return errors.New(fmt.Sprint("Message size for", name, version, "exceeds max length of", RegistrationMessageSizeBytes, "bytes"))
	}

	_, err = conn.Write(dataOut.Bytes())
	return err
}

func (this *multicastApiRegistry) GetAvailableApis() []Api {
	this.mapRWMutex.RLock()
	allApis := make([]Api, 0)

	for curApiName := range this.apiRegs {
		allApis = append(allApis, this.GetApisByApiName(curApiName)...)
	}

	this.mapRWMutex.RUnlock()

	return allApis
}

func (this *multicastApiRegistry) GetApisByApiName(name string) []Api {
	this.mapRWMutex.RLock()
	regs := this.apiRegs[name]
	apis := make([]Api, len(regs))
	for i, curReg := range regs {
		apis[i] = curReg.api
	}
	this.mapRWMutex.RUnlock()

	return apis
}

func (this *multicastApiRegistry) listenMutlicast() {
	readBuff := make([]byte, RegistrationMessageSizeBytes)
	for {
		nRead, rAddr, err := this.mConn.ReadFromUDP(readBuff)
		if err != nil {
			log.Println("Error during multicast read", err)
		} else {
			message := &apiRegisterMessageJSON{}
			err = json.NewDecoder(bytes.NewReader(readBuff[0:nRead])).Decode(message)
			if err != nil {
				log.Println("Error decoding multicast json", err)
			} else {
				api := &apiImpl{name: message.ApiName, version: message.ApiVersion, remoteIP: rAddr.IP, remotePort: rAddr.Port}
				log.Println("Registered", api)
				this.updateApis(api)
			}
		}
	}
}

func (this *multicastApiRegistry) updateApis(api Api) {
	this.mapRWMutex.Lock()
	apisForName, contains := this.apiRegs[api.Name()]
	if !contains {
		apisForName = []*apiRegistration{{api: api, timeRegistered: time.Now()}}
		this.apiRegs[api.Name()] = apisForName
	} else {
		matchReg := getRegMatch(api, apisForName)

		if matchReg != nil {
			matchReg.timeRegistered = time.Now()
		} else {
			matchReg = &apiRegistration{api: api, timeRegistered: time.Now()}
		}
	}
	this.mapRWMutex.Unlock()
}

func getRegMatch(api Api, apis []*apiRegistration) *apiRegistration {
	for _, curReg := range apis {
		if api.Name() == curReg.api.Name() &&
			api.Version() == curReg.api.Version() &&
			api.HostIP().String() == curReg.api.HostIP().String() &&
			api.HostPort() == curReg.api.HostPort() {
			return curReg
		}
	}
	return nil
}
