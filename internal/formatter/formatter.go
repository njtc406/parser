// Package formatter
// @Title  title
// @Description  desc
// @Author  yr  2025/1/22
// @Update  yr  2025/1/22
package formatter

import "fmt"

type Formatter interface {
	Format(data interface{}) ([]byte, error)
}

type Factory struct {
}

func NewFactory() *Factory {
	return &Factory{}
}

func (f *Factory) NewFormatter(fileType string) (Formatter, error) {
	switch fileType {
	case "json":
		return &JsonFormatter{}, nil
	case "yaml":
		return &YamlFormatter{}, nil
	default:
		return nil, fmt.Errorf("不支持的输出文件类型: %s", fileType)
	}
}
