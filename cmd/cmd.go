package cmd

import (
	"fmt"
	"github.com/kN6jq/gatherSearch/module"
	"github.com/kN6jq/gatherSearch/utils"
	"log"
	"os"
	"strings"

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
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVarP(&domain, "domain", "d", "", "Domain name")
	rootCmd.PersistentFlags().StringVarP(&ip, "ip", "i", "", "IP address")
	rootCmd.PersistentFlags().StringVarP(&file, "file", "f", "", "File name")
	rootCmd.PersistentFlags().BoolVarP(&custom, "custom", "c", false, "Custom query")
}

var fofaCmd = &cobra.Command{
	Use:   "fofa",
	Short: "Use fofa module",
	//Args:  cobra.ExactArgs(1),
	Run: runFofa,
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

var hunterCmd = &cobra.Command{
	Use:   "hunter",
	Short: "Use hunter module",
	//Args:  cobra.ExactArgs(1),
	Run: runHunter,
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

var shodanCmd = &cobra.Command{
	Use:   "shodan",
	Short: "Use shodan module",
	//Args:  cobra.ExactArgs(1),
	Run: runShodan,
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

var shodandbCmd = &cobra.Command{
	Use:   "shodandb",
	Short: "Use shodandb module",
	//Args:  cobra.ExactArgs(1),
	Run: runShodandb,
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

var zoneCmd = &cobra.Command{
	Use:   "zone",
	Short: "Use zone module",
	Run:   runZone,
}

func init() {
	zoneCmd.Flags().StringVarP(&name, "name", "n", "", "Enterprise name")
	rootCmd.AddCommand(zoneCmd)
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

func init() {
	rootCmd.AddCommand(fofaCmd, hunterCmd, shodanCmd, shodandbCmd)
}
