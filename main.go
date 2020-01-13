package main

import (
	"time"
)

func main() {
	go WinMain()

	for {
		for i:=uintptr(50);i<500;i=i+50 {
			drawRectangle(i, i+25, 200, 300, 3, 0xffff00)
		}
		time.Sleep(1 * time.Millisecond)
	}
}
