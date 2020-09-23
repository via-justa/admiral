package main

import (
	"github.com/via-justa/admiral/cmd"
)

var appVersion = "undefined"

func main() {
	cmd.AppVersion = appVersion

	cmd.Execute()

}
