package cmd

import (
	"fmt"
	"kubectl-iperf/utils"

	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
)

var (
	kubeClient *kubernetes.Clientset
)
var InstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install Iperf3 (https://iperf.fr/)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		kubeClient = utils.InitKubernetes(&Kubeconfig)
		fmt.Println("Installing Iperf3...")
		if *NamespaceFlag {
			utils.CreateNamespace(Namespace, kubeClient)
			utils.Deploy(Namespace, kubeClient)
		} else {
			utils.Deploy(DefaultNamespace,kubeClient)
		}
	},
}

