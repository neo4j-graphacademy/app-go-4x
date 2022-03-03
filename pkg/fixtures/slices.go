package fixtures

func Slice(slice []map[string]interface{}, skip, limit int) []map[string]interface{} {
	start, end := sliceBounds(slice, skip, limit)
	return slice[start:end]
}

func sliceBounds(slice []map[string]interface{}, skip, limit int) (int, int) {
	maxIndex := len(slice) - 1
	return minInt(skip, minInt(skip, maxIndex)),
		minInt(skip+limit, minInt(skip+limit, maxIndex+1))
}

func minInt(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
