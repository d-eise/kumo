package loopclosure

import "fmt"

func badRangeLoop() {
	slice := []int{1, 2, 3}
	for i, v := range slice {
		go func() { // want `goroutine captures loop variable "i"`
			fmt.Println(i)
		}()
		_ = v
	}
}

func badForLoop() {
	for i := 0; i < 10; i++ {
		go func() { // want `goroutine captures loop variable "i"`
			fmt.Println(i)
		}()
	}
}

func safeRangeLoop() {
	slice := []int{1, 2, 3}
	for _, v := range slice {
		v := v
		go func() {
			fmt.Println(v)
		}()
	}
}

func safeForLoop() {
	for i := 0; i < 10; i++ {
		i := i
		go func() {
			fmt.Println(i)
		}()
	}
}

func safePassByArg() {
	slice := []int{1, 2, 3}
	for i, v := range slice {
		go func(idx, val int) {
			fmt.Println(idx, val)
		}(i, v)
	}
}
