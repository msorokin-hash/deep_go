package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v homework_test.go

type TreeNode struct {
	Key   int
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

type OrderedMap struct {
	root *TreeNode
	size int
}

func NewOrderedMap() OrderedMap {
	return OrderedMap{}
}

func (m *OrderedMap) Insert(key, value int) {
	m.root = m.insertNode(m.root, key, value)
	m.size++
}

func (m *OrderedMap) insertNode(root *TreeNode, key int, value int) *TreeNode {
	if root == nil {
		return &TreeNode{Key: key, Val: value}
	}

	if key < root.Key {
		root.Left = m.insertNode(root.Left, key, value)
	} else if key > root.Key {
		root.Right = m.insertNode(root.Right, key, value)
	} else {
		root.Val = value
		m.size--
	}

	return root
}

func (m *OrderedMap) Erase(key int) {
	if m.Contains(key) {
		m.root = m.deleteNode(m.root, key)
		m.size--
	}
}

func (m *OrderedMap) deleteNode(root *TreeNode, key int) *TreeNode {
	if root == nil {
		return nil
	}

	if key < root.Key {
		root.Left = m.deleteNode(root.Left, key)
	} else if key > root.Key {
		root.Right = m.deleteNode(root.Right, key)
	} else {
		if root.Left == nil {
			return root.Right
		}
		if root.Right == nil {
			return root.Left
		}

		minNode := m.findMin(root.Right)
		root.Key = minNode.Key
		root.Val = minNode.Val
		root.Right = m.deleteNode(root.Right, minNode.Key)
	}

	return root
}

func (m *OrderedMap) findMin(node *TreeNode) *TreeNode {
	current := node
	for current.Left != nil {
		current = current.Left
	}
	return current
}

func (m *OrderedMap) Contains(key int) bool {
	return m.findNode(m.root, key) != nil
}

func (m *OrderedMap) findNode(node *TreeNode, key int) *TreeNode {
	if node == nil || node.Key == key {
		return node
	}
	if key < node.Key {
		return m.findNode(node.Left, key)
	}
	return m.findNode(node.Right, key)
}

func (m *OrderedMap) Size() int {
	return m.size
}

func (m *OrderedMap) ForEach(action func(int, int)) {
	m.inOrder(m.root, action)
}

func (m *OrderedMap) inOrder(node *TreeNode, action func(int, int)) {
	if node != nil {
		m.inOrder(node.Left, action)
		action(node.Key, node.Val)
		m.inOrder(node.Right, action)
	}
}

func TestCircularQueue(t *testing.T) {
	data := NewOrderedMap()
	assert.Zero(t, data.Size())

	data.Insert(10, 10)
	data.Insert(5, 5)
	data.Insert(15, 15)
	data.Insert(2, 2)
	data.Insert(4, 4)
	data.Insert(12, 12)
	data.Insert(14, 14)

	assert.Equal(t, 7, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(3))
	assert.False(t, data.Contains(13))

	var keys []int
	expectedKeys := []int{2, 4, 5, 10, 12, 14, 15}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))

	data.Erase(15)
	data.Erase(14)
	data.Erase(2)

	assert.Equal(t, 4, data.Size())
	assert.True(t, data.Contains(4))
	assert.True(t, data.Contains(12))
	assert.False(t, data.Contains(2))
	assert.False(t, data.Contains(14))

	keys = nil
	expectedKeys = []int{4, 5, 10, 12}
	data.ForEach(func(key, _ int) {
		keys = append(keys, key)
	})

	assert.True(t, reflect.DeepEqual(expectedKeys, keys))
}
