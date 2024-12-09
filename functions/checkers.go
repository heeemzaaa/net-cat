package netcat

func Printable(message string) bool {
	slice := []rune(message)
	for i := 0; i < len(slice); i++ {
		if !(slice[i] >= 32 && slice[i] <= 126) {
			return false
		}
	}
	return true
}

func SpaceName(name string) bool {
	slice := []rune(name)
	for i := 0; i < len(slice); i++ {
		if slice[i] == ' ' {
			return true
		}
	}
	return false
}
