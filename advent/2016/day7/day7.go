package day7

type Day struct{}

func (d *Day) Run(input []byte) (int, int, error) {
	var s1, s2 int

	foundAbba := false
	insideHypernet := false
	for i := 0; i < len(input); i++ {

		if input[i+3] == '\n' {
			// end of the line, update state and skip the next
			if foundAbba {
				s1++
			}

			foundAbba = false
			insideHypernet = false
			i += 3 // skip the last 3 chars
			continue
		}
		if input[i+3] == '[' {
			insideHypernet = true
			i += 3
			continue
		}
		if input[i+3] == ']' {
			insideHypernet = false
			i += 3
			continue
		}
		// look for ABBA patterns
		// where a1 == b1 && a2 == b2
		// a1 b1   b2 a2
		// A  B  = B  A
		a1 := input[i]
		b1 := input[i+1]
		b2 := input[i+2]
		a2 := input[i+3]

		if a1 == a2 && b1 == b2 && a1 != b1 {
			if insideHypernet {
				// this invalidates the ip
				foundAbba = false
				for input[i] != '\n' {
					i++
				}
				continue
			}
			foundAbba = true
		}
	}

	return s1, s2, nil
}
