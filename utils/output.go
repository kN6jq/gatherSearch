package utils

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// PrintTableHeader 打印表格表头
func PrintTableHeader() {
	// 定义列宽 - 增加各列宽度
	domainWidth := 30
	urlWidth := 50
	titleWidth := 40
	ipWidth := 30 // 确保IPv6地址可以完整显示
	portWidth := 8
	countryWidth := 15
	areaWidth := 25

	// 打印表头分隔线
	fmt.Println(strings.Repeat("-", domainWidth+urlWidth+titleWidth+ipWidth+portWidth+countryWidth+areaWidth+14))

	// 打印表头
	fmt.Printf("| %-*s | %-*s | %-*s | %-*s | %-*s | %-*s | %-*s |\n",
		domainWidth-2, "域名",
		urlWidth-2, "URL",
		titleWidth-2, "标题",
		ipWidth-2, "IP",
		portWidth-2, "端口",
		countryWidth-2, "国家",
		areaWidth-2, "地区")

	// 打印表头分隔线
	fmt.Println(strings.Repeat("-", domainWidth+urlWidth+titleWidth+ipWidth+portWidth+countryWidth+areaWidth+14))
}

// PrintTableRow 打印表格行
func PrintTableRow(domain, url, title, ip, port, country, area string) {
	// 定义列宽 - 与表头保持一致
	domainWidth := 30
	urlWidth := 50
	titleWidth := 40
	ipWidth := 30
	portWidth := 8
	countryWidth := 15
	areaWidth := 25

	// 截断并处理可能的中文字符
	domain = formatString(domain, domainWidth-2)
	url = formatString(url, urlWidth-2)
	title = formatString(title, titleWidth-2)
	ip = formatString(ip, ipWidth-2, true) // 特殊处理IP地址
	port = formatString(port, portWidth-2)
	country = formatString(country, countryWidth-2)
	area = formatString(area, areaWidth-2)

	// 打印数据行
	fmt.Printf("| %-*s | %-*s | %-*s | %-*s | %-*s | %-*s | %-*s |\n",
		domainWidth-2, domain,
		urlWidth-2, url,
		titleWidth-2, title,
		ipWidth-2, ip,
		portWidth-2, port,
		countryWidth-2, country,
		areaWidth-2, area)
}

// formatString 格式化字符串，处理中文字符宽度并根据需要截断
func formatString(s string, maxWidth int, isIP ...bool) string {
	if maxWidth <= 0 {
		return ""
	}

	// 如果字符串为空，直接返回
	if s == "" {
		return s
	}

	// 特殊处理IP地址，确保IP地址不被截断
	if (len(isIP) > 0 && isIP[0]) ||
		(len(s) <= 15 && strings.Count(s, ".") == 3) || // IPv4
		strings.Contains(s, ":") { // IPv6
		return s // 不截断IP地址
	}

	// 计算字符串的显示宽度
	width := 0
	var result []rune
	for _, r := range s {
		// 中文字符通常占用2个显示宽度
		var charWidth int
		if r > 0x7F { // 非ASCII字符（如中文）
			charWidth = 2
		} else {
			charWidth = 1
		}

		// 如果添加这个字符会超过最大宽度，则添加省略号并结束
		if width+charWidth > maxWidth {
			// 确保有足够空间放置省略号
			if maxWidth >= 3 {
				// 移除最后几个字符以腾出空间放置省略号
				for utf8.RuneCountInString(string(result)) > 0 && width+3 > maxWidth {
					lastRune := result[len(result)-1]
					result = result[:len(result)-1]
					if lastRune > 0x7F {
						width -= 2
					} else {
						width -= 1
					}
				}
				result = append(result, []rune("...")...)
			}
			break
		}

		// 添加当前字符到结果中
		result = append(result, r)
		width += charWidth
	}

	return string(result)
}

// PrintTableFooter 打印表格底部分隔线
func PrintTableFooter() {
	// 定义列宽 - 与表头保持一致
	domainWidth := 30
	urlWidth := 50
	titleWidth := 40
	ipWidth := 30
	portWidth := 8
	countryWidth := 15
	areaWidth := 25

	// 打印表底分隔线
	fmt.Println(strings.Repeat("-", domainWidth+urlWidth+titleWidth+ipWidth+portWidth+countryWidth+areaWidth+14))
}
