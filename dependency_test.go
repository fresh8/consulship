package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMergeDepsConfig(t *testing.T) {
	var mergeTests = []struct {
		name     string
		base     []DependencyConfig
		local    []DependencyConfig
		expected []DependencyConfig
		err      error
	}{
		{
			"Blank configs",
			[]DependencyConfig{},
			[]DependencyConfig{},
			[]DependencyConfig{},
			nil,
		},
		{
			"Only base config",
			[]DependencyConfig{
				{
					Name: "redis",
					Env:  "local",
					Port: 6379,
				},
			},
			[]DependencyConfig{},
			[]DependencyConfig{
				{
					Name: "redis",
					Env:  "local",
					Port: 6379,
				},
			},
			nil,
		},
		{
			"Only local config",
			[]DependencyConfig{},
			[]DependencyConfig{
				{
					Name: "redis",
					Env:  "local",
					Port: 6379,
				},
			},
			[]DependencyConfig{
				{
					Name: "redis",
					Env:  "local",
					Port: 6379,
				},
			},
			nil,
		},
		{
			"Both configs (same content in same service)",
			[]DependencyConfig{
				{
					Name: "redis",
					Env:  "local",
					Port: 6379,
				},
			},
			[]DependencyConfig{
				{
					Name: "redis",
					Env:  "local",
					Port: 6379,
				},
			},
			[]DependencyConfig{
				{
					Name: "redis",
					Env:  "local",
					Port: 6379,
				},
			},
			nil,
		},
		{
			"Both configs (different content in same service)",
			[]DependencyConfig{
				{
					Name: "redis",
					Env:  "local",
					Port: 6379,
				},
			},
			[]DependencyConfig{
				{
					Name: "redis",
					Env:  "local",
					Port: 6379,
				},
			},
			[]DependencyConfig{
				{
					Name: "redis",
					Env:  "local",
					Port: 6379,
				},
			},
			nil,
		},
		{
			"Both configs (different services)",
			[]DependencyConfig{
				{
					Name: "aerospike",
					Env:  "local",
					Port: 3000,
				},
			},
			[]DependencyConfig{
				{
					Name: "redis",
					Env:  "local",
					Port: 6379,
				},
			},
			[]DependencyConfig{
				{
					Name: "aerospike",
					Env:  "local",
					Port: 3000,
				},
				{
					Name: "redis",
					Env:  "local",
					Port: 6379,
				},
			},
			nil,
		},
		{
			"Both configs (missing fields)",
			[]DependencyConfig{
				{
					Name:    "redis",
					Env:     "local",
					Port:    6379,
					Version: "v1.0.0",
				},
			},
			[]DependencyConfig{
				{
					Name: "redis",
					Env:  "staging",
				},
			},
			[]DependencyConfig{
				{
					Name:    "redis",
					Env:     "staging",
					Port:    6379,
					Version: "v1.0.0",
				},
			},
			nil,
		},
		{
			"Both configs (additional fields)",
			[]DependencyConfig{
				{
					Name: "redis",
					Env:  "local",
				},
			},
			[]DependencyConfig{
				{
					Name:    "redis",
					Env:     "staging",
					Port:    6379,
					Version: "v1.0.0",
				},
			},
			[]DependencyConfig{
				{
					Name:    "redis",
					Env:     "staging",
					Port:    6379,
					Version: "v1.0.0",
				},
			},
			nil,
		},
	}

	for _, test := range mergeTests {
		got, err := mergeDepConfigs(test.base, test.local)
		if err != nil {
			assert.Equal(t, test.err.Error(), err.Error(), test.name)
		} else {
			assert.Equal(t, test.expected, got, test.name)
		}
	}
}
