package actions

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"encoding/json"
	"log"

	"google.golang.org/api/dialogflow/v2"
)

type scaleParameters struct {
	Pods float64 `json:"pods"`
}

// ScaleDeploymentHandler handler for scale deployment event
func ScaleDeploymentHandler(q *dialogflow.GoogleCloudDialogflowV2WebhookRequest) (*dialogflow.GoogleCloudDialogflowV2WebhookResponse, error) {
	var parameters scaleParameters

	if err := json.Unmarshal(q.QueryResult.Parameters, &parameters); err != nil {
		return nil, err
	}

	log.Printf("Scaling Hello World deployment, number of pods scaled: %v", parameters.Pods)

	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	// Use the current context in kubeconfig to build flag
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	// Create the kubernetes client
	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}

	// Try to list the Pods
	pods, err := k8sClient.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Checking kube-system:  %s pods \n", len(pods.Items))
	for index := range pods.Items {
		fmt.Printf("Pod %v: %s \n", index+1, pods.Items[index].Name)
	}

	response := &dialogflow.GoogleCloudDialogflowV2WebhookResponse{
		FulfillmentText: fmt.Sprintf("Scaling the Hello World deployment to %v pods. It's too easy!", parameters.Pods),
	}

	return response, nil
}

func homeDir() string {
	// For Mac and Linux users
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // For Windows users
}
