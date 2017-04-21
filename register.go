package main

import (
	"log"

	"github.com/hashicorp/consul/api"
)

func registerService(reg *api.AgentServiceRegistration) error {
	return consulByEnv["local"].Agent().ServiceRegister(reg)
}

func registerLocalService(dep DependencyConfig) error {
	return registerService(&api.AgentServiceRegistration{
		ID:      dep.Name,
		Name:    dep.Name,
		Port:    dep.Port,
		Address: dep.Address,
		Tags:    dep.Tags,
	})
}

func registerRemoteService(dep DependencyConfig) error {
	sourceConsul, ok := consulByEnv[dep.Env]
	if !ok {
		log.Fatalf("No such consul env (%s) for service %s", dep.Env, dep.Name)
	}
	if len(dep.Tags) == 0 {
		dep.Tags = append(dep.Tags, "")
	}
	for _, tag := range dep.Tags {
		services, _, err := sourceConsul.Catalog().Service(dep.Name, tag, nil)
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
		reg := &api.AgentServiceRegistration{
			ID:      dep.Name,
			Name:    dep.Name,
			Port:    service.ServicePort,
			Address: address,
		}
		// when there is no tag we should pass an empty array
		if tag != "" {
			reg.Tags = append(reg.Tags, tag)
		}
		err = registerService(reg)
		if err != nil {
			return err
		}
	}
	return nil
}

func copyConsulServices(deps []DependencyConfig) {
	for _, dep := range deps {
		var err error
		if dep.Env == "local" {
			err = registerLocalService(dep)
		} else {
			err = registerRemoteService(dep)
		}
		if err != nil {
			log.Fatalf("Cannot register service %s with local consul: %s", dep.Name, err.Error())
		}
	}
}
