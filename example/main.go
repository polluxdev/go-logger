package main

import "github.com/polluxdev/go-logger"

func main() {
	zerolog := logger.NewZerolog("debug")
	zerolog.Debug("This is a debug message")
	zerolog.Error("This is a error message")
	zerolog.Info("This is a info message")
	zerolog.Warn("This is a warning message")
}
