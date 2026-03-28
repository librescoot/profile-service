package names

import (
	"fmt"
	"math/rand/v2"
)

func Generate() string {
	adj := adjectives[rand.IntN(len(adjectives))]
	animal := animals[rand.IntN(len(animals))]
	return fmt.Sprintf("%s %s", adj, animal)
}

func GenerateWithSeed(seed int) string {
	r := rand.New(rand.NewPCG(uint64(seed), uint64(seed*31)))
	adj := adjectives[r.IntN(len(adjectives))]
	animal := animals[r.IntN(len(animals))]
	return fmt.Sprintf("%s %s", adj, animal)
}
