package main

import (
	"errors"
	"testing"

	"i3-auto-mode/internal/config"

	"go.i3wm.org/i3"
)

func newSwitcher() *modeSwitcher {
	return &modeSwitcher{
		classModeMap: map[string]string{
			"App1": "mode1",
			"App2": "mode2",
		},
		defaultMode: "default",
		currentMode: "default",
	}
}

func TestWindowEventMode_FocusKnownClass(t *testing.T) {
	switcher := newSwitcher()
	event := &i3.WindowEvent{
		Change: "focus",
		Container: i3.Node{
			WindowProperties: i3.WindowProperties{Class: "App1"},
		},
	}

	if got := switcher.windowEventMode(event); got != "mode1" {
		t.Errorf("windowEventMode() = %q, want %q", got, "mode1")
	}
}

func TestWindowEventMode_FocusUnknownClassFallsBackToDefault(t *testing.T) {
	switcher := newSwitcher()
	switcher.currentMode = "mode1"
	event := &i3.WindowEvent{
		Change: "focus",
		Container: i3.Node{
			WindowProperties: i3.WindowProperties{Class: "UnknownApp"},
		},
	}

	if got := switcher.windowEventMode(event); got != "default" {
		t.Errorf("windowEventMode() = %q, want %q (fallback to default_mode)", got, "default")
	}
}

func TestWindowEventMode_CloseOfFocusedWindowResetsToDefault(t *testing.T) {
	switcher := newSwitcher()
	switcher.currentMode = "mode1"
	event := &i3.WindowEvent{
		Change:    "close",
		Container: i3.Node{Focused: true},
	}

	if got := switcher.windowEventMode(event); got != "default" {
		t.Errorf("windowEventMode() = %q, want %q", got, "default")
	}
}

func TestWindowEventMode_CloseOfBackgroundWindowLeavesModeUnchanged(t *testing.T) {
	switcher := newSwitcher()
	switcher.currentMode = "mode1"
	event := &i3.WindowEvent{
		Change:    "close",
		Container: i3.Node{Focused: false},
	}

	if got := switcher.windowEventMode(event); got != "mode1" {
		t.Errorf("windowEventMode() = %q, want unchanged %q", got, "mode1")
	}
}

func TestWindowEventMode_TitleChangeLeavesModeUnchanged(t *testing.T) {
	switcher := newSwitcher()
	switcher.currentMode = "mode1"
	event := &i3.WindowEvent{
		Change: "title",
		Container: i3.Node{
			WindowProperties: i3.WindowProperties{Class: "App1"},
		},
	}

	if got := switcher.windowEventMode(event); got != "mode1" {
		t.Errorf("windowEventMode() = %q, want unchanged %q", got, "mode1")
	}
}

func TestWorkspaceEventMode_FocusResetsToDefault(t *testing.T) {
	switcher := newSwitcher()
	switcher.currentMode = "mode1"
	event := &i3.WorkspaceEvent{Change: "focus"}

	if got := switcher.workspaceEventMode(event); got != "default" {
		t.Errorf("workspaceEventMode() = %q, want %q", got, "default")
	}
}

func TestWorkspaceEventMode_OtherChangeLeavesModeUnchanged(t *testing.T) {
	switcher := newSwitcher()
	switcher.currentMode = "mode1"
	event := &i3.WorkspaceEvent{Change: "urgent"}

	if got := switcher.workspaceEventMode(event); got != "mode1" {
		t.Errorf("workspaceEventMode() = %q, want unchanged %q", got, "mode1")
	}
}

func TestEventMode_UnknownEventTypeLeavesModeUnchanged(t *testing.T) {
	switcher := newSwitcher()
	switcher.currentMode = "mode1"

	if got := switcher.eventMode(&i3.BindingEvent{}); got != "mode1" {
		t.Errorf("eventMode() = %q, want unchanged %q", got, "mode1")
	}
}

func TestApplyMode_RunsCommandAndUpdatesCurrentModeWhenModeDiffers(t *testing.T) {
	switcher := newSwitcher()
	var gotCommand string
	switcher.runCommand = func(command string) ([]i3.CommandResult, error) {
		gotCommand = command
		return nil, nil
	}

	switcher.applyMode("mode1")

	if gotCommand != `mode "mode1"` {
		t.Errorf("runCommand called with %q, want %q", gotCommand, `mode "mode1"`)
	}
	if switcher.currentMode != "mode1" {
		t.Errorf("currentMode = %q, want %q", switcher.currentMode, "mode1")
	}
}

func TestApplyMode_SkipsCommandWhenModeUnchanged(t *testing.T) {
	switcher := newSwitcher()
	called := false
	switcher.runCommand = func(string) ([]i3.CommandResult, error) {
		called = true
		return nil, nil
	}

	switcher.applyMode(switcher.currentMode)

	if called {
		t.Error("runCommand was called, want no-op when newMode == currentMode")
	}
}

func TestApplyMode_LeavesCurrentModeUnchangedOnCommandError(t *testing.T) {
	switcher := newSwitcher()
	switcher.runCommand = func(string) ([]i3.CommandResult, error) {
		return nil, errors.New("boom")
	}

	switcher.applyMode("mode1")

	if switcher.currentMode != "default" {
		t.Errorf("currentMode = %q, want unchanged %q after a failed command", switcher.currentMode, "default")
	}
}

func TestCurrentFocusMode_UsesFocusedWindowClass(t *testing.T) {
	switcher := newSwitcher()
	switcher.getTree = func() (i3.Tree, error) {
		return i3.Tree{Root: &i3.Node{
			Focused:          true,
			WindowProperties: i3.WindowProperties{Class: "App1"},
		}}, nil
	}

	if got := switcher.currentFocusMode(); got != "mode1" {
		t.Errorf("currentFocusMode() = %q, want %q", got, "mode1")
	}
}

func TestCurrentFocusMode_FallsBackToDefaultWhenNothingFocused(t *testing.T) {
	switcher := newSwitcher()
	switcher.getTree = func() (i3.Tree, error) {
		return i3.Tree{Root: &i3.Node{}}, nil
	}

	if got := switcher.currentFocusMode(); got != "default" {
		t.Errorf("currentFocusMode() = %q, want %q", got, "default")
	}
}

func TestCurrentFocusMode_FallsBackToDefaultOnTreeError(t *testing.T) {
	switcher := newSwitcher()
	switcher.getTree = func() (i3.Tree, error) {
		return i3.Tree{}, errors.New("boom")
	}

	if got := switcher.currentFocusMode(); got != "default" {
		t.Errorf("currentFocusMode() = %q, want %q", got, "default")
	}
}

func TestBuildClassModeMap(t *testing.T) {
	configuration := &config.Root{
		DefaultMode: "default",
		Modes: map[string]config.ModeConfig{
			"mode1": {ApplicationClasses: []string{"App1"}},
			"mode2": {ApplicationClasses: []string{"App2", "App3"}},
		},
	}

	got := buildClassModeMap(configuration)
	want := map[string]string{
		"App1": "mode1",
		"App2": "mode2",
		"App3": "mode2",
	}
	if len(got) != len(want) {
		t.Fatalf("buildClassModeMap() = %v, want %v", got, want)
	}
	for class, mode := range want {
		if got[class] != mode {
			t.Errorf("buildClassModeMap()[%q] = %q, want %q", class, got[class], mode)
		}
	}
}
