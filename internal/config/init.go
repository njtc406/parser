// Package config
// @Title  title
// @Description  desc
// @Author  yr  2025/1/22
// @Update  yr  2025/1/22
package config

import (
	"github.com/njtc406/chaosutil/validate"
	"github.com/njtc406/viper"
	"os"
	"path/filepath"
)

var (
	runtimeViper = viper.New()
	Conf         = new(conf)
)

const defaultConfPath = "./configs"

func Init(confPath string) {
	// 解析配置
	parseConf(confPath)
	clearOutDir()
	// 初始化目录
	initDir()
}

// parseConf 解析本地配置文件
func parseConf(confPath string) {
	if confPath == "" {
		confPath = defaultConfPath
	}

	// 解析配置路径
	runtimeViper.SetConfigType("yaml")
	runtimeViper.SetConfigName("conf")
	runtimeViper.AddConfigPath(filepath.ToSlash(confPath))

	if err := runtimeViper.ReadInConfig(); err != nil {
		panic(err)
	} else if err = runtimeViper.Unmarshal(Conf); err != nil {
		panic(err)
	} else if err = validate.Struct(Conf); err != nil {
		panic(validate.TransError(err, validate.ZH))
	}
}

// initDir 创建必要的目录
func initDir() {
	createDirIfNotExists(Conf.ClientOutPath)
	createDirIfNotExists(Conf.ServerOutPath)
	createDirIfNotExists(Conf.GMOutPath)
}

// createDirIfNotExists 创建目录
func createDirIfNotExists(dir string) {
	if err := os.MkdirAll(dir, 0644); err != nil {
		panic(err)
	}
}

func clearOutDir() error {
	// 清理输出目录
	if err := os.RemoveAll(Conf.ClientOutPath); err != nil {
		return err
	}
	if err := os.RemoveAll(Conf.ServerOutPath); err != nil {
		return err
	}
	if err := os.RemoveAll(Conf.GMOutPath); err != nil {
		return err
	}
	return nil
}
