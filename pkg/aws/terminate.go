package aws

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
)

// Terminate terminates a specific EC2 instance
func Terminate(name string) {

	instances := describeInstance(name)

	if len(instances) == 0 {
		log.Infof("No instance with name %s found", name)
	}

	for _, instance := range instances {

		log.Infof("Terminating: %s", *instance.InstanceId)

		terminateInstancesInput := &ec2.TerminateInstancesInput{
			InstanceIds: []*string{
				instance.InstanceId,
			},
		}

		_, err := ec2svc.TerminateInstances(terminateInstancesInput)

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
			return
		}
	}








}