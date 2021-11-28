package aws

import (
	"errors"
	"fmt"
)

func GetIp(name string) (ip string, err error) {
	instances := describeInstance(name)

	switch len(instances) {
	case 0:
		return "", errors.New(fmt.Sprintf("no instance with name %s found", name))
	case 2:
		return "", errors.New(fmt.Sprintf("more than one instances with name %s found", name))
	}

	if instances[0].PublicIpAddress == nil {
		return "", errors.New("public IP isn't available yet")
	}

	return *instances[0].PublicIpAddress, nil
}