package advent

func IsDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func ToDigit(c byte) int {
	return int(c - byte('0'))
}
