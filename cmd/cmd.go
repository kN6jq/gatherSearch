package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/kN6jq/gatherSearch/module"
	"github.com/kN6jq/gatherSearch/utils"

	"github.com/spf13/cobra"
)

var (
	domain  string
	ip      string
	file    string
	custom  bool
	name    string
	rootCmd = &cobra.Command{
		Use:   "gatherSearch",
		Short: "gatherSearch is a tool to search information from various platforms",
		// 禁用自动生成标签
		DisableAutoGenTag: true,
		// 不遍历子命令的标志
		TraverseChildren: false,
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var fofaCmd = &cobra.Command{
	Use:   "fofa",
	Short: "Use fofa module",
	Run:   runFofa,
	// 禁用使用行中的标志
	DisableFlagsInUseLine: true,
	// 不遍历子命令的标志
	TraverseChildren: false,
}

var hunterCmd = &cobra.Command{
	Use:   "hunter",
	Short: "Use hunter module",
	Run:   runHunter,
	// 禁用使用行中的标志
	DisableFlagsInUseLine: true,
	// 不遍历子命令的标志
	TraverseChildren: false,
}

var shodanCmd = &cobra.Command{
	Use:   "shodan",
	Short: "Use shodan module",
	Run:   runShodan,
	// 禁用使用行中的标志
	DisableFlagsInUseLine: true,
	// 不遍历子命令的标志
	TraverseChildren: false,
}

var shodandbCmd = &cobra.Command{
	Use:   "shodandb",
	Short: "Use shodandb module",
	Run:   runShodandb,
	// 禁用使用行中的标志
	DisableFlagsInUseLine: true,
	// 不遍历子命令的标志
	TraverseChildren: false,
}

var zoneCmd = &cobra.Command{
	Use:   "zone",
	Short: "Use zone module",
	Run:   runZone,
	// 禁用使用行中的标志
	DisableFlagsInUseLine: true,
	// 不遍历子命令的标志
	TraverseChildren: false,
}

// 初始化命令和标志
func init() {
	// 禁用默认的补全命令
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// 禁用全局标志显示
	rootCmd.DisableFlagsInUseLine = true

	// 设置根命令不传递标志给子命令
	rootCmd.PersistentFlags().SetInterspersed(false)

	// 配置fofa命令的标志
	fofaCmd.Flags().StringVarP(&domain, "domain", "d", "", "Domain name")
	fofaCmd.Flags().StringVarP(&ip, "ip", "i", "", "IP address")
	fofaCmd.Flags().StringVarP(&file, "file", "f", "", "File name")
	fofaCmd.Flags().BoolVarP(&custom, "custom", "c", false, "Custom query")

	// 配置hunter命令的标志
	hunterCmd.Flags().StringVarP(&domain, "domain", "d", "", "Domain name")
	hunterCmd.Flags().StringVarP(&ip, "ip", "i", "", "IP address")
	hunterCmd.Flags().StringVarP(&file, "file", "f", "", "File name")
	hunterCmd.Flags().BoolVarP(&custom, "custom", "c", false, "Custom query")

	// 配置shodan命令的标志
	shodanCmd.Flags().StringVarP(&domain, "domain", "d", "", "Domain name")
	shodanCmd.Flags().StringVarP(&file, "file", "f", "", "File name")

	// 配置shodandb命令的标志
	shodandbCmd.Flags().StringVarP(&ip, "ip", "i", "", "IP address")
	shodandbCmd.Flags().StringVarP(&file, "file", "f", "", "File name")

	// 配置zone命令的标志
	zoneCmd.Flags().StringVarP(&name, "name", "n", "", "Enterprise name")
	zoneCmd.Flags().StringVarP(&domain, "domain", "d", "", "Domain name")

	// 添加所有子命令到根命令
	rootCmd.AddCommand(fofaCmd, hunterCmd, shodanCmd, shodandbCmd, zoneCmd)
}

func runFofa(cmd *cobra.Command, args []string) {
	fofaExcelFile := utils.FofaExcelFile()
	switch {
	case domain != "":
		searchData := fmt.Sprintf(`domain="%s"`, domain)
		module.RunFofa(searchData, fofaExcelFile)
	case ip != "":
		searchData := fmt.Sprintf(`ip="%s"`, ip)
		module.RunFofa(searchData, fofaExcelFile)
	case file != "":
		runFofaWithFile(file, custom, fofaExcelFile)
	default:
		_ = cmd.Help()
	}
}

func runFofaWithFile(filename string, custom bool, excelFile string) {
	lines := readFileLines(filename)
	if custom {
		log.Println("使用自定义搜索语法,请确保文件中每行都是搜索语法")
		for _, line := range lines {
			log.Println("自定义语句查询: ", line)
			module.RunFofa(line, excelFile)
		}
	} else {
		var ipData, domainData []string
		for _, line := range lines {
			if utils.IsIP(line) {
				ipData = append(ipData, fmt.Sprintf(`ip="%s"`, line))
			} else {
				domainData = append(domainData, fmt.Sprintf(`domain="%s"`, line))
			}
		}
		config := utils.GetConfig()
		batchSize := config.Module.Fofa.BatchSize
		runBatchSearch(ipData, batchSize, module.RunFofa, excelFile)
		runBatchSearch(domainData, 1, module.RunFofa, excelFile)
	}
}

func runHunter(cmd *cobra.Command, args []string) {
	hunterExcelFile := utils.HunterExcelFile()
	switch {
	case domain != "":
		searchData := fmt.Sprintf(`domain.suffix="%s"`, domain)
		module.RunHunter(searchData, hunterExcelFile)
	case ip != "":
		searchData := fmt.Sprintf(`ip="%s"`, ip)
		module.RunHunter(searchData, hunterExcelFile)
	case file != "":
		runHunterWithFile(file, custom, hunterExcelFile)
	default:
		_ = cmd.Help()
	}
}

func runHunterWithFile(filename string, custom bool, excelFile string) {
	lines := readFileLines(filename)
	if custom {
		log.Println("使用自定义搜索语法,请确保文件中每行都是搜索语法")
		for _, line := range lines {
			log.Println("正在搜索自定义语法: ", line)
			module.RunHunter(line, excelFile)
		}
	} else {
		var ipData, domainData []string
		for _, line := range lines {
			if utils.IsIP(line) {
				ipData = append(ipData, fmt.Sprintf(`ip="%s"`, line))
			} else {
				domainData = append(domainData, fmt.Sprintf(`domain.suffix="%s"`, line))
			}
		}
		config := utils.GetConfig()
		batchSize := config.Module.Hunter.BatchSize
		runBatchSearch(ipData, batchSize, module.RunHunter, excelFile)
		runBatchSearch(domainData, 1, module.RunHunter, excelFile)
	}
}

func runShodan(cmd *cobra.Command, args []string) {
	shodanExcelFile := utils.ShodanExcelFile()
	switch {
	case domain != "":
		module.RunShodan(domain, shodanExcelFile)
	case file != "":
		lines := readFileLines(file)
		for _, line := range lines {
			module.RunShodan(line, shodanExcelFile)
		}
	default:
		_ = cmd.Help()
	}
}

func runShodandb(cmd *cobra.Command, args []string) {
	shodandbExcelFile := utils.ShodanPortExcelFile()
	switch {
	case ip != "":
		var ipList []string
		if strings.Contains(ip, "/") {
			ipList, _ = utils.GenerateIPList(ip)
		} else {
			ipList = append(ipList, ip)
		}
		for _, ips := range ipList {
			module.RunShodandb(ips, shodandbExcelFile)
		}
	case file != "":
		lines := readFileLines(file)
		for _, line := range lines {
			if utils.IsIP(line) {
				module.RunShodandb(line, shodandbExcelFile)
			}
		}
	default:
		_ = cmd.Help()
	}
}

func runZone(cmd *cobra.Command, args []string) {
	switch {
	case name != "":
		zoneSiteExcelFile := utils.ZoneSiteExcelFile()
		module.RunZoneSite(name, zoneSiteExcelFile)
	case domain != "":
		zoneDomainExcelFile := utils.ZoneDomainExcelFile()
		module.RunZoneDomain(domain, zoneDomainExcelFile)
	default:
		_ = cmd.Help()
	}
}

func readFileLines(filename string) []string {
	lines, err := utils.ReadFileLines(filename)
	if err != nil {
		log.Println("read file error:", err)
		return nil
	}
	return lines
}

func runBatchSearch(data []string, batchSize int, runner func(string, string), excelFile string) {
	batches := utils.SplitAndFormatSlice(data, batchSize)
	for _, batch := range batches {
		log.Println("正在搜索: ", batch)
		runner(batch, excelFile)
	}
}
