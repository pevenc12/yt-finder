package helper

import (
	"errors"
	"strconv"
)

// ParseFlags extracts parameters in flag
func ParseFlags(inputs []string) ([]string, int, error) {
	var terms []string
	var day int = 7
	var setPeriod bool = false
	var prefix string = "--search-period"
	for _, fl := range inputs {
		if setPeriod == true {
			period, errPeriod := strconv.Atoi(fl)
			if errPeriod != nil {
				return []string{}, 0, errors.New("invalid period")
			}
			day = period
			break
		} else {
			// if --search-period flag exists
			if fl == prefix {
				setPeriod = true
				continue
			} else {
				terms = append(terms, fl)
			}
		}
	}

	// Check if search terms are empty
	if len(terms) == 0 {
		return []string{}, 0, errors.New("invalid search term")
	}

	return terms, day, nil
}
