package main

import (
	"fmt"
	"sync"

	"github.com/xxpauloxx/go-simple-debugger/debugger"
)

func main() {
	x := 10
	y := "valor y"

	fmt.Println("Iniciando o exemplo...")
	debugger.Breakpoint(map[string]interface{}{
		"x": &x,
		"y": &y,
	})

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		z := 20
		debugger.GoroutineBreakpoint(map[string]interface{}{
			"z": &z,
		}, &wg)
		fmt.Printf("Value z: %d\n", z)
	}()

	wg.Wait()

	fmt.Printf("Após o debugger: x = %d, y = %s\n", x, y)
}
