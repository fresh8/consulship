package main

import (
	"fmt"
	"os"

	"github.com/connectedventures/gonfigurator"
	"github.com/imdario/mergo"
)

type DependencyConfig struct {
	Name    string   `json:"name"`
	Env     string   `json:"env"`
	Version string   `json:"version"`
	Port    int      `json:"port"`
	Tags    []string `json:"tags"`
}

func parseDepConfigs(baseDeps, localDeps *[]DependencyConfig) error {
	gonfigurator.ParseCustomFlag("/etc/consulship/dependencies.yaml", "baseDeps", baseDeps)

	// Default local configuration to pwd
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	gonfigurator.ParseCustomFlag(fmt.Sprintf("%s/dependencies.yaml", wd), "deps", localDeps)

	return nil
}

func mergeDepConfigs(baseDeps, localDeps []DependencyConfig) ([]DependencyConfig, error) {
	mappedDeps := make(map[string]DependencyConfig)

	for _, baseDep := range baseDeps {
		mappedDeps[baseDep.Name] = baseDep
	}

	for _, localDep := range localDeps {
		val, ok := mappedDeps[localDep.Name]
		if ok {
			err := mergo.Merge(&localDep, val)
			if err != nil {
				return nil, err
			}
		}

		mappedDeps[localDep.Name] = localDep
	}

	deps := []DependencyConfig{}
	for _, dep := range mappedDeps {
		deps = append(deps, dep)
	}

	return deps, nil
}
