package apireg

import (
	"errors"
	"net"

	"github.com/google/uuid"
)

type Api interface {
	//Name of the Api that is being served
	Name() string
	//Version of the Api that the server is serving. Used to allow for services to figure out if they can communicate to server or not
	Version() Version
	//UUID of the application that is serving the Api
	UUID() uuid.UUID
	//IP address that the application is being served on. Used for the client to dial back
	HostIP() net.IP
	//Port that the client should dial the serving application on.
	HostPort() int
	//Equal is used to determine if the two apis are the same
	Equal(Api) bool
	//Environment that the server hosting this api is running in Prod, Non-Prod or ALL
	Environment() Environment
}

type apiImpl struct {
	name       string
	version    Version
	uuid       uuid.UUID
	remoteIP   net.IP
	remotePort int
	env        Environment
}

func NewApi(name string, ver Version, uuid uuid.UUID, env Environment, hostIP net.IP, port int) (Api, error) {
	if name == "" {
		return nil, errors.New("name is required for NewApi")
	} else if ver == nil {
		return nil, errors.New("ver (version) is required for NewApi")
	} else if hostIP == nil {
		return nil, errors.New("hostIP is required for NewApi")
	} else if port <= 0 {
		return nil, errors.New("port must be > 0 for NewApi")
	}

	return &apiImpl{name: name, version: ver, uuid: uuid, env: env, remoteIP: hostIP, remotePort: port}, nil
}

func (this *apiImpl) Name() string {
	return this.name
}
func (this *apiImpl) Version() Version {
	return this.version
}

func (this *apiImpl) UUID() uuid.UUID {
	return this.uuid
}

func (this *apiImpl) HostIP() net.IP {
	return this.remoteIP
}
func (this *apiImpl) HostPort() int {
	return this.remotePort
}
func (this *apiImpl) Equal(other Api) bool {
	return other != nil &&
		this.name == other.Name() &&
		this.version.Equal(other.Version()) &&
		this.remoteIP.Equal(other.HostIP()) &&
		this.uuidEqual(other) &&
		this.remotePort == other.HostPort()
}

func (this *apiImpl) uuidEqual(other Api) bool {
	thisUUIDAsBytes := [16]byte(this.uuid)
	otherUUIDAsBytes := [16]byte(other.UUID())

	for i, thisByte := range thisUUIDAsBytes {
		if thisByte != otherUUIDAsBytes[i] {
			//If we found a byte that doesn't equal then short early
			return false
		}
	}
	//If we made it here then they are all equal
	return true
}

func (this *apiImpl) Environment() Environment {
	return this.env
}
