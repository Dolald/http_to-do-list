package main

import "fmt"

func main() {
	fmt.Println(hasCycle)
}

type ListNode struct {
	Val  int
	Next *ListNode
}

func hasCycle(head *ListNode) bool {
	visited_nodes := make(map[*ListNode]bool)
	current_node := head
	for current_node != nil {
		if visited_nodes[current_node] {
			return true
		}
		visited_nodes[current_node] = true
		current_node = current_node.Next
	}
	return false
}
