package main

import "github.com/bilgehannal/cursor-config/curset/cmd"

// version is set at build time via ldflags.
var version = "dev"

func main() {
	cmd.Execute(version)
}
