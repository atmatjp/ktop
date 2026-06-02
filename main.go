package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := filepath.Join(os.Getenv("HOME"), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatalf("kubeconfig読み込みエラー: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("clientset作成エラー: %v", err)
	}

	/*
	metricsClient, err := metricsv.NewForConfig(config)
	if err != nil {
		log.Fatalf("Metricsクライアント作成エラー: %v", err)
	}
	*/

	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Pod取得エラー: %v", err)
	}


	for i, pod := range pods.Items {
		fmt.Printf("[%d] Namespace: %-15s | Pod名: %s\n", i+1, pod.Namespace, pod.Name)
	}
}
