package test

import (
	"WordKiller/models"
	"WordKiller/ui"
	"WordKiller/wordhandlers"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

func Test() {
	delay := 2
	fmt.Println("[*]测试中...")

	body, err := ioutil.ReadFile("papers/paper_20240414152156.json")
	if err != nil {
		log.Fatalf("读取文件错误: %v", err)
	}

	var paper models.Paper
	if err := json.Unmarshal(body, &paper); err != nil {
		log.Fatalf("解析 JSON 错误: %v", err)
	}

	fmt.Println("[*]等待提交中...")
	go ui.DisplayProgress(time.Duration(delay) * time.Second)
	time.Sleep(time.Duration(delay) * time.Second)

	answer := wordhandlers.GetAnswer(paper)
	fmt.Printf("答案详情: %+v\n", answer)
	fmt.Println()
	fmt.Println("[+]测试完成！")
}
