package main

import (
	"fmt"
	"monkey/repl"
	"os"
)

func main() {
	fmt.Printf("MonkeyLang.\n")
	repl.Start(os.Stdin, os.Stdout)
}
