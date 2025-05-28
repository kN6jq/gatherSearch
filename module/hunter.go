package module

import (
	"encoding/base64"
	"log"
	"strconv"
	"time"

	"github.com/kN6jq/gatherSearch/utils"
)

type hunterResults struct {
	Code int `json:"code"`
	Data struct {
		AccountType string `json:"account_type"`
		Total       int    `json:"total"`
		Time        int    `json:"time"`
		Arr         []struct {
			IsRisk         string `json:"is_risk"`
			URL            string `json:"url"`
			IP             string `json:"ip"`
			Port           int    `json:"port"`
			WebTitle       string `json:"web_title"`
			Domain         string `json:"domain"`
			IsRiskProtocol string `json:"is_risk_protocol"`
			Protocol       string `json:"protocol"`
			BaseProtocol   string `json:"base_protocol"`
			StatusCode     int    `json:"status_code"`
			Component      []struct {
				Name    string `json:"name"`
				Version string `json:"version"`
			} `json:"component"`
			Os        string `json:"os"`
			Company   string `json:"company"`
			Number    string `json:"number"`
			Country   string `json:"country"`
			Province  string `json:"province"`
			City      string `json:"city"`
			UpdatedAt string `json:"updated_at"`
			IsWeb     string `json:"is_web"`
			AsOrg     string `json:"as_org"`
			Isp       string `json:"isp"`
			Banner    string `json:"banner"`
		} `json:"arr"`
		ConsumeQuota string `json:"consume_quota"`
		RestQuota    string `json:"rest_quota"`
		SyntaxPrompt string `json:"syntax_prompt"`
	} `json:"data"`
	Message string `json:"message"`
}

var (
	hunterTestResults hunterResults
	hunterDataResults hunterResults
	hunterdataTotal   int
)

func RunHunter(data string, filename string) {
	searchData := base64.URLEncoding.EncodeToString([]byte(data))
	config := utils.GetConfig()
	hunterUrl := config.Module.Hunter.URL
	hunterKey := config.Module.Hunter.Key

	// 先获取一条数据，获取总数
	hunterReq := hunterUrl + "?api-key=" + hunterKey + "&search=" + searchData + "&page=1&page_size=1"

	response, err := utils.Req().SetSuccessResult(&hunterTestResults).Get(hunterReq)
	if err != nil {
		log.Println("Hunter API request failed")
		log.Println(err)
		return
	}
	if hunterTestResults.Code == 40204 {
		log.Println("Da Niu, your points are used up, try again tomorrow")
		return
	}
	if hunterTestResults.Code == 401 {
		log.Println("Da Niu, Token expired, please check")
		return
	}
	if hunterTestResults.Code == 429 {
		log.Println("Da Niu,your so Fast,wait 10s")
		time.Sleep(time.Second * 10)
		RunHunter(searchData, filename)
		return
	}
	if response.IsSuccessState() {
		hunterdataTotal = hunterTestResults.Data.Total
	}
	log.Printf("共搜索到数据: %d 个", hunterdataTotal)
	// 计算总页数
	if hunterdataTotal > 0 {
		// 根据配置决定获取数量
		maxResults := hunterdataTotal
		if !config.Module.Hunter.All {
			// 如果不获取全部，则限制为配置的size
			if config.Module.Hunter.Size < maxResults {
				maxResults = config.Module.Hunter.Size
				log.Printf("根据配置限制，将只获取 %d 条数据", maxResults)
			}
		}

		// 每页请求数量
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
			var rows [][]string
			time.Sleep(time.Second * 3)

			// 计算当前页应获取的数量
			currentPageSize := pageSize
			if pageIndex == totalPages && lastPageSize > 0 {
				currentPageSize = lastPageSize
			}

			hunterDataReq := hunterUrl + "?api-key=" + hunterKey + "&search=" + searchData + "&page=" + strconv.Itoa(pageIndex) + "&page_size=" + strconv.Itoa(currentPageSize)
			log.Printf("正在请求Hunter API第 %d 页数据", pageIndex)

			hunterResponse, err := utils.Req().SetSuccessResult(&hunterDataResults).Get(hunterDataReq)
			if err != nil {
				log.Println("Hunter API request failed")
				log.Println(err)
			}
			if hunterResponse.IsSuccessState() {
				for i := range hunterDataResults.Data.Arr {
					domain := utils.ToString(hunterDataResults.Data.Arr[i].Domain)
					url := utils.ToString(hunterDataResults.Data.Arr[i].URL)
					webTitle := utils.ToString(hunterDataResults.Data.Arr[i].WebTitle)
					statusCode := hunterDataResults.Data.Arr[i].StatusCode
					ip := utils.ToString(hunterDataResults.Data.Arr[i].IP)
					port := hunterDataResults.Data.Arr[i].Port
					area, country, _ := utils.QueryIp(ip)

					// 使用表格形式输出
					utils.PrintTableRow(domain, url, webTitle, ip, strconv.Itoa(port), country, area)

					// 保存到Excel的数据行，注意顺序要与表头一致
					row := []string{domain, url, webTitle, ip, strconv.Itoa(port), strconv.Itoa(statusCode), country, area}
					rows = append(rows, row)

					processedCount++
					if !config.Module.Hunter.All && processedCount >= config.Module.Hunter.Size {
						// 达到限制数量，提前结束
						log.Printf("已达到配置的最大获取数量: %d", config.Module.Hunter.Size)
						break
					}
				}
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

			// 如果已达到限制数量，提前结束循环
			if !config.Module.Hunter.All && processedCount >= config.Module.Hunter.Size {
				break
			}
		}

		// 打印表格底部分隔线
		utils.PrintTableFooter()
	}
	time.Sleep(time.Second * 3) // 有没有结果都要等待 3s 而不是有结果才延时
}
