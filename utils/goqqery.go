package utils

import (
	"log"
	"os"
	"strings"

	"github.com/kN6jq/gatherSearch/utils/goqqwry"
)

func init() {
	goqqwry.IPData.FilePath = "./qqwry.dat"
}

func QueryIp(ip string) (area string, country string, iscloud bool) {
	// 检查IP是否为空
	if ip == "" {
		log.Println("IP地址为空，跳过地理位置查询")
		return "", "", false
	}

	// 检查文件是否存在
	if _, err := os.Stat("./qqwry.dat"); os.IsNotExist(err) {
		log.Println("qqwry.dat文件不存在，跳过地理位置查询")
		return "", "", false
	}

	// 使用defer-recover防止程序崩溃
	defer func() {
		if r := recover(); r != nil {
			log.Printf("IP地理位置查询发生错误: %v", r)
			area = ""
			country = ""
			iscloud = false
		}
	}()

	// 初始化IP数据库
	res := goqqwry.IPData.InitIPData()
	if _, ok := res.(error); ok {
		log.Println("初始化IP数据库失败，跳过地理位置查询")
		return "", "", false
	}

	// 创建查询实例
	qqWry := goqqwry.NewQQwry()

	// 查询IP
	result := qqWry.Find(ip)

	// 检查结果
	if result.Country == "" && result.Area == "" {
		log.Printf("未找到IP %s 的地理位置信息", ip)
		return "", "", false
	}

	// 判断是否为云服务
	isCloud := strings.Contains(result.Area, "云") ||
		strings.Contains(result.Area, "数据中心") ||
		strings.Contains(result.Area, "cloud")

	return result.Area, result.Country, isCloud
}
