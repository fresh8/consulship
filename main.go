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
	consulAddr := os.Getenv("CONSUL_ADDR")

	consulEnvConfig = append(consulEnvConfig, ConsulConfig{
		Name:    "local",
		Address: consulAddr,
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

func copyConsulServices(deps []DependencyConfig) {
	for _, dep := range deps {
		sourceConsul, ok := consulByEnv[dep.Env]
		if !ok {
			log.Fatalf("No such consul env (%s) for service %s", dep.Env, dep.Name)
		}
		services, _, err := sourceConsul.Catalog().Service(dep.Name, "", nil)
		if err != nil {
			log.Fatalf("Cannot get service %s from consul env (%s)", dep.Name, dep.Env)
		}
		if len(services) == 0 {
			log.Fatalf("Found no entries for service %s in consul env (%s)", dep.Name, dep.Env)
		}

		service := services[0]
		address := service.ServiceAddress
		if address == "" {
			address = service.Address
		}
		err = consulByEnv["local"].Agent().ServiceRegister(&api.AgentServiceRegistration{
			ID:      dep.Name,
			Name:    dep.Name,
			Port:    service.ServicePort,
			Address: address,
		})
		if err != nil {
			log.Fatalf("Cannot register service %s with local consul", dep.Name)
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
