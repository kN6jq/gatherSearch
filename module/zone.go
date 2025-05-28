package module

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kN6jq/gatherSearch/utils"
)

type Icon_md5_base64 struct {
	Base64 string `json:"base64"`
	Md5    string `json:"md5"`
}

type Parse_ip struct {
	Ip   string `json:"ip"`
	Port string `json:"port"`
}

type Data struct {
	_id                  string            `json:"_id"`
	Ip                   string            `json:"ip"`
	Ip_addr              string            `json:"ip_addr"`
	Port                 string            `json:"port"`
	Url                  string            `json:"url"`
	Redirect_url         string            `json:"redirect_url"`
	Title                string            `json:"title"`
	Os                   string            `json:"os"`
	Ping                 interface{}       `json:"ping"`
	Cms                  string            `json:"cms"`
	Banner_os            string            `json:"banner_os"`
	Component            string            `json:"component"`
	Area                 string            `json:"area"`
	City                 string            `json:"city"`
	Continent            string            `json:"continent"`
	Country              string            `json:"country"`
	Device_type          string            `json:"device_type"`
	Lang                 string            `json:"lang"`
	Lang_version         string            `json:"lang_version"`
	Latitude             string            `json:"latitude"`
	Longitude            string            `json:"longitude"`
	Operator             string            `json:"operator"`
	Protection           string            `json:"protection"`
	Protocol             string            `json:"protocol"`
	Province             string            `json:"province"`
	Service              string            `json:"service"`
	Versions             string            `json:"versions"`
	Extra_info           []interface{}     `json:"extra_info"`
	App_name             string            `json:"app_name"`
	App_version          string            `json:"app_version"`
	App_brand            string            `json:"app_brand"`
	Banner               string            `json:"banner"`
	Html_banner          string            `json:"html_banner"`
	Accuracy             int               `json:"accuracy"`
	Group                string            `json:"group"`
	Company              interface{}       `json:"company"`
	Tags                 []interface{}     `json:"tags"`
	Icon_md5_base64      []Icon_md5_base64 `json:"icon_md5_base64"`
	Counterfeit          string            `json:"counterfeit"`
	Status_code          string            `json:"status_code"`
	Risk_score           int               `json:"risk_score"`
	Url_directory        []interface{}     `json:"url_directory"`
	Parse_ip             []Parse_ip        `json:"parse_ip"`
	Toplv_domain         string            `json:"toplv_domain"`
	Server_name          string            `json:"server_name"`
	Server_version       string            `json:"server_version"`
	Server_brand         string            `json:"server_brand"`
	Os_name              string            `json:"os_name"`
	Os_version           string            `json:"os_version"`
	Os_brand             string            `json:"os_brand"`
	Framework_name       string            `json:"framework_name"`
	Framework_version    string            `json:"framework_version"`
	Framework_brand      string            `json:"framework_brand"`
	Js_framework_name    string            `json:"js_framework_name"`
	Js_framework_version string            `json:"js_framework_version"`
	Js_framework_brand   string            `json:"js_framework_brand"`
	Explore_timestamp    string            `json:"explore_timestamp"`
	Beian                string            `json:"beian"`
	Cname                string            `json:"cname"`
	Ssl_hostname         string            `json:"ssl_hostname"`
	Hostname             string            `json:"hostname"`
	Is_cdn               int               `json:"is_cdn"`
	Timestamp            string            `json:"timestamp"`
	Url_analyzer         interface{}       `json:"url_analyzer"`
}

type SiteData struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Sort     string `json:"sort"`
	Page     int    `json:"page"`
	Pagesize int    `json:"pagesize"`
	Total    string `json:"total"`
	Data     []Data `json:"data"`
}

type Msg struct {
	Ip string `json:"ip"`
}

type Datas struct {
	_id          string      `json:"_id"`
	Msg          Msg         `json:"msg"`
	Toplv_domain string      `json:"toplv_domain"`
	Icp          string      `json:"icp"`
	Domain       string      `json:"domain"`
	Company      interface{} `json:"company"`
	Url          string      `json:"url"`
}

type DomainData struct {
	Code     int     `json:"code"`
	Message  string  `json:"message"`
	Sort     string  `json:"sort"`
	Page     int     `json:"page"`
	Pagesize int     `json:"pagesize"`
	Total    string  `json:"total"`
	Data     []Datas `json:"data"`
}

var (
	zoneSiteResult   SiteData
	zoneDomainResult DomainData
)

type zoneResponse struct {
	Code int `json:"code"`
	Data struct {
		Count   int `json:"count"`
		Results []struct {
			URL        string `json:"url"`
			Title      string `json:"title"`
			StatusCode string `json:"status_code"`
			IP         string `json:"ip"`
			Port       string `json:"port"`
		} `json:"results"`
	} `json:"data"`
}

func RunZone(data string, filename string) {
	config := utils.GetConfig()
	zoneUrl := config.Module.Zone.URL
	zoneKey := config.Module.Zone.Key
	zoneSize := config.Module.Zone.Size

	// 构建请求URL，使用name参数
	requestUrl := fmt.Sprintf("%sget?key=%s&name=%s&pagesize=%d", zoneUrl, zoneKey, data, zoneSize)

	// 创建HTTP客户端
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	log.Printf("正在请求Zone API: %s", requestUrl)

	// 发送GET请求
	resp, err := client.Get(requestUrl)
	if err != nil {
		log.Println("Error sending request to 0.zone:", err)
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return
	}

	// 解析JSON响应
	var response zoneResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Println("Error parsing JSON response:", err)
		return
	}

	// 检查响应码
	if response.Code != 0 {
		log.Printf("API returned error code: %d\n", response.Code)
		return
	}

	// 打印结果数量
	log.Printf("共搜索到数据: %d 个", response.Data.Count)

	if response.Data.Count > 0 {
		// 打印表头
		utils.PrintTableHeader()

		var rows [][]string
		for _, result := range response.Data.Results {
			url := result.URL
			title := result.Title
			statusCode := result.StatusCode
			ip := result.IP
			port := result.Port

			// 查询IP地理位置
			area, country, _ := utils.QueryIp(ip)

			// 使用表格输出
			utils.PrintTableRow(extractDomain(url), url, title, ip, port, country, area)

			// 保存到Excel
			row := []string{extractDomain(url), url, title, ip, port, statusCode, country, area}
			rows = append(rows, row)
		}

		if len(rows) > 0 {
			err := utils.WriteDataToExcel(filename, rows)
			if err != nil {
				log.Println("保存文件时发生错误:", err)
				return
			}
		}

		// 打印表格底部分隔线
		utils.PrintTableFooter()
	}
}

// 从URL中提取域名
func extractDomain(url string) string {
	// 移除协议部分
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")

	// 获取域名部分（第一个斜杠之前）
	parts := strings.Split(url, "/")
	if len(parts) > 0 {
		return parts[0]
	}
	return url
}

func RunZoneSite(search string, filename string) {
	searchSite(search, filename)

}

func RunZoneDomain(search string, filename string) {
	searchDomain(search, filename)
}

func searchSite(data string, filename string) {
	config := utils.GetConfig()
	url := config.Module.Zone.URL
	key := config.Module.Zone.Key
	size := config.Module.Zone.Size

	// 构建正确的API请求，使用name参数
	zoneReq := url + "site?key=" + key + "&name=" + data + "&pagesize=" + strconv.Itoa(size)

	response, err := utils.Req().SetSuccessResult(&zoneSiteResult).Get(zoneReq)
	if err != nil {
		log.Println("0.zone request error:", err)
		return
	}
	if response.IsSuccessState() {
		zoneSiteTotal, _ := strconv.Atoi(zoneSiteResult.Total)
		log.Printf("共搜索到数据: %d 个", zoneSiteTotal)
		if zoneSiteTotal > 0 {
			// 打印表头
			utils.PrintTableHeader()

			var rows [][]string
			for i := range zoneSiteResult.Data {
				url := zoneSiteResult.Data[i].Url
				title := zoneSiteResult.Data[i].Title
				statusCode := zoneSiteResult.Data[i].Status_code
				ip := zoneSiteResult.Data[i].Ip
				port := zoneSiteResult.Data[i].Port
				area, country, _ := utils.QueryIp(ip)

				// 使用表格输出
				domain := extractDomain(url)
				utils.PrintTableRow(domain, url, title, ip, port, country, area)

				row := []string{domain, url, title, ip, port, statusCode, country, area}
				rows = append(rows, row)
			}

			if len(rows) > 0 {
				err := utils.WriteDataToExcel(filename, rows)
				if err != nil {
					log.Println("保存文件时发生错误:", err)
					return
				}
			}

			// 打印表格底部分隔线
			utils.PrintTableFooter()
		}
	} else {
		log.Println("0.zone request error:", err)
	}
}

func searchDomain(data string, filename string) {
	config := utils.GetConfig()
	url := config.Module.Zone.URL
	key := config.Module.Zone.Key
	size := config.Module.Zone.Size

	// 构建正确的API请求，使用domain参数
	zoneReq := url + "domain?key=" + key + "&domain=" + data + "&pagesize=" + strconv.Itoa(size)
	log.Printf("正在请求Zone API: %s", zoneReq)

	response, err := utils.Req().SetSuccessResult(&zoneDomainResult).Get(zoneReq)
	if err != nil {
		log.Println("0.zone request error:", err)
		return
	}
	if response.IsSuccessState() {
		zoneDomainTotal, _ := strconv.Atoi(zoneDomainResult.Total)
		log.Printf("共搜索到数据: %d 个", zoneDomainTotal)
		if zoneDomainTotal > 0 {
			// 打印表头
			utils.PrintTableHeader()

			var rows [][]string
			for i := range zoneDomainResult.Data {
				domain := zoneDomainResult.Data[i].Domain
				url := zoneDomainResult.Data[i].Url
				ip := zoneDomainResult.Data[i].Msg.Ip
				area, country, _ := utils.QueryIp(ip)

				// 使用表格输出
				utils.PrintTableRow(domain, url, "", ip, "", country, area)

				row := []string{domain, url, "", ip, "", country, area}
				rows = append(rows, row)
			}

			if len(rows) > 0 {
				err := utils.WriteDataToExcel(filename, rows)
				if err != nil {
					log.Println("保存文件时发生错误:", err)
					return
				}
			}

			// 打印表格底部分隔线
			utils.PrintTableFooter()
		}
	} else {
		log.Println("0.zone request error:", err)
	}
}
