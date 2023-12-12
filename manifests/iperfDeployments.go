package manifests

import (
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ServerDeployment(namespace string) *appsv1.Deployment {

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "iperf-server",
			Namespace: namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name": "iperf-server", "app.kubernetes.io/instance": "iperf",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/name": "iperf-server", "app.kubernetes.io/instance": "iperf",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/name": "iperf-server", "app.kubernetes.io/instance": "iperf",
					},
					Annotations: map[string]string{
						"sidecar.istio.io/inject": "false",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:            "iperf",
							ImagePullPolicy: "Always",
							Image:           "clearlinux/iperf:3",
							Command:         []string{"iperf3"},
							Args:            []string{"-s", "-p", "5201", "--json", "--logfile", "/dev/null"},
							Ports: []apiv1.ContainerPort{
								{
									Name:          "tcp",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 5201,
								}, {
									Name:          "udp",
									Protocol:      apiv1.ProtocolUDP,
									ContainerPort: 5201,
								},
							},
						},
					},
				},
			},
		},
	}
}

func ClientDeployment(namespace string) *appsv1.Deployment {

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "iperf-client",
			Namespace: namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name": "iperf-client", "app.kubernetes.io/instance": "iperf",
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/name": "iperf-client", "app.kubernetes.io/instance": "iperf",
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/name": "iperf-client", "app.kubernetes.io/instance": "iperf",
					},
					Annotations: map[string]string{
						"sidecar.istio.io/inject": "false",
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:            "iperf",
							ImagePullPolicy: "Always",
							Image:           "clearlinux/iperf:3",
							Command:         []string{"sleep"},
							Args:            []string{"infinity"},
						},
					},
				},
			},
		},
	}
}

func Service(namespace string) *apiv1.Service {

	return &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "iperf",
			Namespace: namespace,
			Labels: map[string]string{
				"app.kubernetes.io/instance": "iperf",
			},
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Name:     "tcp",
					Protocol: apiv1.ProtocolTCP,
					Port:     5201,
				}, {
					Name:     "udp",
					Protocol: apiv1.ProtocolUDP,
					Port:     5201,
				},
			},
			Selector: map[string]string{
				"app.kubernetes.io/instance": "iperf", "app.kubernetes.io/name": "iperf-server",
			},
			ClusterIP: "",
		},
	}
}
func int32Ptr(i int32) *int32 { return &i }
