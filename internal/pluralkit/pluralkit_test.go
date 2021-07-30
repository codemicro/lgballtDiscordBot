package pluralkit

import (
	"fmt"
	"os"
	"runtime"
	"runtime/pprof"
	"testing"
	"time"
)

func Test_leaks(t *testing.T) {
	a, _ := os.OpenFile("first.txt", os.O_WRONLY|os.O_CREATE, 0644)
	b, _ := os.OpenFile("second.txt", os.O_WRONLY|os.O_CREATE, 0644)

	pprof.Lookup("goroutine").WriteTo(a, 1)
	fmt.Println(runtime.NumGoroutine())
	MessageById("870605457536978954")
	pprof.Lookup("goroutine").WriteTo(b, 1)

	a.Close()
	b.Close()

	fmt.Println(runtime.NumGoroutine())
	MessageById("870605457536978954")
	fmt.Println(runtime.NumGoroutine())
	MessageById("870605457536978954")
	fmt.Println(runtime.NumGoroutine())
	MessageById("870605457536978954")
	fmt.Println(runtime.NumGoroutine())
	MessageById("870605457536978954")
	fmt.Println(runtime.NumGoroutine())
	MessageById("870605457536978954")
	fmt.Println(runtime.NumGoroutine())
	MessageById("870605457536978954")
	fmt.Println(runtime.NumGoroutine())
	time.Sleep(time.Second * 30)
	fmt.Println(runtime.NumGoroutine())
}
