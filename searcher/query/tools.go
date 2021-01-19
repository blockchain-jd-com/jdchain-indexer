package query

import "strings"

func mapFromToToOffset(from, to, max int64) (count, offset int64) {
	if from < 0 {
		from = 0
	}
	offset = from - 1
	if offset < 0 {
		offset = 0
	}

	count = to - from + 1
	if count < 0 {
		count = 0
	}

	if max > 0 {
		if count > max {
			count = max
		}
	}
	return
}

// split string to single and double charactors, eg. abc -> a b c ab bc
func stringSplit(s string) (ss []string) {
	if len(s) <= 0 {
		return
	}

	singles := strings.Split(s, "")
	var doubles []string
	for i := 0; i < len(singles)-2; i++ {
		doubles = append(doubles, singles[i]+singles[i+1])
	}
	return append(singles, doubles...)
}
