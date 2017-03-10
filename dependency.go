package main

import (
	"fmt"
	"os"

	"github.com/connectedventures/gonfigurator"
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
	for _, localDep := range localDeps {
		for i, baseDep := range baseDeps {
			if baseDep.Name == localDep.Name {
				baseDeps[i] = localDep
			}
		}
	}

	return baseDeps, nil
}
