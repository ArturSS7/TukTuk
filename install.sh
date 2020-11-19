#!/bin/bash
echo "[*] Installing dependencies"
#requied go modules
go get github.com/gorilla/sessions
go get github.com/labstack/echo
go get github.com/labstack/echo-contrib/session
go get github.com/lib/pq
go get github.com/miekg/dns
go get golang.org/x/oauth2
go get golang.org/x/oauth2/google
go get google.golang.org/api/gmail/v1
go get golang.org/x/crypto/acme/autocert
go get github.com/go-telegram-bot-api/telegram-bot-api
go get github.com/aiomonitors/godiscord
go get github.com/emersion/go-sasl
go get github.com/vjeantet/ldapserver

echo "[*] Installing project"

mv ../TukTuk $GOPATH/src/TukTuk
cd $GOPATH/src
go install TukTuk
cd TukTuk
go build

echo "[*] If have received no errors build is successful"
echo "The project has moved to: "$GOPATH/src/TukTuk
echo "In order to configure database run db_init.sh"
echo "After database init and configuring certificates run ./TukTuk and that's all"
echo "The project binary and sources are here: "$GOPATH/src/TukTuk


