// Copyright 2019 Changkun Ou. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package algo

// InsertSort ...
func InsertSort(arr []int) []int {
	sorted := arr
	for i := 1; i < len(sorted); i++ {
		for j := i - 1; j >= 0; j-- {
			if sorted[j] < sorted[j+1] {
				break
			}
			sorted[j], sorted[j+1] = sorted[j+1], sorted[j]
		}
	}
	return sorted
}

// MergeSort ...
func MergeSort(arr []int) []int {
	if len(arr) <= 1 {
		return arr
	}
	num := len(arr) / 2
	l := MergeSort(arr[:num])
	r := MergeSort(arr[num:])
	return merge(l, r)
}

func merge(left []int, right []int) (result []int) {
	l, r := 0, 0
	for l < len(left) && r < len(right) {
		if left[l] < right[r] {
			result = append(result, left[l])
			l++
		} else {
			result = append(result, right[r])
			r++
		}
	}
	result = append(result, left[l:]...)
	result = append(result, right[r:]...)
	return
}
