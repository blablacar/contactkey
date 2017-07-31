package main

import (
	"flag"
	"github.com/remyLemeunier/contactkey/commands"
)

func main() {

	// Needed for glog package
	flag.Parse()
	flag.Lookup("alsologtostderr").Value.Set("true")

	commands.Execute()
}
