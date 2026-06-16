package main

import (
	"i3-auto-mode/internal/config"
	"log"
	"strconv"

	"go.i3wm.org/i3"
)

// switches between i3 modes based on a window in focus while handling i3 events
type modeSwitcher struct {
	classModeMap map[string]string
	defaultMode  string
	currentMode  string
	runCommand   func(string) ([]i3.CommandResult, error)
	getTree      func() (i3.Tree, error)
}

func newModeSwitcher(configuration *config.Root) *modeSwitcher {
	return &modeSwitcher{
		classModeMap: buildClassModeMap(configuration),
		defaultMode:  configuration.DefaultMode,
		runCommand:   i3.RunCommand,
		getTree:      i3.GetTree,
	}
}

func (switcher *modeSwitcher) start() {
	switcher.applyMode(switcher.currentFocusMode())
}

func (switcher *modeSwitcher) handle(event i3.Event) {
	switcher.applyMode(switcher.eventMode(event))
}

func (switcher *modeSwitcher) eventMode(event i3.Event) string {
	switch event := event.(type) {
	case *i3.WindowEvent:
		return switcher.windowEventMode(event)
	case *i3.WorkspaceEvent:
		return switcher.workspaceEventMode(event)
	default:
		log.Printf("ignoring unsupported event type %T", event)
		return switcher.currentMode
	}
}

func (switcher *modeSwitcher) windowEventMode(event *i3.WindowEvent) string {
	switch event.Change {
	case "focus":
		return switcher.modeForClass(event.Container.WindowProperties.Class)
	case "close":
		if !event.Container.Focused {
			return switcher.currentMode
		}
		return switcher.defaultMode
	default:
		return switcher.currentMode
	}
}

func (switcher *modeSwitcher) workspaceEventMode(event *i3.WorkspaceEvent) string {
	if event.Change != "focus" {
		return switcher.currentMode
	}
	return switcher.defaultMode
}

func (switcher *modeSwitcher) currentFocusMode() string {
	tree, err := switcher.getTree()
	if err != nil {
		log.Printf("failed to query i3 tree for startup sync, falling back to default_mode: %v", err)
		return switcher.defaultMode
	}

	focused := tree.Root.FindFocused(func(node *i3.Node) bool { return node.Focused })
	if focused == nil {
		return switcher.defaultMode
	}

	return switcher.modeForClass(focused.WindowProperties.Class)
}

func (switcher *modeSwitcher) applyMode(newMode string) {
	if newMode == switcher.currentMode {
		return
	}
	command := "mode " + strconv.Quote(newMode)
	if _, err := switcher.runCommand(command); err != nil {
		log.Printf("failed to switch to mode '%s': %v", newMode, err)
		return
	}
	switcher.currentMode = newMode
}

func (switcher *modeSwitcher) modeForClass(class string) string {
	if mode, ok := switcher.classModeMap[class]; ok {
		return mode
	}
	return switcher.defaultMode
}

func buildClassModeMap(configuration *config.Root) map[string]string {
	classesMap := make(map[string]string)

	for mode, modeConfig := range configuration.Modes {
		classes := modeConfig.ApplicationClasses
		for _, class := range classes {
			classesMap[class] = mode
		}
	}

	return classesMap
}
