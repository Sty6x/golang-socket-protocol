package server

func CheckExistingUserConnnection(connections []string, target string) bool {
	for _, id := range connections {
		if id == target {
			return true
		}
	}
	return false
}
