package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
)

// Terminate terminates a specific EC2 instance
func Terminate(name string) error {
	instances, err := describeInstances(name)
	if err != nil {
		return err
	}

	for _, instance := range instances {
		if *instance.State.Name == ec2.InstanceStateNameTerminated {
			continue

			//TODO this needs work!
		}

		log.Infof("Terminating: %s", *instance.InstanceId)

		// range over and delete network interfaces
		for niNum, ni := range instance.NetworkInterfaces {
			if niNum == 0 {
				continue
			}

			log.Infof("Deleting ni %s", *ni.NetworkInterfaceId)

			_, err := ec2svc.DetachNetworkInterface(&ec2.DetachNetworkInterfaceInput{
				AttachmentId: ni.Attachment.AttachmentId,
				Force:        aws.Bool(true),
			})
			if err != nil {
				return err
			}

			err = ec2svc.WaitUntilNetworkInterfaceAvailable(&ec2.DescribeNetworkInterfacesInput{
				NetworkInterfaceIds: []*string{
					instance.InstanceId,
				},
			})
			if err != nil {
				return err
			}

			_, err = ec2svc.DeleteNetworkInterface(&ec2.DeleteNetworkInterfaceInput{
				NetworkInterfaceId: ni.NetworkInterfaceId,
			})
			if err != nil {
				return err
			}
		}

		terminateInstancesInput := &ec2.TerminateInstancesInput{
			InstanceIds: []*string{
				instance.InstanceId,
			},
		}

		_, err = ec2svc.TerminateInstances(terminateInstancesInput)

		if err != nil {
			return err
		}

		err = ec2svc.WaitUntilInstanceTerminated(&ec2.DescribeInstancesInput{
			InstanceIds: []*string{
				instance.InstanceId,
			},
		})
		if err != nil {
			return err
		}
	}

	for _, instance := range instances {
		// range over and delete attached security groups
		for _, sg := range instance.SecurityGroups {
			log.Infof("Deleting sg: %s", *sg.GroupId)
			_, err = ec2svc.DeleteSecurityGroup(&ec2.DeleteSecurityGroupInput{
				GroupId: sg.GroupId,
			})
			if err != nil {
				return err
			}
		}
	}

	if len(instances) == 0 {
		log.Infof("No instance with name %s found", name)
		return nil
	}

	return err
}
