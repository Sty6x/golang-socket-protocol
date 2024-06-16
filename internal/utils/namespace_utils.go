package utils

func Check_existing_user_connection(connections []string, target string) bool {
	for _, id := range connections {
		if id == target {
			return true
		}
	}
	return false
}
