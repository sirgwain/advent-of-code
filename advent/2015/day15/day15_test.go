package day15

import "testing"

func Test_calculateScore(t *testing.T) {
	tests := []struct {
		name        string
		ingredients [][5]int
		count       [4]int
		want        int
	}{
		// Butterscotch: capacity -1, durability -2, flavor 6, texture 3, calories 8
		// Cinnamon: capacity 2, durability 3, flavor -2, texture -1, calories 3
		// 44 teaspooons butterscotch, 56 teaspoons cinammon
		{
			name: "44-56",
			ingredients: [][5]int{
				{-1, -2, 6, 3, 8},
				{2, 3, -2, -1, 3},
				{0, 0, 0, 0, 0}, // empty ingredients
				{0, 0, 0, 0, 0},
			},
			count: [4]int{44, 56, 0, 0},
			want:  62842880, // 68 * 80 * 152 * 76
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateScore(tt.ingredients, tt.count); got != tt.want {
				t.Errorf("calculateScore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bestScore(t *testing.T) {
	tests := []struct {
		name             string
		ingredients      [][5]int
		caloriesRequired int
		want             int
	}{
		// Butterscotch: capacity -1, durability -2, flavor 6, texture 3, calories 8
		// Cinnamon: capacity 2, durability 3, flavor -2, texture -1, calories 3
		// 44 teaspooons butterscotch, 56 teaspoons cinammon
		{
			name: "44-56",
			ingredients: [][5]int{
				{-1, -2, 6, 3, 8},
				{2, 3, -2, -1, 3},
				{0, 0, 0, 0, 0}, // empty ingredients
				{0, 0, 0, 0, 0},
			},
			caloriesRequired: -1,
			want:             62842880, // 44a and 56b
		},
		{
			name: "40-60 500 calories",
			ingredients: [][5]int{
				{-1, -2, 6, 3, 8},
				{2, 3, -2, -1, 3},
				{0, 0, 0, 0, 0}, // empty ingredients
				{0, 0, 0, 0, 0},
			},
			caloriesRequired: 500,
			want:             57600000, // 40a and 50b
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bestScore(tt.ingredients, tt.caloriesRequired); got != tt.want {
				t.Errorf("bestScore() = %v, want %v", got, tt.want)
			}
		})
	}
}
