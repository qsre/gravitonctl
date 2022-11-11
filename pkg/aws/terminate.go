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

	if len(instances) == 0 {
		log.Infof("No instance with name %s found", name)
		return nil
	}

	for _, instance := range instances {
		if *instance.State.Name == ec2.InstanceStateNameTerminated {
			continue

			//TODO this needs work!
		}

		// range over and delete network interfaces
		for _, ni := range instance.NetworkInterfaces {
			if *ni.Attachment.DeviceIndex == 0 {
				// root network interface, cannot detach these
				continue
			}

			log.Infof("Detaching network interface: %s", *ni.NetworkInterfaceId)

			_, err := ec2svc.DetachNetworkInterface(&ec2.DetachNetworkInterfaceInput{
				AttachmentId: ni.Attachment.AttachmentId,
				Force:        aws.Bool(true),
			})
			if err != nil {
				return err
			}

			err = ec2svc.WaitUntilNetworkInterfaceAvailable(&ec2.DescribeNetworkInterfacesInput{
				NetworkInterfaceIds: []*string{
					ni.NetworkInterfaceId,
				},
			})
			if err != nil {
				return err
			}

			log.Warnf("Network interface: %s detached, but not deleted!", *ni.NetworkInterfaceId)

			// not deleting network interfaces as they may have elastic IP's attached
			// checking for, and deleting these elastic IP's will be implemented later

			// _, err = ec2svc.DeleteNetworkInterface(&ec2.DeleteNetworkInterfaceInput{
			// 	NetworkInterfaceId: ni.NetworkInterfaceId,
			// })
			// if err != nil {
			// 	return err
			// }
		}

		terminateInstancesInput := &ec2.TerminateInstancesInput{
			InstanceIds: []*string{
				instance.InstanceId,
			},
		}

		log.Infof("Terminating: %s", *instance.InstanceId)

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

	return err
}
