package settings

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

func ImportTOML(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read settings file: %w", err)
	}

	var raw map[string]map[string]interface{}
	if err := toml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse TOML: %w", err)
	}

	flat := make(map[string]string)
	for section, fields := range raw {
		for key, value := range fields {
			flatKey := section + "." + key
			flat[flatKey] = fmt.Sprintf("%v", value)
		}
	}

	return flat, nil
}
