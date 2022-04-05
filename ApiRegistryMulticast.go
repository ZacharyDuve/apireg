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
	MulticastGroupIP             string        = "224.0.0.78"
	MulticastGroupPort           int           = 5324
	RegistrationMessageSizeBytes int           = 1400
	RegistrationLifeSpan         time.Duration = time.Second * 90
	RegistrationUpdateInterval   time.Duration = time.Second * 30
)

type multicastApiRegistry struct {
	mConn            *net.UDPConn
	apisRWMutex      *sync.RWMutex
	apiRegs          map[string][]*apiRegistration
	ownedRegs        map[string]*ownedApi
	ownedRegsRWMutex *sync.RWMutex
}

func NewRegistry() (ApiRegistry, error) {
	r := &multicastApiRegistry{}
	r.apisRWMutex = &sync.RWMutex{}
	r.ownedRegsRWMutex = &sync.RWMutex{}
	r.apiRegs = make(map[string][]*apiRegistration)
	r.ownedRegs = make(map[string]*ownedApi)

	mC, err := net.ListenMulticastUDP("udp", nil, &net.UDPAddr{IP: net.ParseIP(MulticastGroupIP), Port: MulticastGroupPort})

	if err != nil {
		return nil, err
	}
	r.mConn = mC

	go r.listenMutlicast()
	go r.cleanupExpiredRegLoop()
	go r.resendOwnedRegistrationsLoop()
	return r, nil
}

func (this *multicastApiRegistry) ownsApi(name string) bool {
	this.ownedRegsRWMutex.RLock()
	_, contains := this.ownedRegs[name]
	this.ownedRegsRWMutex.RUnlock()

	return contains
}

func (this *multicastApiRegistry) RegisterApi(name string, version string) error {
	if name == "" {
		return errors.New("name was empty and name is a required parameter")
	}
	if version == "" {
		return errors.New("version was empty and version is a required parameter")
	}
	contains := this.ownsApi(name)
	if contains {
		log.Println("Already contains", name, "nothing to do")
		return nil
	}

	log.Println("Registering", name, version)

	err := sendApiRegistration(name, version)

	if err != nil {
		this.ownedRegsRWMutex.Lock()
		this.ownedRegs[name] = &ownedApi{name: name, version: version}
		this.ownedRegsRWMutex.Unlock()
	}
	return err
}

func sendApiRegistration(name, version string) error {
	log.Println("Sending Api Registration for", name, version)
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: net.ParseIP(MulticastGroupIP), Port: MulticastGroupPort})

	if err != nil {
		return err
	}

	message := &apiRegisterMessageJSON{ApiName: name, ApiVersion: version, LifeSpan: RegistrationLifeSpan}

	dataOut := bytes.NewBuffer(make([]byte, 0, RegistrationMessageSizeBytes))

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

func computeOwnedRegKey(name, version string) string {
	return fmt.Sprint(name, ":", version)
}

func (this *multicastApiRegistry) resendOwnedRegistrationsLoop() {
	updateTicker := time.NewTicker(RegistrationUpdateInterval)
	for range updateTicker.C {
		this.processRegResends()
	}
}

func (this *multicastApiRegistry) processRegResends() {
	log.Println("Starting to process Registration Resends")
	this.ownedRegsRWMutex.RLock()
	log.Println("Number of cur owned APIs:", len(this.ownedRegs))
	for _, curOwnedApi := range this.ownedRegs {
		sendApiRegistration(curOwnedApi.name, curOwnedApi.version)
	}
	this.ownedRegsRWMutex.RUnlock()
}

func (this *multicastApiRegistry) GetAvailableApis() []Api {
	this.apisRWMutex.RLock()
	allApis := make([]Api, 0)

	for curApiName := range this.apiRegs {
		allApis = append(allApis, this.GetApisByApiName(curApiName)...)
	}
	this.apisRWMutex.RUnlock()
	return allApis
}

func (this *multicastApiRegistry) GetApisByApiName(name string) []Api {
	this.apisRWMutex.RLock()
	regs := this.apiRegs[name]
	apis := make([]Api, 0)
	now := time.Now()
	for _, curReg := range regs {
		if curReg.timeRegistered.Add(curReg.lifeSpan).After(now) {
			apis = append(apis, curReg.api)
		}
	}
	this.apisRWMutex.RUnlock()

	return apis
}

func (this *multicastApiRegistry) cleanupExpiredRegLoop() {
	cleanupTicker := time.NewTicker(RegistrationLifeSpan)
	for range cleanupTicker.C {
		expiredApiRegs := this.findExpiredApiRegs()
		if len(expiredApiRegs) > 0 {
			this.apisRWMutex.Lock()
			for _, curExpiredReg := range expiredApiRegs {
				origRegsForName := this.apiRegs[curExpiredReg.api.Name()]
				newRegsForName := make([]*apiRegistration, 0)
				for _, curRegsForName := range origRegsForName {
					if curRegsForName.timeRegistered.Add(curRegsForName.lifeSpan).After(time.Now()) {
						newRegsForName = append(newRegsForName, curRegsForName)
					}
				}
				if len(newRegsForName) > 0 {
					this.apiRegs[curExpiredReg.api.Name()] = newRegsForName
				} else {
					delete(this.apiRegs, curExpiredReg.api.Name())
				}
			}
			this.apisRWMutex.Unlock()
		}
	}
}

func (this *multicastApiRegistry) findExpiredApiRegs() []*apiRegistration {
	now := time.Now()
	this.apisRWMutex.RLock()
	expiredApiRegs := make([]*apiRegistration, 0)
	for _, regs := range this.apiRegs {
		for _, curReg := range regs {
			if curReg.timeRegistered.Add(curReg.lifeSpan).Before(now) {
				expiredApiRegs = append(expiredApiRegs, curReg)
			}
		}
	}
	this.apisRWMutex.RUnlock()

	return expiredApiRegs
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
				this.updateApis(api, message.LifeSpan)
			}
		}
	}
}

func (this *multicastApiRegistry) updateApis(api Api, lifespan time.Duration) {
	this.apisRWMutex.Lock()
	apisForName, contains := this.apiRegs[api.Name()]
	if !contains {
		apisForName = []*apiRegistration{{api: api, timeRegistered: time.Now(), lifeSpan: lifespan}}
		this.apiRegs[api.Name()] = apisForName
	} else {
		matchReg := getRegMatch(api, apisForName)

		if matchReg != nil {
			matchReg.timeRegistered = time.Now()
		} else {
			matchReg = &apiRegistration{api: api, timeRegistered: time.Now(), lifeSpan: lifespan}
		}
	}
	this.apisRWMutex.Unlock()
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
