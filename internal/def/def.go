// Package def
// @Title  title
// @Description  desc
// @Author  yr  2025/1/22
// @Update  yr  2025/1/22
package def

const (
	XlsxColTypeInt   = "int"
	XlsxColTypeStr   = "string"
	XlsxColTypeBool  = "bool"
	XlsxColTypeFloat = "float"
)

var TypeMap = map[string]interface{}{
	"int":    0,
	"string": "",
	"bool":   false,
	"float":  float64(0),
}
