package multicast

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ZacharyDuve/apireg"
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
	version apireg.Version
	port    int
}

type multicastApiRegistry struct {
	mAddr *net.UDPAddr
	mConn *net.UDPConn
	//Need to save all of the apis that have been registered externally
	apiRegs *syncApiRegStore
	//Need to know which api registrations are ours so that due to multicast we can double check
	ownedApis          *syncApiStore
	purgeExpiredTicker *time.Ticker
	id                 uuid.UUID
	environment        apireg.Environment
}

func NewMulticastRegistry(lAddr *net.UDPAddr, e apireg.Environment, sId uuid.UUID) (apireg.ApiRegistry, error) {
	//If we are not passed in a lAddr then lets set to defaults
	if lAddr == nil {
		lAddr = &net.UDPAddr{IP: net.ParseIP(DEFAULT_MULTICAST_GROUP_IP), Port: DEFAULT_MULTICAST_GROUP_PORT}
	}

	r := &multicastApiRegistry{}
	r.purgeExpiredTicker = time.NewTicker(registrationPurgeInterval)
	r.apiRegs = newSyncApiRegistrationStore(r.purgeExpiredTicker.C)
	r.id = sId
	r.environment = e
	r.mAddr = lAddr

	r.ownedApis = newSyncApiStore()

	mC, err := net.ListenMulticastUDP("udp", nil, lAddr)

	if err != nil {
		return nil, err
	}
	r.mConn = mC

	go r.listenMutlicast()
	go r.resendOwnedRegistrationsLoop()
	return r, nil
}

func (this *multicastApiRegistry) RegisterApi(name string, version apireg.Version, port int) error {
	if name == "" {
		return errors.New("name was empty and name is a required parameter")
	}
	//We just set a bogus ip as listeners don't get this ip but from the actual packet
	localApi, err := apireg.NewApi(name, version, this.environment, net.ParseIP("0.0.0.0"), port)

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

func (this *multicastApiRegistry) sendApiRegistration(a apireg.Api) error {
	conn, err := net.DialUDP("udp", nil, this.mAddr)

	if err != nil {
		return err
	}

	message := &apiRegisterMessageJSON{
		ApiName:     a.Name(),
		ApiVersion:  &versionJSON{Major: a.Version().Major(), Minor: a.Version().Minor(), BugFix: a.Version().BugFix()},
		ApiPort:     a.HostPort(),
		SenderUUID:  this.id.String(),
		Environment: this.environment}

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

func (this *multicastApiRegistry) GetAvailableApis() []apireg.Api {
	allRegs := this.apiRegs.GetAllRegs()
	allApis := make([]apireg.Api, len(allRegs))
	for i, curReg := range allRegs {
		allApis[i] = curReg.Api()
	}

	return allApis
}

func (this *multicastApiRegistry) GetApisByApiName(name string) []apireg.Api {
	regs := this.apiRegs.GetAllRegsForName(name)
	apis := make([]apireg.Api, len(regs))

	for i, curReg := range regs {
		apis[i] = curReg.Api()
	}
	return apis
}

func (this *multicastApiRegistry) AddEventListener(l apireg.RegistrationListener) {
	this.apiRegs.AddListener(l)
}

func (this *multicastApiRegistry) RemoveEventListener(l apireg.RegistrationListener) {
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
				if message.SenderUUID == ourIDAsString || !shouldProcessMessage(this.environment, message.Environment) {
					continue
				}
				var a apireg.Api
				apiVersion := apireg.NewVersion(message.ApiVersion.Major, message.ApiVersion.Minor, message.ApiVersion.BugFix)
				a, err = apireg.NewApi(message.ApiName, apiVersion, message.Environment, rAddr.IP, message.ApiPort)
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

func shouldProcessMessage(ourEnv, otherEnv apireg.Environment) bool {
	return ourEnv == apireg.All || otherEnv == apireg.All || ourEnv == otherEnv
}

func (this *multicastApiRegistry) updateForApi(a apireg.Api) {
	apisForName := this.apiRegs.GetAllRegsForName(a.Name())

	if len(apisForName) == 0 {
		reg, _ := newApiRegistration(a, time.Now(), registrationLifeSpan)
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
			reg, _ := newApiRegistration(a, time.Now(), registrationLifeSpan)
			this.apiRegs.AddReg(reg)
		}
	}
}
