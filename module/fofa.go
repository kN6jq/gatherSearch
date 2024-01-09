package module

import (
	"encoding/base64"
	"fmt"
	"github.com/kN6jq/gatherSearch/utils"
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

var (
	fofaResults   fofaResponse
	fofadataTotal int
)

func RunFofa(data string, filename string) {
	search := base64.URLEncoding.EncodeToString([]byte(data))
	config := utils.GetConfig()
	fofaurl := config.Module.Fofa.URL
	fofaemail := config.Module.Fofa.Email
	fofatoken := config.Module.Fofa.Token
	fofaReq := fofaurl + "?full=true&fields=domain,host,title,ip,port&full=false&page=1&size=1&email=" + fofaemail + "&key=" + fofatoken + "&qbase64=" + search
	response, err := utils.Req().SetSuccessResult(&fofaResults).Get(fofaReq)
	if err != nil {
		log.Println("fofa request error:", err)
		return
	}
	if response.IsSuccessState() {
		fofadataTotal = fofaResults.Size
	} else {
		log.Println("fofa request error:", err)
		return
	}
	log.Printf("共搜索到数据: %d 个", fofadataTotal)
	if fofadataTotal > 0 {
		pageSize := 100 // 每页处理 100 条数据
		// 这里修改为通过开关设置获取全部还是指定数量
		if utils.Config.Module.Fofa.All == true {
			fofadataTotal = fofadataTotal
		} else if utils.Config.Module.Fofa.All == false {
			fofadataTotal = config.Module.Fofa.Size
		}
		// 计算总页数
		totalPages := (fofadataTotal + pageSize - 1) / pageSize

		// 使用循环逐页处理数据
		for pageIndex := 1; pageIndex <= totalPages; pageIndex++ {
			time.Sleep(time.Second * 3)
			fofaDataReq := fofaurl + "?full=true&fields=domain,host,title,ip,port&full=false&page=" + strconv.Itoa(pageIndex) + "&size=" + strconv.Itoa(pageSize) + "&email=" + fofaemail + "&key=" + fofatoken + "&qbase64=" + search
			fofaResponse, err := utils.Req().SetSuccessResult(&fofaResults).Get(fofaDataReq)
			if err != nil {
				log.Println("Fofa API request failed")
				return
			}

			if fofaResponse.IsSuccessState() {
				var rows [][]string
				// 循环输出 results 列表
				for _, item := range fofaResults.Results {
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
		time.Sleep(time.Second * 3)
	}
}
