package module

import (
	"fmt"
	"github.com/kN6jq/gatherSearch/utils"
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

var (
	shodandbResults ShodandbData
)

func RunShodandb(data string, filename string) {
	url := "https://internetdb.shodan.io/" + data
	response, _ := utils.Req().SetSuccessResult(&shodandbResults).Get(url)
	if response.IsSuccessState() {
		var rows [][]string
		for i := range shodandbResults.Ports {
			ip := shodandbResults.Ip
			port := strconv.Itoa(shodandbResults.Ports[i])
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
