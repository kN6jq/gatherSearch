package module

import (
	"fmt"
	"gatherSearch/utils"
	"github.com/imroc/req/v3"
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

func RunHunter(search string, filename string) {
	var testresults hunterResults
	var dataresults hunterResults
	var dataTotal int
	config := utils.GetConfig()
	hunterUrl := config.Module.Hunter.URL
	hunterKey := config.Module.Hunter.Key
	start_time := utils.GetLastYearDate()
	end_time := utils.GetNowDate()
	// 先获取一条数据，获取总数
	hunterReq := hunterUrl + "?api-key=" + hunterKey + "&search=" + search + "&page=" + "1" + "&page_size=" + "1" + "&is_web=1&port_filter=true" + "&start_time=" + start_time + "&end_time=" + end_time
	client := req.C().R().
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	response, err := client.SetSuccessResult(&testresults).Get(hunterReq)
	if err != nil {
		log.Println("Hunter API request failed")
		return
	}
	if testresults.Code == 40204 {
		log.Println("Da Niu, your points are used up, try again tomorrow")
		return
	}
	if testresults.Code == 401 {
		log.Println("Da Niu, Token expired, please check")
		return
	}
	if testresults.Code == 429 {
		log.Println("Da Niu,your so Fast,wait 10s")
		time.Sleep(time.Second * 10)
		RunHunter(search, filename)
		return
	}
	if response.IsSuccessState() {
		dataTotal = testresults.Data.Total
	}
	log.Printf("共搜索到数据: %d 个", dataTotal)
	// 计算总页数
	if dataTotal > 0 {
		pageSize := 10 // 每页处理 10 条数据
		if dataTotal > 3000 {
			dataTotal = utils.Config.Module.Hunter.Size
		}
		// 计算总页数
		totalPages := (dataTotal + pageSize - 1) / pageSize

		// 使用循环逐页处理数据
		for pageIndex := 1; pageIndex <= totalPages; pageIndex++ {
			var rows [][]string
			time.Sleep(time.Second * 3)
			hunterDataReq := hunterUrl + "?api-key=" + hunterKey + "&search=" + search + "&page=" + strconv.Itoa(pageIndex) + "&page_size=" + strconv.Itoa(pageSize) + "&is_web=1&port_filter=true" + "&start_time=" + start_time + "&end_time=" + end_time
			hunterResponse, err := client.SetSuccessResult(&dataresults).Get(hunterDataReq)
			if err != nil {
				log.Println("Hunter API request failed")
				return
			}
			if hunterResponse.IsSuccessState() {
				for i := range dataresults.Data.Arr {
					domain := utils.ToString(dataresults.Data.Arr[i].Domain)
					url := utils.ToString(dataresults.Data.Arr[i].URL)
					webTitle := utils.ToString(dataresults.Data.Arr[i].WebTitle)
					statusCode := dataresults.Data.Arr[i].StatusCode
					ip := utils.ToString(dataresults.Data.Arr[i].IP)
					port := dataresults.Data.Arr[i].Port
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
