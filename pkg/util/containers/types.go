// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2017 Datadog, Inc.

package containers

import (
	"github.com/DataDog/datadog-agent/pkg/util/containers/metrics"
)

// Known container runtimes
const (
	RuntimeNameDocker     string = "docker"
	RuntimeNameContainerd string = "containerd"
	RuntimeNameCRIO       string = "cri-o"
)

// Supported container states
const (
	ContainerUnknownState    string = "unknown"
	ContainerCreatedState           = "created"
	ContainerRunningState           = "running"
	ContainerRestartingState        = "restarting"
	ContainerPausedState            = "paused"
	ContainerExitedState            = "exited"
	ContainerDeadState              = "dead"
)

// Supported container health
const (
	ContainerUnknownHealth  string = "unknown"
	ContainerStartingHealth        = "starting"
	ContainerHealthy               = "healthy"
	ContainerUnhealthy             = "unhealthy"
)

// Container represents a single container on a machine
// and includes Cgroup-level statistics about the container.
type Container struct {
	Type     string
	ID       string
	EntityID string
	Name     string
	Image    string
	ImageID  string
	Created  int64
	State    string
	Health   string
	Pids     []int32
	Excluded bool

	CPULimit       float64
	SoftMemLimit   uint64
	MemLimit       uint64
	CPUNrThrottled uint64
	CPU            *metrics.CgroupTimesStat
	Memory         *metrics.CgroupMemStat
	IO             *metrics.CgroupIOStat
	Network        metrics.ContainerNetStats
	StartedAt      int64

	// For internal use only
	cgroup *metrics.ContainerCgroup
}
