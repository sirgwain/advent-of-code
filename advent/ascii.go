package advent

func IsDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func ToDigit(c byte) int {
	return int(c - byte('0'))
}

// Atoi progresses 
func Atoi(input []byte) int {
	result := 0
	for i := 0; i < len(input) && IsDigit(input[i]); i++ {
		result = result*10 + int(input[i]-byte('0'))
	}
	return result
}
