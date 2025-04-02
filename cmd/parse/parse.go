package parse

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/lilendian0x00/xray-knife/v2/pkg"
	"github.com/lilendian0x00/xray-knife/v2/pkg/xray"
	"github.com/lilendian0x00/xray-knife/v2/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/xtls/xray-core/infra/conf"
)

var (
	readFromSTDIN   bool
	configLink      string
	configLinksFile string
	coreType        string
	// outputType      string
	outputFile string
	showFailed bool
)

// ParseCmd represents the parse command
var ParseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Gives a detailed info about the config link",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 && readFromSTDIN != true && configLink == "" && configLinksFile == "" {
			cmd.Help()
			return
		}

		core := pkg.NewAutomaticCore(true, true, coreType)
		var links []string

		if readFromSTDIN {
			reader := bufio.NewReader(os.Stdin)
			fmt.Println("Enter your config link:")
			text, _ := reader.ReadString('\n')
			links = append(links, text)
		} else if configLink != "" {
			links = append(links, configLink)
		} else if configLinksFile != "" {
			links = utils.ParseFileByNewline(configLinksFile)
			//fmt.Println(links)
		}

		type FailInfo struct {
			uri string
			err string
		}

		failed_configs := make(map[int]FailInfo)

		xray_conf := conf.Config{}

		fmtInfo := color.New(color.FgCyan, color.Bold)
		fmtError := color.New(color.FgRed)
		for i, link := range links {

			p, err := core.CreateProtocol(link)
			if err != nil {
				if showFailed {
					failed_configs[i] = FailInfo{uri: link, err: err.Error()}
				} else {
					fmt.Println(fmtInfo.Sprint(i), " | ", fmtError.Sprint(err), " | ", link)
				}
				continue
			}

			if outputFile == "" {
				if len(links) > 1 {
					fmtInfo.Printf("Config Number: %d\n", i+1)
				}
				fmt.Println(p.DetailsStr())
			} else {
				if coreType == "xray" {
					outbound, _ := p.(xray.Protocol).BuildOutboundDetourConfig(true)
					xray_conf.OutboundConfigs = append(xray_conf.OutboundConfigs, *outbound)
				} else if coreType == "singbox" {
					//TODO: Not implemented yet
					fmtError.Print("Not implemented: singbox outbound")
					os.Exit(0)
				}
			}

			// time.Sleep(time.Duration(100) * time.Millisecond)
		}

		if outputFile != "" {
			jsonBytes, _ := json.MarshalIndent(xray_conf, "", "  ")

			if err := utils.WriteIntoFile(outputFile, []byte(jsonBytes)); err != nil {
				fmt.Errorf("failed to save configs: %v", err)
			}
		}

		if showFailed {
			for i, failed := range failed_configs {
				fmt.Println(fmtInfo.Sprint(i), " | ", fmtError.Sprint(failed.err), " | ", failed.uri)
			}
		}
	},
}

func init() {
	ParseCmd.Flags().BoolVarP(&readFromSTDIN, "stdin", "i", false, "Read config link from the console")
	ParseCmd.Flags().StringVarP(&configLink, "config", "c", "", "The config link")
	ParseCmd.Flags().StringVarP(&configLinksFile, "file", "f", "", "Read config links from a file")
	ParseCmd.Flags().StringVarP(&coreType, "core", "z", "auto", "Core type forced (auto, xray, singbox)")
	// ParseCmd.Flags().StringVarP(&outputType, "type", "x", "outbound-array", "Output type (outbound-array)")
	ParseCmd.Flags().StringVarP(&outputFile, "out", "o", "", "Output file for parsed config links")
	ParseCmd.Flags().BoolVarP(&showFailed, "show-failed", "", false, "Show failed URIs and their errors in the end")
}
