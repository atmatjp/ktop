package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	/*
	"os/signal"
	"time"
	*/

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	metricsClient "k8s.io/metrics/pkg/client/clientset/versioned"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	kubeconfig := filepath.Join(home, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	// metrics-server 用クライアント
	metricsClientset, err := metricsClient.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// Kubernetes API 用クライアント
	coreClientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}


	nodeMetricsList, err := metricsClientset.
		MetricsV1beta1().
		NodeMetricses().
		List(context.Background(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf(
		"%-5s %5s %5s\n",
		"NODE",
		"CPU(%)",
		"RAM(%)",
	)

	for _, nodeMetric := range nodeMetricsList.Items {
		node, err := coreClientset.
			CoreV1().
			Nodes().
			Get(context.Background(), nodeMetric.Name, metav1.GetOptions{})
		if err != nil {
			fmt.Printf("failed to get node %s: %v\n", nodeMetric.Name, err)
			continue
		}

		// CPU
		usedCPU := nodeMetric.Usage.Cpu().MilliValue()
		totalCPU := node.Status.Capacity.Cpu().MilliValue()

		cpuPercent := 0.0
		if totalCPU > 0 {
			cpuPercent = float64(usedCPU) / float64(totalCPU) * 100
		}

		// Memory
		usedMem := nodeMetric.Usage.Memory().Value()
		totalMem := node.Status.Capacity.Memory().Value()

		memPercent := 0.0
		if totalMem > 0 {
			memPercent = float64(usedMem) / float64(totalMem) * 100
		}


		fmt.Printf(
			"%-5s %5.1f%% %5.1f%%\n",
			nodeMetric.Name,
			cpuPercent,
			memPercent,
		)

	}
}