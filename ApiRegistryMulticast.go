package apireg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ZacharyDuve/apireg/api"
	"github.com/ZacharyDuve/apireg/environment"
	"github.com/ZacharyDuve/apireg/event"
	"github.com/ZacharyDuve/apireg/store"
	"github.com/ZacharyDuve/serverid"
)

const (
	multicastGroupIP             string        = "224.0.0.78"
	multicastGroupPort           int           = 5324
	registrationMessageSizeBytes int           = 1400
	registrationLifeSpan         time.Duration = registrationUpdateInterval * 4
	registrationUpdateInterval   time.Duration = time.Second * 15
	registrationPurgeInterval    time.Duration = time.Second * 30
)

type ownedApi struct {
	name    string
	version *api.Version
	port    int
}

type multicastApiRegistry struct {
	mConn              *net.UDPConn
	apiRegs            store.ApiRegistrationStore
	ownedApis          store.ApiStore
	purgeExpiredTicker *time.Ticker
	id                 string
	environment        environment.Environment
}

func NewRegistry(e environment.Environment) (ApiRegistry, error) {
	var err error
	r := &multicastApiRegistry{}
	r.purgeExpiredTicker = time.NewTicker(registrationPurgeInterval)
	r.apiRegs = store.NewSyncApiRegistrationStore(r.purgeExpiredTicker.C)
	sIdSvc, err := serverid.NewFileServerIdService("")
	if err != nil {
		return nil, err
	}
	r.id = sIdSvc.GetServerId().String()
	r.environment = e
	if err != nil {
		return nil, err
	}

	r.ownedApis = store.NewSyncApiStore()

	mC, err := net.ListenMulticastUDP("udp", nil, &net.UDPAddr{IP: net.ParseIP(multicastGroupIP), Port: multicastGroupPort})

	if err != nil {
		return nil, err
	}
	r.mConn = mC

	go r.listenMutlicast()
	go r.resendOwnedRegistrationsLoop()
	return r, nil
}

func (this *multicastApiRegistry) RegisterApi(name string, version *api.Version, port int) error {
	if name == "" {
		return errors.New("name was empty and name is a required parameter")
	}
	localApi, err := api.NewApi(name, version, this.environment, net.ParseIP("127.0.0.1"), port)

	if err != nil {
		return err
	}
	if this.ownedApis.Contains(localApi) {
		return nil
	}

	err = this.sendApiRegistration(localApi)

	if err == nil {
		this.ownedApis.Add(localApi)
	}
	return err
}

func (this *multicastApiRegistry) sendApiRegistration(a api.Api) error {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{IP: net.ParseIP(multicastGroupIP), Port: multicastGroupPort})

	if err != nil {
		return err
	}

	message := &apiRegisterMessageJSON{ApiName: a.Name(), ApiVersion: a.Version(), ApiPort: a.HostPort(), SenderId: this.id, Environment: this.environment}

	dataOut := bytes.NewBuffer(make([]byte, 0, registrationMessageSizeBytes))

	err = json.NewEncoder(dataOut).Encode(message)

	if err != nil {
		return err
	}

	if dataOut.Len() > registrationMessageSizeBytes {
		return errors.New(fmt.Sprint("Message size for", a.Name(), a.Version(), "exceeds max length of", registrationMessageSizeBytes, "bytes"))
	}

	_, err = conn.Write(dataOut.Bytes())

	return err
}

func (this *multicastApiRegistry) resendOwnedRegistrationsLoop() {
	updateTicker := time.NewTicker(registrationUpdateInterval)
	for range updateTicker.C {
		this.processRegResends()
	}
}

func (this *multicastApiRegistry) processRegResends() {
	for _, curOwnedApi := range this.ownedApis.All() {
		this.sendApiRegistration(curOwnedApi)
	}
}

func (this *multicastApiRegistry) GetAvailableApis() []api.Api {
	allRegs := this.apiRegs.GetAllRegs()
	allApis := make([]api.Api, len(allRegs))
	for i, curReg := range allRegs {
		allApis[i] = curReg.Api()
	}

	return allApis
}

func (this *multicastApiRegistry) GetApisByApiName(name string) []api.Api {
	regs := this.apiRegs.GetAllRegsForName(name)
	apis := make([]api.Api, len(regs))

	for i, curReg := range regs {
		apis[i] = curReg.Api()
	}
	return apis
}

func (this *multicastApiRegistry) AddListener(l event.RegistrationListener) {
	this.apiRegs.AddListener(l)
}

func (this *multicastApiRegistry) RemoveListener(l event.RegistrationListener) {
	this.apiRegs.RemoveListener(l)
}

func (this *multicastApiRegistry) listenMutlicast() {
	readBuff := make([]byte, registrationMessageSizeBytes)
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
				//If we got a message from ourselves or for another environment then ignore it
				if message.SenderId == this.id || !shouldProcessMessage(this.environment, message.Environment) {
					continue
				}
				var a api.Api
				a, err = api.NewApi(message.ApiName, message.ApiVersion, message.Environment, rAddr.IP, message.ApiPort)
				if err != nil {
					log.Println("Error generating new Api from message")
				} else {
					this.updateForApi(a)
				}
			}
		}
	}
}

//Us	| Msg	| pro
// A	| A		| Y
// A	| P		| Y
// A	| N		| Y
// P	| A		| Y
// P	| P		| Y
// P	| N		| N
// N	| A		| Y
// N	| P		| N
// N	| N		| Y

func shouldProcessMessage(ourEnv, otherEnv environment.Environment) bool {
	return ourEnv == environment.All || otherEnv == environment.All || ourEnv == otherEnv
}

func (this *multicastApiRegistry) updateForApi(a api.Api) {
	apisForName := this.apiRegs.GetAllRegsForName(a.Name())

	if len(apisForName) == 0 {
		reg, _ := store.NewApiRegistration(a, time.Now(), registrationLifeSpan)
		this.apiRegs.AddReg(reg)
	} else if len(apisForName) > 0 {
		matched := false
		for _, curReg := range apisForName {
			if curReg.Api().Equal(a) {
				curReg.UpdateTimeRegistered(time.Now())
				matched = true
			}
		}

		if !matched {
			reg, _ := store.NewApiRegistration(a, time.Now(), registrationLifeSpan)
			this.apiRegs.AddReg(reg)
		}
	}
}
