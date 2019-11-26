package util

import "log"

func Logerr(n int, err error) {
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
}
