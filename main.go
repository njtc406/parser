// Package main
// @Title  title
// @Description  desc
// @Author  yr  2025/1/22
// @Update  yr  2025/1/22
package main

import (
	"fmt"
	"os"
	"parser/internal/config"
	"parser/internal/parser"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("======================运行失败=======================")
			fmt.Println("error:", r)
			fmt.Println("==============================================")
		}
	}()
	var confPath string
	if len(os.Args) > 1 {
		confPath = os.Args[1]
	}

	// 解析配置文件
	config.Init(confPath)

	fmt.Println("======================配置信息=======================")
	fmt.Println()
	fmt.Println("xlsx文件路径:", config.Conf.FileDir)
	fmt.Println()
	fmt.Println("客户端输出文件路径:", config.Conf.ClientOutPath)
	fmt.Println()
	fmt.Println("服务端输出文件路径:", config.Conf.ServerOutPath)
	fmt.Println()
	fmt.Println("解析文件后缀:", config.Conf.ParseFileSuffix)
	fmt.Println()
	fmt.Println("解析前缀:", config.Conf.ParserPrefix)
	fmt.Println()
	fmt.Println("输出文件格式:", config.Conf.OutputFileFormat)
	fmt.Println()

	// 创建一个解析器
	p := parser.NewFactory()

	fmt.Println("======================开始解析=======================")

	if err := p.ExecHooks(); err != nil {
		fmt.Println("======================解析失败=======================")
		fmt.Println("error:", err)
		fmt.Println("==============================================")
		goto DoEnd
	}

	// 开始解析
	if err := p.Parse(); err != nil {
		fmt.Println("======================解析失败=======================")
		fmt.Println("error:", err)
		fmt.Println("==============================================")
		goto DoEnd
	}

	fmt.Println("======================解析完成=======================")

DoEnd:
	fmt.Println("回车键退出...")
	_, _ = fmt.Scanln()
}
