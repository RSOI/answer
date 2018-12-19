package utils

import "log"

// DEBUG mode
var DEBUG = false

// LOG message
func LOG(msg string) {
	if DEBUG {
		log.Println(msg)
	}
}
