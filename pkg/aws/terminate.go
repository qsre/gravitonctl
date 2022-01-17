package aws

import (
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
)

// Terminate terminates a specific EC2 instance
func Terminate(name string) error {

	instances, err := describeRunningInstances(name)
	if err != nil {
		return err
	}

	for _, instance := range instances {

		log.Infof("Terminating: %s", *instance.InstanceId)

		terminateInstancesInput := &ec2.TerminateInstancesInput{
			InstanceIds: []*string{
				instance.InstanceId,
			},
		}

		_, err = ec2svc.TerminateInstances(terminateInstancesInput)

		if err != nil {
			return err
		}
	}

	for _, instance := range instances {
		err = ec2svc.WaitUntilInstanceTerminated(&ec2.DescribeInstancesInput{
			InstanceIds: []*string{
				instance.InstanceId,
			},
		})
		if err != nil {
			return err
		}
	}

	sgIds, err := getSecurityGroupIds(securityGroupName(name))
	if err != nil {
		return err
	}

	for _, sgId := range sgIds {
		_, err = ec2svc.DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{
			GroupId: &sgId,
		})
		if err != nil {
			return err
		}
	}

	if len(instances) == 0 {
		log.Infof("No instance with name %s found", name)
		return nil
	}

	return err
}
