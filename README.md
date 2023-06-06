# Readme
## Description
openfortivpn-saml allows the use of SAML authentication when using Fortinet/FortiGate SSLVPN with Okta as IdP.
## Disclaimer
- Due to the nature of SAML, we need to emulate a browser.
- When starting openfortivpn-saml for the first time, it will download dependencies such as Chromium and Playwright.  
This does not happen in subsequent runs and should not take too long.
## Installation
- Get the correct [release](https://git.deribit.internal/deribit/sys-admin/openfortivpn-saml/-/releases) or compile yourself if it isn't precompiled.
- Place `config.yaml` in `/etc/openfortivpn-saml`, or in the same directory as the binary.
##### config.yaml contents
```
---
# When totp: true, openfortivpn-saml will run in CLI-only mode
# and ask for a MFA/TOTP
totp: true

# Okta credentials
username:
password:
```
## Usage
openfortivpn-saml can run both with or without a browser.  
- When `config.yaml` contains `totp: false`, it is possible to use YubiKey (WebAuthn) in an interactive browser.  
- If instead `config.yaml` contains `totp: true`, the CLI will ask for a TOTP generated by tools such as Authy, Google Authenticator or 1Password.

Depending on how you currently use openfortivpn, your run command may vary:  
```
openfortivpn-saml | sudo openfortivpn -c /etc/openfortivpn/config --cookie-on-stdin
```
