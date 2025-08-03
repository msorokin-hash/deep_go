package main

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func Defragment(memory []byte, pointers []unsafe.Pointer) {
	// нечего дефрагментировать
	if len(memory) == 0 || len(pointers) == 0 {
		return
	}

	// считаем количество занятых ячеек памяти
	occupied := make(map[int]bool)
	for _, ptr := range pointers {
		ptrAddr := uintptr(ptr)
		baseAddr := uintptr(unsafe.Pointer(&memory[0]))
		index := int(ptrAddr - baseAddr)
		if index >= 0 && index < len(memory) {
			occupied[index] = true
		}
	}

	newPos := 0
	posMap := make(map[int]int)

	// копируем занятые ячейки в начало массива
	for i := 0; i < len(memory); i++ {
		if occupied[i] {
			if i != newPos {
				memory[newPos] = memory[i]
				posMap[i] = newPos
			} else {
				posMap[i] = i
			}
			newPos++
		}
	}

	// зануляем оставшиеся ячейки
	for i := newPos; i < len(memory); i++ {
		memory[i] = 0
	}

	// пересчитываем указатели
	for i, ptr := range pointers {
		ptrAddr := uintptr(ptr)
		baseAddr := uintptr(unsafe.Pointer(&memory[0]))
		oldIndex := int(ptrAddr - baseAddr)
		if newIndex, exists := posMap[oldIndex]; exists {
			pointers[i] = unsafe.Pointer(&memory[newIndex])
		}
	}
}

func TestDefragmentation(t *testing.T) {
	var fragmentedMemory = []byte{
		0xFF, 0x00, 0x00, 0x00,
		0x00, 0xFF, 0x00, 0x00,
		0x00, 0x00, 0xFF, 0x00,
		0x00, 0x00, 0x00, 0xFF,
	}

	var fragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[5]),
		unsafe.Pointer(&fragmentedMemory[10]),
		unsafe.Pointer(&fragmentedMemory[15]),
	}

	var defragmentedPointers = []unsafe.Pointer{
		unsafe.Pointer(&fragmentedMemory[0]),
		unsafe.Pointer(&fragmentedMemory[1]),
		unsafe.Pointer(&fragmentedMemory[2]),
		unsafe.Pointer(&fragmentedMemory[3]),
	}

	var defragmentedMemory = []byte{
		0xFF, 0xFF, 0xFF, 0xFF,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
	}

	Defragment(fragmentedMemory, fragmentedPointers)
	assert.True(t, reflect.DeepEqual(defragmentedMemory, fragmentedMemory))
	assert.True(t, reflect.DeepEqual(defragmentedPointers, fragmentedPointers))
}
