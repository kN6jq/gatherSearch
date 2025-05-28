package module

import (
	"log"
	"strconv"

	"github.com/kN6jq/gatherSearch/utils"
)

type ShodandbData struct {
	Cpes      []string      `json:"cpes"`
	Hostnames []interface{} `json:"hostnames"`
	Ip        string        `json:"ip"`
	Ports     []int         `json:"ports"`
	Tags      []string      `json:"tags"`
	Vulns     []string      `json:"vulns"`
}

var (
	shodandbResults ShodandbData
)

func RunShodandb(data string, filename string) {
	url := "https://internetdb.shodan.io/" + data

	response, err := utils.Req().SetSuccessResult(&shodandbResults).Get(url)
	if err != nil {
		log.Println("shodandb request error:", err)
		return
	}

	if !response.IsSuccessState() {
		log.Println("shodandb API请求失败")
		return
	}

	// 检查IP是否为空
	if shodandbResults.Ip == "" {
		log.Printf("未找到IP %s 的相关数据", data)
		return
	}

	// 检查端口列表是否为空
	if len(shodandbResults.Ports) == 0 {
		log.Printf("IP %s 没有开放端口信息", data)
		return
	}

	// 打印表头
	utils.PrintTableHeader()

	var rows [][]string
	for i := range shodandbResults.Ports {
		ip := shodandbResults.Ip
		port := strconv.Itoa(shodandbResults.Ports[i])

		// 获取IP地理位置，忽略错误
		area, country, _ := utils.QueryIp(ip)

		// 从主机名列表中获取域名
		domain := getDomainFromHostnames(shodandbResults.Hostnames)

		// 使用表格输出
		utils.PrintTableRow(domain, "", "", ip, port, country, area)

		rows = append(rows, []string{domain, "", "", ip, port, country, area})
	}

	if len(rows) > 0 {
		// 保存数据到 Excel 文件
		err := utils.WriteDataToExcel(filename, rows)
		if err != nil {
			log.Println("保存文件时发生错误:", err)
			return
		}
	} else {
		log.Println("没有有效数据可保存")
	}

	// 打印表格底部分隔线
	utils.PrintTableFooter()
}

// 从主机名列表中获取域名
func getDomainFromHostnames(hostnames []interface{}) string {
	if len(hostnames) == 0 {
		return ""
	}

	// 尝试获取第一个主机名
	if hostname, ok := hostnames[0].(string); ok {
		return hostname
	}

	return ""
}
