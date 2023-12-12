package utils

import (
	"bytes"
	"context"
	"fmt"
	"kubectl-iperf/manifests"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

func CreateNamespace(namespace string, kubeClient *kubernetes.Clientset) {
	fmt.Printf("Creating \"%s\" namespace...\n", namespace)
	nsName := &apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: namespace,
		},
	}
	_, err := kubeClient.CoreV1().Namespaces().Create(context.TODO(), nsName, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Namespace", namespace, "created")
	}
}

func Deploy(namespace string, kubeClient *kubernetes.Clientset) {
	fmt.Println("Creating deployments...")
	deploymentsClient := kubeClient.AppsV1().Deployments(namespace)
	serviceClient := kubeClient.CoreV1().Services(namespace)

	serverDeployment := manifests.ServerDeployment(namespace)
	clientDeployment := manifests.ClientDeployment(namespace)
	serviceDeployment := manifests.Service(namespace)

	resultServer, err := deploymentsClient.Create(context.TODO(), serverDeployment, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Created deployment %q.\n", resultServer.GetObjectMeta().GetName())
	}
	resultClient, err := deploymentsClient.Create(context.TODO(), clientDeployment, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Created deployment %q.\n", resultClient.GetObjectMeta().GetName())
	}
	resultService, err := serviceClient.Create(context.TODO(), serviceDeployment, metav1.CreateOptions{})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Created service %q.\n", resultService.GetObjectMeta().GetName())
	}
}
func GetPod(ns string, listOptions metav1.ListOptions, kubeClient *kubernetes.Clientset) (*apiv1.PodList, error) {
	return kubeClient.CoreV1().Pods(ns).List(context.TODO(), listOptions)
}
func DeleteNamespace(namespace string, kubeClient *kubernetes.Clientset) {
	err := kubeClient.CoreV1().Namespaces().Delete(context.TODO(), namespace, metav1.DeleteOptions{})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Namespace", namespace, "deleted.")
	}
}
func ExecuteRemoteCommand(pod *apiv1.Pod, command string) (string, string, error) {
	kubeCfg := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{},
	)
	restCfg, err := kubeCfg.ClientConfig()
	if err != nil {
		return "", "", err
	}
	coreClient, err := kubernetes.NewForConfig(restCfg)
	if err != nil {
		return "", "", err
	}

	buf := &bytes.Buffer{}
	errBuf := &bytes.Buffer{}
	request := coreClient.CoreV1().RESTClient().
		Post().
		Namespace(pod.Namespace).
		Resource("pods").
		Name(pod.Name).
		SubResource("exec").
		VersionedParams(&apiv1.PodExecOptions{
			Container: "iperf",
			Command:   []string{"/bin/sh", "-c", "iperf3 " + command},
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
			TTY:       true,
		}, scheme.ParameterCodec)
	exec, err := remotecommand.NewSPDYExecutor(restCfg, "POST", request.URL())
	if err != nil {
		return "", "", fmt.Errorf("%w Failed executing command %s on %v/%v", err, command, pod.Namespace, pod.Name)
	}
	err = exec.StreamWithContext(context.TODO(), remotecommand.StreamOptions{
		Stdout: buf,
		Stderr: errBuf,
	})
	if err != nil {
		return buf.String(), errBuf.String(), fmt.Errorf("%w Failed executing command %s on %v/%v", err, command, pod.Namespace, pod.Name)
	}

	return buf.String(), errBuf.String(), nil
}
