package aws

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func GetIp(name string) (ip string, err error) {
	instances, err := describeInstances(name)

	if err != nil {
		return ip, err
	}

	if len(instances) == 0 {
		return "", errors.New(fmt.Sprintf("no instance with name %s found", name))
	} else {
		for _, instance := range instances {
			if *instance.State.Name == ec2.InstanceStateNameRunning || *instance.State.Name == ec2.InstanceStateNamePending {
				if instances[0].PublicIpAddress == nil {
					return "", errors.New("public IP isn't available yet")
				}

				return *instances[0].PublicIpAddress, nil
			}
		}
	}

	return "", errors.New(fmt.Sprintf("multiple but no valid instances with name %s found", name))

}
