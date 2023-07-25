package module

import (
	"fmt"
	"gatherSearch/utils"
	"log"
	"strconv"
)

type ShodandbData struct {
	Cpes      []string      `json:"cpes"`
	Hostnames []interface{} `json:"hostnames"`
	Ip        string        `json:"ip"`
	Ports     []int         `json:"ports"`
	Tags      []string      `json:"tags"`
	Vulns     []string      `json:"vulns"`
}

func RunShodandb(search string, filename string) {
	url := "https://internetdb.shodan.io/" + search
	var results ShodandbData
	response, _ := utils.Req().SetSuccessResult(&results).Get(url)
	if response.IsSuccessState() {
		var rows [][]string
		for i := range results.Ports {
			ip := results.Ip
			port := strconv.Itoa(results.Ports[i])
			rows = append(rows, []string{ip + ":" + port})
			fmt.Printf("%s:%s\n", ip, port)
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
	}
}
