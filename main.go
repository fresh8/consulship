package main

import (
	"fmt"
	"log"
	"os"

	"github.com/connectedventures/gonfigurator"
	"github.com/hashicorp/consul/api"
)

type ConsulConfig struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

const defaultLocalConsulAddr = "localhost:8500"

var (
	consulByEnv = make(map[string]*api.Client)

	workDir string
)

func init() {
	var err error
	// Default local configuration to pwd
	workDir, err = os.Getwd()
	if err != nil {
		log.Fatal("Cannot get working directory")
	}
}

func createConsulClients(consulEnvConfig []ConsulConfig) {
	localConsulAddr := os.Getenv("LOCAL_CONSUL_ADDR")
	if localConsulAddr == "" {
		localConsulAddr = defaultLocalConsulAddr
	}

	consulEnvConfig = append(consulEnvConfig, ConsulConfig{
		Name:    "local",
		Address: localConsulAddr,
	})

	var err error
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

func main() {
	var consulEnvConfig []ConsulConfig
	var baseDeps, localDeps []DependencyConfig

	err := parseDepConfigs(&baseDeps, &localDeps)
	if err != nil {
		log.Fatal(err)
	}

	gonfigurator.ParseCustomFlag(fmt.Sprintf("%s/.consulship/consul-env.yaml", workDir), "consulEnv", &consulEnvConfig)

	createConsulClients(consulEnvConfig)
	err = gonfigurator.Load()

	if err != nil {
		log.Fatalf("Cannot read yaml configurations: %s", err.Error())
	}

	depConfig, err := mergeDepConfigs(baseDeps, localDeps)
	if err != nil {
		log.Fatal(err)
	}

	createConsulClients(consulEnvConfig)
	copyConsulServices(depConfig)
}
