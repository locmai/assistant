package actions

import (
	"encoding/json"
	"fmt"
	"google.golang.org/api/dialogflow/v2"
	"log"
)

type createParameters struct {
	Nodes float64 `json:"nodes"`
}

func CreateCluster(q *dialogflow.GoogleCloudDialogflowV2WebhookRequest) (*dialogflow.GoogleCloudDialogflowV2WebhookResponse, error){
	var parameters createParameters

	if err := json.Unmarshal(q.QueryResult.Parameters, &parameters); err != nil {
		return nil, err
	}

	log.Printf("Create Kubernetes, number of node provisioned: %v", parameters.Nodes)

	response := &dialogflow.GoogleCloudDialogflowV2WebhookResponse{
		FulfillmentText: fmt.Sprintf("Creating a Kubernetes cluster with %v nodes. That's what I called \"Kubernetes the easy way.\"" , parameters.Nodes),
	}
	return response, nil
}
