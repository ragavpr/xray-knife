package parse

import (
	"bufio"
	"fmt"
	"os"

	"github.com/lilendian0x00/xray-knife/v2/pkg"
	"github.com/lilendian0x00/xray-knife/v2/utils"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	readFromSTDIN   bool
	configLink      string
	configLinksFile string
	showFailed      bool
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

		core := pkg.NewAutomaticCore(true, true)
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

		fmtInfo := color.New(color.FgCyan, color.Bold)
		fmtError := color.New(color.FgRed)
		for i, link := range links {
			if len(links) > 1 {
				fmtInfo.Printf("Config Number: %d\n", i+1)
			}

			p, err := core.CreateProtocol(link)
			if err != nil {
				if showFailed {
					failed_configs[i] = FailInfo{uri: link, err: err.Error()}
				} else {
					fmt.Fprintf(os.Stderr, "Skipped on protocol creation: %v", err)
				}
				continue
			}

			fmt.Println(p.DetailsStr())

			// time.Sleep(time.Duration(100) * time.Millisecond)
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
	ParseCmd.Flags().BoolVarP(&showFailed, "show-failed", "", false, "Show failed URIs and their errors in the end")
}
