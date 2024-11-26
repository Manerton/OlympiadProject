package parsing

import (
	"fmt"
	"strconv"
)

// Returns offset, limit, error
func ParseOffsetLimit(offsetStr, limitStr string) (*int, *int, error) {
	var offsetRes *int = nil
	var limitRes *int = nil
	if offsetStr != "" {
		tempOffset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse offset: %w", err)
		}
		offsetRes = &tempOffset
	}

	if limitStr != "" {
		tempLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse limit: %w", err)
		}
		limitRes = &tempLimit
	}
	return offsetRes, limitRes, nil
}
