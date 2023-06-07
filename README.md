# Readme
## Description
openfortivpn-saml allows the use of SAML authentication when using Fortinet/FortiGate SSLVPN with Okta as IdP.
## Installation
- Get the correct [release](https://git.deribit.internal/deribit/sys-admin/openfortivpn-saml/-/releases) or compile yourself if it isn't precompiled.
- Run the application and go through the setup wizard.
  - It will prompt for a master password. This is used to encrypt your credentials when saving to disk.
  - The master password is not stored anywhere and is only known by you, just like with a password manager.
- If your Okta credentials have changed or you forgot your master password, remove the config file and start the application again to reinitialize.
  - The config file location is shown when running the application.
## Usage
openfortivpn-saml can run both with or without a browser.  
- During the first run, the application will ask if you intend to use TOTP or not.
  - Answering yes here will not show a browser, keeping it strictly CLI.
  - Answering no here will allow the use of YubiKey and other authentication methods that require a browser.
- If you change your mind later, you can edit the config file and change `totp:`

#### Running
Depending on how you currently use openfortivpn, your run command may vary:  
```
openfortivpn-saml | sudo openfortivpn -c /etc/openfortivpn/config --cookie-on-stdin
```
#### MacOS Gatekeeper
As Sentillia is not a licensed Mac/iOS developer, this software is not officially signed.  
This means you need to allow openfortivpn-saml after the first run via  
`System Preferences > Security & Privacy > Gatekeeper`  
On some systems, it can be found under:  
`System Preferences > Security & Privacy > General`

## Disclaimer
- Due to the nature of SAML, we need to emulate a browser.
- When starting openfortivpn-saml for the first time, it will download dependencies such as Chromium and Playwright.  
  This does not happen in subsequent runs and should not take too long.
