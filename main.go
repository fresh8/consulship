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
		_, err = consulByEnv["local"].Catalog().Register(&api.CatalogRegistration{ID: dep.Name, Address: address}, nil)
		if err != nil {
			log.Fatalf("Cannot register service %s with local consul", dep.Name)
		}
	}
}

func main() {
	var consulEnvConfig []ConsulConfig
	var depConfig []DependencyConfig

	gonfigurator.ParseCustomFlag("/etc/consulship/consul-env.yaml", "consulEnv", &consulEnvConfig)
	gonfigurator.ParseCustomFlag("/etc/consulship/dependencies.yaml", "deps", &depConfig)

	createConsulClients(consulEnvConfig)
	err := gonfigurator.Load()

	if err != nil {
		log.Fatal("Cannot read yaml configurations")
	}

	createConsulClients(consulEnvConfig)
	copyConsulServices(depConfig)
}
