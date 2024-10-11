package multicast

import "github.com/ZacharyDuve/apireg"

type apiRegisterMessageJSON struct {
	ApiName     string             `json:"api-name"`
	ApiVersion  *versionJSON       `json:"api-version"`
	ApiPort     int                `json:"api-port"`
	SenderUUID  string             `json:"sender-uuid"`
	Environment apireg.Environment `json:"env"`
}
