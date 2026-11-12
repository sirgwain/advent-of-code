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

// AtoiFixed parses a fixed-width, space-padded unsigned number such as " 6" or "41".
// Every byte must be a digit or a space: no signs, and no validation.
// Other bytes produce garbage (b&0xF), e.g. 'A' (0x41) is read as 1.
func AtoiFixed(b []byte) int {
	n := 0
	for _, c := range b {
		n = n*10 + int(c&0xF)
	}
	return n
}

// AtoiFixed2 is AtoiFixed for exactly 2 bytes (b[0] and b[1]), unrolled for speed.
// Use it for 2-column fields; it is much faster than the AtoiFixed loop.
func AtoiFixed2(b []byte) int {
	return int(b[0]&0xF)*10 + int(b[1]&0xF)
}
