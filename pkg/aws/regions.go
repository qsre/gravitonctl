package aws

func GetRegions() (regions []string, err error) {
	resultRegions, err := ec2svc.DescribeRegions(nil)
	if err != nil {
		return nil, err
	}

	for _, region := range resultRegions.Regions {
		regions = append(regions, *region.RegionName)
	}

	return regions, err
}
