package collection

import (
	"container/list"
	"errors"
	"math"
	"sync"
)

var (
	ErrEmpty = errors.New("没有数据")
)

// Stack 栈定义
type Stack interface {
	Push(data any) error
	Pop() (any, error)
	Top() (any, error)
	Len() int
	Cap() int
}

// ArrayStack 栈的数组实现
type ArrayStack struct {
	mux sync.Mutex
	arr []any
	p   int
}

func NewArrayStack(cap int) Stack {
	return &ArrayStack{arr: make([]any, cap), p: -1, mux: sync.Mutex{}}
}

func (a *ArrayStack) Push(data any) error {
	a.mux.Lock()
	defer a.mux.Unlock()
	a.p++

	if cap(a.arr)-1 < a.p {
		a.arr = append(a.arr, data)
	} else {
		a.arr[a.p] = data
	}
	return nil
}

func (a *ArrayStack) Pop() (any, error) {
	a.mux.Lock()
	defer a.mux.Unlock()
	if a.p < 0 {
		return nil, ErrEmpty
	}
	d := a.arr[a.p]
	a.p--
	return d, nil
}

func (a *ArrayStack) Top() (any, error) {
	if a.p < 0 {
		return nil, ErrEmpty
	}
	return a.arr[a.p], nil
}

func (a *ArrayStack) Len() int {
	return a.p + 1
}

func (a *ArrayStack) Cap() int {
	return cap(a.arr)
}

// LinkedStack 栈的链表实现
type LinkedStack struct {
	mux  sync.Mutex
	list *list.List
}

func NewLinkedStack() Stack {
	return &LinkedStack{mux: sync.Mutex{}, list: list.New()}
}

func (l *LinkedStack) Push(data any) error {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.list.PushBack(data)
	return nil
}

func (l *LinkedStack) Pop() (any, error) {
	l.mux.Lock()
	defer l.mux.Unlock()
	return l.list.Remove(l.list.Back()), nil
}

func (l *LinkedStack) Top() (any, error) {
	return l.list.Back().Value, nil
}

func (l *LinkedStack) Len() int {
	return l.list.Len()
}

func (l *LinkedStack) Cap() int {
	return math.MaxInt
}
