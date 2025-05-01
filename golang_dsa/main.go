package main

import (
	"fmt"
	"sort"
)

func singleNumber(nums []int) int {
	if nums[0] != nums[1] {
		return nums[0]
	}
	if len(nums) == 1 {
		return nums[0]
	}
	if nums[len(nums)-1] != nums[len(nums)-2] {
		return nums[len(nums)-1]
	}
    var uniqueElem int
    sort.Ints(nums)
    for i := 1; i < len(nums)-1 ; i++ {
        if nums[i] != nums[i-1] && nums[i] != nums[i+1] {
			uniqueElem = nums[i]
        }else {
			continue
		}
    }
    return uniqueElem
}

func main() {
	sl1 := []int{1,2,2,3,3}
	uniqueElem := singleNumber(sl1)
	fmt.Print(uniqueElem)
}

