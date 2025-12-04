package main

/*
#题目描述
实现一个 HTTP 接口，接收一个 JSON 格式的请求体，包含一个整数数组 nums 和一个目标值 target，在 nums 中找到两个数，使它们的和等于 target，并返回它们的下标。
#要求:
	•	你可以使用任意熟悉的http框架，服务运行在 8080端口
	•	这是一个完整的项目，应包括但不限于代码的分层设计、代码组织方式、参数校验、状态码设计、错误处理、单元测试等
	•	允许使用搜索引擎或查阅官方文档，但禁止使用AI，一经发现直接淘汰
*/

func test(nums []int, target int) []int {
	hashtable := make(map[int]int, 0)
	for i, num := range nums {
		if j, ok := hashtable[target-num]; ok {
			return []int{i, j}
		}
		hashtable[num] = i
	}
	return nil
}

func main() {
	//arr := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	//target := 7
	//fmt.Println(test(arr, target))

}
