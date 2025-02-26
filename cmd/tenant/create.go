// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation and Dapr Contributors.
// Licensed under the MIT License.
// ------------------------------------------------------------

package tenant

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var username string
var password string
var remark string
var TenantCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "create tenant.",
	Example: TenantHelpExample,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			print.PendingStatusEvent(os.Stdout, "tenantTitle not fount ...\n # auth plugins. in Kubernetes mode \n tkeel auth createtenant -k tenantTitle adminName adminPassword")
			return
		}
		title := args[0]
		err := kubernetes.TenantCreate(title, remark, username, password)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}

		print.SuccessStatusEvent(os.Stdout, "Success! ")
	},
}

func init() {
	TenantCreateCmd.Flags().BoolP("help", "h", false, "Print this help message")
	TenantCreateCmd.Flags().StringVarP(&username, "username", "u", "", "username of tenant")
	TenantCreateCmd.Flags().StringVarP(&password, "password", "p", "", "password of tenant")
	TenantCreateCmd.Flags().StringVarP(&remark, "remark", "r", "", "remark of tenant")
	TenantCreateCmd.MarkFlagRequired("username")
	TenantCreateCmd.MarkFlagRequired("password")
	TenantCmd.AddCommand(TenantCreateCmd)
}
