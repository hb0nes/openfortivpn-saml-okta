# Readme
## Description
openfortivpn-saml allows the use of SAML authentication when using Fortinet/FortiGate SSLVPN with Okta as IdP.
## Prerequisites
- openfortivpn-saml requires `openfortivpn` to be installed.
  - If the `--cookie-on-stdin` option is not known, your openfortivpn version is too old.
  - Mac
    - `brew install openfortivpn`
  - Linux
    - Recommended to build yourself from https://github.com/adrienverge/openfortivpn
- openfortivpn requires root privileges to set up the tunnel.  
In order to run it with sudo non-interactively, run this command:  
`echo "${USER} ALL = (root) NOPASSWD: $(which openfortivpn)" | sudo tee /etc/sudoers.d/openfortivpn`

## Installation and configuration
- Get the correct release at `Deployments > Releases` or compile yourself if it isn't precompiled.
- Move the binary to e.g.: `/usr/local/bin`
- Run the binary (see `Running` section below) and go through the setup wizard, which happens when running it for the first time.
  - It will prompt for a master password. This is used to encrypt your credentials.
    - _The master password is not stored anywhere and is only known by you, just like with a password manager._
  - It will ask if you intend to use **Okta FastPass**, TOTP (MFA), Okta Verify, Webauthn (YubiKey) or none of the above.
    - _Note: **Okta FastPass** is highly recommended as it requires no typing or remembering whatsoever._
    - _Note: TouchID won't work with webauthn, YubiKey will._
  - Answering _yes_ to any of these questions will **not** show a browser, keeping openfortivpn-saml strictly CLI.
  - Answering _no_ **will** show a browser, allowing the use of other authentication methods or manual actions.
- If your Okta credentials have changed or you forgot your master password, remove the config file and start the application again to reinitialize.
  - The config file location is shown when running the application.
## Usage
openfortivpn-saml can run both with or without a browser.  

- If you change your mind later, you can edit the config file and change `totp:`

#### Running
Depending on how you currently use openfortivpn, your run command may vary:  
```
sudo openfortivpn -c /etc/openfortivpn/config --cookie="$(openfortivpn-saml)"
```
#### MacOS Gatekeeper
As openfortivpn-saml is not written by a licensed Mac/iOS developer, it is not officially signed.  
This means you need to allow openfortivpn-saml after the first run via  
`System Preferences > Security & Privacy > Gatekeeper`  
On some systems, it can be found under:  
`System Preferences > Security & Privacy > General`

## Disclaimer
- Due to the nature of SAML, we need to emulate a browser.
- When starting openfortivpn-saml for the first time, it will download dependencies such as Chromium and Playwright.  
  This does not happen in subsequent runs and should not take too long.
