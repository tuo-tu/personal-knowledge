package main

import (
	"fmt"
	"math"
)

func test1(A, B []int) []int {
	hash := map[int]bool{}
	res := []int{}

	for _, v := range B {
		hash[v] = true
	}

	for i, v := range A {
		if _, ok := hash[v]; ok {
			res = append(res, i)
		}
	}
	return res
}

type Node struct {
	ID  string
	PID string
}

func test2(nodes []Node) []string {
	hash := map[string]string{}
	res := []string{}
	for _, node := range nodes {
		hash[node.ID] = node.PID
	}

	var path func(string) string
	path = func(id string) string {
		pid := hash[id]
		if pid == "-1" {
			return id
		}
		return path(pid) + "/" + id
	}

	for _, node := range nodes {
		res = append(res, "/"+path(node.ID)+"/")
	}
	return res
}

func test3(matrix [][]int) int {
	if len(matrix) == 0 {
		return 0
	}
	m, n := len(matrix), len(matrix[0])

	dp := make([][]int, m)
	for i := range dp {
		dp[i] = make([]int, n)
	}

	for i := 0; i < n; i++ {
		dp[0][i] = matrix[0][i]
	}

	// minPre := math.MaxInt
	for i := 1; i < m; i++ {

		for j := 0; j < n; j++ {
			minPre := math.MaxInt
			if j > 0 {
				minPre = min(minPre, dp[i-1][j-1])
			}
			minPre = min(minPre, dp[i-1][j])
			if j+1 < n {
				minPre = min(minPre, dp[i-1][j+1])
			}

			dp[i][j] = matrix[i][j] + minPre
		}
	}

	res := math.MaxInt
	for i := 0; i < n; i++ {
		res = min(res, dp[m-1][i])
	}
	return res
}

//func min(x, y int) int {
//	if x > y {
//		return y
//	}
//	return x
//}

func main() {
	matrix := [][]int{
		{5, 8, 1, 2},
		{4, 1, 7, 3},
		{3, 6, 2, 9},
	}
	fmt.Println(test3(matrix))
}
