package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"os"
)

func main() {
	allRegions, err := GetRegions()
	if err != nil {
		panic(err)
	}

	for _, r := range allRegions {
		allVpcs, err := GetVpcs(r)
		if err != nil {
			panic(err)
		}
		for _, v := range allVpcs {
			instanceCount, err := countInstancesInVpc(v, r)
			if err != nil {
				panic(err)
			}
			fmt.Println("\nRegion: ", r)
			fmt.Println("VPC: ", v)
			fmt.Println("Instance count: ", instanceCount)
		}
	}
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func createEc2Client(region string) *ec2.EC2 {
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		exitErrorf("failed to load config, %v", err)
	}
	cfg.Region = region

	ec2Svc := ec2.New(cfg)
	return ec2Svc
}

func countInstancesInVpc(vpc string, region string) (int, error) {

	instanceCount := 0
	ec2Client := createEc2Client(region)
	params := &ec2.DescribeInstancesInput{
		Filters: []ec2.Filter{
			{
				Name:   aws.String("vpc-id"),
				Values: []string{vpc},
			},
		},
	}

	request := ec2Client.DescribeInstancesRequest(params)
	response, err := request.Send()
	if err != nil {
		panic(err)
	}
	for _, r := range response.Reservations {
		for _, _ = range r.Instances {
			instanceCount += 1
		}
	}
	return instanceCount, nil
}

func GetRegions() ([]string, error) {

	regions := make([]string, 0)
	ec2Client := createEc2Client("ap-southeast-1")

	params := &ec2.DescribeRegionsInput{}

	request := ec2Client.DescribeRegionsRequest(params)
	response, err := request.Send()
	if err != nil {
		panic(err)
	}

	for _, region := range response.Regions {
		regions = append(regions, *region.RegionName)
	}
	return regions, nil
}

func GetVpcs(region string) ([]string, error) {

	vpcs := make([]string, 0)
	ec2Client := createEc2Client(region)

	params := &ec2.DescribeVpcsInput{}

	request := ec2Client.DescribeVpcsRequest(params)
	response, err := request.Send()
	if err != nil {
		panic(err)
	}

	for _, i := range response.Vpcs {
		vpcs = append(vpcs, *i.VpcId)
	}
	return vpcs, nil
}
