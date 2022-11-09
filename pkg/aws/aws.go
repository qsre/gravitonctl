package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var REGION = "eu-west-3"
var sess *session.Session

var ec2svc *ec2.EC2

func init() {
	setSessions()
}

func setSessions() {
	sess = session.Must(session.NewSession(&aws.Config{
		Region: aws.String(REGION),
	}))

	ec2svc = ec2.New(sess)
}

func ReInitWithRegion(region string) {
	REGION = region
	setSessions()
}

func describeAllInstances() (instances []*ec2.Instance, err error) {
	describeInstancesInput := ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:created-by"),
				Values: []*string{
					aws.String("gravitonctl"),
				},
			},
		},
	}

	describeInstancesResult, err := ec2svc.DescribeInstances(&describeInstancesInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return instances, err
			}
		} else {
			return instances, err
		}
	}

	for _, reservations := range describeInstancesResult.Reservations {
		instances = append(instances, reservations.Instances...)
	}

	return instances, nil
}

func DescribeAllRunningInstances() (instances []*ec2.Instance, err error) {
	allInstances, err := describeAllInstances()
	if err != nil {
		return instances, err
	}

	for _, instance := range allInstances {
		if *instance.State.Name == ec2.InstanceStateNameRunning {
			instances = append(instances, instance)
		}
	}

	return instances, err
}

func DescribeAllStoppedInstances() (instances []*ec2.Instance, err error) {
	allInstances, err := describeAllInstances()
	if err != nil {
		return instances, err
	}

	for _, instance := range allInstances {
		if *instance.State.Name == ec2.InstanceStateNameStopped {
			instances = append(instances, instance)
		}
	}

	return instances, err
}

func filterInstancesByName(unfilteredInstances []*ec2.Instance, name string) (instances []*ec2.Instance) {
	for _, instance := range unfilteredInstances {
		for _, tag := range instance.Tags {
			if *tag.Key == "Name" && *tag.Value == name {
				instances = append(instances, instance)
			}
		}
	}

	return instances
}

func describeInstances(name string) (instances []*ec2.Instance, err error) {
	unfilteredInstances, err := describeAllInstances()
	if err != nil {
		return instances, err
	}

	return filterInstancesByName(unfilteredInstances, name), err
}

func describeRunningInstances(name string) (instances []*ec2.Instance, err error) {
	unfilteredInstances, err := DescribeAllRunningInstances()
	if err != nil {
		return instances, err
	}

	return filterInstancesByName(unfilteredInstances, name), err
}

func describeStoppedInstances(name string) (instances []*ec2.Instance, err error) {
	unfilteredInstances, err := DescribeAllStoppedInstances()
	if err != nil {
		return instances, err
	}

	return filterInstancesByName(unfilteredInstances, name), err
}
