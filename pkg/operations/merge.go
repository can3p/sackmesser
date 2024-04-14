package operations

import (
	"github.com/can3p/sackmesser/pkg/traverse/types"
	"github.com/pkg/errors"
)

func Merge(root types.Node, path []string, args ...any) error {
	if len(args) != 1 {
		return errors.Errorf("set operation expects one argument")
	}

	value, ok := args[0].(map[string]any)
	if !ok {
		return errors.Errorf("Merge expects a json as an argument")
	}

	if len(path) == 1 {
		fieldName := path[0]
		fieldVal, err := root.GetField(fieldName)

		if err == types.ErrFieldMissing {
			return root.SetField(fieldName, value)
		}

		if err != nil {
			return err
		}

		fieldVal = mergeObject(fieldVal, value)
		return root.SetField(fieldName, fieldVal)
	}

	node, err := root.Visit(path[0])

	if err != nil {
		return err
	}

	return Merge(node, path[1:], value)
}

func mergeObject(existingValue any, value map[string]any) any {
	typed, ok := existingValue.(map[string]any)

	if !ok {
		return value
	}

	for fieldName, value := range value {
		subfieldValue, exists := typed[fieldName]

		if !exists {
			typed[fieldName] = value
			continue
		}

		typedSubfieldValue, subfieldIsMap := subfieldValue.(map[string]any)
		typedNewValue, newValueIsMap := value.(map[string]any)

		if !subfieldIsMap || !newValueIsMap {
			typed[fieldName] = value
			continue
		}

		typed[fieldName] = mergeObject(typedSubfieldValue, typedNewValue)
	}

	return typed
}
