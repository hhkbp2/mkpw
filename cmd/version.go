package cmd

import "fmt"

var Version string
var GitHash string

func printVersionInfo() {
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("Git Commit Hash: %s\n", GitHash)
}
