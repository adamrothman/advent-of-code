package main

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

func extractDigits(x int) (digits []int) {
	var digitCount int
	if x == 0 {
		digitCount = 1
	} else {
		digitCount = int(math.Log10(float64(x))) + 1
	}

	digits = make([]int, digitCount)
	for i := 0; i < digitCount; i++ {
		power := int(math.Pow10(i))
		digit := (x / power) % 10
		digits[i] = digit
	}
	if len(digits) > 1 {
		reverseSlice(digits)
	}

	return
}

func reverseSlice(s []int) {
	for i := len(s)/2 - 1; i >= 0; i-- {
		opp := len(s) - 1 - i
		s[i], s[opp] = s[opp], s[i]
	}
}

func formatSlice(s []int) string {
	builder := strings.Builder{}
	for _, digit := range s {
		builder.WriteString(strconv.FormatInt(int64(digit), 10))
	}
	return builder.String()
}

func formatScoreboard(sb []int, elfA, elfB int) string {
	builder := strings.Builder{}
	for i, score := range sb {
		var format string
		if i == elfA {
			format = "(%d)"
		} else if i == elfB {
			format = "[%d]"
		} else {
			format = " %d "
		}
		builder.WriteString(fmt.Sprintf(format, score))
	}
	return builder.String()
}

type StopFunc func(sb []int) bool

func simulateRecipesUntilStop(stop StopFunc) (scoreboard []int) {
	scoreboard = make([]int, 2)
	scoreboard[0] = 3
	scoreboard[1] = 7

	elfARecipe, elfBRecipe := 0, 1

outer:
	for {
		elfAScore, elfBScore := scoreboard[elfARecipe], scoreboard[elfBRecipe]

		// Get new scores
		for _, newScore := range extractDigits(elfAScore + elfBScore) {
			scoreboard = append(scoreboard, newScore)
			if stop(scoreboard) {
				break outer
			}
		}

		// Move elves forward
		elfARecipe = (elfARecipe + 1 + elfAScore) % len(scoreboard)
		elfBRecipe = (elfBRecipe + 1 + elfBScore) % len(scoreboard)
	}

	return
}

func main() {
	input := 110201

	scoreboard := simulateRecipesUntilStop(func(sb []int) bool {
		return len(sb) == input+10
	})
	fmt.Printf("After %d recipes, the scores of the next 10 are: %s\n", input, formatSlice(scoreboard[input:input+10]))

	inputDigits := extractDigits(input)
	inputDigitCount := len(inputDigits)
	scoreboard = simulateRecipesUntilStop(func(sb []int) bool {
		if len(sb) >= inputDigitCount {
			tail := sb[len(sb)-inputDigitCount:]
			if reflect.DeepEqual(tail, inputDigits) {
				return true
			}
		}
		return false
	})
	fmt.Printf("%d first appears after %d recipes\n", input, len(scoreboard)-inputDigitCount)
}
