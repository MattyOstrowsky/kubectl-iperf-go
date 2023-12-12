/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kubectl-iperf/utils"

)

// runCmd represents the run command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete Iperf3",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Deleting Iperf3...")
		kubeClient = utils.InitKubernetes(&Kubeconfig)
		if *NamespaceFlag {
			utils.DeleteNamespace(Namespace, kubeClient)
		} else {
			deploymentsClient := kubeClient.AppsV1().Deployments(DefaultNamespace)
			serviceClient := kubeClient.CoreV1().Services(DefaultNamespace)

			if err := deploymentsClient.Delete(context.TODO(), "iperf-client", metav1.DeleteOptions{}); err != nil {
				fmt.Println(err)
				
			}
			
			if err := deploymentsClient.Delete(context.TODO(), "iperf-server", metav1.DeleteOptions{});err != nil {
				fmt.Println(err)
			}

			
			if err := serviceClient.Delete(context.TODO(), "iperf", metav1.DeleteOptions{});err != nil {
				fmt.Println(err)
			}
			
		}
		fmt.Println("Done")
	}}

