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
		return errors.New("this instance already exists")
	}

	// Security Groups
	// TODO: needs to be in gravitonctl/pkg/aws

	sgName := fmt.Sprintf("%s-sg", name)
	var sg ec2.SecurityGroup

	sgOutput, _ := ec2svc.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		GroupNames: []*string{
			&sgName,
		},
	})

	if sgOutput.SecurityGroups == nil {
		_, err := ec2svc.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
			Description: &sgName,
			GroupName:   &sgName,
		})
		if err != nil {
			return err
		}

		err = ec2svc.WaitUntilSecurityGroupExists(&ec2.DescribeSecurityGroupsInput{
			GroupNames: []*string{
				&sgName,
			},
		})
		if err != nil {
			return err
		}

		describeSgOutput, err := ec2svc.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
			GroupNames: []*string{
				&sgName,
			},
		})

		if err != nil {
			return err
		}

		sg = *describeSgOutput.SecurityGroups[0]

		_, err = ec2svc.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
			GroupId: describeSgOutput.SecurityGroups[0].GroupId,
			IpPermissions: []*ec2.IpPermission{
				{
					FromPort:   aws.Int64(22),
					ToPort:     aws.Int64(22),
					IpProtocol: aws.String("tcp"),
					IpRanges: []*ec2.IpRange{
						{
							CidrIp:      aws.String("0.0.0.0/0"),
							Description: aws.String("Inbound ssh access"),
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}

		// outbound
		_, err = ec2svc.UpdateSecurityGroupRuleDescriptionsEgress(&ec2.UpdateSecurityGroupRuleDescriptionsEgressInput{
			GroupId: describeSgOutput.SecurityGroups[0].GroupId,
			IpPermissions: []*ec2.IpPermission{
				{
					IpProtocol: aws.String("-1"),
					IpRanges: []*ec2.IpRange{
						{
							CidrIp:      aws.String("0.0.0.0/0"),
							Description: aws.String("Outbound access"),
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}

	} else {
		sg = *sgOutput.SecurityGroups[0]
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
			sg.GroupId,
		},

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
	}

	log.Infof("Instance launched: %s", *result.Instances[0].InstanceId)

	return nil
}
