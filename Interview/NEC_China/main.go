package main

import "fmt"

//func formatSize(bytes int64) string {
//	str := []string{"KB", "MB", "GB", "TB", "PB", "EB"}
//	const unit = 1024
//	if bytes < unit {
//		return fmt.Sprintf("%d B", bytes)
//	}
//	div, index := int64(unit), 0
//	for n := bytes / unit; n >= unit; n /= unit {
//		div *= unit
//		index++
//	}
//	// bytes不会变，变的是div
//	return fmt.Sprintf("%.2f %s", float64(bytes)/float64(div), str[index])
//}

func main() {
	//fmt.Println(formatSize(1023))
	//fmt.Println(formatSize(1025))
	//fmt.Println(formatSize(10467028))
	//fmt.Println(formatSize(1024 * 10467028))
	//buf := make([]int, 4, 5)
	buf := []int{1, 2, 3, 4, 5}
	fmt.Println(cap(buf))
	buf = buf[5:]
	fmt.Println(buf)
}
