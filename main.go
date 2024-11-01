package main

import (
	"errors"
	"log/slog"
	"os"
)

/*
Inspired by Mat Ryer's "How I Write HTTP Web Services After n Years" series

Defer he heavy lifting into a function that returns an error if something breaks.
You can easily log that error at one place and return it in the setup for main.
This avoids a lot of the annoying 'if err != nil { println(err); return }' we've
all seen out there
*/
func main() {
	var logger = slog.Default()
	if err := run(logger); err != nil {
		logger.Error("Server failed while running.", slog.Any("result", err))
		os.Exit(2)
	}
}

func run(logger *slog.Logger) error {
	return errors.New("NYI")
}
