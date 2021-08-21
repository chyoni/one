package utils

import "log"

// HandleErr is cause panic if err is not nil.
func HandleErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
