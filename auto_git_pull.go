package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"hook_test2/utils"
)

const (
	// 配置文件
	HOOK_TEST1_CONF_FILE string = "auto_pull_conf"
)

func main() {
	host := "localhost"
	if runtime.GOOS == "linux" {
		host = "linuxHost"
	}
	http.HandleFunc("/hook_test1", GitHubHookTest1)
	err := http.ListenAndServe(host+":5572", nil)
	if utils.IsErr(err) {
		return
	}
}

func GitHubHookTest1(w http.ResponseWriter, r *http.Request) {
	var err error
	utils.Info("========= 收到一次请求 ==========")
	defer func() {
		// 回到原来目录
		os.Chdir("../hook_test2")
		utils.Info("====== this is hook_test2 =========")
		utils.Info("============ end =============\n")
		if err != nil {
			fmt.Fprintf(w, "err >> "+err.Error())
		} else {
			fmt.Fprintf(w, "success")
		}
	}()
	// 读取文件的信息
	bytes, err := ioutil.ReadFile(HOOK_TEST1_CONF_FILE)
	if err != nil {
		utils.Info("read conf file err >> ", err)
		return
	}
	// 按照换行符分割
	text := string(bytes)
	cmdarr := strings.Split(text, "\r\n")
	// 是否调到新文件
	isBegin := true
	for _, val := range cmdarr {
		tmpval := strings.TrimSpace(val)
		//新的命令开始 --> 切换目录
		if tmpval != "" && isBegin {
			os.Chdir(tmpval)
			utils.Info("====== this is ", tmpval, " =========")
		} else if tmpval != "" {
			utils.Info("--------- a new cmd --------")
			// //分割命令
			// cmds := strings.Split(tmpval, " ")
			// //命令
			// command := cmds[0]
			// //参数
			// params := cmds[1:]
			//执行命令
			// err = execCommand(command, params, tmpval)
			err = execCommand(tmpval)
			if err != nil {
				return
			}
		}
		if tmpval == "" {
			isBegin = true
			//空行跳回原来目录
			os.Chdir("../hook_test2")
			utils.Info("====== this is hook_test2 =========")
		} else {
			isBegin = false
		}
	}
	return
}

//执行命令
// func execCommand(command string, params []string, allCommand string) (err error) {
// cmd := exec.Command(command, params...)
func execCommand(command string) (err error) {
	//合成命令
	cmd := exec.Command("/bin/sh", "-c", command)
	//显示运行的命令
	utils.Info(cmd.Args)
	//返回一个连接命令标准输出的管道
	stdout, err := cmd.StdoutPipe()
	if utils.IsErr(err) {
		return err
	}
	//执行命令,但不等待命令完成
	err = cmd.Start()
	if utils.IsErr(err) {
		return err
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
		return err
	}
	return nil
}
