package main

import (
	"log/slog"

	"log"

	"github.com/antsrp/house_service/internal/setup"
)

func main() {
	h, _, dbConnection, logger, err := setup.Setup("DB")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	defer dbConnection.Close()
	if err := h.Run(); err != nil {
		logger.Error("cannot run rest server", slog.Any("error", err.Error()))
	}
}
