package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
)

func getMemAvailable(s *bufio.Scanner) (int, error) {
	memAvailable := 0
	for s.Scan() {
		if found, _ := fmt.Sscanf(s.Text(), "MemAvailable: %d kB", &memAvailable); found > 0 {
			break
		}
	}
	return memAvailable, s.Err()
}

func getMetadata(metadataClient *imds.Client, path string) (string, error) {
	var err error
	if res, err := metadataClient.GetMetadata(context.Background(), &imds.GetMetadataInput{Path: path}); err == nil {
		if resBytes, err := io.ReadAll(res.Content); err == nil {
			return string(resBytes), err
		}
	}
	return "", err
}

func main() {
	procMeminfo, err := os.Open("/proc/meminfo")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to open /proc/meminfo,", err)
		os.Exit(1)
	}
	meminfoScanner := bufio.NewScanner(procMeminfo)

	metadataClient := imds.New(imds.Options{})
	credentialsProvider := ec2rolecreds.New() // Only use EC2 instance credentials

	region, err := getMetadata(metadataClient, "placement/availability-zone")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get the region from EC2 metadata,", err)
		os.Exit(1)
	}
	region = region[:len(region)-1] // Drop AZ to get the region

	instanceId, err := getMetadata(metadataClient, "instance-id")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to get the instance id from EC2 metadata,", err)
		os.Exit(1)
	}

	fmt.Println("Got metadata, region:", region, "instance-id:", instanceId)

	cloudwatchClient := cloudwatch.New(cloudwatch.Options{
		Credentials: aws.NewCredentialsCache(credentialsProvider),
		Region:      region,
	})

	ticker := time.Tick(time.Minute)
	for timestamp := time.Now(); ; timestamp = <-ticker {
		procMeminfo.Seek(0, 0)
		memAvailable, err := getMemAvailable(meminfoScanner)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Getting RAM info failed:", err)
			os.Exit(1)
		}
		fmt.Println("Sending the metric, value:", memAvailable/1024)
		_, err = cloudwatchClient.PutMetricData(context.Background(), &cloudwatch.PutMetricDataInput{
			Namespace: aws.String("Custom/EC2"),
			MetricData: []types.MetricDatum{
				{
					MetricName: aws.String("MemAvailable"),
					Dimensions: []types.Dimension{{Name: aws.String("InstanceId"), Value: &instanceId}},
					Timestamp:  &timestamp,
					Value:      aws.Float64(float64(memAvailable / 1024)),
					Unit:       types.StandardUnitMegabytes,
				},
			},
		})
		if err != nil {
			fmt.Fprintln(os.Stderr, "Sending the metric failed, ", err)
			os.Exit(1)
		}
	}
}
