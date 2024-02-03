package utils

import "sort"

// Filter applies a callback function to each item in a slice and returns a new slice containing only the items for which the callback returns true.
// The callback function receives an item and its index as arguments and returns a boolean and an error.
// If the callback function returns an error, the filtering process stops and the error is returned.
func Filter[T any](items []T, callbackFn func(item T, index int) (bool, error)) ([]T, error) {
	rv := []T{}

	for i, item := range items {
		result, err := callbackFn(item, i)
		if err != nil {
			return nil, err
		}
		if result {
			rv = append(rv, item)
		}
	}

	return rv, nil
}

// Map applies a callback function to each item in a slice and returns a new slice of the results.
// The callback function receives an item and its index as arguments and returns a new value and an error.
// If the callback function returns an error, the mapping process stops and the error is returned.
func Map[T any, R any](items []T, callbackFn func(item T, index int) (R, error)) ([]R, error) {
	rv := []R{}

	for i, item := range items {
		result, err := callbackFn(item, i)
		if err != nil {
			return nil, err
		}
		rv = append(rv, result)
	}

	return rv, nil
}

// Keys extracts the keys from a map and returns them as a slice.
// The order of the keys in the returned slice is not specified and may vary.
func Keys[K comparable, V any](items map[K]V) []K {
	keys := []K{}

	for key := range items {
		keys = append(keys, key)
	}

	return keys
}

// MapToSlice converts a map with integer keys to a slice.
// It sorts the keys and populates the slice with values in that order.
func MapToSlice[V any](m map[int]V) ([]V, error) {
	keys := Keys(m) // キーの取得
	sort.Ints(keys) // キーをソート

	// ソートされたキーに基づいてスライスを生成
	return Map(keys, func(k int, _ int) (V, error) {
		return m[k], nil
	})
}
