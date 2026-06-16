package main

import (
	"i3-auto-mode/internal/cli"
	"i3-auto-mode/internal/config"
	"log"

	"github.com/go-playground/validator/v10"
	"go.i3wm.org/i3"
)

func main() {
	cli := cli.New()
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterStructValidation(config.Validate, config.Root{})
	configuration, err := config.FromFile(cli.ConfigPath, validate)
	if err != nil {
		log.Fatal(err)
	}

	switcher := newModeSwitcher(configuration)
	switcher.start()

	receiver := i3.Subscribe(
		i3.WindowEventType,
		i3.WorkspaceEventType,
	)

	for receiver.Next() {
		switcher.handle(receiver.Event())
	}

	// receiver.Next() returns false in case of a fatal connection loss
	if err := receiver.Close(); err != nil {
		log.Fatalf("lost connection to i3: %v", err)
	}
}
