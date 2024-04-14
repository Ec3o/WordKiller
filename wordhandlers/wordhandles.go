package wordhandlers

import (
	"WordKiller/models"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// CleanString 去除字符串中的空格和特殊字符
func CleanString(input string) string {
	// 小写转换，移除空格、标点符号等
	input = strings.ToLower(input)
	return strings.NewReplacer(" ", "", ".", "", "，", "", ",", "", ";", "", ":", "").Replace(input)
}
func GetAnswer(paper models.Paper) models.Answer {
	matched := 0
	fmt.Println("\n[*]开始查找答案...")

	file, err := os.ReadFile("cet-4.json")
	if err != nil {
		log.Fatalf("无法读取题库文件: %v", err)
	}

	var bank []models.AnswerBank
	if err := json.Unmarshal(file, &bank); err != nil {
		log.Fatalf("解析题库数据失败: %v", err)
	}

	ans := models.Answer{
		PaperId: paper.PaperId,
		Type:    paper.Type,
		List:    make([]models.AnswerDetail, 0),
	}

	for index, question := range paper.List {
		detail := models.AnswerDetail{
			Input:         "A", // 默认选项
			PaperDetailId: question.PaperDetailId,
		}
		title := CleanString(question.Title)
		fmt.Println("\n[*]正在处理题目:第", index+1, "题", title)
		// 收集所有可能的匹配项
		var possibleMatches []string
		match := false
		for _, obj := range bank {
			if strings.Contains(obj.Mean, title) || strings.Contains(title, obj.Word) {
				match = true
				possibleMatches = append(possibleMatches, obj.Word, obj.Mean)
				fmt.Println("[*]匹配到可能项:", obj.Word, obj.Mean)
			}
		}
		if !match {
			fmt.Println("[x]未匹配到任何项")
		}
		maxMatchCount := 0
		answers := map[string]string{
			"A": question.AnswerA,
			"B": question.AnswerB,
			"C": question.AnswerC,
			"D": question.AnswerD,
		}
		for k, v := range answers {
			cleanV := CleanString(v)
			fmt.Println("[*]正在处理选项:", k, cleanV)
			matchCount := 0
			for _, match := range possibleMatches {
				if strings.Contains(match, cleanV) || strings.Contains(cleanV, match) {
					matchCount++
				}
			}
			if matchCount > maxMatchCount {
				matched++
				maxMatchCount = matchCount
				detail.Input = k[len(k)-1:] // 取答案的最后一个字符作为选项标记（A/B/C/D）
				fmt.Println("[+]最佳匹配项:", k, cleanV, "匹配次数:", matchCount)
			}
		}

		ans.List = append(ans.List, detail)
	}

	fmt.Println("[+]查找结束！共匹配到", matched, "个答案")
	return ans
}
