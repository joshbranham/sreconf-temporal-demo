package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/joshbranham/sreconf-temporal-demo/pkg/activities"
	"github.com/joshbranham/sreconf-temporal-demo/pkg/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

const TaskQueueName = "default"

func main() {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalln("Unable to load AWS config", err)
	}

	// Create the Temporal client
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// Create the Temporal worker
	w := worker.New(c, TaskQueueName, worker.Options{})

	// Setup dependencies and inject our activities
	ec2Client := ec2.NewFromConfig(cfg)
	ec2Activities := activities.NewActivities(ec2Client)
	w.RegisterActivity(ec2Activities)

	// Register our one workflow
	w.RegisterWorkflow(workflows.TagAllSubnets)

	// Start the Worker
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start Temporal worker", err)
	}
}
