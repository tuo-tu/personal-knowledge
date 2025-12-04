package main

import "fmt"

// 题目1：一圈孩子报数，叫到的出圈，要求你输出所有孩子的编号
func josephusArray(n, target int) []int {
	// people用来存储孩子的编号，从1开始
	people := make([]int, n)
	for i := range people {
		people[i] = i + 1
	}

	res := []int{}
	index := 0 // 当前报数的起始位置
	for len(people) > 0 {
		index = (index + target - 1) % len(people)
		res = append(res, people[index])
		// 包前不包后
		people = append(people[:index], people[index+1:]...)
	}

	return res
}

// 破冰游戏
func iceBreakingGame(num int, target int) int {
	if num == 1 {
		return 0
	}
	prevRemaining := iceBreakingGame(num-1, target)
	index := (prevRemaining + target) % num
	return index
}

// 题目2：两个列表的交集
func intersection(nums1, nums2 []int) []int {
	exist, ans := map[int]struct{}{}, []int{}
	for _, x := range nums1 {
		exist[x] = struct{}{}
	}

	for _, x := range nums2 {
		if _, ok := exist[x]; ok {
			delete(exist, x)
			ans = append(ans, x)
		}
	}
	return ans
}

func main() {
	//fmt.Println(josephusArray(7, 4))
	//fmt.Println(iceBreakingGame(6, 4))
	fmt.Println(intersection([]int{1, 2, 3, 6}, []int{2, 6, 3, 4}))
}
