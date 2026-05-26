// Package goroutine contains test cases for the goroutine analyzer.
package goroutine

import "fmt"

func init() {
	go func() { // want `goroutine launched inside init function`
		fmt.Println("running in goroutine")
	}()
}

func init() {
	// no goroutine here — should be fine
	fmt.Println("safe init")
}

func safeFunc() {
	// launching a goroutine outside init is fine
	go func() {
		fmt.Println("safe goroutine")
	}()
}

func anotherInit() {
	// not named init, so it's fine
	go func() {
		fmt.Println("goroutine in non-init func")
	}()
}
