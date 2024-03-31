package main

import (
	"gointer/repl"
	"os"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
