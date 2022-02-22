package utils

func IndexOf(e string, slice []string) int {
	for k, v := range slice {
		if e == v {
			return k
		}
	}
	return -1 //not found.
}
