package cmd

import (
	apiv1 "k8s.io/api/core/v1"

	"context"
	"fmt"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"strings"
	"kubectl-iperf/utils"

)

var command string

var runCmd = &cobra.Command{
	Use:   "run -- <command>",
	Short: "Run Iperf3 command",
	Long: `example:
	kubectl iperf run -- --help
	kubectl iperf run -- -t 90s -b 1K -c iperf -p 5201 `,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}else {
			command = strings.Join(args, " ")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		clientOptions := metav1.ListOptions{
			LabelSelector: "app.kubernetes.io/name=iperf-client",
		}
		kubeClient = utils.InitKubernetes(&Kubeconfig)
		var pod *apiv1.PodList
		var err error
		if *NamespaceFlag {
			pod, err = kubeClient.CoreV1().Pods(Namespace).List(context.TODO(), clientOptions)

		} else {
			pod, err = kubeClient.CoreV1().Pods(DefaultNamespace).List(context.TODO(), clientOptions)
		}
		if err != nil{
			fmt.Println(err)
		}else if len(pod.Items) == 0{
			fmt.Println("Iperf not installed.")
		}
		fmt.Println("Executing command:","iperf3", command)
		stdout, stderr, err := utils.ExecuteRemoteCommand(&pod.Items[0], command)
		fmt.Println("Output:")
		if err != nil {
			fmt.Println(stderr)
			fmt.Println(stdout)
			fmt.Println(err)
		}else{
			fmt.Println(stdout)
		}

	},
}
