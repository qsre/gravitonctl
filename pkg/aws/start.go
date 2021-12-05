package aws

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"

	log "github.com/sirupsen/logrus"
)

// Start starts an EC2 Graviton instance
func Start(name string) error {

	instances, err := describeRunningInstances(name)
	if err != nil {
		return err
	}

	if len(instances) > 0 {
		return errors.New("this instance already exists")
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
		ImageId:      aws.String("ami-0b5d05f884ff8bd47"),
		InstanceType: aws.String(ec2.InstanceTypeT4gMicro),
		KeyName:      aws.String("berty_key"),
		MaxCount:     aws.Int64(1),
		MinCount:     aws.Int64(1),

		// SecurityGroupIds: []*string{
		// 	aws.String("sg-064d01f01ebe545ce"),
		// },
		// SubnetId: aws.String("subnet-51519f2a"),

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
		return err
	}

	log.Infof("Instance launched: %s", *result.Instances[0].InstanceId)

	return nil
}
