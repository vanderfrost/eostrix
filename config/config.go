package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	SecurityToken  string `json:"securityToken"`
	DefaultChannel string `json:"defaultChannel"`
	LeetcodeRoleID string `json:"leetcodeRoleID"`
}

func ParseConfig() *Config {
	jsonFile, err := os.Open("config/config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	jsonDecoder := json.NewDecoder(jsonFile)
	config := &Config{}
	if err := jsonDecoder.Decode(config); err != nil {
		log.Fatal(err)
	}

	return config
}
