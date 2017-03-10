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

type DependencyConfig struct {
	Name    string   `json:"name"`
	Env     string   `json:"env"`
	Version string   `json:"version"`
	Tags    []string `json:"tags"`
}

var (
	consulByEnv = make(map[string]*api.Client)
)

func main() {
	var consulEnvConfig []ConsulConfig
	consulAddr := os.Getenv("CONSUL_ADDR")

	var depConfig []DependencyConfig

	gonfigurator.ParseCustomFlag("/etc/consulship/consul-env.yaml", "consulEnv", &consulEnvConfig)
	gonfigurator.ParseCustomFlag("/etc/consulship/dependencies.yaml", "deps", &depConfig)
	err := gonfigurator.Load()

	if err != nil {
		log.Fatal("Cannot read yaml configurations")
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
