package generator

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"time"
)

func GenerateUniqueNumbers(n, digits int) ([]int, error) {
	max := int(math.Pow10(digits))

	if n > max {
		return nil, fmt.Errorf("to many participants: %d > %d", n, max)
	}

	numbers := make([]int, max)
	for i := 0; i < max; i++ {
		numbers[i] = i
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(numbers), func(i, j int) {
		numbers[i], numbers[j] = numbers[j], numbers[i]
	})

	log.Println("NUMBERS", numbers[:n+10])

	return numbers[:n], nil
}
