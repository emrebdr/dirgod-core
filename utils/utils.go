package utils

func Contains(array []any, value any) bool {
	for _, v := range array {
		if v == value {
			return true
		}
	}
	return false
}
