package activities

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

type EC2Activities struct {
	ec2Client *ec2.Client
}

func NewActivities(ec2Client *ec2.Client) *EC2Activities {
	return &EC2Activities{
		ec2Client: ec2Client,
	}
}

func (a *EC2Activities) FetchVPCIds(ctx context.Context) ([]string, error) {
	vpcs, err := a.ec2Client.DescribeVpcs(ctx, &ec2.DescribeVpcsInput{})
	if err != nil {
		return nil, err
	}

	var vpcIds []string
	for _, vpc := range vpcs.Vpcs {
		vpcIds = append(vpcIds, *vpc.VpcId)
	}

	return vpcIds, nil
}

func (a *EC2Activities) FetchSubnetIds(ctx context.Context, vpcId string) ([]string, error) {
	subnets, err := a.ec2Client.DescribeSubnets(ctx, &ec2.DescribeSubnetsInput{
		Filters: []types.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vpcId},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	var subnetIds []string
	for _, subnet := range subnets.Subnets {
		subnetIds = append(subnetIds, *subnet.SubnetId)
	}

	return subnetIds, nil
}

func (a *EC2Activities) TagSubnet(ctx context.Context, subnetId string, key string, value string) error {
	_, err := a.ec2Client.CreateTags(ctx, &ec2.CreateTagsInput{
		Resources: []string{subnetId},
		Tags: []types.Tag{
			{
				Key:   aws.String(key),
				Value: aws.String(value),
			},
		},
	})
	return err
}
