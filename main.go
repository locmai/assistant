package main

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/locmai/assistant/actions"
	"github.com/locmai/assistant/server"
)

var (
	addr string
)

func main() {

	// var kubeconfig *string
	// if home := homeDir(); home != "" {
	// 	kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	// } else {
	// 	kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	// }
	// flag.Parse()

	// // Use the current context in kubeconfig to build flag
	// config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Create the kubernetes client
	// k8sClient, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// Try to list the Pods
	// pods, err := k8sClient.CoreV1().Pods("").List(metav1.ListOptions{})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Checking kube-system: %v pods \n", len(pods.Items))
	// for index := range pods.Items {
	// 	fmt.Printf("Pod %v: %s \n", index+1, pods.Items[index].Name)
	// }

	// Load env var from .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	flag.StringVar(&addr, "http", os.Getenv("ADDR"), "HTTP listen address")
	flag.Parse()

	fs := server.NewServer()
	fs.Addr = addr
	fs.DisableBasicAuth = true

	fs.Actions.Set("create_cluster", actions.CreateClusterHandler)
	fs.Actions.Set("scale_deployment", actions.ScaleDeploymentHandler)
	log.Println("Server is started")

	if err := fs.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

func homeDir() string {
	// For Mac and Linux users
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // For Windows users
}
