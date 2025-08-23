package utils

import (
	"encoding/json"
	"sort"
)

func sortMapRecursively(input any) any {
	inputMap, ok := input.(map[string]any)
	if !ok {
		return input
	}

	sorted := make(map[string]any)
	keys := make([]string, 0, len(inputMap))

	// Collect and sort keys
	for k := range inputMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Rebuild map in sorted order
	for _, k := range keys {
		v := inputMap[k]
		switch nested := v.(type) {
		case map[string]any:
			sorted[k] = sortMapRecursively(nested)
		default:
			sorted[k] = v
		}
	}

	return sorted
}

// MarshalSortedJSON recursively sorts keys in a struct
func MarshalSortedJSON(v any) ([]byte, error) {
	unsortedJSONBytes, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var unsortedData any
	if err := json.Unmarshal(unsortedJSONBytes, &unsortedData); err != nil {
		return nil, err
	}

	sortedData := sortMapRecursively(unsortedData)

	sortedJSONBytes, err := json.MarshalIndent(sortedData, "", "\t")
	if err != nil {
		return nil, err
	}

	// add the new line
	sortedJSONBytes = append(sortedJSONBytes, []byte("\n")...)

	return sortedJSONBytes, nil
}
