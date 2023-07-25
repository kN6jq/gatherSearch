package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"gatherSearch/module"
	"gatherSearch/utils"
	"log"
	"strings"
)

var (
	domain    string // 域名
	ip        string // ip
	file      string // 文件
	name      string
	platform  string
	customize bool
)

func init() {
	flag.StringVar(&domain, "d", "", "domain")
	flag.StringVar(&ip, "i", "", "ip")
	flag.StringVar(&file, "f", "", "file")
	flag.StringVar(&name, "n", "", "0zone 对应的企业名字")
	flag.StringVar(&platform, "p", "", "platform ( fofo | hunter | shodan | shodandb | 0zone  )")
	flag.BoolVar(&customize, "c", false, "customize")
	flag.Parse()
}

func main() {
	config := utils.GetConfig()

	if platform == "fofa" {
		// 先判断是否做了fofa的参数配置
		if config.Module.Fofa.Email == "" || config.Module.Fofa.Token == "" || config.Module.Fofa.URL == "" {
			fmt.Println("fofa config error")
			return
		}
		filename := utils.FofaExcelFile()
		if domain != "" {
			log.Println("正在搜索domain: ", domain)
			searchDomain := fmt.Sprintf(`domain="%s"`, domain)
			search := base64.URLEncoding.EncodeToString([]byte(searchDomain))
			module.RunFofa(search, filename)
		} else if ip != "" {
			log.Println("正在搜索ip: ", ip)
			searchIp := fmt.Sprintf(`ip="%s"`, ip)
			search := base64.URLEncoding.EncodeToString([]byte(searchIp))
			module.RunFofa(search, filename)
		} else if utils.FileExists(file) {
			lines, err := utils.ReadFileLines(file)
			if err != nil {
				log.Println("read file error:", err)
			}
			if customize {
				for _, line := range lines {
					search := base64.URLEncoding.EncodeToString([]byte(line))
					module.RunFofa(search, filename)
				}
			} else {
				for _, line := range lines {
					if utils.IsIP(line) {
						log.Println("正在搜索ip: ", line)
						searchIp := fmt.Sprintf(`ip="%s"`, line)
						search := base64.URLEncoding.EncodeToString([]byte(searchIp))
						module.RunFofa(search, filename)
					} else {
						log.Println("正在搜索domain: ", line)
						searchDomain := fmt.Sprintf(`domain="%s"`, line)
						search := base64.URLEncoding.EncodeToString([]byte(searchDomain))
						module.RunFofa(search, filename)
					}
				}
			}

		} else {
			flag.Usage()
		}
	} else if platform == "hunter" {
		// 先判断是否做了hunter的参数配置
		if config.Module.Hunter.Key == "" || config.Module.Hunter.URL == "" {
			log.Println("hunter配置文件错误, 请检查配置文件")
			return
		}
		if domain != "" {
			log.Println("正在搜索domain: ", domain)
			hunterExcelFile := utils.HunterExcelFile()
			searchDomain := fmt.Sprintf(`domain.suffix="%s"`, domain)
			search := base64.URLEncoding.EncodeToString([]byte(searchDomain))
			module.RunHunter(search, hunterExcelFile)
		} else if ip != "" {
			log.Println("正在搜索ip: ", ip)
			hunterExcelFile := utils.HunterExcelFile()
			searchIp := fmt.Sprintf(`ip="%s"`, ip)
			search := base64.URLEncoding.EncodeToString([]byte(searchIp))
			module.RunHunter(search, hunterExcelFile)
		} else if utils.FileExists(file) {
			lines, err := utils.ReadFileLines(file)
			if err != nil {
				log.Println("read file error:", err)
			}
			hunterExcelFile := utils.HunterExcelFile()
			if customize {
				log.Println("使用自定义搜索语法,请确保文件中每行都是搜索语法")
				for _, line := range lines {
					log.Println("正在搜索自定义语法: ", line)
					search := base64.URLEncoding.EncodeToString([]byte(line))
					module.RunHunter(search, hunterExcelFile)
				}
			} else {
				for _, line := range lines {
					if utils.IsIP(line) {
						log.Println("正在搜索ip: ", line)
						searchIp := fmt.Sprintf(`ip="%s"`, line)
						search := base64.URLEncoding.EncodeToString([]byte(searchIp))
						module.RunHunter(search, hunterExcelFile)
					} else {
						log.Println("正在搜索domain: ", line)
						searchDomain := fmt.Sprintf(`domain.suffix="%s"`, line)
						search := base64.URLEncoding.EncodeToString([]byte(searchDomain))
						module.RunHunter(search, hunterExcelFile)
					}
				}
			}

		} else {
			flag.Usage()
		}
	} else if platform == "shodan" {
		// 先判断是否做了shodan的参数配置
		if config.Module.Shodan.Key == "" || config.Module.Shodan.URL == "" {
			log.Println("shodan配置文件错误, 请检查配置文件")
			return
		}
		shodanExcelFile := utils.ShodanExcelFile()
		if domain != "" {
			log.Println("正在搜索domain: ", domain)
			module.RunShodan(domain, shodanExcelFile)
		} else {
			flag.Usage()
		}
	} else if platform == "shodandb" {
		shodandbExcelFile := utils.ShodanPortExcelFile()
		if ip != "" {
			ipList := []string{}
			if strings.Contains(ip, "/") {
				ipList, _ = utils.GenerateIPList(ip)

			} else {
				ipList = append(ipList, ip)
			}
			for _, ips := range ipList {
				log.Println("正在搜索ip: ", ips)
				module.RunShodandb(ips, shodandbExcelFile)
			}

		} else if file != "" {
			lines, err := utils.ReadFileLines(file)
			if err != nil {
				log.Println("read file error:", err)
			}
			for _, line := range lines {
				log.Println("正在搜索ip: ", line)
				module.RunShodandb(line, shodandbExcelFile)
			}
		} else {
			log.Println("shodandb 搜索参数错误, 请检查参数")
		}
	} else if strings.Contains(platform, "0zone") {
		if config.Module.Zone.URL == "" || config.Module.Zone.Key == "" {
			log.Println("0zone配置文件错误, 请检查配置文件")
			return
		}
		zoneSiteExcelFile := utils.ZoneSiteExcelFile()
		zoneDomainExcelFile := utils.ZoneDomainExcelFile()
		if name != "" {
			log.Println("正在搜索企业: ", name)
			module.RunZoneSite(name, zoneSiteExcelFile)
		} else if domain != "" {
			log.Println("正在搜索domain: ", domain)
			module.RunZoneDomain(domain, zoneDomainExcelFile)
		}
	}

}
