package utils

import (
	"github.com/kN6jq/gatherSearch/utils/goqqwry"
	"log"
	"strings"
)

func init() {
	goqqwry.IPData.FilePath = "./qqwry.dat"
}

func QueryIp(ip string) (area string, country string, iscloud bool) {
	res := goqqwry.IPData.InitIPData()
	if v, ok := res.(error); ok {
		log.Panic(v)
	}
	qqWry := goqqwry.NewQQwry()

	if ip == "" {
		log.Fatalln("请输入ip地址")
	}

	rs := map[string]goqqwry.ResultQQwry{}
	rs[ip] = qqWry.Find(ip)
	if strings.Contains(rs[ip].Area, "云") || strings.Contains(rs[ip].Area, "数据中心") || strings.Contains(rs[ip].Area, "cloud") {
		return rs[ip].Area, rs[ip].Country, true
	} else {
		return rs[ip].Area, rs[ip].Country, false
	}
}
