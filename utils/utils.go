package utils

import (
	"bufio"
	"fmt"
	"github.com/tealeg/xlsx"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		// 文件存在
		return true
	} else if os.IsNotExist(err) {
		// 文件不存在
		return false
	} else {
		// 发生错误
		return false
	}
}
func ReadFileLines(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func HunterExcelFile() string {
	headerRow := []string{"Domain", "URL", "WebTitle", "StatusCode", "IP", "Port"}
	filename, err := CreateExcelFile("hunter", headerRow)
	if err != nil {
		log.Fatalln("创建文件时发生错误:", err)
	}
	return filename
}
func FofaExcelFile() string {
	headerRow := []string{"Domain", "Host", "Title", "Ip", "Port"}
	filename, err := CreateExcelFile("fofa", headerRow)
	if err != nil {
		log.Fatalln("创建文件时发生错误:", err)
	}
	return filename
}
func ShodanExcelFile() string {
	headerRow := []string{"Domain", "Subdomain", "Ip"}
	filename, err := CreateExcelFile("shodan", headerRow)
	if err != nil {
		log.Fatalln("创建文件时发生错误:", err)
	}
	return filename
}

func CreateExcelFile(sheetName string, headerRow []string) (string, error) {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet(sheetName)
	if err != nil {
		return "", fmt.Errorf("无法创建工作表: %w", err)
	}

	// 插入表头
	row := sheet.AddRow()
	row.SetHeightCM(1) // 设置表头行高
	for _, header := range headerRow {
		cell := row.AddCell()
		cell.Value = header
	}

	// 生成文件名
	filename := fmt.Sprintf("%s.xlsx", GetCurrentTimestamp())

	// 保存文件
	err = file.Save(filename)
	if err != nil {
		return "", fmt.Errorf("无法保存文件: %w", err)
	}

	return filename, nil
}

func WriteDataToExcel(filename string, data [][]string) error {
	file, err := xlsx.OpenFile(filename)
	if err != nil {
		return fmt.Errorf("无法打开文件: %w", err)
	}

	sheet := file.Sheets[0] // 假设只有一个工作表

	// 保存数据
	for _, row := range data {
		rowObj := sheet.AddRow()
		for _, value := range row {
			cell := rowObj.AddCell()
			cell.Value = value
		}
	}

	err = file.Save(filename)
	if err != nil {
		return fmt.Errorf("无法保存文件: %w", err)
	}

	return nil
}

func ToString(value interface{}) string {
	str, ok := value.(string)
	if !ok {
		return strings.Repeat(" ", 50)
	}
	return limitString(str, 50)
}

func limitString(str string, width int) string {
	if len(str) > width {
		str = str[:width] // 截断字符串，确保不超过指定的宽度
	}
	return str
}

func IsIP(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}

func GetNowDate() string {
	return time.Now().Format("2006-01-02")
}
func GetLastYearDate() string {
	currentTime := time.Now()
	lastYear := currentTime.Year() - 1
	lastYearDate := time.Date(lastYear, currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())
	return lastYearDate.Format("2006-01-02")
}
func GetCurrentTimestamp() string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%d", timestamp)
}
