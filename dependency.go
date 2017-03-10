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

func parseDepConfigs() ([]DependencyConfig, []DependencyConfig, error) {
	var baseDeps, localDeps []DependencyConfig
	gonfigurator.ParseCustomFlag("/etc/consulship/dependencies.yaml", "baseDeps", &baseDeps)

	// Default local configuration to pwd
	wd, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}

	gonfigurator.ParseCustomFlag(fmt.Sprintf("%s/dependencies.yaml", wd), "deps", &localDeps)

	return baseDeps, localDeps, nil
}

func mergeDepConfigs(baseDeps, localDeps []DependencyConfig) ([]DependencyConfig, error) {
	return baseDeps, nil
}
