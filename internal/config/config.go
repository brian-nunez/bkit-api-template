package config

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/brian-nunez/bconfig"
	"github.com/brian-nunez/bconfig/drivers/file"
)

type customEnvSource struct{}

func EnvSource() bconfig.Source {
	return &customEnvSource{}
}

func (e *customEnvSource) Name() string {
	return "custom-env-nested"
}

func (e *customEnvSource) Load(ctx context.Context) (map[string]any, error) {
	result := make(map[string]any)
	environ := os.Environ()

	for _, envVar := range environ {
		pair := strings.SplitN(envVar, "=", 2)
		if len(pair) < 2 {
			continue
		}
		key := pair[0]
		val := pair[1]

		// Support both underscore and dot-separated names
		normalizedKey := strings.ReplaceAll(key, ".", "_")
		parts := strings.Split(normalizedKey, "_")
		if len(parts) == 0 {
			continue
		}

		firstPart := strings.ToLower(parts[0])
		// Only process variables corresponding to our config sections to avoid polluting with other env vars
		if firstPart != "server" && firstPart != "telemetry" && firstPart != "kv" && firstPart != "db" {
			continue
		}

		for i, part := range parts {
			parts[i] = strings.ToLower(part)
		}

		// Traverse and build deep nesting
		curr := result
		for i := 0; i < len(parts)-1; i++ {
			part := parts[i]
			if _, exists := curr[part]; !exists {
				curr[part] = make(map[string]any)
			}

			// Override if the existing value is not a map
			nextMap, ok := curr[part].(map[string]any)
			if !ok {
				nextMap = make(map[string]any)
				curr[part] = nextMap
			}
			curr = nextMap
		}

		// Assign the parsed value to the leaf key
		leafKey := parts[len(parts)-1]
		curr[leafKey] = parseValue(val)
	}

	return result, nil
}

// parseValue converts environment value strings into basic types
func parseValue(val string) any {
	if val == "true" {
		return true
	}
	if val == "false" {
		return false
	}
	if i, err := strconv.Atoi(val); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(val, 64); err == nil {
		return f
	}
	return val
}

// Load reads config from config.yaml and applies overrides from env vars without BKIT_ prefix
func Load(ctx context.Context) (*bconfig.Config, error) {
	cfg, err := bconfig.Load(
		ctx,
		file.Source("config.yaml"), // Load defaults first
		EnvSource(),                // Override with custom nested env variables (e.g. SERVER_PORT=9090 or SERVER.PORT=9090)
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	return cfg, nil
}
