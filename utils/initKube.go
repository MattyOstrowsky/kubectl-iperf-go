package utils
import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

)

func InitKubernetes(kubeconfig *string) *kubernetes.Clientset {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	if err != nil {
		panic(err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	return clientset
}