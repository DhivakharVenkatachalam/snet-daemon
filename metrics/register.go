// Copyright 2018 singularitynet foundation.
// All rights reserved.
// <<add licence terms for code reuse>>

// package for monitoring and reporting the daemon metrics
package metrics

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/singnet/snet-daemon/config"
	log "github.com/sirupsen/logrus"
)

/*
Daemon Registration is required for the metrics services to uniquely identify every daemon and store the
metrics accordingly. To register with metrics service, every daemon must have a Unique Identity i.e. DaemonID.
Daemon ID is a SHA256 hash value generated by using - Org Id, Service ID, Group Id of Daemon (Derived from from service metadata) and Daemon Endpoint.
Post beta, this ID will be used to enable Token based authentication for accessing metrics services.*/

// generates DaemonID nad returns i.e. DaemonID = HASH (Org Name, Service Name, daemon endpoint)
func GetDaemonID() string {
	rawID := config.GetString(config.OrganizationId) + config.GetString(config.ServiceId) + daemonGroupId + config.GetString(config.DaemonEndPoint) + config.GetString(config.RegistryAddressKey)
	//get hash of the string id combination
	hasher := sha256.New()
	hasher.Write([]byte(rawID))
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash
}

type RegisterDaemonPayload struct {
	NetworkID int    `json:"netId"`
	DaemonID  string `json:"daemonId"`
}
type TokenGenerated struct {
	Status string `json:"status"`
	Data   struct {
		Token string `json:"token"`
	} `json:"data"`
}

var daemonGroupId string

var daemonAuthorizationToken string

// setter method for daemonGroupID
func SetDaemonGrpId(grpId string) {
	daemonGroupId = grpId
}

// New Daemon registration. Generates the DaemonID and use that as getting access token
func RegisterDaemon(serviceURL string) bool {
	daemonID := GetDaemonID()
	status := false
	// call the service and get the result
	status = callRegisterService(daemonID, serviceURL)
	if !status {
		log.Infof("Daemon unable to register with the monitoring service. ")
	} else {
		log.Infof("Daemon successfully registered with the monitoring service. ")
	}
	return status
}

/*
service request
{"daemonID":"3a4ebeb75eace1857a9133c7a50bdbb841b35de60f78bc43eafe0d204e523dfe"}

service output
true/false
*/