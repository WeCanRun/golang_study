package _sort

import "testing"

func quickSort(nums []int) {
	helper(nums, 0, len(nums)-1)
}

func helper(nums []int, low, high int) {
	if low >= high {
		return
	}
	index := partition(nums, low, high)
	helper(nums, 0, index-1)
	helper(nums, index+1, high)
}

func partition(nums []int, low, high int) int {
	i, j := low+1, high
	for {
		for nums[i] > nums[low] && i < high {
			i++
		}
		for nums[j] < nums[low] && j > low {
			j--
		}
		if i >= j {
			break
		}
		swap(nums, i, j)
	}
	swap(nums, low, j)
	return j
}

func swap(nums []int, i, j int) {
	nums[i], nums[j] = nums[j], nums[i]
}

func TestQuickSort(t *testing.T) {
	nums := []int{14, 24, 5, 5, 524, 23, 434}
	quickSort(nums)
	t.Log(nums)
}
