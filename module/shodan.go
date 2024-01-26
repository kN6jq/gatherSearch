package module

import (
	"fmt"
	"github.com/kN6jq/gatherSearch/utils"
	"log"
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
	response, err := utils.Req().SetSuccessResult(&shodanResults).Get(shodanReq)
	if err != nil {
		log.Println("shodan request error:", err)
		return
	}
	if response.IsSuccessState() {
		var rows [][]string
		for i := range shodanResults.Data {
			subdomain := ""
			if shodanResults.Data[i].Subdomain == "" {
				subdomain = "*." + data
			} else {
				subdomain = shodanResults.Data[i].Subdomain + "." + data
			}
			ip := shodanResults.Data[i].Value
			area, country, _ := utils.QueryIp(ip)
			rows = append(rows, []string{data, subdomain, ip, area, country})
			fmt.Printf("%-20s %-20s %-20s %-20s %-20s\n", data, subdomain, ip, area, country)
		}
		if len(rows) > 0 {
			// 保存数据到 Excel 文件
			err := utils.WriteDataToExcel(filename, rows)
			if err != nil {
				log.Println("保存文件时发生错误:", err)
				return
			}
			rows = nil
		}
	} else {
		log.Println("shodan request error:", err)
		return
	}
}
