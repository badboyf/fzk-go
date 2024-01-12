package main

import (
	"fmt"
)

type ListNode struct {
	Val  int
	Next *ListNode
}

func initListNode(param ...int) *ListNode {
	var head *ListNode
	var rs *ListNode
	for _, v := range param {
		if rs == nil {
			rs = &ListNode{Val: v}
			head = rs
		} else {
			rs.Next = &ListNode{Val: v}
			rs = rs.Next
		}
	}
	return head
}
func printListNode(p *ListNode)  {
	h := p
	for h != nil {
		fmt.Printf("%v  ", h.Val)
		h = h.Next
	}
	fmt.Printf("\n")
}

func max(v1 int, v2 int) int {
	if v1 > v2 {
		return v1
	}
	return v2
}
func min(v1, v2 int) int  {
	if v1 > v2 {
		return v2
	}
	return v1
}

/* 1. 两数之和
给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出 和为目标值 target  的那 两个 整数，并返回它们的数组下标。
你可以假设每种输入只会对应一个答案。但是，数组中同一个元素在答案里不能重复出现。
你可以按任意顺序返回答案。*/
func twoSum(nums []int, target int) []int {
	if nums == nil || len(nums) < 2 {
		return nil
	}
	m := map[int]int{}
	for i, x := range nums {
		if p, ok := m[target-x]; ok {
			return []int{p, i}
		}
		m[x] = i
	}
	return nil
}

/**
2. 两数相加
给你两个 非空 的链表，表示两个非负的整数。它们每位数字都是按照 逆序 的方式存储的，并且每个节点只能存储 一位 数字。请你将两个数相加，并以相同形式返回一个表示和的链表。
你可以假设除了数字 0 之外，这两个数都不会以 0 开头。
示例 1：
输入：l1 = [2,4,3], l2 = [5,6,4]   输出：[7,0,8]       解释：342 + 465 = 807.
*/
func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {
	var tail *ListNode
	var head *ListNode
	carry, v1, v2 := 0, 0, 0
	for l1 != nil || l2 != nil {
		if l1 != nil {
			v1 = l1.Val
			l1 = l1.Next
		}
		if l2 != nil {
			v2 = l2.Val
			l2 = l2.Next
		}
		sum := v1 + v2 + carry
		sum, carry = sum % 10, sum / 10
		if tail == nil {
			tail = &ListNode{Val: sum}
			head = tail
		} else {
			tail.Next = &ListNode{Val: sum}
			tail = tail.Next
		}
	}
	if carry > 0 && tail != nil {
		tail.Next = &ListNode{Val: carry}
	}
	printListNode(head)
	return head
}
/*
3. 无重复字符的最长子串
给定一个字符串 s ，请你找出其中不含有重复字符的 最长子串 的长度。
示例 1:
输入: s = "abcabcbb" 输出: 3
解释: 因为无重复字符的最长子串是 "abc"，所以其长度为 3。
*/
func lengthOfLongestSubstring(s string) int {
	rk, rs := 0, 0
	m := map[byte]int{}
	for i := range s {
		if i != 0 {
			delete(m, s[i - 1])
		}
		for rk < len(s) && m[s[rk]] == 0 {
			m[s[rk]]++
			rk++
		}
		rs = max(rs, rk - i)
	}
	return rs
}

/*  4. 寻找两个正序数组的中位数
给定两个大小分别为 m 和 n 的正序（从小到大）数组 nums1 和 nums2。请你找出并返回这两个正序数组的 中位数 。
示例 1：
输入：nums1 = [1,3], nums2 = [2]
输出：2.00000
解释：合并数组 = [1,2,3] ，中位数 2
*/
//func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
//
//}

/*5. 最长回文子串
	给你一个字符串 s，找到 s 中最长的回文子串。
	示例 1：
	输入：s = "babad"
	输出："bab"
	解释："aba" 同样是符合题意的答案。
	https://leetcode-cn.com/problems/longest-palindromic-substring/
*/
func longestPalindrome(s string) string {
	if len(s) == 0 {
		return s
	}
	start := 0
	end := 0
	for i := range s {
		left1, right1 := longestPalindromeLength(s, i, i)
		left2, right2 := longestPalindromeLength(s, i, i + 1)
		if (right2 - left2) - (right1 - left1) > 0 {
			if (right2 - left2) > (end - start) {
				start = left2
				end = right2
			}
		} else {
			if (right1 - left1) > (end - start) {
				start = left1
				end = right1
			}
		}
	}
	return s[start:end+1]
}
func longestPalindromeLength(s string, start int, end int) (int, int) {
	for start >= 0 && end < len(s) && s[start] == s[end] {
		start--
		end++
	}
	start++
	end--
	return start, end
}

/*10. 正则表达式匹配
给你一个字符串 s 和一个字符规律 p，请你来实现一个支持 '.' 和 '*' 的正则表达式匹配。
'.' 匹配任意单个字符
'*' 匹配零个或多个前面的那一个元素
所谓匹配，是要涵盖 整个 字符串 s的，而不是部分字符串。
示例 1：
输入：s = "aa" p = "a"
输出：false
解释："a" 无法匹配 "aa" 整个字符串。
示例 5：
输入：s = "mississippi" p = "mis*is*p*."
输出：false
*/
func isMatch(s string, p string) bool {
	return false
}

/*11. 盛最多水的容器    https://leetcode-cn.com/problems/container-with-most-water/
给你 n 个非负整数 a1，a2，...，an，每个数代表坐标中的一个点 (i, ai) 。在坐标内画 n 条垂直线，垂直线 i 的两个端点分别为 (i, ai) 和 (i, 0) 。找出其中的两条线，使得它们与 x 轴共同构成的容器可以容纳最多的水。
输入：[1,8,6,2,5,4,8,3,7]
输出：49
解释：图中垂直线代表输入数组 [1,8,6,2,5,4,8,3,7]。在此情况下，容器能够容纳水（表示为蓝色部分）的最大值为 49。
*/
func maxArea(height []int) int {
	left := 0
	right := len(height) - 1
	m := 0
	for left < right {
		m = max(m, (right - left) * min(height[left], height[right]))
		if height[left] > height[right] {
			right--
		} else {
			left++
		}
	}
	return m
}

/*17. 电话号码的字母组合
给定一个仅包含数字 2-9 的字符串，返回所有它能表示的字母组合。答案可以按 任意顺序 返回。
给出数字到字母的映射如下（与电话按键相同）。注意 1 不对应任何字母。
示例 1：
输入：digits = "23"
输出：["ad","ae","af","bd","be","bf","cd","ce","cf"]
*/
func letterCombinations(digits string) []string {
	letterMap := map[string]string{
		"2": "abc",
		"3": "def",
		"4": "ghi",
		"5": "jkl",
		"6": "mno",
		"7": "pqrs",
		"8": "tuv",
		"9": "wxyz",
	}
	rs := make([]string, 0, 1)
	backtrack(&digits, 0, "", &letterMap, &rs)
	return rs
}
func backtrack(digits *string, index int, item string, letterMap *map[string]string, rs *[]string)  {
	if index == len((*digits)) {
		*rs = append(*rs, item)
	} else {
		letters := (*letterMap)[string((*digits)[index])]
		for _, v := range letters {
			item = item + string(v)
			backtrack(digits, index+1, item, letterMap, rs)
			item = item[0:len(item)-1]
		}
	}
}

/*
19. 删除链表的倒数第 N 个结点
给你一个链表，删除链表的倒数第 n 个结点，并且返回链表的头结点。
输入：head = [1,2,3,4,5], n = 2
输出：[1,2,3,5]
*/
func removeNthFromEnd(head *ListNode, n int) *ListNode {
	rs, left, right := head, head, head
	i := 0
	for right.Next != nil {
		right = right.Next
		if i >= n {
			left = left.Next
		}
		i++
	}
	if i < n {

	} else     if rs == left {
		rs = rs.Next
	} else {
		left.Next = left.Next.Next
	}

	return rs
}

/* https://leetcode-cn.com/problems/merge-two-sorted-lists/
21. 合并两个有序链表
将两个升序链表合并为一个新的 升序 链表并返回。新链表是通过拼接给定的两个链表的所有节点组成的。
输入：l1 = [1,2,4], l2 = [1,3,4]
输出：[1,1,2,3,4,4]
*/
func mergeTwoLists(l1 *ListNode, l2 *ListNode) *ListNode {
	rs := &ListNode{}
	head := rs

	for l1 != nil && l2 != nil {
		if l1.Val > l2.Val {
			head.Next = l2
			l2 = l2.Next
		} else {
			head.Next = l1
			l1 = l1.Next
		}
		head = head.Next
	}
	if l2 == nil {
		head.Next = l1
	} else {
		head.Next = l2
	}

	return rs.Next
}

/*https://leetcode-cn.com/problems/generate-parentheses/
22. 括号生成
数字 n 代表生成括号的对数，请你设计一个函数，用于能够生成所有可能的并且 有效的 括号组合。
有效括号组合需满足：左括号必须以正确的顺序闭合。

示例 1：
输入：n = 3
输出：["((()))","(()())","(())()","()(())","()()()"]
*/
func generateParenthesis(n int) []string {
	if n == 0 {
		return []string{}
	}
	rs := &[]string{}
	generateParenthesisBackTrack("", n, 0, 0, rs)
	return *rs
}

func generateParenthesisBackTrack(itm string, cnt int, open int, close int, rs *[]string)  {
	if len(itm) == cnt * 2 {
		*rs = append(*rs, itm)
		return
	}
	if open < cnt {
		tmpItm := itm
		tmpItm += "("
		generateParenthesisBackTrack(tmpItm, cnt, open + 1, close, rs)
	}
	if close < open {
		tmpItm := itm
		tmpItm += ")"
		generateParenthesisBackTrack(tmpItm, cnt, open, close + 1, rs)
	}
}
/*  https://leetcode-cn.com/problems/merge-k-sorted-lists/
23. 合并K个升序链表
给你一个链表数组，每个链表都已经按升序排列。
请你将所有链表合并到一个升序链表中，返回合并后的链表。
示例 1：

输入：lists = [[1,4,5],[1,3,4],[2,6]]
输出：[1,1,2,3,4,4,5,6]
解释：链表数组如下：
[
  1->4->5,
  1->3->4,
  2->6
]
将它们合并到一个有序链表中得到。
1->1->2->3->4->4->5->6
*/
func mergeKLists(lists []*ListNode) *ListNode {

	return mergeKListsMerge(lists, 0, len(lists) - 1)
}
func mergeKListsMerge(lists []*ListNode, left int, right int) *ListNode {
	if left == right {
		return lists[left]
	} else if left > right {
		return nil
	}
	mid := (left + right) / 2
	return mergeTwoLists(mergeKListsMerge(lists, left, mid), mergeKListsMerge(lists, mid + 1, right))
}
func mergeKListsMergeTwo(l1, l2 *ListNode) *ListNode {
	rs := &ListNode{}
	head := rs
	for l1 != nil && l2 != nil {
		if l1.Val > l2.Val {
			head.Next = l2
			l2 = l2.Next
		} else {
			head.Next = l1
			l1 = l1.Next
		}
		head = head.Next
	}
	if l1 != nil {
		head.Next = l1
	}
	if l2 != nil {
		head.Next = l2
	}
	return rs.Next
}

/* https://leetcode-cn.com/problems/next-permutation/
31. 下一个排列
实现获取 下一个排列 的函数，算法需要将给定数字序列重新排列成字典序中下一个更大的排列（即，组合出下一个更大的整数）。
如果不存在下一个更大的排列，则将数字重新排列成最小的排列（即升序排列）。
必须 原地 修改，只允许使用额外常数空间。
示例 1：
输入：nums = [1,2,3]
输出：[1,3,2]
*/
func nextPermutation(nums []int) {
	left, right := 0, 0
	v := -1
	for i := len(nums) - 1; i > 0; i-- {
		if nums[i-1] < nums[i] {
			left = i - 1
			v = nums[i - 1]
			break
		}
	}

	for i := len(nums) - 1; i > left; i-- {
		if nums[i] > v {
			right = i
			break
		}
	}
	if left > -1 {
		swap(nums, left, right)
	}
	idx := 1
	for i := len(nums) - 1; i > (left + len(nums)) / 2 - 1; i-- {
		swap(nums, left + idx, len(nums) - idx)
		idx++
	}
	fmt.Printf("%v\n", nums)
}
func swap(nums []int, left int, right int)  {
	tmp := nums[left]
	nums[left] = nums[right]
	nums[right] = tmp
}

/*39. 组合总和
给定一个无重复元素的正整数数组 candidates 和一个正整数 target ，找出 candidates 中所有可以使数字和为目标数 target 的唯一组合。
candidates 中的数字可以无限制重复被选取。如果至少一个所选数字数量不同，则两种组合是唯一的。
对于给定的输入，保证和为 target 的唯一组合数少于 150 个。

示例 1：
输入: candidates = [2,3,6,7], target = 7
输出: [[7],[2,2,3]]
*/
func combinationSum(candidates []int, target int) (ans [][]int) {
	comb := []int{}

	var dfs func(target, idx int)
	dfs = func(target int, idx int) {
		fmt.Printf("%v %v %v \n", target, idx, comb)
		if idx >= len(candidates) {
			return
		}
		if target == 0 {
			ans = append(ans, append([]int(nil), comb...))
			return
		}
		dfs(target, idx + 1)
		if target - candidates[idx] >= 0 {
			comb = append(comb, candidates[idx])
			dfs(target - candidates[idx], idx)
			comb = comb[:len(comb) - 1]
		}
	}

	dfs(target, 0)
	return
}

/* https://leetcode-cn.com/problems/permutations/
46. 全排列
给定一个不含重复数字的数组 nums ，返回其 所有可能的全排列 。你可以 按任意顺序 返回答案。

示例 1：
输入：nums = [1,2,3]
输出：[[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]
*/
func permute(nums []int) (result [][]int) {
	var dfs func(index int)
	dfs = func(index int) {
		if index == len(nums) {
			result = append(result, append([]int(nil), nums...))
			return
		}

		for i := index; i < len(nums); i++ {
			swap(nums, index, i)
			dfs(index + 1)
			swap(nums, index, i)
		}
		fmt.Printf("%v \n", nums)
	}

	dfs(0)
	return
}



func main() {
	//fmt.Printf("%v", twoSum([]int{2, 7, 11, 15}, 9))
	//fmt.Printf("%v", addTwoNumbers(initListNode(2, 4, 3), initListNode(5, 6, 4)))
	//fmt.Printf("%v\n", lengthOfLongestSubstring("abcdbcbb"))
	//fmt.Printf("%s\n", longestPalindrome(""))
	//fmt.Printf("%v\n", maxArea([]int{1,8,6,2,5,4,8,3,7}))
	//fmt.Printf("%v\n", letterCombinations("ab"))
	//printListNode(removeNthFromEnd(initListNode(1, 2), 3))
	//printListNode(mergeTwoLists(initListNode(1, 3, 5), initListNode(2, 4, 6)))
	//fmt.Printf("%v\n", generateParenthesis(3))
	//printListNode(mergeKLists([]*ListNode{
	//	initListNode(1, 4, 5),
	//	initListNode(1, 3, 4),
	//	initListNode(2, 6)}))
	//nextPermutation([]int{4,5,2,6,3,1}) // 3,2,1
	//fmt.Printf("%v \n", combinationSum([]int{2, 3, 5}, 8))
	fmt.Printf("%v \n", permute([]int{1, 2, 3, 4}))
}
