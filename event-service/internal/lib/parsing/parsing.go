package parsing

import (
	"fmt"
	"strconv"
)

// Returns offset, limit, error
func ParsePageLimitToOffsetLimit(pageStr, limitStr string) (*int, *int, error) {
	offsetRes := new(int)
	var limitRes *int = nil
	var pageRes *int = nil

	if pageStr != "" {
		tempPage, err := strconv.Atoi(pageStr)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse offset: %w", err)
		}
		pageRes = &tempPage
	}

	if limitStr != "" {
		tempLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to parse limit: %w", err)
		}
		limitRes = &tempLimit
	}

	if pageRes != nil && limitRes != nil {
		*offsetRes = (*pageRes - 1) * (*limitRes)
	}

	return offsetRes, limitRes, nil
}
