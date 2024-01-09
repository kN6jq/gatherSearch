package cmd

import (
	"flag"
	"fmt"
	"github.com/kN6jq/gatherSearch/module"
	"github.com/kN6jq/gatherSearch/utils"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

type FofaFlags struct {
	Domain string
	Ip     string
	File   string
	Custom bool
}

type HunterFlags struct {
	Domain string
	Ip     string
	File   string
	Custom bool
}

type ShodanFlags struct {
	Domain string
	File   string
}

type ShodandbFlags struct {
	Ip   string
	File string
}

type ZoneFlags struct {
	Name   string
	Domain string
}

var (
	fofa     FofaFlags
	hunter   HunterFlags
	shodan   ShodanFlags
	shodandb ShodandbFlags
	zone     ZoneFlags
)

var rootCmd = &cobra.Command{
	Use:  "learnCobra",
	Long: "learnCobra is a tool to detect",
}

var fofaCmd = &cobra.Command{
	Use:   "fofa",
	Short: "[*] Use fofa module. -h or --help for help.",
	Run: func(cmd *cobra.Command, args []string) {
		if fofa.Domain == "" && fofa.Ip == "" && fofa.File == "" {
			flag.Usage()
			os.Exit(1)
		}
		fofaExcelFile := utils.FofaExcelFile()
		if fofa.Domain != "" {
			log.Println("正在搜索domain: ", fofa.Domain)
			searchData := fmt.Sprintf(`domain="%s"`, fofa.Domain)
			module.FofaSearch(searchData, fofaExcelFile)
		} else if fofa.Ip != "" {
			log.Println("正在搜索ip: ", fofa.Ip)
			searchData := fmt.Sprintf(`ip="%s"`, fofa.Ip)
			module.FofaSearch(searchData, fofaExcelFile)
		} else if fofa.File != "" {
			lines, err := utils.ReadFileLines(fofa.File)
			if err != nil {
				log.Println("read file error:", err)
			}
			if fofa.Custom {
				log.Println("使用自定义搜索语法,请确保文件中每行都是搜索语法")
				for _, line := range lines {
					log.Println("自定义语句查询: ", line)
					module.FofaSearch(line, fofaExcelFile)
				}
			} else {
				// 如果一个文件内都是ip,那么将ip拼接为指定个数进行查询，快
				var allData []string
				for _, line := range lines {
					if utils.IsIP(line) {
						searchData := fmt.Sprintf(`ip="%s"`, line)
						allData = append(allData, searchData)
					}
				}
				config := utils.GetConfig()
				searchData := utils.SplitAndFormatSlice(allData, config.Module.Fofa.BatchSize)
				for _, datum := range searchData {
					log.Println("正在搜索: ", datum)
					module.FofaSearch(datum, fofaExcelFile)
				}

				for _, line := range lines {
					if !utils.IsIP(line) {
						searchData := fmt.Sprintf(`domain="%s"`, line)
						module.FofaSearch(searchData, fofaExcelFile)
					}
				}

			}
		} else {
			_ = cmd.Help()
			return
		}
	},
}

var hunterCmd = &cobra.Command{
	Use:   "hunter",
	Short: "[*] Use hunter module. -h or --help for help.",
	Run: func(cmd *cobra.Command, args []string) {
		if hunter.Domain == "" && hunter.Ip == "" && hunter.File == "" {
			flag.Usage()
			os.Exit(1)
		}
		hunterExcelFile := utils.HunterExcelFile()
		if hunter.Domain != "" {
			log.Println("正在搜索domain: ", hunter.Domain)
			searchData := fmt.Sprintf(`domain.suffix="%s"`, hunter.Domain)
			module.RunHunter(searchData, hunterExcelFile)
		} else if hunter.Ip != "" {
			log.Println("正在搜索ip: ", hunter.Ip)
			searchData := fmt.Sprintf(`ip="%s"`, hunter.Ip)
			module.RunHunter(searchData, hunterExcelFile)
		} else if hunter.File != "" {
			lines, err := utils.ReadFileLines(hunter.File)
			if err != nil {
				log.Println("read file error:", err)
			}
			if hunter.Custom {
				log.Println("使用自定义搜索语法,请确保文件中每行都是搜索语法")
				for _, line := range lines {
					log.Println("正在搜索自定义语法: ", line)
					module.RunHunter(line, hunterExcelFile)
				}
			} else {
				for _, line := range lines {
					if utils.IsIP(line) {
						log.Println("正在搜索ip: ", line)
						searchData := fmt.Sprintf(`ip="%s"`, line)
						module.RunHunter(searchData, hunterExcelFile)
					} else {
						log.Println("正在搜索domain: ", line)
						searchData := fmt.Sprintf(`domain.suffix="%s"`, line)
						module.RunHunter(searchData, hunterExcelFile)
					}
				}
			}

		} else {
			flag.Usage()
		}
	},
}

var shodanCmd = &cobra.Command{
	Use:   "shodan",
	Short: "[*] Use shodan module. -h or --help for help.",
	Run: func(cmd *cobra.Command, args []string) {
		if shodan.Domain == "" && shodan.File == "" {
			flag.Usage()
			os.Exit(1)
		}
		shodanExcelFile := utils.ShodanExcelFile()
		if shodan.Domain != "" {
			log.Println("正在搜索domain: ", shodan.Domain)
			module.RunShodan(shodan.Domain, shodanExcelFile)
		} else if shodan.File != "" {
			lines, err := utils.ReadFileLines(shodan.File)
			if err != nil {
				log.Println("read file error:", err)
			}
			for _, line := range lines {
				log.Println("正在搜索domain: ", line)
				module.RunShodan(line, shodanExcelFile)
			}
		} else {
			_ = cmd.Help()
			return
		}
	},
}

var showandbCmd = &cobra.Command{
	Use:   "shodandb",
	Short: "[*] Use shodandb module. -h or --help for help.",
	Run: func(cmd *cobra.Command, args []string) {
		if shodandb.Ip == "" && shodandb.File == "" {
			flag.Usage()
			os.Exit(1)
		}
		shodandbExcelFile := utils.ShodanPortExcelFile()
		if shodandb.Ip != "" {
			ipList := []string{}
			if strings.Contains(shodandb.Ip, "/") {
				ipList, _ = utils.GenerateIPList(shodandb.Ip)
			} else {
				ipList = append(ipList, shodandb.Ip)
			}
			for _, ips := range ipList {
				log.Println("正在搜索ip: ", ips)
				module.RunShodandb(ips, shodandbExcelFile)
			}

		} else if shodandb.File != "" {
			lines, err := utils.ReadFileLines(shodandb.File)
			if err != nil {
				log.Println("read file error:", err)
			}
			for _, line := range lines {
				if utils.IsIP(line) {
					log.Println("正在搜索ip: ", line)
					module.RunShodandb(line, shodandbExcelFile)
				}
			}
		} else {
			_ = cmd.Help()
			return
		}
	},
}

var zoneCmd = &cobra.Command{
	Use:   "zone",
	Short: "[*] Use zone module. -h or --help for help.",
	Run: func(cmd *cobra.Command, args []string) {
		if zone.Name == "" && zone.Domain == "" {
			flag.Usage()
			os.Exit(1)
		}
		if zone.Name != "" {
			zoneSiteExcelFile := utils.ZoneSiteExcelFile()
			log.Println("正在搜索企业: ", zone.Name)
			module.RunZoneSite(zone.Name, zoneSiteExcelFile)
		} else if zone.Domain != "" {
			zoneDomainExcelFile := utils.ZoneDomainExcelFile()
			log.Println("正在搜索domain: ", zone.Domain)
			module.RunZoneDomain(zone.Domain, zoneDomainExcelFile)
		} else {
			_ = cmd.Help()
			return
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(fofaCmd)
	fofaCmd.Flags().StringVarP(&fofa.Domain, "domain", "d", "", "Domain name")
	fofaCmd.Flags().StringVarP(&fofa.Ip, "ip", "i", "", "IP address")
	fofaCmd.Flags().StringVarP(&fofa.File, "file", "f", "", "File name")
	fofaCmd.Flags().BoolVarP(&fofa.Custom, "custom", "c", false, "Custom query")

	rootCmd.AddCommand(hunterCmd)
	hunterCmd.Flags().StringVarP(&hunter.Domain, "domain", "d", "", "Domain name")
	hunterCmd.Flags().StringVarP(&hunter.Ip, "ip", "i", "", "IP address")
	hunterCmd.Flags().StringVarP(&hunter.File, "file", "f", "", "File name")
	hunterCmd.Flags().BoolVarP(&hunter.Custom, "custom", "c", false, "Custom query")

	rootCmd.AddCommand(shodanCmd)
	shodanCmd.Flags().StringVarP(&shodan.Domain, "domain", "d", "", "Domain name")
	shodanCmd.Flags().StringVarP(&shodan.File, "file", "f", "", "File name")

	rootCmd.AddCommand(showandbCmd)
	showandbCmd.Flags().StringVarP(&shodandb.Ip, "ip", "i", "", "IP address")
	showandbCmd.Flags().StringVarP(&shodandb.File, "file", "f", "", "File name")

	rootCmd.AddCommand(zoneCmd)
	zoneCmd.Flags().StringVarP(&zone.Name, "name", "n", "", "Zone name")
	zoneCmd.Flags().StringVarP(&zone.Domain, "domain", "d", "", "Domain name")

}
