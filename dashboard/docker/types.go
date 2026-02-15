package docker

import "github.com/moby/moby/api/types/container"

type ContainerStats struct {
	CPU    float64
	Memory float64
}

type errMsg struct {
	Err error
}

type dockerMsg struct {
	Containers []container.Summary
}
