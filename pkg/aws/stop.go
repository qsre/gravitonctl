package aws

import (
	"errors"

	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
)

// Stop stops an EC2 Graviton instance
func Stop(name string, keyName string) error {

	instances, err := describeRunningInstances(name)
	if err != nil {
		return err
	}

	if len(instances) == 0 {
		return errors.New("this instance is already stopped")
	}

	log.Infof("stopping %s\n", keyName)

	var instanceIds []*string
	for _, instance := range instances {
		instanceIds = append(instanceIds, instance.InstanceId)
	}

	_, err = ec2svc.StopInstances(&ec2.StopInstancesInput{
		InstanceIds: instanceIds,
	})
	if err != nil {
		return err
	}

	err = ec2svc.WaitUntilInstanceStopped(&ec2.DescribeInstancesInput{
		InstanceIds: instanceIds,
	})
	if err != nil {
		return err
	}

	log.Infof("stopped %s\n", keyName)

	return nil
}
