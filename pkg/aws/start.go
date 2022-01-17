package aws

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"

	log "github.com/sirupsen/logrus"
)

// Start starts an EC2 Graviton instance
func Start(name string, keyName string) error {

	instances, err := describeRunningInstances(name)
	if err != nil {
		return err
	}

	if len(instances) > 0 {
		return errors.New("this instance is already started")
	}

	sgName := securityGroupName(name)
	var sgId string

	groupIds, err := getSecurityGroupIds(sgName)
	if err != nil {
		return err
	}

	fmt.Println(groupIds)

	// check if security group exists
	if len(groupIds) == 0 {

		// create security group
		groupId, err := createSecurityGroup(sgName)
		if err != nil {
			return err
		}

		// create SSH rules
		err = createSSHRules(groupId)
		if err != nil {
			return err
		}

		sgId = groupId

	} else {
		sgId = groupIds[0]
	}

	// bare minimum input
	input := &ec2.RunInstancesInput{
		BlockDeviceMappings: []*ec2.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/sdh"),
				Ebs: &ec2.EbsBlockDevice{
					VolumeSize: aws.Int64(16),
				},
			},
		},
		ImageId:      aws.String("ami-0a7559a0ef82639f2"),
		InstanceType: aws.String(ec2.InstanceTypeT4gMicro),
		KeyName:      &keyName,
		MaxCount:     aws.Int64(1),
		MinCount:     aws.Int64(1),

		SecurityGroupIds: []*string{
			&sgId,
		},

		TagSpecifications: []*ec2.TagSpecification{
			{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					{
						Key:   aws.String("created-by"),
						Value: aws.String("gravitonctl"),
					},
					{
						Key:   aws.String("Name"),
						Value: &name,
					},
				},
			},
		},
	}

	result, err := ec2svc.RunInstances(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			default:
				return err
			}
		} else {
			return err
		}
	}

	log.Infof("Instance launched: %s", *result.Instances[0].InstanceId)

	return nil
}
