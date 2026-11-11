package advent

func IsDigit(b byte) bool {
	return b >= '0' && b <= '9'
}

func ToDigit(c byte) int {
	return int(c - byte('0'))
}

// Atoi processes a byte stream for numbers returning the number and bytes processed
func Atoi(input []byte) (int, int) {
	result := 0
	i := 0
	for i = 0; i < len(input) && IsDigit(input[i]); i++ {
		result = result*10 + int(input[i]-byte('0'))
	}
	return result, i
}
