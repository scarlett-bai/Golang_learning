package main

import (
	"context"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func main() {
	c, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	resp, err := c.ContainerCreate(ctx, &container.Config{
		Image: "mongo:4.4",
		ExposedPorts: nat.PortSet{
			"27017/tcp": {},
		},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			"27017/tcp": []nat.PortBinding{
				{
					HostIP: "127.0.0.1", // 只接受本地的请求，且不需要权限  0.0.0.0 是可以接受所有的请求
					// HostPort: "27018",
					HostPort: "0", // 可以动态分配一个闲置的端口使用，而不需要特意分配
				},
			},
		},
	}, nil, "")

	if err != nil {
		panic(err)
	}

	err = c.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("container started")
	time.Sleep(10 * time.Second)

	inspRes, err := c.ContainerInspect(ctx, resp.ID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("listening at %+v\n", inspRes.NetworkSettings.Ports["27017/tcp"][0])
	err = c.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		panic(err)
	}
}
