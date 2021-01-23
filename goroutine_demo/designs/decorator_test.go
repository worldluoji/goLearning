package designs

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func Hello(s string) {
	fmt.Println("Hello Golang " + s)
}

func decorator(f func(s string)) func(s string) {
	return func(s string) {
		fmt.Println("before...")
		f(s)
		fmt.Println("after...")
	}
}

func TestDecorator(t *testing.T) {
	d := decorator(Hello)
	d("zhazhahui")
}

type SumFunc func(int64, int64) int64

func getFunctionName(i interface{}) string {
	// 反射机制来获取函数名；
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func timedSumFunc(f SumFunc) SumFunc {
	return func(start, end int64) int64 {
		defer func(t time.Time) {
			fmt.Printf("--- Time Elapsed (%s): %v ---\n", getFunctionName(f), time.Since(t))
		}(time.Now())
		return f(start, end)
	}
}

func Sum1(start, end int64) int64 {
	var sum int64
	sum = 0
	if start > end {
		start, end = end, start
	}
	for i := start; i <= end; i++ {
		sum = sum + i
	}
	return sum
}

func Sum2(start, end int64) int64 {
	if start > end {
		start, end = end, start
	}
	return (start + end) * (end - start + 1) / 2
}

func TestDecorator2(t *testing.T) {
	f1 := timedSumFunc(Sum1)
	f2 := timedSumFunc(Sum2)
	fmt.Printf("%d, %d\n", f1(-10000, 10000000), f2(-10000, 10000000))
}
