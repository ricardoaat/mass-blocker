package blocker

import log "github.com/sirupsen/logrus"

func checkErr(err error, message string) {
	if err != nil {
		log.Panic(message, err)
	}
}
