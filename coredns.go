package main

//go:generate go run directives_generate.go

import (
	"coredns/coremain"

	// Plug in CoreDNS
	_ "coredns/core/plugin"
)

func main() {
	coremain.Run()
}
