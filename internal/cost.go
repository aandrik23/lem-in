package internal

// capacityForTurn returns how many total ants can finish by turn T,
// given each path’s length L_i (number of edges).
func capacityForTurn(costs []int, T int) int {
	sum := 0
	for _, L := range costs {
		if T > L {
			sum += T - L
		}
	}
	return sum
}

// computeAntsPerPath allocates exactly totalAnts across paths with costs L_i
// so that they all finish in the minimum T turns.
func ComputeAntsPerPath(costs []int, totalAnts int) []int {
	k := len(costs)
	antsPerPath := make([]int, k)

	// Find minimal T by binary‐search
	minL, maxL := costs[0], costs[0]
	for _, L := range costs {
		if L < minL {
			minL = L
		}
		if L > maxL {
			maxL = L
		}
	}
	// Lower bound for T is minL+1, upper bound minL+totalAnts (worst case)
	lo, hi := minL+1, minL+totalAnts
	for lo < hi {
		mid := (lo + hi) / 2
		if capacityForTurn(costs, mid) >= totalAnts {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	T := lo

	// Now assign Ai = max(0, T - Li)
	sum := 0
	for i, L := range costs {
		cap := T - L
		if cap < 0 {
			cap = 0
		}
		antsPerPath[i] = cap
		sum += cap
	}

	// If we allocated too many, remove the excess from the longest paths
	excess := sum - totalAnts
	for i := k - 1; excess > 0 && i >= 0; i-- {
		if antsPerPath[i] > 0 {
			antsPerPath[i]--
			excess--
		}
	}
	return antsPerPath
}
