// Package config
// @Title  title
// @Description  desc
// @Author  yr  2025/1/22
// @Update  yr  2025/1/22
package config

type conf struct {
	FileDir          string   `binding:"required"` // xlsx文件所在目录
	ClientOutPath    string   `binding:"required"` // 客户端输出目录
	ServerOutPath    string   `binding:"required"` // 服务端输出目录
	GMOutPath        string   `binding:"required"` // gm输出目录
	ParseFileSuffix  []string `binding:""`         // 解析器需要解析的文件后缀(默认xlsx)
	ParserPrefix     string   `binding:""`         // 解析器需要解析的表前缀(如果没有前缀,默认只解析第一个文件)
	OutputFileFormat string   `binding:""`         // 输出文件格式(默认json)
}
