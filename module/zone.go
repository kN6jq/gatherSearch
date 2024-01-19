package module

import (
	"fmt"
	"github.com/kN6jq/gatherSearch/utils"
	"log"
	"strconv"
	"time"
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
	data = fmt.Sprintf("(company==%s)||(title==%s)||(banner==%s)||(html_banner==%s)||(component==%s)||(ssl_info.detail==%s)", data, data, data, data, data, data)
	response, err := utils.Req().SetHeader("Content-Type", "application/json").
		SetBody(fmt.Sprintf(`{"query":"%s", "query_type":"site", "page":1, "pagesize":1, "zone_key_id": "%s"}`, data, key)).
		SetSuccessResult(&zoneSiteResult).Post(url)
	if err != nil {
		log.Println("0.zone request error:", err)
	}
	var Total = 0
	if response.IsSuccessState() {
		Total, _ = strconv.Atoi(zoneSiteResult.Total)
	} else {
		log.Println("暂未发现信息系统")
	}
	time.Sleep(time.Second * 2)
	if Total > 0 {
		if Total > 10000 {
			Total = utils.Config.Module.Zone.Size
		}
		pageSize := 40 // 每页处理 40 条数据
		totalPages := (Total + pageSize - 1) / pageSize
		for pageIndex := 1; pageIndex <= totalPages; pageIndex++ {
			var rows [][]string
			response, err := utils.Req().SetHeader("Content-Type", "application/json").
				SetBody(fmt.Sprintf(`{"query":"%s", "query_type":"site", "page":%d, "pagesize":%d, "zone_key_id":"%s"}`, data, pageIndex, pageSize, key)).
				SetSuccessResult(&zoneSiteResult).Post(url)
			if err != nil {
				log.Println("0.zone request error:", err)
			}
			if response.IsSuccessState() {
				for _, v := range zoneSiteResult.Data {
					url := v.Url
					title := utils.RemoveSpaces(v.Title)
					status_code := v.Status_code
					ip := v.Ip
					port := v.Port
					rows = append(rows, []string{url, title, status_code, ip, port})
					fmt.Printf("%-20s %-30s %-20s %-20s %-20s\n", v.Url, utils.RemoveSpaces(v.Title), v.Status_code, v.Ip, v.Port)
				}
			}
			if len(rows) > 0 {
				err := utils.WriteDataToExcel(filename, rows)
				if err != nil {
					log.Println("写入excel失败:", err)
				}
				rows = nil
			}
			time.Sleep(time.Second * 2)
		}
	}
	time.Sleep(time.Second * 2)
}

func searchDomain(data string, filename string) {
	config := utils.GetConfig()
	url := config.Module.Zone.URL
	key := config.Module.Zone.Key
	response, err := utils.Req().SetHeader("Content-Type", "application/json").
		SetBody(fmt.Sprintf(`{"query":"%s", "query_type":"domain", "page":1, "pagesize":1, "zone_key_id": "%s"}`, data, key)).
		SetSuccessResult(&zoneDomainResult).Post(url)
	if err != nil {
		log.Println("0.zone request error:", err)
	}
	var Total = 0
	if response.IsSuccessState() {
		Total, _ = strconv.Atoi(zoneDomainResult.Total)
	} else {
		log.Println("暂未发现信息系统")
	}
	time.Sleep(time.Second * 2)
	if Total > 0 {
		if Total > 10000 {
			Total = config.Module.Zone.Size
		}
		pageSize := 40 // 每页处理 40 条数据
		totalPages := (Total + pageSize - 1) / pageSize

		for pageIndex := 1; pageIndex <= totalPages; pageIndex++ {
			var rows [][]string
			response, err := utils.Req().SetHeader("Content-Type", "application/json").
				SetBody(fmt.Sprintf(`{"query":"%s", "query_type":"domain", "page":%d, "pagesize":%d, "zone_key_id":"%s"}`, data, pageIndex, pageSize, key)).
				SetSuccessResult(&zoneDomainResult).Post(url)
			if err != nil {
				log.Println("0.zone request error:", err)
			}
			if response.IsSuccessState() {
				for _, v := range zoneDomainResult.Data {
					url := v.Url
					ip := v.Msg.Ip
					rows = append(rows, []string{url, ip})
					fmt.Printf("%-20s %-30s\n", url, ip)
				}
			}
			if len(rows) > 0 {
				err := utils.WriteDataToExcel(filename, rows)
				if err != nil {
					log.Println("写入excel失败:", err)
				}
				rows = nil
			}
			time.Sleep(time.Second * 2)
		}
	}

	time.Sleep(time.Second * 2)
}
