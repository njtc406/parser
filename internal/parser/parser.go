// Package parser
// @Title  title
// @Description  desc
// @Author  yr  2025/1/22
// @Update  yr  2025/1/22
package parser

import (
	"fmt"
	"io/fs"
	"parser/internal/config"
	"parser/internal/utils"
	"path/filepath"
	"strings"
)

type Parser interface {
	Parse(filePath string) error
}

// parserMap 解析器集合(后续来支持不同格式)
var parserMap = map[string]Parser{
	".xlsx": NewXlsxParser(),
	//".xls":  NewXlsParser(),
	//".csv":  NewCsvParser(),
}

type Hook func() error

var beginHooks = []Hook{
	ParseLanguage,
}

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) ExecHooks() error {
	for _, hook := range beginHooks {
		if err := hook(); err != nil {
			return err
		}
	}
	return nil
}

func (f *Factory) Parse() error {
	return filepath.Walk(config.Conf.FileDir, func(filePath string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if utils.SliceContains(config.Conf.ParseFileSuffix, strings.ToLower(filepath.Ext(filePath))) {
			x, ok := parserMap[filepath.Ext(filePath)]
			if !ok {
				return fmt.Errorf("not support file type: %s", filepath.Ext(filePath))
			}
			fmt.Println("解析文件: [", filePath, "]")
			return x.Parse(filePath)
		}
		return nil
	})
}
