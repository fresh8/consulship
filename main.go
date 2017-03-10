package main

import (
	"log"
	"os"

	"github.com/connectedventures/gonfigurator"
	"github.com/hashicorp/consul/api"
)

type ConsulConfig struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

var (
	consulByEnv = make(map[string]*api.Client)
)

func main() {
	var consulEnvConfig []ConsulConfig
	consulAddr := os.Getenv("CONSUL_ADDR")

	baseDeps, localDeps, err := parseDepConfigs()
	if err != nil {
		log.Fatal(err)
	}

	gonfigurator.ParseCustomFlag("/etc/consulship/consul-env.yaml", "consulEnv", &consulEnvConfig)
	err = gonfigurator.Load()
	if err != nil {
		log.Fatal(err)
	}

	deps, err := mergeDepConfigs(baseDeps, localDeps)
	if err != nil {
		log.Fatal(err)
	}

	for _, dep := range deps {
		log.Println(dep.Name)
	}

	consulEnvConfig = append(consulEnvConfig, ConsulConfig{
		Name:    "local",
		Address: consulAddr,
	})

	// Create consul clients from consul env config
	for _, consulEnv := range consulEnvConfig {
		config := api.DefaultConfig()
		config.Address = consulEnv.Address
		consulByEnv[consulEnv.Name], err = api.NewClient(config)
		if err != nil {
			log.Fatalf("Cannot create consul client for env %s: %s", consulEnv.Name, err.Error())
		}
	}
}
