include .envrc

.PHONY: app
app:
	cd main && go build -ldflags "-X main.TOKEN=${TOKEN} -X main.DOMAIN=${DOMAIN} -X main.AUTH_EMAIL=${AUTH_EMAIL}" -o ../bin/update-dns
