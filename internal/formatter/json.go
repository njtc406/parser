// Package formatter
// @Title  title
// @Description  desc
// @Author  yr  2025/1/22
// @Update  yr  2025/1/22
package formatter

import "encoding/json"

type JsonFormatter struct {
}

func (j *JsonFormatter) Format(data interface{}) ([]byte, error) {
	if data == nil {
		return nil, nil
	}
	return json.MarshalIndent(data, "", "  ")
}
