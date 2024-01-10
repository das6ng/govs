package main

import (
	"fmt"
	"os"
)

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Printf("[ERROR] %s\n", err.Error())
		os.Exit(1)
	}
}
