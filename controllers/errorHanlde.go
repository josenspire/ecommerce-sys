package controllers

import "log"

func CheckError(err error) {
	if err != nil {
		log.Fatalf("Error happend: %s", err.Error())
		return
	}
}
