package util

func StringSearch(data []string, search string) int {
	for index, value := range data {
		if search == value {
			return index
		}
	}
	return -1
}
