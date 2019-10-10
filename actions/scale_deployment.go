package actions

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
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

	deploymentClient := k8sClient.AppsV1().Deployments(apiv1.NamespaceDefault)

	retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
		// Retrieve the latest version of Deployment before attempting update
		// RetryOnConflict uses exponential backoff to avoid exhausting the apiserver
		result, getErr := deploymentClient.Get("helloworld", metav1.GetOptions{})

		if getErr != nil {
			panic(fmt.Errorf("Failed to get latest version of Deployment: %v", getErr))
		}

		result.Spec.Replicas = int32Ptr(int32(parameters.Pods)) // update replica count
		_, updateErr := deploymentClient.Update(result)
		return updateErr
	})

	if retryErr != nil {
		panic(fmt.Errorf("Update failed: %v", retryErr))
	}

	fmt.Println("Updated deployment...")

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

func int32Ptr(i int32) *int32 { return &i }
