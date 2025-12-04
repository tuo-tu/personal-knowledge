# 案例

## 并发目录大小统计

### 业务逻辑

统计目录的文件数量和大小（或其他信息）。示例输出：

```shell
// 某个目录：
2637 files 1149.87 MB
```

### 实现思路

- 给定一个或多个目录，并发的统计每个目录的size，最后累加到一起。
- 当目录中存在子目录时，递归的统计。
- 每个目录的统计都由独立的Goroutine完成
- 累计总Size由独立的Goroutine完成
- 使用Channel传递获取的文件大小
- 使用WaitGroup调度

### 核心代码

```go
// 读取目录内容
// os.ReadDir
func ReadDir(name string) ([]DirEntry, error)
entries, err := os.ReadDir(dir)

// 取得文件信息
info, err := entry.Info()

//判定是否为目录
entry.IsDir()
```

### 编码实现

```go
func WalkDir(dirs ...string) string {
     wg := &sync.WaitGroup{}
    
    // 用户未指定目录（dirs 为空），默认遍历当前目录（.）
    if len(dirs) == 0 {
        dirs = []string{"."}
    }
    
    // 在goroutine中传递文件大小，一次只能放入一个文件（存储文件的大小）
    filesizeCh := make(chan int64, 1)

    // 每个目录启动一个独立的 goroutine 进行遍历：
    for _, dir := range dirs {
        wg.Add(1)
        go walkDir(dir, filesizeCh, wg)
    }
    
    // 监控子goroutine是否完成，随后关闭通道
    go func(wg *sync.WaitGroup) {
        wg.Wait()
        close(filesizeCh)
    }(wg)

    // 不停从filesizeCh通道中读取文件大小，累加文件数量和总大小。
    var fileNum, sizeTotal int64
    for filesize := range filesizeCh {
        fileNum++ // 每获取一个，文件数量就加1
        sizeTotal += filesize
    }

    // 格式化输出文件数量及总大小，单位为MB
    return fmt.Sprintf("%d files %.2f MB\n", fileNum, float64(sizeTotal)/1e6)
}

// 将某个目录下的条所有文件的大小挨个放入fileSizes通道里，方便后续读取
func walkDir(dir string, fileSizes chan<- int64, wg *sync.WaitGroup) {
    defer wg.Done()
    
    // 遍历dir目录下的所有条目
    for _, fileinfo := range fileInfos(dir) {
        // 如果文件是文件夹，则递归调用 walkDir
        if fileinfo.IsDir() {
            // 将多个路径片段拼接成一个完整的路径。
            subDir := filepath.Join(dir, fileinfo.Name())
            wg.Add(1)
            // 继续递归子目录
            go walkDir(subDir, fileSizes, wg)
        } else {
            // 如果是普通文件，则将文件大小发送到通道
            fileSizes <- fileinfo.Size()
        }
    }
}

// 获取某个目录下所有条目的详细信息
func fileInfos(dir string) []fs.FileInfo {
    // 1.读取目录中的所有条目（文件和子目录）
    entries, err := os.ReadDir(dir)
    if err != nil {
        // 打印错误信息到标准错误输出
        fmt.Fprintf(os.Stderr, "walkdir: %v\n", err)
        return []fs.FileInfo{}
    }
    
    // 初始化FileInfo切片，用于存储某个目录中的所有条目信息
    infos := make([]fs.FileInfo, 0, len(entries))
    // 2.遍历条目并获取详细信息（如文件名、大小、权限等），将获取到的info追加到infos切片中。
    for _, entry := range entries {
        info, err := entry.Info()
        if err != nil {
            continue
        }
        infos = append(infos, info)
    }
    return infos
}
```

### 测试执行

```shell
> go test -run=WalkDir
70 files 0.09 MB

PASS
ok      goConcurrency   0.321s
```

## 快速排序的并发编程实现（后续再看）

### 典型的单线程快速排序实现

这是**基于单枢轴分区**的快速排序算法

```go
func QuickSortSingle(arr []int) []int {
    // 确保arr中至少存在2个或以上元素
    if arr == nil || len(arr) < 2 {
        return arr
    }
    // 执行排序
    quickSortSingle(arr, 0, len(arr)-1)
    return arr
}

func quickSortSingle(arr []int, l, r int) {
    // 判定待排序范围是否合法
    if l < r {
        // 获取参考元素位置索引
        mid := partition(arr, l, r)
        // 递归排序左边
        quickSortSingle(arr, l, mid-1)
        // 递归排序右边
        quickSortSingle(arr, mid+1, r)
    }
}

// 大小分区，将小于等于枢轴的元素移动到左侧，并返回参考元素索引
func partition(arr []int, l, r int) int {
    // 标记当前数组中小于等于枢轴元素的子数组的边界，表示当前没有任何元素被归类为“小于或等于枢轴”。
    p := l - 1
    for i := l; i <= r; i++ {
        if arr[i] <= arr[r] {
            // 必须要满足条件，p才会自增，更不会执行swap
            p++
            swap(arr, p, i)
        }
    }
    return p
}

// 交换arr中i和j元素
func swap(arr []int, i, j int) {
    t := arr[i]
    arr[i] = arr[j]
    arr[j] = t
}
```

### 并发编程实现思路

- **使用独立的Goroutine完成arr中左右部分的排序**
- WaitGroup完成等待阻塞同步

### 编码实现

```go
// QuickSortConcurrency 快速排序调用函数
func QuickSortConcurrency(arr []int) []int {
	// 一、校验arr是否满足排序需要，至少要有2个元素
	if arr == nil || len(arr) < 2 {
		return arr
	}

	// 同步的控制
	wg := &sync.WaitGroup{}
	// 二、执行排序
	// 初始排序整体[0, len(arr)-1]
	wg.Add(1)
	go quickSortConcurrency(arr, 0, len(arr)-1, wg)
	wg.Wait()

	// 三：返回结果
	return arr
}

// 实现递归快排的核心函数
// 接收arr，和排序区间的索引位置[l, r]
func quickSortConcurrency(arr []int, l, r int, wg *sync.WaitGroup) {
	// 一、-1wg的计数器
	defer wg.Done()

	// 二、判定是否需要排序， l < r
	if l < r {
		// 三、大小分区元素，并获取参考元素索引
		mid := partition(arr, l, r)

		// 四、并发对左部分排序
		wg.Add(1)
		go quickSortConcurrency(arr, l, mid-1, wg)

		// 五、并发的对右部分排序
		wg.Add(1)
		go quickSortConcurrency(arr, mid+1, r, wg)
	}
}
```

partition 和 swap 部分不变。

### 测试执行

```go
func TestQuickSortConcurrency(t *testing.T) {
	randArr := GenerateRandArr(1000)
	sortArr := QuickSortConcurrency(randArr)
	fmt.Println(sortArr)
}

// 生成大的随机数组
func GenerateRandArr(l int) []int {
	// 生产大量的随机数
	arr := make([]int, l)
	rand.Seed(time.Now().UnixMilli())
	for i := 0; i < l; i++ {
		arr[i] = int(rand.Int31n(int32(l * 5)))
	}

	return arr
}
```

```shell
> go test -run=QuickSortConcurrency
```
