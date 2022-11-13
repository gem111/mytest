package queue

import (
	"context"
)

type Queue[T any] interface {
	// Enqueue 定义方法
	//入队和出队两个方法
	EnQueue(ctx context.Context, data T) error
	DeQueue(ctx context.Context) (T, error)

	IsFull() bool
	IsEmpty() bool
	Len() uint64
}

func Enqueue() {
}
