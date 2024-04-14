package nethandlers

import (
	"WordKiller/models"
	"WordKiller/ui"
	"WordKiller/wordhandlers"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func GetHeaders(token string) http.Header {
	ticket := generateTicket(21) // 自定义函数，模拟 JavaScript 中的 ticket 函数
	headers := http.Header{}
	headers.Set("Skl-Ticket", ticket)
	headers.Set("X-Auth-Token", token)
	headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	headers.Set("Accept", "application/json, text/plain, */*")
	headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
	headers.Set("Connection", "keep-alive")
	headers.Set("Referer", "https://skl.hdu.edu.cn/")
	return headers
}
func PostHeader(token string) http.Header {
	headers := http.Header{}
	headers.Set("Host", "skl.hdu.edu.cn")
	headers.Set("Sec-Ch-Ua", `"Chromium";v="117", "Not;A=Brand";v="8"`)
	headers.Set("Skl-Ticket", generateTicket(21))
	headers.Set("Sec-Ch-Ua-Mobile", "?0")
	headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.5938.132 Safari/537.36")
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept", "application/json, text/plain, */*")
	headers.Set("X-Auth-Token", token)
	headers.Set("Sec-Ch-Ua-Platform", `"Windows"`)
	headers.Set("Origin", "https://skl.hduhelp.com")
	headers.Set("Sec-Fetch-Site", "cross-site")
	headers.Set("Sec-Fetch-Mode", "cors")
	headers.Set("Sec-Fetch-Dest", "empty")
	headers.Set("Referer", "https://skl.hduhelp.com/")
	headers.Set("Accept-Encoding", "gzip, deflate, br")
	headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
	return headers
}
func generateTicket(length int) string {
	const NL = "useandom-26T198340PX75pxJACKVERYMINDBUSHWOLF_GQZbfghjklqvwyzrict"
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	ticket := make([]byte, length)
	for i, b := range bytes {
		ticket[i] = NL[b&63] // 与 `& 63` 相同，确保索引在 0-63 范围内
	}
	return string(ticket)
}
func Exam(token string, week int, mode string, delay int) {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error retrieving working directory:", err)
		return
	}
	client := &http.Client{}
	startTime := time.Now().UnixMilli()
	url := fmt.Sprintf("https://skl.hdu.edu.cn/api/paper/new?type=%s&week=%d&startTime=%d", mode, week, startTime)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header = GetHeaders(token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	//将body保存到json文件中
	filename := fmt.Sprintf("paper_%s.json", time.Now().Format("20060102150405"))
	fullPath := filepath.Join(wd, "papers", filename)
	err = ioutil.WriteFile(fullPath, body, 0644)
	fmt.Println("[*]存储试题信息中...")
	if err != nil {
		fmt.Println("[x]文件存储错误:", err)
		return
	}

	var paper models.Paper
	if err := json.Unmarshal(body, &paper); err != nil {
		log.Fatal(err)
	}

	fmt.Println("[*]等待提交中...")
	go ui.DisplayProgress(time.Duration(delay) * time.Second)
	time.Sleep(time.Duration(delay) * time.Second)

	ans := wordhandlers.GetAnswer(paper)
	data, err := json.Marshal(ans)
	fmt.Println(string(data))
	if err != nil {
		log.Fatal(err)
	}

	saveURL := "https://skl.hdu.edu.cn/api/paper/save"
	req, err = http.NewRequest("POST", saveURL, bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}
	req.Header = PostHeader(token)
	resp, err = client.Do(req)
	body, err = ioutil.ReadAll(resp.Body)
	//fmt.Println("[*]Response status:", resp.Status)
	//fmt.Println("[*]Response body:", string(body))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	fmt.Println("[+]提交成功！")
	time.Sleep(2 * time.Second)
	fmt.Println("[+]系统处理中...")
	time.Sleep(2 * time.Second)
	fmt.Println("[+]正在查看成绩...")
	detailURL := fmt.Sprintf("https://skl.hdu.edu.cn/api/paper/detail?paperId=%s", paper.PaperId)
	req, err = http.NewRequest("GET", detailURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header = GetHeaders(token)

	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var score models.ExamScore
	if err := json.Unmarshal(body, &score); err != nil {
		log.Fatal(err)
	}
	//fmt.Println("[*]Response status:", resp.Status)
	//fmt.Println("[*]Response body:", string(body))
	fmt.Printf("[+]本次成绩: %d\n", score.Mark)

}
