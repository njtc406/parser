// Package formatter
// @Title  title
// @Description  desc
// @Author  yr  2025/1/22
// @Update  yr  2025/1/22
package formatter

import "gopkg.in/yaml.v3"

type YamlFormatter struct {
}

func (y *YamlFormatter) Format(data interface{}) ([]byte, error) {
	if data == nil {
		return nil, nil
	}
	return yaml.Marshal(data)
}
