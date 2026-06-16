package config

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func newValidator() *validator.Validate {
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterStructValidation(Validate, Root{})
	return validate
}

func TestValidate_AcceptsUniqueApplicationClasses(t *testing.T) {
	root := Root{
		DefaultMode: "default",
		Modes: map[string]ModeConfig{
			"mode1": {ApplicationClasses: []string{"App1"}},
			"mode2": {ApplicationClasses: []string{"App2"}},
		},
	}

	if err := newValidator().Struct(&root); err != nil {
		t.Errorf("Struct() = %v, want no error", err)
	}
}

func TestValidate_RejectsDuplicateApplicationClass(t *testing.T) {
	root := Root{
		DefaultMode: "default",
		Modes: map[string]ModeConfig{
			"mode1": {ApplicationClasses: []string{"App1"}},
			"mode2": {ApplicationClasses: []string{"App1"}},
		},
	}

	err := newValidator().Struct(&root)
	if err == nil {
		t.Fatal("Struct() = nil, want a validation error for the duplicate class")
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok || len(validationErrors) != 1 {
		t.Fatalf("Struct() = %v, want exactly one validator.ValidationErrors entry", err)
	}
	if validationErrors[0].Tag() != ErrorApplicationClassDuplicate {
		t.Errorf("validation tag = %q, want %q", validationErrors[0].Tag(), ErrorApplicationClassDuplicate)
	}
}

func TestValidate_RejectsEmptyModes(t *testing.T) {
	root := Root{DefaultMode: "default"}

	if err := newValidator().Struct(&root); err == nil {
		t.Error("Struct() = nil, want an error because modes is empty")
	}
}

func TestValidate_RejectsModeWithoutApplicationClasses(t *testing.T) {
	root := Root{
		DefaultMode: "default",
		Modes: map[string]ModeConfig{
			"mode1": {},
		},
	}

	if err := newValidator().Struct(&root); err == nil {
		t.Error("Struct() = nil, want an error because the mode has no application_classes")
	}
}

func TestValidate_RequiresDefaultMode(t *testing.T) {
	root := Root{
		Modes: map[string]ModeConfig{
			"mode1": {ApplicationClasses: []string{"App1"}},
		},
	}

	if err := newValidator().Struct(&root); err == nil {
		t.Error("Struct() = nil, want an error because default_mode is missing")
	}
}

func TestGetCommonElement(t *testing.T) {
	tests := []struct {
		name      string
		slice1    []string
		slice2    []string
		wantFound bool
		wantValue string
	}{
		{name: "no overlap", slice1: []string{"a", "b"}, slice2: []string{"c"}, wantFound: false},
		{name: "overlap", slice1: []string{"a", "b"}, slice2: []string{"b", "c"}, wantFound: true, wantValue: "b"},
		{name: "empty slices", slice1: nil, slice2: nil, wantFound: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, found := getCommonElement(test.slice1, test.slice2)
			if found != test.wantFound {
				t.Errorf("getCommonElement() found = %v, want %v", found, test.wantFound)
			}
			if found && value != test.wantValue {
				t.Errorf("getCommonElement() value = %q, want %q", value, test.wantValue)
			}
		})
	}
}
