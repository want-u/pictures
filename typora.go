package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	timestemp   = ""                                // 时间标签
	client      = &http.Client{}                    // 创建一个客户端
	apiUrl      = "https://gitee.com/api/v5/repos/" // Request URL
	contentType = "application/json;charset=UTF-8"  // 定义网络文件的类型和网页的编码
)

type GiteeResponse struct {
	Message string                 `json:"message"` // 请求失败的消息
	Content map[string]interface{} `json:"content"` // 请求成功的内容
	Commit  map[string]interface{} `json:"commit"`  // 请求成功的commit
}

type GiteeRequest struct {
	Access_token string `json:"access_token"` // 用户授权码
	Owner        string `json:"owner"`        // 仓库所属空间地址(企业、组织或个人的地址path)
	Repo         string `json:"repo"`         // 仓库路径(path)
	Branch       string `json:"branch"`       // 分支名称。默认为仓库对默认分支
	Path         string `json:"path"`         // 文件的路径(目录+文件名，这里我默认目录为空，文件直接放在仓库根目录)
	Content      string `json:"content"`      // 文件内容, 要用 base64 编码
	Message      string `json:"message"`      // 提交信息
}

func (g *GiteeRequest) putPics(picSlice []string) {
	// 遍历文件列表
	for _, v := range picSlice {
		// 上传文件，返回文件url
		g.postOne(&v)
	}
}

func (g *GiteeRequest) postOne(pic *string) {
	// 读取文件
	fileByte, err := os.ReadFile(*pic)
	if err != nil {
		fmt.Println("打开文件失败:", err)
		return
	}
	// content "base64编码后的字符串"
	g.Content = base64.StdEncoding.EncodeToString(fileByte)
	// Path "上传文件路径"
	if g.Path == "" {
		g.Path = filepath.Base(*pic)
	} else {
		g.Path = strings.Trim(g.Path, "/") + "/" + filepath.Base(*pic)
	}
	// message "Upload 文件名 by upPic"
	timestemp = time.Now().Format("2006-01-02 15:04:05")
	g.Message = "Upload " + g.Path + " by upPic - " + timestemp
	// url "https://gitee.com/api/v5/repos/luoxian1011/pictures/contents/pic.test3"
	postUrl := apiUrl + g.Owner + "/" + g.Repo + "/contents/" + g.Path
	// 序列化请求参数
	data, err := json.Marshal(g)
	if err != nil {
		fmt.Println("请求数据序列化失败:", err)
		return
	}
	// 路径置空 --
	// 处理的小bug，因为path没有置空，多文件时会将前面的文件名更改为前缀路径，上传显示成功，实际成功了一个寂寞（给文件夹做了个提交），Content也为nil
	g.Path = ""
	// 开始上传文件
	response, err := client.Post(postUrl, contentType, bytes.NewReader(data))
	if err != nil {
		fmt.Println("上传文件失败:", err)
		return
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("读取响应失败! 响应码:", response.StatusCode, err)
		return
	}
	defer response.Body.Close() // 关闭
	// 如果状态码不是200，就是响应错误
	if response.StatusCode != 201 {
		fmt.Println("请求失败! 响应码:", response.StatusCode, string(body))
		return
	}
	// 序列化响应体
	var giteeResponse GiteeResponse
	err = json.Unmarshal(body, &giteeResponse)
	if err != nil {
		fmt.Println("序列化响应体失败:", err)
	}
	// https://gitee.com/luoxian1011/pictures/raw/master/pic.test
	fmt.Println("Upload Success:")
	// 输出文件url
	fmt.Println(giteeResponse.Content["download_url"])
}

func main() {
	// gitee pic
	// 命令行参数: 从第五个参数开始传入文件路径
	argsLen := len(os.Args)
	if argsLen < 6 {
		fmt.Println("参数输入有误:")
		fmt.Println("Usage: upPic.exe access_token owner repo branch path file...")
		return
	}
	// 创建请求结构体
	giteeRequest := &GiteeRequest{
		Access_token: os.Args[1],
		Owner:        os.Args[2],
		Repo:         os.Args[3],
		Branch:       os.Args[4],
		Path:         os.Args[5],
	}
	// 拿到文件路径切片
	picSlice := os.Args[6:]
	// 上传图片
	giteeRequest.putPics(picSlice)
}
