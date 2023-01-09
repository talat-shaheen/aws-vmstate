package main

import (
	"context"
	"errors"
	"time"

	//"flag"

	"fmt"

	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/smithy-go"
)

var client *ec2.Client

// EC2StartInstancesAPI defines the interface for the StartInstances function.
// We use this interface to test the function using a mocked service.
type EC2StartInstancesAPI interface {
	StartInstances(ctx context.Context,
		params *ec2.StartInstancesInput,
		optFns ...func(*ec2.Options)) (*ec2.StartInstancesOutput, error)
}

// StartInstance starts an Amazon Elastic Compute Cloud (Amazon EC2) instance.
// Inputs:
//
//	c is the context of the method call, which includes the AWS Region.
//	api is the interface that defines the method call.
//	input defines the input arguments to the service call.
//
// Output:
//
//	If success, a StartInstancesOutput object containing the result of the service call and nil.
//	Otherwise, nil and an error from the call to StartInstances.
func StartInstances(c context.Context, api EC2StartInstancesAPI, input *ec2.StartInstancesInput) (*ec2.StartInstancesOutput, error) {
	resp, err := api.StartInstances(c, input)

	var apiErr smithy.APIError
	if errors.As(err, &apiErr) && apiErr.ErrorCode() == "DryRunOperation" {
		fmt.Println("User has permission to start an instance.")
		input.DryRun = aws.Bool(false)
		return api.StartInstances(c, input)
	}
	return resp, err
}

func StartInstancesCmd(client EC2StartInstancesAPI, instanceIds []string) {

	fmt.Println(instanceIds)
	input := &ec2.StartInstancesInput{
		InstanceIds: instanceIds,
		DryRun:      aws.Bool(true),
	}
	_, err := StartInstances(context.TODO(), client, input)
	if err != nil {
		fmt.Println("Got an error starting the instance")
		fmt.Println(err)
		//return
	}
	fmt.Println("Started instances with IDs " + strings.Join(instanceIds, ","))
}

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}
	client = ec2.NewFromConfig(cfg)

}
func main() {

	for {

		var instanceIds = make([]string, 0)
		input1 := &ec2.DescribeInstancesInput{
			Filters: []types.Filter{
				{
					Name: aws.String("tag:POC"),
					Values: []string{
						"GolangOperator",
					},
				},
			},
		}

		result, err := client.DescribeInstances(context.TODO(), input1)
		if err != nil {
			fmt.Println("Got an error fetching the status of the instance")
			fmt.Println(err)
		} else {

			fmt.Println(result)
			for _, r := range result.Reservations {
				fmt.Println("Reservation ID: " + *r.ReservationId)
				fmt.Println("Instance IDs:")
				for _, i := range r.Instances {
					fmt.Println("   " + *i.InstanceId)
					//value := *i.InstanceId
					instanceIds = append(instanceIds, *i.InstanceId)
				}

				fmt.Println(instanceIds)
			}

			input := &ec2.DescribeInstanceStatusInput{
				InstanceIds:         instanceIds,
				IncludeAllInstances: aws.Bool(true),
			}
			output, err := client.DescribeInstanceStatus(context.TODO(), input)
			if err != nil {
				fmt.Println("Got an error fetching the status of the instance")
				fmt.Println(err)
			} else {
				fmt.Println(output)
				if len(output.InstanceStatuses) != 1 {
					fmt.Println("The total number of instances did not match the request")
				}
				//////////////////////////////////////////////////////////////////////////////

				for _, instanceStatus := range output.InstanceStatuses {
					fmt.Println("+++++++++++++++++++++++++++++++++++++++++")
					fmt.Println("status check loop\n")
					fmt.Println(*instanceStatus.InstanceId, instanceStatus.InstanceState.Name)
					for key, value := range instanceIds {
						if *instanceStatus.InstanceId == value {
							fmt.Println("instance is found in config file")
							fmt.Printf(" %v : %v \n", key, value)
							if instanceStatus.InstanceState.Name == "running" {
								fmt.Println("instance is running\n")
							} else {
								fmt.Println("instance is not running\n")
								StartInstancesCmd(client, []string{*instanceStatus.InstanceId})
							}
						}
					} //key search ends
				} //instance id check ends
			} //aws-sdk call ends
			fmt.Println("+++++++++++++++++++++++++++++++++++++++++")
		} //instance exists
		time.Sleep(60 * time.Second)

	} //for loop ends
}
