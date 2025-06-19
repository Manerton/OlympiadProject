package parser

import (
	"fmt"
	"strconv"
)

func ParsePageLimit(pageStr string, limitStr string) (*int, *int, error) {
	var page *int = nil
	var limit *int = nil

	if pageStr != "" {
		tempPage, err := strconv.Atoi(pageStr)
		if err != nil {
			return nil, nil, fmt.Errorf("failed parse page: %w", err)
		}
		page = &tempPage
	}

	if limitStr != "" {
		tempLimit, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, nil, fmt.Errorf("failed parse limit: %w", err)
		}
		limit = &tempLimit
	}

	return page, limit, nil
}
