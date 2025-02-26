/*
Copyright 2021 The tKeel Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

//
//          developer         +        paas manager         +     tantent manager
//        						|                             |
//           +------------+       |       +-----------+         |      +----------+
//           |            |       |       |           |         |      |          |
//           | developing |       |       | published |         |      | disabled |
//           |            |       |       |           |         |      |          |
//           +----+-------+       |       +---+-------+         |      +---+------+
//        		|               |   install |                 |          |
//        		|  ^            |           v   ^             |          | ^
//        		|  |            |               | uninstall   |          | |
//        		|  |            |       +-------+---+         |          | |
//      release |  |            |       |           |         |   enable | |
//        		|  | upgrade    |       | installed |         |          | | disable
//        		|  |            |       |           |         |          | |
//        		|  |            |       +---+-------+         |          | |
//        		|  |            |  register |                 |          | |
//        		v  |            |           v  ^              |          v |
//        		   |            |              | remove       |            |
//           +-------+----+       |       +------+----+         |      +-----+----+
//           |            |       |       |           |         |      |          |
//           |  release   |       |       |registered |         |      | enabled  |
//           |            |       |       |           |         |      |          |
//           +------------+       +       +-----------+         +      +----------+

package plugin

import (
	"fmt"
	"os"

	"github.com/dapr/cli/utils"
	"github.com/gocarina/gocsv"
	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/print"
)

var PluginHelpExample = `
# Get status of tKeel plugins from Kubernetes
tkeel plugin list
tkeel plugin list -r <repo>
tkeel plugin install <repo>/<plugin> <pluginID>
tkeel plugin install <repo>/<plugin>@v0.1.0 <pluginID>
tkeel plugin uninstall <pluginID>
tkeel plugin show <pluginID>
tkeel plugin enable <pluginID> -t <tenantId>
tkeel plugin disable <pluginID> -t <tenantId>
`

var PluginCmd = &cobra.Command{
	Use:     "plugin",
	Short:   "manage plugins.",
	Example: PluginHelpExample,
	Run: func(cmd *cobra.Command, args []string) {
		// Prompt help information If there is no parameter
		if len(args) == 0 {
			cmd.Help()
			return
		}
	},
}

func init() {
	PluginCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "The output format of the list. Valid values are: json, yaml, or table (default)")
	PluginCmd.PersistentFlags().BoolP("help", "h", false, "Print this help message")
}

func outputList(list interface{}, length int) {
	if outputFormat == "json" || outputFormat == "yaml" {
		err := utils.PrintDetail(os.Stdout, outputFormat, list)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}
	} else {
		table, err := gocsv.MarshalString(list)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}

		// Standalone mode displays a separate message when no instances are found.
		if length == 0 {
			fmt.Println("No Dapr instances found.")
			return
		}

		utils.PrintTable(table)
	}
}
