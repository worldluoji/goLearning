package goroutine

import (
	"fmt"
	"testing"
	"time"
)

func hello(msg string) {
	fmt.Println("Hello " + msg)
}

func TestGoroutineCase1(t *testing.T) {
	//在新的协程中执行hello方法
	go hello("World")
	//等待100毫秒让协程执行结束, 否则可能退出后，goroutine还没有结束
	time.Sleep(100 * time.Millisecond)
}

// 在容量为C的channel上的第k个接收先行发生于从这个channel上的第k+C次发送完成。
func TestGoroutineCase2(t *testing.T) {
	start := time.Now()
	ch := make(chan int, 2)
	for i := 0; i < 6; i++ {
		go func(num int) {
			ch <- num
			t.Logf("发送第%d个元素\n", num)
		}(i)
	}
	// 该例子可以看到，接收的顺序，不一定就是发送的顺序。
	for i := 0; i < 6; i++ {
		go func(i int) {
			o := <-ch
			t.Logf("收到的第%d个元素为%d\n", i, o)
		}(i)
	}
	elapse := time.Since(start)
	t.Log("耗时时间为：", elapse)
	time.Sleep(100 * time.Millisecond)
}

var cht = make(chan int)
var a string

func setVal() {
	a = "hello golang"
	cht <- 9
}

func TestGoroutineCase3(t *testing.T) {
	go setVal()
	<-cht
	// 无缓冲channel的接收先行发生于发送完成，因此能正确打印出hello golang
	fmt.Println(a)
}

func setVal2() {
	a = "hello golang2"
	close(cht)
}

// 对channel的关闭先行发生于接收到值，因为channel已经被关闭了
func TestGoroutineCase4(t *testing.T) {
	go setVal2()
	<-cht
	// channel接收一定在发送之前，因此能正确打印出hello golang
	fmt.Println(a)
}

func TestSelectChannelCase(t *testing.T) {
	ch := make(chan string)
	go func(chan string) {
		time.Sleep(101 * time.Millisecond)
		ch <- "data...data...data"
	}(ch)
	select {
	case msg := <-ch:
		t.Log("从ch读取到数据：", msg)
	case <-time.After(100 * time.Millisecond):
		t.Log("已超时")
		// default:
		// 	t.Log("什么也没找到")
	}
}

func echo(nums []int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

func sq(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			out <- n * n
		}
		close(out)
	}()
	return out
}

func odd(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			if n%2 == 1 {
				out <- n
			}
		}
		close(out)
	}()
	return out
}

// pipeline实际就是用channel将结果连接起来
func TestPipeline1(t *testing.T) {
	nums := []int{1, 2, 3, 4, 5, 6}
	// in fact , it is a pipeline
	for r := range odd(sq(echo(nums))) {
		fmt.Println(r)
	}
}

type EchoFunc func([]int) <-chan int
type PipelineFunc func(<-chan int) <-chan int

// 利用函数式编程对Pipeline进行封装
func Pipeline(nums []int, echoFunc EchoFunc, pipelineFuncs ...PipelineFunc) <-chan int {
	ch := echoFunc(nums)
	for _, f := range pipelineFuncs {
		ch = f(ch)
	}
	return ch
}

func TestPipeline2(t *testing.T) {
	out := Pipeline([]int{1, 2, 3, 4, 5, 6}, echo, sq, odd)
	for r := range out {
		fmt.Println(r)
	}
}
