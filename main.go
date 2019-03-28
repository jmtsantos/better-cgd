package main

import (
	"os"

	"github.com/shopspring/decimal"

	cgdbreak "better-cgd/cgdbreak"

	log "github.com/sirupsen/logrus"
)

func main() {
	var (
		balance  decimal.Decimal
		username string
		password string
		err      error
	)

	log.Println("[observer][cgd][break] starting")

	username = os.Getenv("USER")
	password = os.Getenv("PASS")

	if username == "" || password == "" {
		log.Errorln("error: USER or PASS not set")
		return
	}

	clientBreak := cgdbreak.New(username, password)

	if balance, err = clientBreak.CheckBreakBalance(); err != nil {
		log.WithError(err).Errorln("[observer][cgd][break] error processing balance")
		return
	}

	log.Println("[observer][cgd][break] current balance for Caixa Break", balance.String())
}
