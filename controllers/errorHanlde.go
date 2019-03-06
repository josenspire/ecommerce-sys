package controllers

import "log"

func ErrorHandle(err error) {
	if err != nil {
		log.Fatalf("Error happend: %s", err.Error())
		return
	}
}
