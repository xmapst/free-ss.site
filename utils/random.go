package utils

import (
	"math/rand"
	"regexp"
	"time"
)

func Random(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func ExtractParenthesisContent(str string) (result []string) {
	rex := regexp.MustCompile(`\(([^)]+)\)`)
	out := rex.FindAllStringSubmatch(str, -1)
	for _, i := range out {
		result = append(result, i[1])
	}
	return
}
