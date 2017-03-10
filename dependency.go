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
		inBase := false
		for i, baseDep := range baseDeps {
			if baseDep.Name == localDep.Name {
				baseDeps[i] = mergeDepConfig(baseDep, localDep)
				inBase = true
				break
			}
		}

		if !inBase {
			baseDeps = append(baseDeps, localDep)
		}
	}

	return baseDeps, nil
}

func mergeDepConfig(base, local DependencyConfig) DependencyConfig {
	if local.Env != "" {
		base.Env = local.Env
	}

	if local.Version != "" {
		base.Version = local.Version
	}

	if local.Port != 0 {
		base.Port = local.Port
	}

	// TODO: Work out the best way to do tag merging - should local override regardless of if set or not?

	return base
}
