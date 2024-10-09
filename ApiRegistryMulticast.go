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
	"github.com/ZacharyDuve/apireg/apievent"
	"github.com/ZacharyDuve/apireg/environment"
	"github.com/ZacharyDuve/apireg/store"
	"github.com/google/uuid"
)

const (
	DEFAULT_MULTICAST_GROUP_IP   string        = "224.0.0.78"
	DEFAULT_MULTICAST_GROUP_PORT int           = 5324
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
	mAddr              *net.UDPAddr
	mConn              *net.UDPConn
	apiRegs            store.ApiRegistrationStore
	ownedApis          store.ApiStore
	purgeExpiredTicker *time.Ticker
	id                 uuid.UUID
	environment        environment.Environment
}

func NewMulticastRegistry(lAddr *net.UDPAddr, e environment.Environment, sId uuid.UUID) (ApiRegistry, error) {
	//If we are not passed in a lAddr then lets set to defaults
	if lAddr == nil {
		lAddr = &net.UDPAddr{IP: net.ParseIP(DEFAULT_MULTICAST_GROUP_IP), Port: DEFAULT_MULTICAST_GROUP_PORT}
	}

	r := &multicastApiRegistry{}
	r.purgeExpiredTicker = time.NewTicker(registrationPurgeInterval)
	r.apiRegs = store.NewSyncApiRegistrationStore(r.purgeExpiredTicker.C)
	r.id = sId
	r.environment = e
	r.mAddr = lAddr

	r.ownedApis = store.NewSyncApiStore()

	mC, err := net.ListenMulticastUDP("udp", nil, lAddr)

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
	//We just set a bogus ip as listeners don't get this ip but from the actual packet
	localApi, err := api.NewApi(name, version, this.environment, net.ParseIP("0.0.0.0"), port)

	if err != nil {
		return err
	}
	//If we already know that we have registered this api from us then don't re-register it
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
	conn, err := net.DialUDP("udp", nil, this.mAddr)

	if err != nil {
		return err
	}

	message := &apiRegisterMessageJSON{ApiName: a.Name(), ApiVersion: a.Version(), ApiPort: a.HostPort(), SenderId: this.id.String(), Environment: this.environment}

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

func (this *multicastApiRegistry) AddListener(l apievent.RegistrationListener) {
	this.apiRegs.AddListener(l)
}

func (this *multicastApiRegistry) RemoveListener(l apievent.RegistrationListener) {
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
				ourIDAsString := this.id.String()
				//If we got a message from ourselves or for another environment then ignore it
				if message.SenderId == ourIDAsString || !shouldProcessMessage(this.environment, message.Environment) {
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
