package main

import "math"

func minFallingPathSum(matrix [][]int) int {
	if len(matrix) == 0 || len(matrix[0]) == 0 {
		return 0
	}
	// m,n 分别代表行数和列数
	row, col := len(matrix), len(matrix[0])
	// dp[i][j]代表自上而下的节点的最短路径
	dp := make([][]int, row)
	for i := range dp {
		dp[i] = make([]int, col)
	}
	// 第一行沿用 matrix 的数字
	for i := 0; i < col; i++ {
		dp[0][i] = matrix[0][i]
	}

	// 开始填充剩下的位置
	for i := 1; i < row; i++ {
		for j := 0; j < col; j++ {
			minPre := math.MaxInt
			// 左上方的
			if j > 0 {
				minPre = min(minPre, dp[i-1][j-1])
			}
			// 正上方的
			minPre = min(minPre, dp[i-1][j])
			// 右上方的
			if j+1 < col {
				minPre = min(minPre, dp[i-1][j+1])
			}
			// 不要忘记赋值给dp[i][j]
			dp[i][j] = matrix[i][j] + minPre
		}
	}
	// 直接获取最后一行的结果即可
	minPath := math.MaxInt
	for i := 0; i < col; i++ {
		minPath = min(minPath, dp[row-1][i])
	}
	return minPath
}
