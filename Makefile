build:
	GOOS=windows GOARCH=amd64 go build -o openfortivpn-saml
	tar -cvzf openfortivpn-saml-windows-amd64.tar.gz openfortivpn-saml
	GOOS=linux GOARCH=amd64 go build -o openfortivpn-saml
	tar -cvzf openfortivpn-saml-linux-amd64.tar.gz openfortivpn-saml
	GOOS=darwin GOARCH=arm64 go build -o openfortivpn-saml
	tar -cvzf openfortivpn-saml-darwin-arm64.tar.gz openfortivpn-saml
	rm -v openfortivpn-saml

clean:
	rm -v openfortivpn-saml*
