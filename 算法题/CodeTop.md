#### 3.无重复字符的最长子串

给定一个字符串 `s` ，请你找出其中不含有重复字符的 **最长子串** 的长度。

```go
// 滑动窗口 + map
func lengthOfLongestSubstring(s string) int {
	buf := make(map[byte]int)
	left := 0
	maxLen := 0
	// 不能用 for range，会导致类型不匹配
	for right := 0; right < len(s); right++ {
		if idx, ok := buf[s[right]]; ok && idx >= left {
			// idx此时还未更新
			left = idx + 1
		}
		// 更新字符的位置
		buf[s[right]] = right
		maxLen = max(maxLen, right-left+1)
	}
	return maxLen
}
```

#### 146.LRU 缓存

请你设计并实现一个满足 [LRU (最近最少使用) 缓存](https://baike.baidu.com/item/LRU) 约束的数据结构。

实现 `LRUCache` 类：

- `LRUCache(int capacity)` 以 **正整数** 作为容量 `capacity` 初始化 LRU 缓存
- `int get(int key)` 如果关键字 `key` 存在于缓存中，则返回关键字的值，否则返回 `-1` 。
- `void put(int key, int value)` 如果关键字 `key` 已经存在，则变更其数据值 `value` ；如果不存在，则向缓存中插入该组 `key-value` 。如果插入操作导致关键字数量超过 `capacity` ，则应该 **逐出** 最久未使用的关键字。

函数 `get` 和 `put` 必须以 `O(1)` 的平均时间复杂度运行。

```go
type Node struct {
    Key   int
    Value int
    Pre   *Node
    Next  *Node
}

// 注意，哨兵头节点是不会变的，他们只是边界标记
type LRUCache struct {
    Cap   int
    Cache map[int]*Node
    Head  *Node
    Tail  *Node
}

func Constructor(capacity int) LRUCache {
    Cache := make(map[int]*Node)
    head := &Node{}
    tail := &Node{}
    head.Next = tail
    tail.Pre = head
    return LRUCache{
       Cap:   capacity,
       Cache: Cache,
       Head:  head,
       Tail:  tail,
    }
}

func (l *LRUCache) Get(key int) int {
    if node, ok := l.Cache[key]; ok {
       // node移到头部
       l.MoveToHead(node)
       return node.Value
    }
    return -1
}

func (l *LRUCache) Put(key int, value int) {
    if node, ok := l.Cache[key]; ok {
       l.Cache[key].Value = value
       l.MoveToHead(node)
       // 如果不用else，这里要return
    } else {
       // 添加新节点到头部
       newNode := &Node{
          Key:   key,
          Value: value,
       }
       l.AddToHead(newNode)
       if len(l.Cache) > l.Cap {
          // 删除尾巴
          l.RemoveTail()
       }
    }
}

func (l *LRUCache) MoveToHead(node *Node) {
    node.Pre.Next = node.Next
    node.Next.Pre = node.Pre
    node.Pre = l.Head
    node.Next = l.Head.Next
    l.Head.Next.Pre = node // 关键步骤
    l.Head.Next = node
}

func (l *LRUCache) AddToHead(node *Node) {
    l.Cache[node.Key] = node
    node.Next = l.Head.Next
    node.Pre = l.Head
    l.Head.Next.Pre = node
    l.Head.Next = node
}

func (l *LRUCache) RemoveTail() {
    tail := l.Tail.Pre
    tail.Pre.Next = l.Tail
    l.Tail.Pre = tail.Pre
    tail.Pre = nil
    tail.Next = nil
    delete(l.Cache, tail.Key)
}
```

#### 206. 反转链表

给你单链表的头节点 `head` ，请你反转链表，并返回反转后的链表。

```go
func reverseList(head *ListNode) *ListNode {
    // 不能写成 pre := &ListNode{}，因为这里不需要哨兵节点，如果这样写，最后尾巴会多一个0值节点
    var pre *ListNode
    cur := head
    for cur != nil {
       next := cur.Next // 先记录下一个节点的位置
       cur.Next = pre   // 反转当前节点
       pre = cur        // pre 先前进
       cur = next       // cur 再前进
    }
    return pre
}
```

#### 215. 数组中的第K个最大元素

给定整数数组 `nums` 和整数 `k`，请返回数组中第 `k` 个最大的元素。

请注意，你需要找的是数组排序后的第 `k` 个最大的元素，而不是第 `k` 个不同的元素。

你必须设计并实现时间复杂度为 `O(n)` 的算法解决此问题。

```go
// 快速排序 + for循环
func findKthLargest(nums []int, k int) int {
    // 第 k 大的元素的下标
    target := len(nums) - k
    left := 0
    right := len(nums) - 1
    for {
       index := partion(nums, left, right)
       if index == target {
          return nums[index]
       } else if index < target { // 说明目标在右边
          left = index + 1
       } else { // 目标在左边
          right = index - 1
       }
    }
}

// 返回分区的分分界点下标
func partion(nums []int, left, right int) int {
    // 不要写成div := right，会增加复杂度
    div_value := nums[right]
    i := left
    for j := left; j < right; j++ {
       // 小于才交换，大于等下一次交换
       if nums[j] < div_value { // 注意是比较值，不是比较索引
          nums[i], nums[j] = nums[j], nums[i]
          // 这时候i直接往前推进
          i++
       }
    }
    // 将right放到正确的位置
    nums[i], nums[right] = nums[right], nums[i]
    return i
}
```

#### 25. K 个一组翻转链表

给你链表的头节点 `head` ，每 `k` 个节点一组进行翻转，请你返回修改后的链表。

`k` 是一个正整数，它的值小于或等于链表的长度。如果节点总数不是 `k` 的整数倍，那么请将最后剩余的节点保持原有顺序。

你不能只是单纯的改变节点内部的值，而是需要实际进行节点交换。

```go
// hair + dummy 节点
func reverseKGroup(head *ListNode, k int) *ListNode {
    // hair 永远不会动,指向头
    hair := &ListNode{Next: head}
    // 一开始的dummy也指向头，但是后续会变
    dummy := hair
    // 需要找到每一段的 head、tail节点,每一组作为 for 的一个循环
    for head != nil {
       // 需要通过 dummy 来找到这一组的尾巴
       tail := dummy
       for i := 0; i < k; i++ {
          tail = tail.Next
          // 如果这组元素不足，直接返回结果
          if tail == nil {
             return hair.Next
          }
       }
       // 记录下一组的开头
       nextHead := tail.Next
       // 翻转，返回新的头和尾
       newHead, newTail := reverse(head, tail)
       // 整合连接上
       dummy.Next = newHead
       newTail.Next = nextHead
       // 下一组开始
       // dummy等于上一组的尾节点
       dummy = newTail
       head = nextHead
    }
    return hair.Next
}

func reverse(head, tail *ListNode) (*ListNode, *ListNode) {
    // 记录本组的尾节点的下一个节点，从这里开始翻转
    pre := tail.Next
    cur := head
    for pre != tail {
       next := cur.Next
       // 当前节点指向前驱
       cur.Next = pre
       pre = cur
       cur = next
    }
    return tail, head
}
```

#### 15. 三数之和

给你一个整数数组 `nums` ，判断是否存在三元组 `[nums[i], nums[j], nums[k]]` 满足 `i != j`、`i != k` 且 `j != k` ，同时还满足 `nums[i] + nums[j] + nums[k] == 0` 。请你返回所有和为 `0` 且不重复的三元组。

**注意：**答案中不可以包含重复的三元组。

```go
// 先排序 + 双指针
func threeSum(nums []int) [][]int {
    // 首先进行排序
    sort.Ints(nums)
    n := len(nums)
    result := [][]int{}
    // i < n-2是保证 i 的右边还有两个数可以凑成3个数
    for i := 0; i < n-2; i++ {
       // 注意去重。i > 0 是为了防止越界，因为后面有 i-1
       if i > 0 && nums[i] == nums[i-1] {
          continue
       }
       l, r := i+1, n-1
       for l < r {
          sum := nums[i] + nums[l] + nums[r]
          if sum == 0 {
             result = append(result, []int{nums[i], nums[l], nums[r]})
             // 注意去重，缩小范围
             // 当 nums[l] == nums[l+1] 时，说明 nums[l] 和下一个元素 nums[l+1] 是相同的。
             for l < r && nums[l] == nums[l+1] {
                l++
             }
             for l < r && nums[r] == nums[r-1] {
                r--
             }
             l++
             r--
          } else if sum < 0 {
             l++
          } else {
             r--
          }
       }
    }
    return result
}
```

#### 53. 最大子数组和

给你一个整数数组 `nums` ，请你找出一个具有最大和的连续子数组（子数组最少包含一个元素），返回其最大和。

**子数组**是数组中的一个连续部分。

```go
// 空间优化，两个 max
func maxSubArray(nums []int) int {
    n := len(nums)
    // 以当前下标 i 为结尾的子数组的最大和
    cur := nums[0]
    maxSum := nums[0]
    for i := 1; i < n; i++ {
       cur = max(nums[i], cur+nums[i])
       maxSum = max(cur, maxSum)
    }
    return maxSum
}
```

#### 912. 排序数组

给你一个整数数组 `nums`，请你将该数组升序排列。

你必须在 **不使用任何内置函数** 的情况下解决问题，时间复杂度为 `O(nlog(n))`，并且空间复杂度尽可能小。

```go
// 归并排序
func sortArray(nums []int) []int {
    if len(nums) <= 1 {
       return nums
    }
    mid := len(nums) / 2
    // 开始递归
    left := sortArray(nums[:mid])
    right := sortArray(nums[mid:])
    return merge(left, right)
}

func merge(left, right []int) []int {
    result := []int{}
    i, j := 0, 0
    for i < len(left) && j < len(right) {
       if left[i] <= right[j] {
          result = append(result, left[i])
          i++
       } else {
          result = append(result, right[j])
          j++
       }
       // i++ ,j++写到这里是错的
    }
    // 处理剩下的元素,左右一定有一边是空的，注意这两行代码要放在循环之外
    result = append(result, left[i:]...)
    result = append(result, right[j:]...)
    return result
}
```

#### 21. 合并两个有序链表

将两个升序链表合并为一个新的 **升序** 链表并返回。新链表是通过拼接给定的两个链表的所有节点组成的。

```go
// 哨兵节点不动，cur节点移动
func mergeTwoLists(list1 *ListNode, list2 *ListNode) *ListNode {
    // 初始化哨兵节点,不需要增加 cur1，cur2
    dummy := &ListNode{}
    // cur 和 dummy 指向同一个节点,但是cur会往前移动，dummy不会
    cur := dummy
    for list1 != nil && list2 != nil {
       if list1.Val < list2.Val {
          cur.Next = list1
          // 注意链表要往前滑动
          list1 = list1.Next
       } else {
          cur.Next = list2
          list2 = list2.Next
       }
       // cur 也要往前滑动
       cur = cur.Next
    }
    // 还要拼接剩下的链表
    if list1 != nil {
       cur.Next = list1
    }
    if list2 != nil {
       cur.Next = list2
    }
    return dummy.Next
}
```

#### 5. 最长回文子串

给你一个字符串 `s`，找到 `s` 中最长的 回文 子串。

```go
// 中心扩展法，分奇数、偶数回文长度处理
func longestPalindrome(s string) string {
    if len(s) == 0 {
       return s
    }

    start, end := 0, 0
    left, right := 0, 0
    for i := 0; i < len(s); i++ {
       // 奇数长度的回文处理
       left, right = expand(s, i, i)
       if right-left > end-start {
          start, end = left, right
       }
       // 偶数长度的回文处理
       left, right = expand(s, i, i+1)
       if right-left > end-start {
          start, end = left, right
       }
    }
    return s[start : end+1]
}

func expand(s string, left, right int) (int, int) {
    for left >= 0 && right < len(s) && s[left] == s[right] {
       left--
       right++
    }
    // 返回正确的坐标
    return left + 1, right - 1
}
```

#### 102. 二叉树的层序遍历

给你二叉树的根节点 `root` ，返回其节点值的 **层序遍历** 。 （即逐层地，从左到右访问所有节点）。

```go
// DFS + 递归
// 深度优先搜索（DFS）的典型特征：一条路走到黑，再回头走别的路。
func levelOrder(root *TreeNode) [][]int {
	// result用于存储结果，二维切片的每个一维索引代表一层
	res := [][]int{}
	
	var dfs func(*TreeNode, int)
	dfs = func(node *TreeNode, level int) {
		if node == nil {
			return
		}
		// 如果当前层不存在，先新建一层
		// len(res)表示已经存放了几层
		// 比如目前已经存放了根节点（0层），此时len(res) == 1，但是第1层还不存在，因此要创建[]int{}
		if len(res) == level {
			res = append(res, []int{})
		}
		
		res[level] = append(res[level], node.Val)
		// 递归处理左右子树
		dfs(node.Left, level+1)
		dfs(node.Right, level+1)
	}

	dfs(root, 0)
	return res
}
```

#### 1. 两数之和

给定一个整数数组 `nums` 和一个整数目标值 `target`，请你在该数组中找出 **和为目标值** *`target`* 的那 **两个** 整数，并返回它们的数组下标。

你可以假设每种输入只会对应一个答案，并且你不能使用两次相同的元素。

你可以按任意顺序返回答案。

```go
// 使用map
func twoSum(nums []int, target int) []int {
    hashmap := map[int]int{}
    for i := 0; i < len(nums); i++ {
       if j, ok := hashmap[target-nums[i]]; ok {
          return []int{i, j}
       }
       // 不能写成 hashmap[target-nums[i]] = i
       hashmap[nums[i]] = i
    }
    return nil
}
```

#### 33. 搜索旋转排序数组

整数数组 `nums` 按升序排列，数组中的值 **互不相同** 。

在传递给函数之前，`nums` 在预先未知的某个下标 `k`（`0 <= k < nums.length`）上进行了 **向左旋转**，使数组变为 `[nums[k], nums[k+1], ..., nums[n-1], nums[0], nums[1], ..., nums[k-1]]`（下标 **从 0 开始** 计数）。例如， `[0,1,2,4,5,6,7]` 下标 `3` 上向左旋转后可能变为 `[4,5,6,7,0,1,2]` 。

给你 **旋转后** 的数组 `nums` 和一个整数 `target` ，如果 `nums` 中存在这个目标值 `target` ，则返回它的下标，否则返回 `-1` 。

你必须设计一个时间复杂度为 `O(log n)` 的算法解决此问题。

**注意：**`target` 就是一个**目标数字**，题目给你这个数字，让你判断它是否在旋转后的数组里。

```go
// for循环 + 二分查找.
// 注意：target是目标值，而返回的是下标
func search(nums []int, target int) int {
    left, right := 0, len(nums)-1
    // 开启循环，以左右边界为条件
    for left <= right {
       mid := (left + right) / 2
       if nums[mid] == target {
          return mid
       }

       // 一定注意，一定有一边严格递增，左或者右。
       // 这里不能用 <,一定要用 <=
       if nums[left] <= nums[mid] { // 只能说明左半边严格递增
          // 再判断 target 是否在左半边
          if target >= nums[left] && target < nums[mid] { // target 在左半边
             // 缩小左边范围
             right = mid - 1
          } else { // target 在右半边，跳到右半边找
             // 注意不能写成left++
             left = mid + 1
          }
       } else { // 只能说明右边严格递增了
          if target > nums[mid] && target <= nums[right] { // target 在右半边
             left = mid + 1
          } else { // target 在左半边，跳到左半边找
             right = mid - 1
          }
       }
    }
    return -1
}
```

#### 200. 岛屿数量

给你一个由 `'1'`（陆地）和 `'0'`（水）组成的的二维网格，请你计算网格中岛屿的数量。

岛屿总是被水包围，并且每座岛屿只能由水平方向和/或竖直方向上相邻的陆地连接形成。

此外，你可以假设该网格的四条边均被水包围。

```go
// 深度优先（DFS）+ 递归
// 不断淹没已经访问过的岛屿
func numIslands(grid [][]byte) int {
    // 特殊情况
    if len(grid) == 0 {
       return 0
    }
    count := 0
    // m,n 分别用于表示岛的宽、长
    m, n := len(grid), len(grid[0])
    
    // 每次统计到一个岛（值为1），就递归淹没整个岛屿
    // 注意不能直接写成 dfs := func(i, j int) {}，原因是 Go 的匿名函数在声明自身时不能直接递归
    var dfs func(i, j int)
    dfs = func(i, j int) {
       // 控制边界
       if i < 0 || i >= m || j < 0 || j >= n || grid[i][j] == '0' {
          // 碰到边界直接结束
          return
       }
       // 淹没当前位置，标记为'0'
       grid[i][j] = '0'
       // 当前地点的上下左右也要淹没
       dfs(i, j+1)
       dfs(i, j-1)
       dfs(i+1, j)
       dfs(i-1, j)
    }

    for i := 0; i < m; i++ {
       for j := 0; j < n; j++ {
          // 发现了岛屿
          if grid[i][j] == '1' {
             count++
             dfs(i, j)
          }
       }
    }
    return count
}
```

#### 46. 全排列

给定一个不含重复数字的数组 `nums` ，返回其 *所有可能的全排列* 。你可以 **按任意顺序** 返回答案。

```go
// 回溯法，需要使用一个递归函数
func permute(nums []int) [][]int {
    // 存储最终结果
    res := [][]int{}
    // 回溯递归函数
    var dfs func(int)
    dfs = func(start int) {
       // 递归终止条件：当 start 走到数组末尾，代表得到了一组结果
       if len(nums) == start {
          buf := make([]int, len(nums))
          copy(buf, nums)
          res = append(res, buf)
       }

       // 枚举：把 [start..len-1] 这些“尚未固定”的元素，依次和 start 的位置交换
       for i := start; i < len(nums); i++ {
          // 交换，注意一开始的 start和 i是相等的，相当于交换自身
          nums[start], nums[i] = nums[i], nums[start]
          dfs(start + 1)
          nums[start], nums[i] = nums[i], nums[start] // 回溯
       }
    }
    dfs(0)
    return res
}
```

**回溯法的核心思想：**

回溯法是一种 **搜索算法**，用于 **在问题的解空间中寻找所有可行解**。

它的核心思路是：

1. **选择**：从当前状态选择一个可能的选项。
2. **递归**：进入下一个状态继续选择。
3. **撤销选择（回溯）**：当发现当前选择不行，或者完成一个解后，要撤销之前的选择，尝试其他选项。

简而言之，就是 **试错 + 撤回 + 继续尝试**。

#### 88. 合并两个有序数组

给你两个按 **非递减顺序** 排列的整数数组 `nums1` 和 `nums2`，另有两个整数 `m` 和 `n` ，分别表示 `nums1` 和 `nums2` 中的元素数目。

请你 **合并** `nums2` 到 `nums1` 中，使合并后的数组同样按 **非递减顺序** 排列。

**注意：**最终，合并后数组不应由函数返回，而是存储在数组 `nums1` 中。为了应对这种情况，`nums1` 的初始长度为 `m + n`，其中前 `m` 个元素表示应合并的元素，后 `n` 个元素为 `0` ，应忽略。`nums2` 的长度为 `n` 。

```go
// 逆向双指针
func merge(nums1 []int, m int, nums2 []int, n int) {
    p1, p2 := m-1, n-1
    p := m + n - 1
    for p1 >= 0 && p2 >= 0 {
       if nums1[p1] > nums2[p2] {
          nums1[p] = nums1[p1]
          p1--
       } else {
          nums1[p] = nums2[p2]
          p2--
       }
       p--
    }

    // 注意 nums1 比 nums2 长
    // 因此最后只有可能nums2的剩余元素需要被处理，
    // 如果nums1剩余，不需要被处理
    // 如果 nums2 还有剩余，直接拷贝过去
    for p2 >= 0 {
       nums1[p] = nums2[p2]
       p2--
       p--
    }
}
```

#### 20. 有效的括号

给定一个只包括 `'('`，`')'`，`'{'`，`'}'`，`'['`，`']'` 的字符串 `s` ，判断字符串是否有效。

有效字符串需满足：

1. 左括号必须用相同类型的右括号闭合。
2. 左括号必须以正确的顺序闭合。
3. 每个右括号都有一个对应的相同类型的左括号。

```go
// 栈
func isValid(s string) bool {
    // 定义一个栈来存储左括号。左括号遇到右括号时会弹出，因此栈里面不会存储右括号
    stack := []rune{}
    // 每个右括号对应一个左括号
    pairs := map[rune]rune{
       ')': '(',
       '}': '{',
       ']': '[',
    }
    // 遍历字符串
    for _, ch := range s {
       // 先判断是左括号还是右括号
       // if left, ok := pairs[ch]; ok {
       if ch == ')' || ch == '}' || ch == ']' { // 遇到右括号
          // 已经有了右括号,栈为空和栈顶不匹配都表明栈里面没有和 val 匹配的左括号
          if len(stack) == 0 || stack[len(stack)-1] != pairs[ch] {
             return false // 直接结束
          }
          // 如果有匹配到的左括号，弹出
          stack = stack[:len(stack)-1]
       } else { // 遇到左括号
          // 左括号直接入栈
          // 不能写成 stack[len(stack) - 1] = val,会报错
          stack = append(stack, ch)
       }
    }
    return len(stack) == 0
}
```

#### 121. 买卖股票的最佳时机

给定一个数组 `prices` ，它的第 `i` 个元素 `prices[i]` 表示一支给定股票第 `i` 天的价格。

你只能选择 **某一天** 买入这只股票，并选择在 **未来的某一个不同的日子** 卖出该股票。设计一个算法来计算你所能获取的最大利润。

返回你可以从这笔交易中获取的最大利润。如果你不能获取任何利润，返回 `0` 。

```go
// 贪心 / 动态规划（O(n)
func maxProfit(prices []int) int {
    minPrice := int(^uint(0) >> 1)
    maxProfit := 0
    for _, price := range prices {
       // 找出最小的price
       minPrice = min(price, minPrice)
       // 找出最大的收益
       maxProfit = max(maxProfit, price-minPrice)
    }
    return maxProfit
}
```

#### 103. 二叉树的锯齿形层序遍历

给你二叉树的根节点 `root` ，返回其节点值的 **锯齿形层序遍历** 。（即先从左往右，再从右往左进行下一层遍历，以此类推，层与层之间交替进行）。

```go
// DFS + 就地头插/尾插
func zigzagLevelOrder(node *TreeNode) [][]int {
    // 存储结果
    res := [][]int{}
    // 创建dfs递归函数
    var dfs func(*TreeNode, int)
    // level 代表当前是第几层
    dfs = func(node *TreeNode, level int) {
       // 递归结束条件
       if node == nil {
          return
       }

       // 首次到达该层，需要新创建一个切片[]int{}
       if len(res) == level {
          res = append(res, []int{})
       }
       if level%2 == 0 { // 偶数层(0,2,4...) 左->右: 尾插
          res[level] = append(res[level], node.Val)
       } else { // 奇数层(1,3,5...) 右->左: 头插
          res[level] = append([]int{node.Val}, res[level]...)
       }

       // 递归处理左右子树
       dfs(node.Left, level+1)
       dfs(node.Right, level+1)
    }
    dfs(node, 0)
    return res
}
```

#### 236. 二叉树的最近公共祖先

给定一个二叉树, 找到该树中两个指定节点的最近公共祖先。

[百度百科](https://baike.baidu.com/item/最近公共祖先/8918834?fr=aladdin)中最近公共祖先的定义为：“对于有根树 T 的两个节点 p、q，最近公共祖先表示为一个节点 x，满足 x 是 p、q 的祖先且 x 的深度尽可能大（**一个节点也可以是它自己的祖先**）。”

```go
// 递归
func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
    // 都没找到或者找到了p或者q，就返回
    if root == nil || root == p || root == q {
       // 直接返回root
       return root
    }
    // 在左右子树中找
    left := lowestCommonAncestor(root.Left, p, q)
    right := lowestCommonAncestor(root.Right, p, q)

    // 注意是&&，左右子树都找到了
    if left != nil && right != nil {
       return root
    }
    // 只在左子树找到
    if left != nil {
       return left
    }
    // 只在右子树找到
    return right
}
```

#### 92. 反转链表 II

给你单链表的头指针 `head` 和两个整数 `left` 和 `right` ，其中 `left <= right` 。请你反转从位置 `left` 到位置 `right` 的链表节点，返回 **反转后的链表** 。

```go
// 一次遍历「穿针引线」反转链表（头插法）
func reverseBetween(head *ListNode, left int, right int) *ListNode {
    // 设置哨兵节点
    dummy := &ListNode{Next: head}
    pre := dummy
    // 获取left的前一个pre节点
    for i := 0; i < left-1; i++ {
       pre = pre.Next
    }

    // 开始翻转，头插法
    cur := pre.Next
    for i := 0; i < right-left; i++ {
       next := cur.Next     // 1.选取要头插的节点
       cur.Next = next.Next // 2.从原位置删掉 next
       next.Next = pre.Next // 3.将 next 指向翻转部分的头（永远是 pre.next）,千万不能写成 cur
       pre.Next = next      // 4.pre 指向 next，完成插入。千万不能写成 cur = next
    }
    return dummy.Next
}
```

#### 141. 环形链表

给你一个链表的头节点 `head` ，判断链表中是否有环。

如果链表中有某个节点，可以通过连续跟踪 `next` 指针再次到达，则链表中存在环。 为了表示给定链表中的环，评测系统内部使用整数 `pos` 来表示链表尾连接到链表中的位置（索引从 0 开始）。**注意：`pos` 不作为参数进行传递** 。仅仅是为了标识链表的实际情况。

*如果链表中存在环* ，则返回 `true` 。 否则，返回 `false` 。

```go
// 哈希表
func hasCycle(head *ListNode) bool {
    visited := map[*ListNode]bool{}
    cur := head
    for cur != nil {
       if visited[cur] {
          return true
       }
       visited[cur] = true
       cur = cur.Next
    }
    return false
}
```

**tips：**

```go
if _,ok := visited[cur];ok {
    
}
// 可以简化成 
if visited[cur] {
    
}
```

#### 300. 最长递增子序列

给你一个整数数组 `nums` ，找到其中最长严格递增子序列的长度。

**子序列** 是由数组派生而来的序列，删除（或不删除）数组中的元素而不改变其余元素的顺序。例如，`[3,6,2,7]` 是数组 `[0,3,1,6,2,2,7]` 的子序列。

```go
// 动态规划：单切片
func lengthOfLIS(nums []int) int {
    n := len(nums)
    // 特殊情况
    if n == 0 {
       return 0
    }
    // dp[i] 存储以 nums[i] 结尾的字符串的严格递增子序列长度
    dp := make([]int, n)
    for i := 0; i < n; i++ {
       // 长度至少是1
       dp[i] = 1
    }
    // 记录结果，长度至少是1
    res := 1
    // 开始处理字符串
    for i := 1; i < n; i++ {
       // 处理从0到i的位置的字符串
       for j := 0; j < i; j++ {
          if nums[j] < nums[i] {
             // 取较长的一个
             dp[i] = max(dp[i], dp[j]+1)
          }
       }
       res = max(res, dp[i])
    }
    return res
}
```

**二分查找**

```go
// 二分查找
func lengthOfLIS1(nums []int) int {
    tails := []int{}
    for _, num := range nums {
       // 找到第一个 >= num 的位置，当tails为空时，i返回0
       i := sort.SearchInts(tails, num)
       // 表示 num 比 tails 里所有数都大
       if i == len(tails) {
          tails = append(tails, num)
       } else {
          tails[i] = num // 替换，保证结尾更小
       }
    }
    return len(tails)
}
```

#### 54. 螺旋矩阵

给你一个 `m` 行 `n` 列的矩阵 `matrix` ，请按照 **顺时针螺旋顺序** ，返回矩阵中的所有元素。

```go
// 始终沿着当前方向前进。如果下一个格子越界或已访问，就右转
func spiralOrder(matrix [][]int) []int {
    // 特殊情况
    if len(matrix) == 0 {
       return []int{}
    }
    // 分别代表行数、列数
    m, n := len(matrix), len(matrix[0])
    // 结果集
    res := make([]int, 0, m*n)
    // 定义 4 个方向：右、下、左、上
    dirs := [][]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
    // 从左上角 (0,0) 开始，方向（0=右，1=下，2=左，3=上）
    x, y, d := 0, 0, 0
    v := math.MinInt // 用于填充每个走过的地方
    // 开始处理所有节点
    for i := 0; i < m*n; i++ {
       res = append(res, matrix[x][y]) // 把当前节点值加到结果中
       matrix[x][y] = v                // 标记已经走过
       xn, yn := x+dirs[d][0], y+dirs[d][1]
       // 判断是否越界或遍历过,判断xn、yn而不是x
       if xn < 0 || xn >= m || yn < 0 || yn >= n || matrix[xn][yn] == v {
          // 转向
          d = (d + 1) % 4
          // 重新赋值
          xn, yn = x+dirs[d][0], y+dirs[d][1]
       }
       // 下一个节点
       x, y = xn, yn
    }
    return res
}
```

#### 143. 重排链表

给定一个单链表 `L` 的头节点 `head` ，单链表 `L` 表示为：

```
L0 → L1 → … → Ln - 1 → Ln
```

请将其重新排列后变为：

```
L0 → Ln → L1 → Ln - 1 → L2 → Ln - 2 → …
```

不能只是单纯的改变节点内部的值，而是需要实际的进行节点交换。

```go
// 数组辅助
func reorderList(head *ListNode) {
    // 特殊情况
    if head == nil {
       return
    }

    nodes := []*ListNode{}
    // 存储链表的每个节点
    for cur := head; cur != nil; cur = cur.Next {
       nodes = append(nodes, cur)
    }

    // 双指针遍历切片
    left, right := 0, len(nodes)-1
    for left < right {
       nodes[left].Next = nodes[right]
       left++
       // 避免重复连接导致成环
       if left == right {
          break
       }
       nodes[right].Next = nodes[left]
       right--
    }
    // 把链表尾节点的 Next 指针置为 nil
    nodes[right].Next = nil
}
```
