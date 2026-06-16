package config

import (
	"fmt"
	"sort"

	"github.com/go-playground/validator/v10"
)

const (
	ErrorApplicationClassDuplicate = "validation.error.application_class_duplicate"
)

// Custom config validation in addition to one specified by annotation at schema.go.
func Validate(structureValidator validator.StructLevel) {
	validateUniqueApplicationClasses(structureValidator)
}

func validateUniqueApplicationClasses(structureValidator validator.StructLevel) {
	var classes []string
	config := structureValidator.Current().Interface().(Root)

	modes := make([]string, 0, len(config.Modes))
	for mode := range config.Modes {
		modes = append(modes, mode)
	}
	sort.Strings(modes)

	for _, mode := range modes {
		modeConfig := config.Modes[mode]
		duplicateClass, isFound := getCommonElement(classes, modeConfig.ApplicationClasses)
		if isFound {
			path := fmt.Sprintf("Modes[%s].ApplicationClasses", mode)
			structureValidator.ReportError(
				duplicateClass,
				path,
				path,
				ErrorApplicationClassDuplicate,
				duplicateClass,
			)
		}
		classes = append(classes, modeConfig.ApplicationClasses...)
	}
}

func getCommonElement[SliceType comparable](slice1 []SliceType, slice2 []SliceType) (result SliceType, isFound bool) {
	slice1Map := make(map[SliceType]bool, len(slice1))
	for _, slice1Value := range slice1 {
		slice1Map[slice1Value] = true
	}

	for _, slice2Value := range slice2 {
		if slice1Map[slice2Value] {
			return slice2Value, true
		}
	}
	return result, false
}
