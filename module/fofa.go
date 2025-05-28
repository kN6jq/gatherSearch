package module

import (
	"encoding/base64"
	"log"
	"strconv"
	"time"

	"github.com/kN6jq/gatherSearch/utils"
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
	fofatoken := config.Module.Fofa.Token
	fofaReq := fofaurl + "?full=true&fields=domain,host,title,ip,port&full=false&page=1&size=1&key=" + fofatoken + "&qbase64=" + search
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
		// 根据配置决定获取数量
		maxResults := fofadataTotal
		if !config.Module.Fofa.All {
			// 如果不获取全部，则限制为配置的size
			if config.Module.Fofa.Size < maxResults {
				maxResults = config.Module.Fofa.Size
				log.Printf("根据配置限制，将只获取 %d 条数据", maxResults)
			}
		}

		// 每页请求数量，确保不超过API限制
		pageSize := 100

		// 计算需要请求的页数
		totalPages := (maxResults + pageSize - 1) / pageSize

		// 计算最后一页需要的数量
		lastPageSize := maxResults % pageSize
		if lastPageSize == 0 && maxResults > 0 {
			lastPageSize = pageSize
		}

		// 打印表头
		utils.PrintTableHeader()

		// 使用循环逐页处理数据
		var processedCount int = 0
		for pageIndex := 1; pageIndex <= totalPages; pageIndex++ {
			time.Sleep(time.Second * 3)

			// 计算当前页应获取的数量
			currentPageSize := pageSize
			if pageIndex == totalPages && lastPageSize > 0 {
				currentPageSize = lastPageSize
			}

			fofaDataReq := fofaurl + "?full=true&fields=domain,host,title,ip,port&full=false&page=" + strconv.Itoa(pageIndex) + "&size=" + strconv.Itoa(currentPageSize) + "&key=" + fofatoken + "&qbase64=" + search
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
						area, country, _ := utils.QueryIp(ip)
						row := []string{domain, url, webTitle, ip, port, country, area}
						rows = append(rows, row)

						// 使用表格形式输出
						utils.PrintTableRow(domain, url, webTitle, ip, port, country, area)

						processedCount++
						if !config.Module.Fofa.All && processedCount >= config.Module.Fofa.Size {
							// 达到限制数量，提前结束
							log.Printf("已达到配置的最大获取数量: %d", config.Module.Fofa.Size)
							break
						}
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

				// 如果已达到限制数量，提前结束循环
				if !config.Module.Fofa.All && processedCount >= config.Module.Fofa.Size {
					break
				}
			} else {
				log.Println("Fofa API request failed")
			}
		}

		// 打印表格底部分隔线
		utils.PrintTableFooter()

		time.Sleep(time.Second * 3)
	}
}
