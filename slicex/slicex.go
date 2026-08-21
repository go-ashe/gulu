package slicex

import "github.com/samber/lo"

// MergeUniq 按顺序合并多个集合并去重
func MergeUniq[T comparable](slices ...[]T) []T {
	var merged []T
	for _, s := range slices {
		merged = append(merged, s...)
	}
	return lo.Uniq(merged)
}
