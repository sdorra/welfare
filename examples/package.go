package main

import (
	"fmt"

	"os"

	"github.com/sdorra/welfare/packages"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage: package [present|absent] packagename")
		os.Exit(1)
	}

	action := os.Args[1]
	var state packages.State
	if action == "present" {
		state = packages.StatePresent
	} else if action == "absent" {
		state = packages.StateAbsent
	} else {
		fmt.Printf("unknown action: %s\n", action)
	}

	pkg := os.Args[2]
	changed, err := packages.NewAptModule(pkg, state).Run()
	if err != nil {
		panic(err)
	}

	if changed {
		fmt.Printf("changed state of package %s\n", pkg)
	} else {
		fmt.Println("nothing todo")
	}
}
