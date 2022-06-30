package eos_test

import (
	"fmt"
	"os"
	"runtime"
)

func getAPIURL() string {
	apiURL := os.Getenv("EOS_GO_API_URL")
	if apiURL != "" {
		return apiURL
	}

	return "https://api.eosn.io/"
}

func Ensure(condition bool, message string, args ...interface{}) {
	if !condition {
		Quit(message, args...)
	}
}

func NoError(err error, message string, args ...interface{}) {
	if err != nil {
		Quit(message+": "+err.Error(), args...)
	}
}

func Quit(message string, args ...interface{}) {
	fmt.Printf(message+"\n", args...)
	runtime.Goexit()
}
