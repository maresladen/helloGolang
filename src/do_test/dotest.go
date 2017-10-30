package do_test

import (
	"fmt"
)

func fibonacci(c, quit chan int) {
	x, y := 0, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}

//Dotest1 这是一个测试方法
func Dotest1() {
	c := make(chan int, 5)
	quit := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
		}
		quit <- 0
	}()
	fibonacci(c, quit)
}

//NewTestFun aaa
func NewTestFun() {
	ch := make(chan int, 100)
	for b := 0; b < 100; b++ {
		ch <- b
	}
	for {
		c := <-ch
		fmt.Println(c)
		if len(ch) == 0 {
			return
		}
	}
}

//DoTestFun 这是一个测试方法
func DoTestFun() {

	// var regArr []int

	// abc(&regArr)
	// fmt.Println(len(regArr))

	bb := 666
	aa := &bb

	bbc(aa)
	fmt.Println(bb)
}

func abc(regArr *[]int) {

	for a := 0; a < 10; a++ {
		*regArr = append((*regArr), a)
	}
}

func bbc(aa *int) {
	*aa = 123
}
