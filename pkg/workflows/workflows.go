package workflows

import (
	"fmt"
	"time"

	"github.com/joshbranham/sreconf-temporal-demo/pkg/activities"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type TagAllSubnetsInput struct {
	Key   string
	Value string
}

func TagAllSubnets(ctx workflow.Context, input *TagAllSubnetsInput) ([]string, error) {
	// Define the activity options, including the retry policy
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second, //amount of time that must elapse before the first retry occurs
			MaximumInterval:    time.Minute, //maximum interval between retries
			BackoffCoefficient: 2,           //how much the retry interval increases
			// MaximumAttempts: 5, // Uncomment this if you want to limit attempts
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var ec2Activities *activities.EC2Activities

	var vpcIds []string
	err := workflow.ExecuteActivity(ctx, ec2Activities.FetchVPCIds).Get(ctx, &vpcIds)
	if err != nil {
		return []string{}, fmt.Errorf("failed to get VPC IDs: %s", err)
	}

	var taggedSubnetIds []string
	for _, vpcId := range vpcIds {
		var subnetIds []string
		err = workflow.ExecuteActivity(ctx, ec2Activities.FetchSubnetIds, vpcId).Get(ctx, &subnetIds)
		if err != nil {
			return taggedSubnetIds, fmt.Errorf("failed to get subnet IDs: %s", err)
		}

		for _, subnetId := range subnetIds {
			err = workflow.ExecuteActivity(ctx, ec2Activities.TagSubnet, subnetId, input.Key, input.Value).Get(ctx, nil)
			if err != nil {
				return taggedSubnetIds, fmt.Errorf("failed to tag subnet: %s", err)
			}

			taggedSubnetIds = append(taggedSubnetIds, subnetId)
		}

	}

	return taggedSubnetIds, nil
}
