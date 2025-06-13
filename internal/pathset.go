package internal

func computeMakespan(costs, quotas []int) int {
	max := 0
	for i, q := range quotas {
		if q == 0 {
			continue
		}
		t := costs[i] + q
		if t > max {
			max = t
		}
	}
	return max
}

func allDisjointSets(paths [][]string, k int) [][][]string {
	var results [][][]string
	var helper func(start int, combo [][]string, used map[string]bool)
	helper = func(start int, combo [][]string, used map[string]bool) {
		if len(combo) == k {
			cp := make([][]string, k)
			copy(cp, combo)
			results = append(results, cp)
			return
		}
		for i := start; i < len(paths); i++ {
			p := paths[i]
			overlap := false
			for _, room := range p[1 : len(p)-1] {
				if used[room] {
					overlap = true
					break
				}
			}
			if overlap {
				continue
			}
			// mark rooms as used
			for _, room := range p[1 : len(p)-1] {
				used[room] = true
			}
			helper(i+1, append(combo, p), used)
			// unmark
			for _, room := range p[1 : len(p)-1] {
				used[room] = false
			}
		}
	}
	helper(0, [][]string{}, map[string]bool{})
	return results
}

func SelectOptimalDisjointPathSet(ants int, bestStepDisjointPaths [][]string) ([][]string, []int) {
	bestSteps := 1 << 30 // some large number
	var bestSet [][]string
	var bestQuota []int

	// Try all numbers of paths from 1 up to len(allPaths)
	for k := 1; k <= len(bestStepDisjointPaths); k++ {
		sets := allDisjointSets(bestStepDisjointPaths, k)
		for _, set := range sets {
			costs := make([]int, len(set))
			for i, p := range set {
				costs[i] = len(p) - 1
			}
			quota := ComputeAntsPerPath(costs, ants)
			steps := computeMakespan(costs, quota)
			if steps < bestSteps {
				bestSteps = steps
				bestSet = set
				bestQuota = quota
			}
		}
	}
	return bestSet, bestQuota
}
