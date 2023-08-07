package module

import (
	"encoding/base64"
	"fmt"
	"github.com/kN6jq/gatherSearch/utils"
	"log"
	"strconv"
	"time"
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
	start_time := utils.GetLastYearDate()
	end_time := utils.GetNowDate()
	// 先获取一条数据，获取总数
	hunterReq := hunterUrl + "?api-key=" + hunterKey + "&search=" + searchData + "&page=" + "1" + "&page_size=" + "1" + "&is_web=1&port_filter=true" + "&start_time=" + start_time + "&end_time=" + end_time
	response, err := utils.Req().SetSuccessResult(&hunterTestResults).Get(hunterReq)
	if err != nil {
		log.Println("Hunter API request failed")
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
		pageSize := 100 // 每页处理 10 条数据
		if hunterdataTotal > 3000 {
			hunterdataTotal = utils.Config.Module.Hunter.Size
		}
		// 计算总页数
		totalPages := (hunterdataTotal + pageSize - 1) / pageSize

		// 使用循环逐页处理数据
		for pageIndex := 1; pageIndex <= totalPages; pageIndex++ {
			var rows [][]string
			time.Sleep(time.Second * 3)
			hunterDataReq := hunterUrl + "?api-key=" + hunterKey + "&search=" + searchData + "&page=" + strconv.Itoa(pageIndex) + "&page_size=" + strconv.Itoa(pageSize) + "&is_web=1&port_filter=true" + "&start_time=" + start_time + "&end_time=" + end_time
			hunterResponse, err := utils.Req().SetSuccessResult(&hunterDataResults).Get(hunterDataReq)
			if err != nil {
				log.Println("Hunter API request failed")
				return
			}
			if hunterResponse.IsSuccessState() {
				for i := range hunterDataResults.Data.Arr {
					domain := utils.ToString(hunterDataResults.Data.Arr[i].Domain)
					url := utils.ToString(hunterDataResults.Data.Arr[i].URL)
					webTitle := utils.ToString(hunterDataResults.Data.Arr[i].WebTitle)
					statusCode := hunterDataResults.Data.Arr[i].StatusCode
					ip := utils.ToString(hunterDataResults.Data.Arr[i].IP)
					port := hunterDataResults.Data.Arr[i].Port
					fmt.Printf("%-20s %-30s %-40s %-20d %-20s %-20d\n", domain, url, webTitle, statusCode, ip, port)
					row := []string{domain, url, webTitle, strconv.Itoa(statusCode), ip, strconv.Itoa(port)}
					rows = append(rows, row)
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
		}

	}
	time.Sleep(time.Second * 3) // 有没有结果都要等待 3s 而不是有结果才延时
}
