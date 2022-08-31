package main

import (
	"fmt"
	"sync"
)

type MyMap[KEY int | string, VALUE float32 | float64] map[KEY]VALUE

type MyStruct[T int | string] struct {
	Name string
	Data T
}

type IPrintData[T int | float32 | string] interface {
	Print(data T)
}

type MyChan[T int | string] chan T

type WowStruct[T int | float32, S []T] struct {
	Data     S
	MaxValue T
	MinValue T
}

type NewType[T interface{ *int }] []T

type Slice1[T int | string | float32 | float64] []T

type FloatSlice[T float32 | float64] Slice1[T]

type IntAndStringSlice[T int | string] Slice1[T]

type IntSlice[T int] IntAndStringSlice[T]

type WowMap[T int | string] map[string]Slice1[T]

type WowMap2[T Slice1[int] | Slice1[string]] map[string]T

type Info[T int | string, T2 Slice1[string] | int] map[T]T2

func main() {
	var a Slice[int] = []int{1, 2, 3}
	fmt.Printf("Type Name: %T\n", a) //输出：Type Name: Slice[int]

	var m MyMap[string, float64] = map[string]float64{
		"bacon": 6.66,
		"pb":    8.88,
	}
	fmt.Printf("Type MyMap: %T\n", m)

	var ws = WowStruct[int, []int]{Data: []int{3, 4}, MaxValue: 4, MinValue: 3}
	fmt.Printf("Type ws: %T\n", ws)

	var info = make(Info[int, int])
	info[1] = 100
	info[2] = 200
	fmt.Println("info = ", info)

	var info1 = make(Info[int, Slice1[string]])
	info1[1] = []string{"123", "456"}

	var s Slice[int] = []int{1, 2, 3, 4}
	fmt.Println(s.sum()) // 输出：10

	var q1 Queue[int]
	q1.Put(1)
	q1.Put(2)
	fmt.Println(q1)

	//var q2 Queue[int]
	//go func() {
	//	for {
	//		q2.Put(rand.Int())
	//		time.Sleep(500 * time.Millisecond)
	//	}
	//}()
	//
	//go func() {
	//	for {
	//		if ret, b := q2.Pop(); b {
	//			fmt.Println(ret)
	//		} else {
	//			fmt.Println("queue is empty")
	//			return
	//		}
	//		time.Sleep(500 * time.Millisecond)
	//	}
	//}()
	//
	//select {}

	sumInt := Add(1, 2)
	fmt.Println("sumInt: ", sumInt)

	sumFloat64 := Add[float64](1, 3.1)
	fmt.Println("sumFloat64: ", sumFloat64)

	subInt := Sub[int](1, 2)
	fmt.Println("subInt: ", subInt)
}

type Slice[T int | float32 | float64] []T

func (s Slice[T]) sum() T {
	var sumVal T
	for _, v := range s {
		sumVal += v
	}
	return sumVal
}

type Queue[T any] struct {
	elements []T
}

var mux = sync.Mutex{}

// Put put data to the end of Queue
func (q *Queue[T]) Put(data T) {
	mux.Lock()
	q.elements = append(q.elements, data)
	mux.Unlock()
}

// Pop pop data from the head of Queue
func (q *Queue[T]) Pop() (T, bool) {
	mux.Lock()
	defer mux.Unlock()
	var val T
	if len(q.elements) == 0 {
		return val, true
	}
	val = q.elements[0]
	q.elements = q.elements[1:]
	return val, len(q.elements) == 0
}

// Size get the size of Queue
func (q *Queue[T]) Size() int {
	return len(q.elements)
}

func Add[T int | float64 | float32](a T, b T) T {
	return a + b
}

func Sub[T int | float64](a T, b T) T {
	fn := func(i, j T) T {
		return i - j
	}
	return fn(a, b)
}
