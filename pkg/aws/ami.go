package aws

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var AMIPrefix = "amzn2-ami-kernel-5.10-*"

func GetAMI() (image string, err error) {
	images, err := ec2svc.DescribeImages(&ec2.DescribeImagesInput{
		IncludeDeprecated: aws.Bool(false),
		Filters: []*ec2.Filter{
			{
				Name: aws.String("architecture"),
				Values: []*string{
					aws.String(ec2.ArchitectureTypeArm64),
				},
			},
			{
				Name: aws.String("owner-alias"),
				Values: []*string{
					aws.String("amazon"),
				},
			},
			{
				Name: aws.String("name"),
				Values: []*string{
					aws.String(AMIPrefix),
				},
			},
		},
	})

	if err != nil {
		return "", err
	}

	if len(images.Images) > 0 {
		return *images.Images[0].ImageId, nil
	}

	return "", errors.New("no amazon-linux-2 image found")
}
