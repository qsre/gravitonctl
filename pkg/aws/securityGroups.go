package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func securityGroupName(name string) string {
	return fmt.Sprintf("%s-sg", name)
}

func getSecurityGroupIds(name string) (groupIds []string, err error) {
	// describe security group
	sgOutput, err := ec2svc.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		GroupNames: []*string{
			&name,
		},
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "InvalidGroup.NotFound":
				return groupIds, nil
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			return groupIds, err
		}
	}

	for _, sg := range sgOutput.SecurityGroups {
		groupIds = append(groupIds, *sg.GroupId)
	}

	return groupIds, err
}

// creates SSH rules for security group
// takes groupId
func createSSHRules(groupId string) (err error) {
	// inbound
	_, err = ec2svc.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
		GroupId: &groupId,
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
		GroupId: &groupId,
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

	return nil
}

func createSecurityGroup(sgName string) (groupId string, err error) {
	createSgOutput, err := ec2svc.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
		Description: &sgName,
		GroupName:   &sgName,
	})
	if err != nil {
		return "", err
	}

	err = ec2svc.WaitUntilSecurityGroupExists(&ec2.DescribeSecurityGroupsInput{
		GroupNames: []*string{
			&sgName,
		},
	})
	if err != nil {
		return "", err
	}

	return *createSgOutput.GroupId, nil
}
