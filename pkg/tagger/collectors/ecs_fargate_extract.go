// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2017 Datadog, Inc.

// +build docker

package collectors

import (
	"fmt"
	"strings"
	"time"

	"github.com/DataDog/datadog-agent/pkg/tagger/utils"
	"github.com/DataDog/datadog-agent/pkg/util/containers"
	"github.com/DataDog/datadog-agent/pkg/util/docker"
	"github.com/DataDog/datadog-agent/pkg/util/ecs"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)

// parseMetadata parses the task metadata and its container list, and returns a list of TagInfo for the new ones.
// It also updates the lastSeen cache of the ECSFargateCollector and return the list of dead containers to be expired.
func (c *ECSFargateCollector) parseMetadata(meta ecs.TaskMetadata, parseAll bool) ([]*TagInfo, error) {
	var output []*TagInfo
	now := time.Now()

	if meta.KnownStatus != "RUNNING" {
		return output, fmt.Errorf("Task %s is in %s status, skipping", meta.Family, meta.KnownStatus)
	}

	for _, ctr := range meta.Containers {
		if c.expire.Update(ctr.DockerID, now) || parseAll {
			tags := utils.NewTagList()

			// cluster
			tags.AddLow("cluster_name", parseECSClusterName(meta.ClusterName))

			// task
			tags.AddLow("task_family", meta.Family)
			tags.AddLow("task_version", meta.Version)

			// container
			tags.AddLow("ecs_container_name", ctr.Name)
			tags.AddHigh("container_id", ctr.DockerID)
			tags.AddHigh("container_name", ctr.DockerName)

			// container image
			tags.AddLow("docker_image", ctr.Image)
			imageName, shortImage, imageTag, err := containers.SplitImageName(ctr.Image)
			if err != nil {
				log.Debugf("Cannot split %s: %s", ctr.Image, err)
			} else {
				tags.AddLow("image_name", imageName)
				tags.AddLow("short_image", shortImage)
				if imageTag == "" {
					imageTag = "latest"
				}
				tags.AddLow("image_tag", imageTag)
			}

			for labelName, labelValue := range ctr.Labels {
				if tagName, found := c.labelsAsTags[strings.ToLower(labelName)]; found {
					tags.AddAuto(tagName, labelValue)
				}
			}

			low, high := tags.Compute()
			info := &TagInfo{
				Source:       ecsFargateCollectorName,
				Entity:       docker.ContainerIDToEntityName(string(ctr.DockerID)),
				HighCardTags: high,
				LowCardTags:  low,
			}
			output = append(output, info)
		}
	}

	return output, nil
}

// parseECSClusterName allows to handle user-friendly values and arn values
func parseECSClusterName(value string) string {
	if strings.Contains(value, "/") {
		parts := strings.Split(value, "/")
		return parts[len(parts)-1]
	} else {
		return value
	}
}
