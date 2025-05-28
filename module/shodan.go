package module

import (
	"log"

	"github.com/kN6jq/gatherSearch/utils"
)

type shodanResponse struct {
	Domain     string     `json:"domain"`
	Tags       []string   `json:"tags"`
	Subdomains []string   `json:"subdomains"`
	Data       []DataItem `json:"data"`
	More       bool       `json:"more"`
}

type DataItem struct {
	Tags      []string `json:"tags"`
	Subdomain string   `json:"subdomain"`
	Type      string   `json:"type"`
	Ports     []int    `json:"ports"`
	Value     string   `json:"value"`
	LastSeen  string   `json:"last_seen"`
}

var (
	shodanResults shodanResponse
)

func RunShodan(data string, filename string) {
	config := utils.GetConfig()
	shodanUrl := config.Module.Shodan.URL
	shodanKey := config.Module.Shodan.Key
	shodanReq := shodanUrl + "/dns/domain/" + data + "?key=" + shodanKey

	log.Printf("正在请求Shodan API: %s", shodanReq)
	response, err := utils.Req().SetSuccessResult(&shodanResults).Get(shodanReq)
	if err != nil {
		log.Println("shodan request error:", err)
		return
	}

	if !response.IsSuccessState() {
		log.Println("shodan API请求失败")
		return
	}

	// 检查是否有数据
	if len(shodanResults.Data) == 0 {
		log.Printf("未找到域名 %s 的相关数据", data)
		return
	}

	// 打印表头
	utils.PrintTableHeader()

	var rows [][]string
	for i := range shodanResults.Data {
		// 检查数据有效性
		if shodanResults.Data[i].Value == "" {
			log.Printf("第 %d 条数据的IP地址为空，跳过", i+1)
			continue
		}

		subdomain := ""
		if shodanResults.Data[i].Subdomain == "" {
			subdomain = "*." + data
		} else {
			subdomain = shodanResults.Data[i].Subdomain + "." + data
		}

		ip := shodanResults.Data[i].Value

		// 查询IP地理位置，忽略错误
		area, country, _ := utils.QueryIp(ip)

		// 使用表格输出
		utils.PrintTableRow(data, subdomain, "", ip, "", country, area)

		rows = append(rows, []string{data, subdomain, "", ip, "", country, area})
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
