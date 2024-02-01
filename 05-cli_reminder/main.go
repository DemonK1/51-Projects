package main

import (
	"fmt"
	"github.com/gen2brain/beeep"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/common"
	"github.com/olebedev/when/rules/en"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	// 环境变量
	markName  = "GOLANG_CLI_REMINDER"
	markValue = "1"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <hh:mm> <text message\n>", os.Args[0])
		os.Exit(1)
	}

	now := time.Now()
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(common.All...)

	t, err := w.Parse(os.Args[1], now)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	if t == nil {
		fmt.Println("无法解析时间")
		os.Exit(2)
	}
	if now.After(t.Time) {
		fmt.Println("请输入一个未来时间")
		os.Exit(3)
	}

	diff := t.Time.Sub(now)
	// 在当前进程设置环境变量
	//err = os.Setenv("GOLANG_CLI_REMINDER", "1")
	//if err != nil {
	//	fmt.Println("无法设置环境变量：", err)
	//	return
	//}
	if os.Getenv(markName) == markValue {
		time.Sleep(diff)
		err = beeep.Alert("Reminder", strings.Join(os.Args[2:], " "), "assets/information.png")
		if err != nil {
			fmt.Println(err)
			os.Exit(4)
		}
	} else {
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Env = append(os.Environ(), fmt.Sprintf("%s=%s", markName, markValue))
		if err = cmd.Start(); err != nil {
			fmt.Println(err)
			os.Exit(5)
		}
		fmt.Println("之后会收到提醒 ", diff.Round(time.Second))
		os.Exit(0)
	}
}
