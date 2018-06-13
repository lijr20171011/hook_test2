package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"hook_test2/utils"
)

const (
	// 配置文件
	CONF_FILE string = "auth_pull_conf"
)

func main() {
	http.HandleFunc("/pull", GitHubHookTest)
	err := http.ListenAndServe("localhost:5572", nil)
	if utils.IsErr(err) {
		return
	}
}

func GitHubHookTest(w http.ResponseWriter, r *http.Request) {
	utils.Info("========= 收到一次请求 ==========")
	defer func() {
		fmt.Fprintf(w, "success")
	}()
	// 读取文件的信息
	bytes, err := ioutil.ReadFile(CONF_FILE)
	if err != nil {
		utils.Info("read conf file err >> ", err)
		return
	}
	// 按照换行符分割
	text := string(bytes)
	cmdarr := strings.Split(text, "\r\n")
	// 是否新的开始
	isBegin := true
	for _, val := range cmdarr {
		tmpval := strings.TrimSpace(val)
		//新的命令开始 --> 切换目录
		if tmpval != "" && isBegin {
			os.Chdir(tmpval)
			utils.Info("====== this is "+tmpval, " =========")
		} else if tmpval != "" {
			utils.Info("==============")
			//分割命令
			cmds := strings.Split(tmpval, " ")
			//命令
			command := cmds[0]
			//参数
			params := cmds[1:]
			//执行命令
			isSuccess := execCommand(command, params)
			if !isSuccess {
				return
			}
		}
		if tmpval == "" {
			isBegin = true
		} else {
			isBegin = false
		}
	}
	os.Chdir("../hook_test2")
	utils.Info("====== this is hook_test2 =========")
	utils.Info("============ end =============\n")
	return
}

//执行命令
func execCommand(command string, params []string) bool {
	//合成命令
	cmd := exec.Command(command, params...)
	//显示运行的命令
	utils.Info(cmd.Args)
	//返回一个连接命令标准输出的管道
	stdout, err := cmd.StdoutPipe()
	if utils.IsErr(err) {
		return false
	}
	//执行命令,但不等待命令完成
	err = cmd.Start()
	if utils.IsErr(err) {
		return false
	}
	reader := bufio.NewReader(stdout)
	//实时读取
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}
		if utils.IsErr(err) {
			break
		}
		fmt.Print(line)
	}
	fmt.Println()
	//在start之后使用,等待标准输入/标准输出/标准错误复制完成,释放与cmd相关的资源
	err = cmd.Wait()
	if utils.IsErr(err) {
		return false
	}
	return true
}
