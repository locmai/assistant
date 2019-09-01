package actions

import (
	container "cloud.google.com/go/container/apiv1"
	"encoding/json"
	"fmt"
	containerpb "google.golang.org/genproto/googleapis/container/v1"
	"golang.org/x/net/context"
	"google.golang.org/api/dialogflow/v2"
	"log"
)

type createParameters struct {
	Nodes float64 `json:"nodes"`
}

func CreateClusterHandler(q *dialogflow.GoogleCloudDialogflowV2WebhookRequest) (*dialogflow.GoogleCloudDialogflowV2WebhookResponse, error){
	var parameters createParameters

	if err := json.Unmarshal(q.QueryResult.Parameters, &parameters); err != nil {
		return nil, err
	}

	log.Printf("Create Kubernetes, number of node provisioned: %v", parameters.Nodes)

	ctx := context.Background()
	client, err := container.NewClusterManagerClient(ctx)

	if err != nil {
		log.Fatal(err)
	}

	//Some default values, TO-DO: Update to configure the value
	defaultLocation := "asia-southeast1-a"
	projectId := "techcon"

	clusterRequest := containerpb.CreateClusterRequest{
		Cluster:              NewCluster(int32(parameters.Nodes),defaultLocation),
		Parent:               fmt.Sprintf("projects/%s/location/%s",projectId, defaultLocation),
	}

	operation,err := client.CreateCluster(ctx,&clusterRequest)

	if err != nil {
		log.Fatal(err)
	}

	log.Print(operation.Name)
	log.Print(operation.StatusMessage)

	response := &dialogflow.GoogleCloudDialogflowV2WebhookResponse{
		FulfillmentText: fmt.Sprintf("Creating a Kubernetes cluster with %v nodes. That's what I called \"Kubernetes the easy way.\"" , parameters.Nodes),
	}

	return response, nil
}

func NewCluster(nodes int32, location string) *containerpb.Cluster {
	cluster := containerpb.Cluster{
		Name:             "techcon",
		Description:      "Demo Techcon Cluster",
		InitialNodeCount: nodes,
		NodeConfig:       DefaultNodeConfig(),
		Location:         location,
	}
	return &cluster
}

func DefaultNodeConfig() *containerpb.NodeConfig {
	nodeConfig := containerpb.NodeConfig{
		MachineType: "g1-small",
		DiskSizeGb:  30,
		Metadata:    nil,
		Labels:      nil,
		Preemptible: true,
	}
	return &nodeConfig
}
