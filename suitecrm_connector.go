// +build !b24

package main

import (
	"github.com/asm2me/asterlink/connect"
	"github.com/asm2me/asterlink/connect/suitecrm"
	"gopkg.in/yaml.v2"
)

func init() {
	connectors = append(connectors, newSuiteCRMConnector)
}

func newSuiteCRMConnector(cfgBytes []byte) (connecter connect.Connecter, err error) {
	var config struct {
		SuiteCRM suitecrm.Config `yaml:"suitecrm"`
	}

	if err = yaml.Unmarshal(cfgBytes, &config); err != nil {
		return
	}

	if config.SuiteCRM.URL == "" {
		return
	}

	connecter = suitecrm.NewSuiteCRMConnector(&config.SuiteCRM)
	return
}
