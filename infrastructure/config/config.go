package config

import (
	// golang package
	"os"
	"strings"

	// internal package
	"github.com/andrew-susanto/go-sample-arch/infrastructure/log"

	// external package
	"gopkg.in/yaml.v3"
)

var (
	basePath = ""
)

// ParseConfig parse config file with given filename
// returns config and nil error if success
// otherwise return empty config and non nil error
func ParseConfig(environment string) Config {
	filepath := basePath

	environment = strings.ToLower(environment)
	if environment == "dev" {
		filepath += "config/cxp-crmapp.development.yaml"
	} else if environment == "staging" {
		filepath += "config/cxp-crmapp.staging.yaml"
	} else if environment == "production" {
		filepath += "config/cxp-crmapp.production.yaml"
	}

	yamlFile, err := os.ReadFile(filepath)
	if err != nil {
		log.Fatal(err, map[string]interface{}{
			"filepath": filepath,
		}, "ioutil.ReadFile() got error - ParseConfig")
		return Config{}
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		log.Fatal(err, map[string]interface{}{
			"filepath": filepath,
		}, "yaml.Unmarshal() got error - ParseConfig")
		return Config{}
	}

	return config
}
