package utils

import (
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

type tempRank struct {
	list.Rank
	Score float64
}

// custom filtering
func CustomSubstringFilter(term string, targets []string) []list.Rank {
	if term == "" {
		ranks := make([]list.Rank, len(targets))
		for i := range targets {
			ranks[i] = list.Rank{Index: i}
		}
		return ranks
	}

	termLower := strings.ToLower(term)
	var matchingTempRanks []tempRank

	for i, target := range targets {
		targetLower := strings.ToLower(target)
		currentScore := 0.0
		matchedIndexes := []int{}

		// --- 1. Check for Exact Match ---
		if targetLower == termLower {
			currentScore = 2.0

			for k := 0; k < len(target); k++ {
				matchedIndexes = append(matchedIndexes, k)
			}
		} else if strings.Contains(targetLower, termLower) {
			// --- 2. Check for Substring Containment ---
			currentScore = 1.0

			start := 0
			for {
				idx := strings.Index(targetLower[start:], termLower)
				if idx == -1 {
					break
				}
				actualIdx := start + idx
				for k := 0; k < len(termLower); k++ {
					matchedIndexes = append(matchedIndexes, actualIdx+k)
				}
				start = actualIdx + len(termLower)
				if start >= len(targetLower) {
					break
				}
			}
		} else {
			continue
		}

		matchingTempRanks = append(matchingTempRanks, tempRank{
			Rank: list.Rank{
				Index:          i,
				MatchedIndexes: matchedIndexes,
			},
			Score: currentScore,
		})
	}

	// Sort: Exact matches (Score 2.0) first, then substring matches (Score 1.0).
	// Within the same score, maintain original input order (via Index).
	sort.Slice(matchingTempRanks, func(i, j int) bool {
		if matchingTempRanks[i].Score != matchingTempRanks[j].Score {
			return matchingTempRanks[i].Score > matchingTempRanks[j].Score // Higher score first
		}
		return matchingTempRanks[i].Index < matchingTempRanks[j].Index // Original index for tie-breaking
	})

	// Convert tempRank slice back to list.Rank slice
	finalRanks := make([]list.Rank, len(matchingTempRanks))
	for i, tr := range matchingTempRanks {
		finalRanks[i] = tr.Rank
	}

	return finalRanks
}
