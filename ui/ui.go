package ui

import (
	"fmt"
	"time"
)

// displayProgress显示进度条
func DisplayProgress(duration time.Duration) {
	progressBarLength := 50 // 进度条的长度
	interval := duration / time.Duration(progressBarLength)
	fmt.Print("[*]答题进度:[")
	for i := 0; i < progressBarLength-1; i++ {
		time.Sleep(interval)
		fmt.Print("=")
	}
	fmt.Println("=]")
}
