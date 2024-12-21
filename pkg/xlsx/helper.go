package xlsx

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func GetByHeader[T any](sheet [][]string, header string) ([]T, error) {
	var index int = -1
	if len(sheet) > 0 {
		// Find the column index for the header
		for i, h := range sheet[0] {
			if strings.TrimSpace(h) == header {
				index = i
				break
			}
		}
	}

	if index == -1 {
		return nil, errors.New("did not find the header: " + header)
	}

	var resList []T
	for i, row := range sheet {
		if i == 0 {
			continue // Skip header row
		}
		if len(row) > index { // Ensure row has enough columns
			val := strings.TrimSpace(row[index])
			if val != "" {
				var typedVal T
				switch any(typedVal).(type) {
				case string:
					resList = append(resList, any(val).(T))
				case int:
					intVal, err := strconv.Atoi(val)
					if err != nil {
						return nil, fmt.Errorf("failed to convert value %q to int: %w", val, err)
					}
					resList = append(resList, any(intVal).(T))
				default:
					return nil, errors.New("type conversion not supported")
				}
			}
		}
	}
	return resList, nil
}
