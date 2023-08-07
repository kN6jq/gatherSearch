package utils

import (
	"bufio"
	"fmt"
	"github.com/imroc/req/v3"
	"github.com/tealeg/xlsx"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"
)

func Req() *req.Request {
	return req.C().R().SetHeader("User-Agent", RandomUserAgent())
}

func SplitAndFormatSlice(list []string, batchSize int) []string {
	totalLength := len(list)
	segments := (totalLength + batchSize - 1) / batchSize

	searchData := make([]string, 0)

	for i := 0; i < segments; i++ {
		result := ""
		startIndex := i * batchSize
		endIndex := (i + 1) * batchSize
		if endIndex > totalLength {
			endIndex = totalLength
		}
		segment := list[startIndex:endIndex]
		for index, idata := range segment {

			data := fmt.Sprintf("%s", idata)
			if index < len(segment)-1 {
				data += " || "
			}
			result += data
		}
		searchData = append(searchData, result)
	}

	return searchData
}

func GenerateIPList(ipSegment string) ([]string, error) {
	// 解析IP段
	ip, ipNet, err := net.ParseCIDR(ipSegment)
	if err != nil {
		return nil, err
	}

	// 递增IP地址函数
	incIP := func(ip net.IP) {
		for j := len(ip) - 1; j >= 0; j-- {
			ip[j]++
			if ip[j] > 0 {
				break
			}
		}
	}

	// 计算IP地址范围
	var ipList []string
	for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); incIP(ip) {
		ipList = append(ipList, ip.String())
	}

	return ipList, nil
}

func RandomUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; AS; rv:11.0) like Gecko",
		"Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36 Edge/16.16299",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:58.0) Gecko/20100101 Firefox/58.0",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_3) AppleWebKit/604.5.6 (KHTML, like Gecko) Version/11.0.3 Safari/604.5.6",
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(userAgents))
	return userAgents[randomIndex]
}

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

func RemoveSpaces(input string) string {
	t1 := strings.ReplaceAll(input, " ", "")
	t2 := strings.ReplaceAll(t1, "\n", "")
	t3 := strings.ReplaceAll(t2, "\t", "")
	t4 := strings.ReplaceAll(t3, "｜", "|")
	return t4
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
func ShodanPortExcelFile() string {
	headerRow := []string{"port"}
	filename, err := CreateExcelFile("shodandb", headerRow)
	if err != nil {
		log.Fatalln("创建文件时发生错误:", err)
	}
	return filename
}
func ZoneSiteExcelFile() string {
	headerRow := []string{"Url", "title", "状态码", "IP", "端口"}
	filename, err := CreateExcelFile("zone", headerRow)
	if err != nil {
		log.Fatalln("创建文件时发生错误:", err)
	}
	return filename
}

func ZoneDomainExcelFile() string {
	headerRow := []string{"domain", "ip"}
	filename, err := CreateExcelFile("domain", headerRow)
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
