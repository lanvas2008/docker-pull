package docker_tools

import (
	"bufio"
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"io"
	"os"
)

func ImagePull(imageName string, pullOptions image.PullOptions) bool {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println("Error creating client: %v\n", err)
		return false
	}
	pullResponse, err := cli.ImagePull(ctx, imageName, pullOptions)

	// 开始拉取镜像
	if err != nil {
		fmt.Println("Error pulling image: %v\n", err)
		return false
	}
	var reader = io.TeeReader(pullResponse, os.Stdout) // 将读取到的数据同时输出到标准输出
	var bufReader = bufio.NewReader(reader)
	// 打印拉取进度
	for {
		line, err := bufReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading line:", err)
			break
		}
		fmt.Println("pull:", string(line))
	}
	return true
}

// ContainerList 获取容器列表 /**
func ContainerList(listOptions container.ListOptions) []types.Container {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println("Error creating client: %v\n", err)
		return []types.Container{}
	}
	containers, err := cli.ContainerList(ctx, listOptions)
	if err != nil {
		fmt.Println("Error listing containers: %v\n", err)
		return []types.Container{}
	}
	return containers
}

func ContainerRemove(containerID string, options container.RemoveOptions) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println("Error creating client: %v\n", err)
		return err
	}
	err = cli.ContainerRemove(ctx, containerID, options)
	return err
}

func ContainerCreate(config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *ocispec.Platform, containerName string) (container.CreateResponse, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println("Error creating client: %v\n", err)
		return container.CreateResponse{}, err
	}
	return cli.ContainerCreate(ctx, config, hostConfig, networkingConfig, platform, containerName)
}

func ContainerStart(containerID string, options container.StartOptions) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println("Error creating client: %v\n", err)
		return err
	}
	err = cli.ContainerStart(ctx, containerID, options)
	return err
}

func ContainerStop(containerID string, options container.StopOptions) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println("Error creating client: %v\n", err)
		return err
	}
	err = cli.ContainerStop(ctx, containerID, options)
	return err
}

func ContainerPause(containerID string, signal string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println("Error creating client: %v\n", err)
		return err
	}
	err = cli.ContainerPause(ctx, containerID)
	return err
}

func ContainerKill(containerID string, signal string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println("Error creating client: %v\n", err)
		return err
	}
	err = cli.ContainerKill(ctx, containerID, signal)
	return err
}

func ContainerStats(containerID string, stream bool) (container.StatsResponseReader, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println("Error creating client: %v\n", err)
		return container.StatsResponseReader{}, err
	}
	return cli.ContainerStats(ctx, containerID, stream)
}

func ImageList(options image.ListOptions) ([]image.Summary, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println("Error creating client: %v\n", err)
		return []image.Summary{}, err
	}
	return cli.ImageList(ctx, options)
}

func ImageLoad(FileName string) bool {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println("Error creating client: %v\n", err)
		return false
	}

	file, _ := os.Open(FileName)
	defer file.Close()

	response, err := cli.ImageLoad(ctx, file, false)
	defer response.Body.Close()
	// 开始拉取镜像
	if err != nil {
		fmt.Println("Error load image: %v\n", err)
		return false
	}

	var reader = io.TeeReader(response.Body, os.Stdout) // 将读取到的数据同时输出到标准输出
	var bufReader = bufio.NewReader(reader)
	//打印拉取进度
	for {
		_, err := bufReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading line:", err)
			return false
		}
	}
	return true

}

func GetImage(imageId string) image.Summary {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println("error : ", err)
		return image.Summary{}
	}
	options := image.ListOptions{}
	options.All = true
	options.Filters = filters.NewArgs()
	ret, _ := cli.ImageList(ctx, options)
	for i := range ret {
		if ret[i].ID == imageId {
			return ret[i]
		}
	}
	return image.Summary{}
}

func ImageTag(imageId string, imageName string, tag string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		fmt.Println("error : ", err)
		return err
	}
	options := image.ListOptions{}
	options.All = true
	options.Filters = filters.NewArgs()
	return cli.ImageTag(ctx, imageId, imageName+":"+tag)
}
