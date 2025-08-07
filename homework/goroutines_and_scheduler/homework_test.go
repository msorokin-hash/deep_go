package main

import (
	"testing"

	"container/heap"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	priorityQueue *PriorityQueue
	taskRegistry  map[int]*QueueItem
}

func NewScheduler() Scheduler {
	return Scheduler{
		priorityQueue: &PriorityQueue{},
		taskRegistry:  make(map[int]*QueueItem),
	}
}

func (s *Scheduler) AddTask(task Task) {
	if _, exists := s.taskRegistry[task.Identifier]; exists {
		return
	}

	newItem := &QueueItem{
		task:  &task,
		index: -1,
	}
	s.taskRegistry[task.Identifier] = newItem
	heap.Push(s.priorityQueue, newItem)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	if item, exists := s.taskRegistry[taskID]; exists {
		item.task.Priority = newPriority
		heap.Fix(s.priorityQueue, item.index)
	}
}

func (s *Scheduler) GetTask() Task {
	if s.priorityQueue.Len() == 0 {
		return Task{}
	}

	item := heap.Pop(s.priorityQueue).(*QueueItem)
	delete(s.taskRegistry, item.task.Identifier)
	return *item.task
}

type QueueItem struct {
	task  *Task
	index int
}

type PriorityQueue []*QueueItem

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].task.Priority > pq[j].task.Priority
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*QueueItem)
	item.index = n
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1
	*pq = old[0 : n-1]
	return item
}

func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}
	task6 := Task{Identifier: 1, Priority: 100}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)

	task = scheduler.GetTask()
	assert.Equal(t, task6, task)

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}
