package cluster

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/resourcegroupstaggingapi"
)

const (
	TagKeyNodeGroupName     = "tf-eksctl/node-group"
	TagKeyClusterNamePrefix = "tf-eksctl/cluster"
)

func getTargetGroupARNs(sess *session.Session, clusterNamePrefixy string) ([]string, error) {
	api := resourcegroupstaggingapi.New(sess)

	var token *string

	var arns []string

	for {
		log.Printf("getting tagged resources for %s", clusterNamePrefixy)

		res, err := api.GetResources(&resourcegroupstaggingapi.GetResourcesInput{
			PaginationToken:     token,
			ResourceTypeFilters: aws.StringSlice([]string{"elasticloadbalancing:targetgroup"}),
			TagFilters: []*resourcegroupstaggingapi.TagFilter{
				{
					Key:    aws.String(TagKeyClusterNamePrefix),
					Values: aws.StringSlice([]string{clusterNamePrefixy}),
				},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("getting tagged resources for %s: %w", clusterNamePrefixy, err)
		}

		for _, m := range res.ResourceTagMappingList {
			if arn := m.ResourceARN; arn != nil {
				arns = append(arns, *arn)
			}
		}

		token = res.PaginationToken
		if token == nil || *token == "" {
			break
		}
	}

	return arns, nil
}
