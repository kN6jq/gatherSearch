package module

import (
	"fmt"
	"gatherSearch/utils"
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

func RunShodan(search string, filename string) {
	var results shodanResponse
	config := utils.GetConfig()
	shodanUrl := config.Module.Shodan.URL
	shodanKey := config.Module.Shodan.Key
	shodanReq := shodanUrl + "/dns/domain/" + search + "?key=" + shodanKey
	response, err := utils.Req().SetSuccessResult(&results).Get(shodanReq)
	if err != nil {
		log.Println("shodan request error:", err)
		return
	}
	if response.IsSuccessState() {
		var rows [][]string
		for i := range results.Data {
			subdomain := ""
			if results.Data[i].Subdomain == "" {
				subdomain = "*." + search
			} else {
				subdomain = results.Data[i].Subdomain + "." + search
			}
			ip := results.Data[i].Value
			rows = append(rows, []string{search, subdomain, ip})
			fmt.Printf("%-20s %-20s %-20s \n", search, subdomain, ip)
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
