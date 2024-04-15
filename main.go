package main

import (
	"fmt"
	hducashelper "github.com/Ec3o/WordKiller/hducashelper"
	"github.com/Ec3o/WordKiller/nethandlers"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"time"
)

func main() {
	var username, password string
	var week int
	var mode string
	var delay int
	fmt.Print("[*]请输入 CAS 账号:")
	fmt.Scanln(&username)
	fmt.Print("[*]请输入 CAS 密码:")
	bytePassword, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatalln("密码读取失败:", err)
	}
	password = string(bytePassword)
	println("\n[*]正在登录...")
	time.Sleep(2 * time.Second)
	ticker := hducashelper.CasPasswordLogin(username, password) // 杭电 CAS 账号密码
	sklLogin := hducashelper.SklLogin(ticker)
	if sklLogin.Error() != nil {
		log.Fatalln(sklLogin.Error())
	}
	token := sklLogin.GetToken()
	fmt.Printf("[*]登录成功, token: %s\n", token)
	fmt.Println("[*]请输入测试周数:")
	fmt.Scanf("%d\n", &week)
	fmt.Println("[*]请选择模式:自测[0]/考试[1]")
	fmt.Scanln(&mode)
	fmt.Println("[*]请输入做题时间(单位：s,推荐时长:450s):")
	fmt.Scanln(&delay)
	switch mode {
	case "0":
		fmt.Println("[*]开始自测...")
	case "1":
		fmt.Println("[*]开始考试...")
	default:
		log.Fatalln("[-]模式错误")
	}
	nethandlers.Exam(token, week, mode, delay)
	//test.Test()
	fmt.Println("[*]运行结束.按任意键退出...")
	fmt.Scanln()
}
