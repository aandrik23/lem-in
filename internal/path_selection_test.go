package internal

import (
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestSelectOptimalDisjointPathSet(t *testing.T) {
	// 2 disjoint simple paths, 1 overlapping path
	allPaths := [][]string{
		{"start", "A", "end"},
		{"start", "B", "end"},
		{"start", "A", "B", "end"}, // overlaps with above
	}

	ants := 4
	// The optimal is splitting ants between two disjoint paths:
	wantPaths := [][]string{
		{"start", "A", "end"},
		{"start", "B", "end"},
	}
	wantQuota := []int{2, 2} // order can be swapped

	gotPaths, gotQuota := SelectOptimalDisjointPathSet(ants, allPaths)

	// Compare ignoring order
	if !comparePathSets(gotPaths, wantPaths) {
		t.Errorf("SelectOptimalDisjointPathSet() got paths %v, want %v", gotPaths, wantPaths)
	}
	if !sameQuotaUnordered(gotQuota, wantQuota) {
		t.Errorf("SelectOptimalDisjointPathSet() got quota %v, want %v", gotQuota, wantQuota)
	}
}

func comparePathSets(a, b [][]string) bool {
	normalize := func(paths [][]string) []string {
		var out []string
		for _, path := range paths {
			out = append(out, strings.Join(path, ","))
		}
		sort.Strings(out)
		return out
	}
	return reflect.DeepEqual(normalize(a), normalize(b))
}

func sameQuotaUnordered(a, b []int) bool {
	aCopy := append([]int{}, a...)
	bCopy := append([]int{}, b...)
	sort.Ints(aCopy)
	sort.Ints(bCopy)
	return reflect.DeepEqual(aCopy, bCopy)
}
