package module

import (
	"fmt"
	"gatherSearch/utils"
	"github.com/imroc/req/v3"
	"log"
	"strconv"
	"time"
)

type fofaResponse struct {
	Error           bool          `yaml:"error"`
	ConsumedFpoint  int           `yaml:"consumed_fpoint"`
	RequiredFpoints int           `yaml:"required_fpoints"`
	Size            int           `yaml:"size"`
	Page            int           `yaml:"page"`
	Mode            string        `yaml:"mode"`
	Query           string        `yaml:"query"`
	Results         []interface{} `yaml:"results"`
}

func RunFofa(search string, filename string) {
	var results fofaResponse
	var dataTotal int
	config := utils.GetConfig()
	fofaurl := config.Module.Fofa.URL
	fofaemail := config.Module.Fofa.Email
	fofatoken := config.Module.Fofa.Token
	fofaReq := fofaurl + "?full=true&fields=domain,host,title,ip,port&full=false&page=1&size=1&email=" + fofaemail + "&key=" + fofatoken + "&qbase64=" + search
	client := req.C()
	response, err := client.R().SetSuccessResult(&results).Get(fofaReq)
	time.Sleep(time.Second * 2)
	if err != nil {
		log.Println("fofa request error:", err)
		return
	}
	if response.IsSuccessState() {
		dataTotal = results.Size
	} else {
		log.Println("fofa request error:", err)
		return
	}
	if dataTotal > 0 {
		pageSize := 100 // 每页处理 100 条数据
		if dataTotal > 3000 {
			dataTotal = utils.Config.Module.Fofa.Size
		}
		// 计算总页数
		totalPages := (dataTotal + pageSize - 1) / pageSize

		// 使用循环逐页处理数据
		for pageIndex := 1; pageIndex <= totalPages; pageIndex++ {
			time.Sleep(time.Second * 2)
			fofaDataReq := fofaurl + "?full=true&fields=domain,host,title,ip,port&full=false&page=" + strconv.Itoa(pageIndex) + "&size=" + strconv.Itoa(pageSize) + "&email=" + fofaemail + "&key=" + fofatoken + "&qbase64=" + search
			fofaResponse, err := client.R().SetSuccessResult(&results).Get(fofaDataReq)
			if err != nil {
				log.Println("Fofa API request failed")
				return
			}

			if fofaResponse.IsSuccessState() {
				var rows [][]string
				// 循环输出 results 列表
				for _, item := range results.Results {
					result, ok := item.([]interface{})
					if !ok {
						log.Println("无效的结果项")
						continue
					}

					if len(result) >= 5 {
						domain := utils.ToString(result[0])
						url := utils.ToString(result[1])
						webTitle := utils.ToString(result[2])
						ip := utils.ToString(result[3])
						port := utils.ToString(result[4])
						row := []string{domain, url, webTitle, ip, port}
						rows = append(rows, row)
						fmt.Printf("%-20s %-50s %-50s %-30s %-30s\n", domain, url, webTitle, ip, port)
					}
				}
				if len(rows) > 0 {
					err := utils.WriteDataToExcel(filename, rows)
					if err != nil {
						log.Println("保存文件时发生错误:", err)
						return
					}
					rows = nil
				}
			} else {
				log.Println("Fofa API request failed")
			}

		}

	}
}
