package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
)

var REGION = "eu-west-3"
var sess *session.Session

var ec2svc *ec2.EC2

func init() {
	sess = session.Must(session.NewSession(&aws.Config{
		Region: aws.String(REGION),
	}))

	ec2svc = ec2.New(sess)
}

func describeAllInstances() (instances []*ec2.Instance){
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
				log.Error(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			log.Error(err.Error())
		}
		return instances
	}

	for _, reservations := range describeInstancesResult.Reservations {
		instances = append(instances, reservations.Instances...)
	}

	return instances
}

func describeInstance(name string) (instances []*ec2.Instance) {
	unfilteredInstances := describeAllInstances()

	for _, instance := range unfilteredInstances {
		for _, tag := range instance.Tags {
			if *tag.Key == "Name" && *tag.Value == name {
				// check for multiple instances with the same name tag
				instances = append(instances, instance)
			}
		}
	}

	return instances
}