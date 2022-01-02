package aws

import "github.com/aws/aws-sdk-go/service/ec2"

func GetKeyNames() (keys []string, err error) {
	output, err := ec2svc.DescribeKeyPairs(&ec2.DescribeKeyPairsInput{})
	if err != nil {
		return keys, err
	}

	for _, keyPair := range output.KeyPairs {
		keys = append(keys, *keyPair.KeyName)
	}

	return keys, err
}
