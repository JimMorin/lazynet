package main

import (
	"github.com/sureshkrishnan-v/lazynet/internal/tui"
)

// main is the entry point of the application. It creates a new App
// and runs it, handling the entire TUI lifecycle.
func main() {
	app := tui.NewApp()
	if err := app.Run(); err != nil {
		panic(err)
	}
}
