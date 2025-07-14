package utils

import "log"

func CheckErrorAndCrash(err error, msg string) {
	if err == nil {
		return
	}

	log.Fatal(msg)
}
