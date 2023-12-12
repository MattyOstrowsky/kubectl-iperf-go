package cmd

import (
	"os"
	"path/filepath"
	"k8s.io/client-go/util/homedir"
	"github.com/spf13/cobra"	
	apiv1 "k8s.io/api/core/v1"
)

var (
	NamespaceFlag    *bool
	Namespace        string = "iperf-workload"
	DefaultNamespace string = apiv1.NamespaceDefault
	Kubeconfig       string
)
var rootCmd = &cobra.Command{
	Use:   "kubectl-iperf",
	Short: "Iperf is a kubectl plugin for performance tool Iperf3",
	Long: `Example:
	kubectl iperf install -n 
	kubectl iperf status`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
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
	NamespaceFlag = rootCmd.PersistentFlags().BoolP("namespace", "n", false, "Create a new namespace \"iperf-workload\". The default is \"default\" namespace.")
	if home := homedir.HomeDir(); home != "" {
		InstallCmd.Flags().StringVarP(&Kubeconfig, "kubeconfig", "k", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")

	} else {
		InstallCmd.PersistentFlags().StringVarP(&Kubeconfig, "kubeconfig", "k", "", "absolute path to the kubeconfig file")

	}
	rootCmd.AddCommand(InstallCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(runCmd)
}
