package main

import (
	"fmt"

	"github.com/sdorra/welfare"
	"github.com/sdorra/welfare/packages"
)

func main() {
	aptKey := packages.NewAptKeyModule("D742B261", packages.Present)
	exec(aptKey, "add scm-manager key")

	aptRepo := packages.NewAptRepositoryModule("scm-manager", packages.Present)
	aptRepo.Repository = "deb http://maven.scm-manager.org/nexus/content/repositories/releases ./"
	exec(aptRepo, "add scm-manager repository")

	apt := packages.NewAptModule("openjdk-8-jre", packages.Present)
	exec(apt, "install java")

	apt = packages.NewAptModule("scm-server", packages.Present)
	exec(apt, "install scm-server")
}

func exec(module welfare.Module, action string) {
	changed, err := module.Run()
	if err != nil {
		panic(err)
	}
	fmt.Println(action, ":", changed)
}
