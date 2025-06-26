package main

import (
	prettylog "go_tinker/pkg"
	"log/slog"
)


func main() {
    prettylog.NewMinimalPrettyLogger(nil)
    slog.Info("Server started", "port", 8080)
    slog.Warn("High latency", "duration", "400ms")
    slog.Error("Something went wrong", "err", "bamboo")
}

