/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"context"
	"github.com/spf13/cobra"
	"kubectl-iperf/utils"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Display Iperf3 status",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		kubeClient = utils.InitKubernetes(&Kubeconfig)
		if *NamespaceFlag {
			fmt.Println("Namespace:", Namespace)
			displayStatus(Namespace, kubeClient)
		} else {
			fmt.Println("Namespace:", DefaultNamespace)
			displayStatus(DefaultNamespace, kubeClient)
		}
	}}

func displayStatus(ns string, kubeClient *kubernetes.Clientset) {
	serverPod, errServer := utils.GetPod(ns, metav1.ListOptions{LabelSelector: "app.kubernetes.io/name=iperf-server"}, kubeClient)
	clientPod, errClient := utils.GetPod(ns, metav1.ListOptions{LabelSelector: "app.kubernetes.io/name=iperf-client"}, kubeClient)
	service, errService := kubeClient.CoreV1().Services(ns).List(context.TODO(), metav1.ListOptions{LabelSelector: "app.kubernetes.io/instance=iperf"})
	if len(serverPod.Items) == 1 && len(clientPod.Items) == 1 && len(service.Items) == 1 {
		fmt.Println("Iperf status: Installed")
	}else{
		fmt.Println("Iperf status: Not installed")

	}
	if errServer != nil {
		fmt.Println(errServer)
	} else if len(serverPod.Items) == 0 {
		fmt.Printf("Server Not found\n")
	} else {
		displayDeploymentStatus("Server", serverPod.Items[0])
	}
	if errClient != nil {
		fmt.Println(errClient)
	} else if len(clientPod.Items) == 0 {
		fmt.Printf("Client Not found\n")
	} else {
		displayDeploymentStatus("Client", clientPod.Items[0])
	}
	if errService != nil {
		fmt.Println(errService)
	} else if len(service.Items) == 0 {
		fmt.Printf("Service Not found\n")
	} else {
		displaySeviceStatus(service.Items[0])
	}

}

func displayDeploymentStatus(label string, pod v1.Pod) {
	fmt.Printf("%s:\n", label)
	fmt.Println("	Name:", pod.Name)
	fmt.Println("	Status:", pod.Status.Phase)
	fmt.Println("	Namespace:", pod.Namespace)
	fmt.Println("	Node:", pod.Spec.NodeName)

}
func displaySeviceStatus(svc v1.Service) {
	fmt.Println("Service:")
	fmt.Println("	Name:", svc.Name)
	fmt.Println("	Namespace:", svc.Namespace)
	fmt.Println("	ClusterIP:", svc.Spec.ClusterIP)
	fmt.Println("	Ports:")
	for _, port := range svc.Spec.Ports {
		fmt.Println("		Name:", port.Name)
		fmt.Println("		Protocol:", port.Protocol)
		fmt.Println("		Port:", port.Port)
	}
}
