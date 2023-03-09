package mongotesting

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

const (
	image         = "mongo:4.4"
	containerPort = "27017/tcp"
)

// RunWithMongoInDocker runs the tests with a docker container
func RunWithMongoInDocker(m *testing.M, mongoURI *string) int {
	c, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	resp, err := c.ContainerCreate(ctx, &container.Config{
		Image: image,
		ExposedPorts: nat.PortSet{
			containerPort: {},
		},
	}, &container.HostConfig{
		PortBindings: nat.PortMap{
			containerPort: []nat.PortBinding{
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

	defer func() {
		err := c.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{
			Force: true,
		})
		panic(err)
	}()
	containerID := resp.ID
	err = c.ContainerStart(ctx, containerID, types.ContainerStartOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Println("container started")
	time.Sleep(5 * time.Second)

	inspRes, err := c.ContainerInspect(ctx, containerID)
	if err != nil {
		panic(err)
	}
	hostPort := inspRes.NetworkSettings.Ports[containerPort][0]
	*mongoURI = fmt.Sprintf("mongodb://%s:%s", hostPort.HostIP, hostPort.HostPort)
	fmt.Printf("listening at %+v\n", inspRes.NetworkSettings.Ports["27017/tcp"][0])
	return m.Run()
}
