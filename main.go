package main

import (
	"fmt"
	"os"
)

func run() error {
	app := setupApplication()
	return app.Run(os.Args)
}

func setupApplication() *Application {
	a := &Application{}
	a.commands()
	return a
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
