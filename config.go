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
	Username       string `yaml:"username"`
	Password       string `yaml:"password"`
	PasswordMaster string `yaml:"-"`
	Totp           bool   `yaml:"totp"`
	Verify         bool   `yaml:"verify"`
	Webauthn       bool   `yaml:"webauthn"`
	FastPass       bool   `yaml:"fast_pass"`
	Headless       *bool  `yaml:"headless,omitempty"`
}

func getConfigPath() (configPath string) {
	dir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal("Could not determine user directory.")
	}
	return filepath.Join(dir, configName)
}

func configInitWizard() (config *Config, err error) {
	var pwdMaster []byte
	var pwdOkta []byte
	var userOkta string
	var verifyStr string
	var totpStr string
	var webauthnStr string
	var fastPassStr string
	config = &Config{
		Username:       "",
		Password:       "",
		PasswordMaster: "",
		Totp:           false,
		Verify:         false,
		Webauthn:       false,
		FastPass:       false,
		Headless:       nil,
	}
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
	config.PasswordMaster = string(pwdMaster)

	log.Println("Please enter your Okta username.")
	_, err = fmt.Scanf("%s", &userOkta)
	if err != nil {
		log.Fatalf("Error while reading input: %v", err)
	}
	config.Username = userOkta

	log.Println("Please enter your Okta password.")
	pwdOkta, err = term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return nil, fmt.Errorf("error while reading input: %v", err)
	}
	config.Password = string(pwdOkta)

	log.Println("Do you intend to use Okta FastPass (y/n)?")
	_, err = fmt.Scanf("%s", &fastPassStr)
	if err != nil {
		return nil, fmt.Errorf("error while reading input: %v", err)
	}
	if config.FastPass = strings.ToLower(fastPassStr) == "y" || strings.ToLower(fastPassStr) == "yes"; config.FastPass {
		return
	}

	log.Println("Do you intend to use TOTP/MFA such as Google Authenticator (y/n)?")
	_, err = fmt.Scanf("%s", &totpStr)
	if err != nil {
		return nil, fmt.Errorf("error while reading input: %v", err)
	}
	if config.Totp = strings.ToLower(totpStr) == "y" || strings.ToLower(totpStr) == "yes"; config.Totp {
		return
	}

	log.Println("Do you intend to use Okta verify? (y/n)?")
	_, err = fmt.Scanf("%s", &verifyStr)
	if err != nil {
		return nil, fmt.Errorf("error while reading input: %v", err)
	}
	if config.Verify = strings.ToLower(verifyStr) == "y" || strings.ToLower(verifyStr) == "yes"; config.Verify {
		return
	}

	log.Println("Do you intend to use Webauthn (YubiKey)? (y/n)?")
	_, err = fmt.Scanf("%s", &webauthnStr)
	if err != nil {
		return nil, fmt.Errorf("error while reading input: %v", err)
	}
	return
}

func initConfig() (config *Config, err error) {
	config, err = configInitWizard()
	if err != nil {
		return
	}
	configEncrypted := *config
	configEncrypted.Password = encrypt(config.Username, config.Username)
	configEncrypted.Password = encrypt(config.PasswordMaster, config.Password)
	configEncryptedBytes, err := yaml.Marshal(configEncrypted)
	if err != nil {
		return nil, err
	}
	configFile, err := os.Create(getConfigPath())
	if err != nil {
		return nil, fmt.Errorf("could not open file %s for writing: %v", configName, err)
	}
	defer configFile.Close()
	configFile.WriteString(string(configEncryptedBytes))
	return config, err
}

func configRead() (config *Config) {
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
