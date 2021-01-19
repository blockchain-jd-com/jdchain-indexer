package handler

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

func validateFromTo(from, to, maxCount int64) (resultFrom, resultTo int64, ok bool) {
	if from < 0 ||
		to < 0 ||
		to < from {
		return
	}
	resultFrom = from
	if to > (resultFrom + maxCount) {
		resultTo = resultFrom + maxCount
	} else {
		resultTo = to
	}
	ok = true
	return
}
