package main

import (
	"bytes"
	"fmt"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

const configName = "openfortivpn-saml.yaml"

type Config struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Totp     bool   `yaml:"totp"`
	Verify   bool   `yaml:"verify"`
	Headless *bool  `yaml:"headless,omitempty"`
}

func getConfigPath() (configPath string) {
	dir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("Could not determine user directory.")
	}
	return filepath.Join(dir, configName)
}

func initConfig() (config Config, err error) {
	var pwdMaster []byte
	var pwdOkta []byte
	var userOkta string
	var verifyStr string
	var verify bool
	var totpStr string
	var totp bool

	log.Println("It appears this is your first time running openfortivpn-saml. Let's configure it.")
	log.Println("Please enter a master password that you can remember. This password is not stored anywhere.")
	pwdMaster, err = term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Printf("Error while reading input: %v", err)
	}
	log.Println("Please enter it again.")
	pwdMaster2, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Printf("Error while reading input: %v", err)
	}
	if bytes.Compare(pwdMaster, pwdMaster2) != 0 {
		log.Fatalf("Didn't enter matching master passwords.")
	}
	log.Println("Please enter your Okta username.")
	_, err = fmt.Scanf("%s", &userOkta)
	if err != nil {
		log.Printf("Error while reading input: %v", err)
	}
	log.Println("Please enter your Okta password.")
	pwdOkta, err = term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Printf("Error while reading input: %v", err)
	}
	log.Println("Do you intend to use TOTP/MFA such as Google Authenticator (y/n)?")
	_, err = fmt.Scanf("%s", &totpStr)
	if err != nil {
		log.Printf("Error while reading input: %v", err)
	}
	totp = strings.ToLower(totpStr) == "y" || strings.ToLower(totpStr) == "yes"
	if !totp {
		log.Println("Do you intend to use Okta verify? (y/n)?")
		_, err = fmt.Scanf("%s", &verifyStr)
		if err != nil {
			log.Printf("Error while reading input: %v", err)
		}
		verify = strings.ToLower(verifyStr) == "y" || strings.ToLower(verifyStr) == "yes"
	}
	configFile, err := os.Create(getConfigPath())
	if err != nil {
		log.Fatalf("Could not open file %s for writing: %v", configName, err)
	}
	defer configFile.Close()
	config = Config{
		Username: userOkta,
		Password: string(pwdOkta),
		Totp:     totp,
		Verify:   verify,
	}
	configEncrypted := Config{
		Username: encrypt(string(pwdMaster), userOkta),
		Password: encrypt(string(pwdMaster), string(pwdOkta)),
		Totp:     totp,
		Verify:   verify,
	}
	configEncryptedBytes, err := yaml.Marshal(configEncrypted)
	if err != nil {
		log.Fatalf("%v", err)
	}
	configFile.WriteString(string(configEncryptedBytes))
	return config, err
}

func configRead() (config Config) {
	log.Printf("Loading config from: %s...", getConfigPath())
	var file *os.File
	file, _ = os.Open(getConfigPath())
	// Initialize config if none is found
	if file == nil {
		var err error
		config, err = initConfig()
		if err != nil {
			log.Fatalf("Could not init config: %v", err)
		}
		return
	}
	// Decode config if found
	if file != nil {
		defer file.Close()
		decoder := yaml.NewDecoder(file)
		if err := decoder.Decode(&config); err != nil {
			log.Fatalf("Could not decode config at %v: %v", getConfigPath(), err)
		}
	}
	log.Println("Please enter your master password.")
	pwdMaster, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Printf("Error while reading input: %v", err)
	}
	config.Username = decrypt(string(pwdMaster), config.Username)
	config.Password = decrypt(string(pwdMaster), config.Password)
	return config
}
