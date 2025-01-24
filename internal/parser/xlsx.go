// Package parser
// @Title  title
// @Description  desc
// @Author  yr  2025/1/22
// @Update  yr  2025/1/22
package parser

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"parser/internal/config"
	"parser/internal/def"
	"parser/internal/formatter"
	"parser/internal/writer"
	"path"
	"strconv"
	"strings"
)

type headerInfo struct {
	index  int    // 列
	desc   string // 描述
	belong string // 所属(csgn,c 为client s 为server g为gm n表示不会解析)
	name   string // 字段名
	tp     string // 类型
}

type xlsxParser struct {
	headers      []*headerInfo
	resultClient []map[string]interface{}
	resultServer []map[string]interface{}
	resultGm     []map[string]interface{}
}

func NewXlsxParser() *xlsxParser {
	return &xlsxParser{}
}

func (x *xlsxParser) write(fileName string) error {
	fileName = strings.TrimPrefix(fileName, config.Conf.ParserPrefix) // 去掉文件名的前缀
	// 写入数据
	data, err := formatter.NewFactory().NewFormatter(config.Conf.OutputFileFormat)
	if err != nil {
		return err
	}
	if len(x.resultClient) > 0 {
		clientData, err := data.Format(x.resultClient)
		if err != nil {
			return err
		}
		if len(clientData) > 0 {
			if err = writer.Write(path.Join(config.Conf.ClientOutPath, fileName+"."+config.Conf.OutputFileFormat), clientData); err != nil {
				return err
			}
		}
	}

	if len(x.resultServer) > 0 {
		serverData, err := data.Format(x.resultServer)
		if err != nil {
			return err
		}
		if len(serverData) > 0 {
			if err = writer.Write(path.Join(config.Conf.ServerOutPath, fileName+"."+config.Conf.OutputFileFormat), serverData); err != nil {
				return err
			}
		}
	}

	if len(x.resultGm) > 0 {
		gmData, err := data.Format(x.resultGm)
		if err != nil {
			return err
		}
		if len(gmData) > 0 {
			if err = writer.Write(path.Join(config.Conf.GMOutPath, fileName+"."+config.Conf.OutputFileFormat), gmData); err != nil {
				return err
			}
		}
	}
	return nil
}

func (x *xlsxParser) Parse(filePath string) error {
	// 解析 xlsx 文件
	f, err := xlsx.OpenFile(filePath)
	if err != nil {
		return err
	}

	var parseCount int
	if config.Conf.ParserPrefix == "" {
		parseCount = 1 // 默认解析第一个 sheet
	} else {
		parseCount = len(f.Sheets)
	}

	for i := 0; i < parseCount; i++ {
		sheet := f.Sheets[i]
		if config.Conf.ParserPrefix == "" || strings.HasPrefix(sheet.Name, config.Conf.ParserPrefix) {
			//labelName := strings.TrimPrefix(sheet.Name, config.Conf.ParserPrefix)
			fmt.Println("  -->开始解析表:", sheet.Name)
			if err = x.parseSheet(sheet); err != nil {
				return err
			}
			fmt.Println("  -->解析完成表:", sheet.Name)
			fmt.Println("    -->>开始写入:", sheet.Name)
			if err = x.write(sheet.Name); err != nil {
				return err
			}
			fmt.Println("    -->>写入完成:", sheet.Name)
		}
	}
	return nil
}

func (x *xlsxParser) parseSheet(sheet *xlsx.Sheet) error {
	// 解析表头
	for idx, row := range sheet.Rows[:4] {
		for colIndex, cell := range row.Cells {
			// 如果当前列的表头信息未初始化，则初始化
			if colIndex >= len(x.headers) {
				x.headers = append(x.headers, &headerInfo{index: colIndex})
			}

			// 根据行索引填充表头信息
			switch idx {
			case 0:
				x.headers[colIndex].desc = cell.String()
			case 1:
				x.headers[colIndex].name = cell.String()
			case 2:
				x.headers[colIndex].tp = cell.String()
			case 3:
				x.headers[colIndex].belong = cell.String()
			}
		}
	}

	// 打印表头信息
	//for _, header := range x.headers {
	//	fmt.Printf("Header: %+v\n", header)
	//}

	// 解析数据
	for rowIndex, row := range sheet.Rows[4:] {
		clientRowData := make(map[string]interface{})
		serverRowData := make(map[string]interface{})
		gmRowData := make(map[string]interface{})
		// 遍历每一列
		for colIndex := 0; colIndex < len(x.headers); colIndex++ {
			var cellStr string
			// 这里是防止这行配置不全,导致一些列读取不到
			if colIndex >= len(row.Cells) {
				cellStr = ""
			} else {
				cell := row.Cells[colIndex]
				cellStr = cell.String()
			}

			header := x.headers[colIndex]
			if header.belong == "n" {
				continue // 跳过不需要解析的列
			}

			// 解析单元格值
			value, err := x.parseCellValue(header.tp, cellStr)
			if err != nil {
				return fmt.Errorf("解析表[%s] 行:%d, 列:%d error:%v", sheet.Name, rowIndex+5, colIndex+1, err)
			}
			// 将解析后的值存储
			if strings.Contains(header.belong, "c") {
				clientRowData[header.name] = value
			}
			if strings.Contains(header.belong, "s") {
				serverRowData[header.name] = value
			}
			if strings.Contains(header.belong, "g") {
				gmRowData[header.name] = value
			}
		}
		x.resultClient = append(x.resultClient, clientRowData)
		x.resultServer = append(x.resultServer, serverRowData)
		x.resultGm = append(x.resultGm, gmRowData)
	}

	// 打印解析结果
	//for _, rowData := range x.resultClient {
	//	fmt.Println("clientRowData:", rowData)
	//}
	//for _, rowData := range x.resultServer {
	//	fmt.Println("serverRowData:", rowData)
	//}
	//for _, rowData := range x.resultGm {
	//	fmt.Println("gmRowData:", rowData)
	//}
	return nil
}

func GetTypeDefaultValue(typeName string) (interface{}, error) {
	// 基础类型
	if value, ok := def.TypeMap[typeName]; ok {
		return value, nil
	}

	// 处理map类型
	if strings.HasPrefix(typeName, "map[") && strings.Contains(typeName, "]") {
		keyValueType := strings.Split(strings.TrimSuffix(strings.TrimPrefix(typeName, "map["), "]"), "]")
		if len(keyValueType) != 2 {
			return nil, fmt.Errorf("invalid map type format: %s", typeName)
		}
		keyType := keyValueType[0]

		// 根据 keyType 和 valueType 返回空的 map
		switch keyType {
		case def.XlsxColTypeStr:
			return map[string]interface{}{}, nil
		case def.XlsxColTypeInt:
			return map[int]interface{}{}, nil
		case def.XlsxColTypeFloat:
			return map[float64]interface{}{}, nil
		case def.XlsxColTypeBool:
			return map[bool]interface{}{}, nil
		default:
			return nil, fmt.Errorf("不支持作为 map 的 key 的类型: %s", keyType)
		}
	}

	// 处理数组类型
	if strings.HasSuffix(typeName, "[]") {
		return []interface{}{}, nil
	}

	// 处理其他复杂类型（如自定义结构体等）

	// 这里可以根据需要扩展
	return nil, fmt.Errorf("不支持的数据类型: %s", typeName)
}

// 定义每层的分隔符
var delimiters = []string{"|", ";", ":", ","} // 顺序是从最外的一层分隔符到最内的一层分隔符

func (x *xlsxParser) parseCellValue(tp string, value string) (interface{}, error) {
	return x.parseCellValueWithDelimiter(tp, value, 0)
}

func (x *xlsxParser) parseCellValueWithDelimiter(tp string, value string, depth int) (interface{}, error) {
	// 去掉value中的空格
	value = strings.TrimSpace(value)
	if value == "" {
		return GetTypeDefaultValue(tp)
	}

	if depth > len(delimiters) {
		return nil, fmt.Errorf("超过数据最大嵌套层数: %s", tp)
	}

	// 处理 map 类型
	if strings.HasPrefix(tp, "map[") && strings.Contains(tp, "]") {
		var delimiter1 string
		var delimiter2 string
		cnt := strings.Count(tp, "map")
		// 每层map需要两个分隔符
		if depth > len(delimiters)-cnt*2 {
			return nil, fmt.Errorf("超过 map 最大嵌套层数: %s", tp)
		}
		if cnt > 1 {
			delimiter1 = delimiters[len(delimiters)-cnt*2:][depth]
			delimiter2 = delimiters[len(delimiters)-cnt*2+1:][depth]
		} else {
			delimiter1 = delimiters[len(delimiters)-2]
			delimiter2 = delimiters[len(delimiters)-1]
		}
		keyType, valueType, err := parseMapType(tp)
		if err != nil {
			return nil, err
		}

		keyValuePairs := strings.Split(value, delimiter1) // 使用当前层级的分隔符
		result := make(map[string]interface{})            // 由于json格式的map的key必须是string类型，所以这里需要将key转换为string类型
		for _, pair := range keyValuePairs {
			kv := strings.Split(pair, delimiter2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("invalid map key-value pair: %s", pair)
			}
			parsedKey, err := x.parseCellValueWithDelimiter(keyType, kv[0], depth+1)
			if err != nil {
				return nil, err
			}
			parsedValue, err := x.parseCellValueWithDelimiter(valueType, kv[1], depth+1)
			if err != nil {
				return nil, err
			}
			result[fmt.Sprintf("%v", parsedKey)] = parsedValue
		}
		return result, nil
	}

	// 处理数组类型
	if strings.HasSuffix(tp, "[]") {
		cnt := strings.Count(tp, "[]")
		var delimiter string
		if depth > len(delimiters)-cnt-1 {
			return nil, fmt.Errorf("超过数组最大支持维数: %s", tp)
		}
		if cnt > 1 {
			delimiter = delimiters[len(delimiters)-cnt:][depth]
		} else {
			delimiter = delimiters[len(delimiters)-1]
		}

		elementType := strings.TrimSuffix(tp, "[]")
		values := strings.Split(value, delimiter) // 使用当前层级的分隔符
		var result []interface{}
		for _, v := range values {
			parsedValue, err := x.parseCellValueWithDelimiter(elementType, v, depth)
			if err != nil {
				return nil, err
			}
			result = append(result, parsedValue)
		}
		return result, nil
	}

	// 基础类型
	switch tp {
	case def.XlsxColTypeInt:
		return strconv.Atoi(value)
	case def.XlsxColTypeStr:
		return value, nil
	case def.XlsxColTypeBool:
		return strconv.ParseBool(value)
	case def.XlsxColTypeFloat:
		v, err := strconv.ParseFloat(value, 64)
		return v, err
	default:
		return nil, fmt.Errorf("不支持的数据类型: %s", tp)
	}
}

func parseMapType(typeStr string) (keyType string, valueType string, err error) {
	// 去掉开头的 "map[" 和结尾的 "]"
	typeStr = strings.TrimPrefix(typeStr, "map[")
	typeStr = strings.TrimSuffix(typeStr, "]")

	// 分割键类型和值类型
	parts := strings.SplitN(typeStr, "]", 2)
	if len(parts) != 2 {
		err = fmt.Errorf("invalid map type format: %s", typeStr)
		return
	}

	keyType = parts[0]
	valueType = parts[1]

	return
}
