# Readme
- Place config.yaml in /etc/openfortivpn-saml, or in the same directory as the binary.
```
# config.yaml contents
---
# When totp: true, openfortivpn-saml will run in CLI-only mode
# and ask for a MFA/TOTP
totp: true

# Okta credentials
username:
password:
```
- Depending on how you currently use openfortivpn, your run command may vary:  
```
openfortivpn-saml | sudo openfortivpn -c /etc/openfortivpn/config --cookie-on-stdin
```
