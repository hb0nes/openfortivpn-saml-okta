package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Totp     bool   `yaml:"totp"`
}

func readConfig() Config {
	ex, _ := os.Executable()
	exPath := filepath.Dir(ex)
	dirs := []string{"/etc/openfortivpn-saml", exPath}

	// Open YAML file
	var file *os.File
	var config Config

	for _, dir := range dirs {
		file, _ = os.Open(fmt.Sprintf("%s/config.yaml", dir))
		if file != nil {
			break
		}
	}
	if file == nil {
		log.Fatalf("Could not find any config.yaml file at locations: %s", strings.Join(dirs, ", "))
	}
	defer file.Close()

	// Decode YAML file to struct
	if file != nil {
		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(&config); err != nil {
			log.Fatal(err.Error())
		}
	}
	return config
}
