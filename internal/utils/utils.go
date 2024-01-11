package utils

import (
	"fmt"

	"github.com/agnivade/levenshtein"
)

var DataKey []string

// GetBestMatch return bestMatch with bestDistance
func GetBestMatch(searchTerm string) (string, int) {
	bestMatch := ""
	bestDistance := -1 // Initialisation avec une valeur impossible

	for _, str := range DataKey {
		distance := levenshtein.ComputeDistance(searchTerm, str)
		if bestDistance == -1 || distance < bestDistance {
			bestDistance = distance
			bestMatch = str
		}
	}
	fmt.Println(bestDistance)
	return bestMatch, bestDistance
}
