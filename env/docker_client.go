package env

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"regexp"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)


type DockerClient struct {
	Cli      *client.Client
	Address  string
	Ctx      context.Context
	Reader   io.ReadCloser
	Response container.CreateResponse
}


func (dc *DockerClient) Close() {
	dc.Reader.Close()
}

func NewDockerClientEnv(imageName string, containerName string, portMapping string) (*DockerClient, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %v", err)
	}

	existingContainer, err := cli.ContainerInspect(context.Background(), containerName)
	if err == nil {
		if existingContainer.State.Running {
			if err := cli.ContainerStop(context.Background(), containerName, container.StopOptions{}); err != nil {
				return nil, fmt.Errorf("failed to stop existing container: %v", err)
			}
		}
		if err := cli.ContainerRemove(context.Background(), containerName, types.ContainerRemoveOptions{}); err != nil {
			return nil, fmt.Errorf("failed to remove existing container: %v", err)
		}
	} else if !client.IsErrNotFound(err) {
		return nil, fmt.Errorf("error inspecting existing container: %v", err)
	}

	ctx := context.Background()
	reader, err := cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to pull Docker image: %v", err)
	}
	time.Sleep(2 * time.Second)

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: imageName,
		Tty:   true,
	}, &container.HostConfig{
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeVolume, // Or mount.TypeBind for bind mounts
				Source: "MyGoApp",     // Omit or specify volume name for Docker-managed volume
				Target: "/app",
			},
		},
	}, &network.NetworkingConfig{}, nil, containerName)
	if err != nil {
		return nil, fmt.Errorf("failed to create container: %v", err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, fmt.Errorf("failed to start container: %v", err)
	}

	time.Sleep(2 * time.Second)

	inspect, err := cli.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to inspect container: %v", err)
	}
	ipAddress := inspect.NetworkSettings.IPAddress

	dc := &DockerClient{
		Reader:   reader,
		Ctx:      ctx,
		Cli:      cli,
		Response: resp,
		Address:  fmt.Sprintf("%s:%s", ipAddress, portMapping),
	}
	fmt.Println(dc.Address)
	return dc, nil
}

func (dc *DockerClient) ExecuteCode(code string) (string, error) {
	execConfig := types.ExecConfig{
		Cmd:          []string{"sh", "-c", code},
		AttachStdout: true,
		AttachStderr: true,
	}
	execResp, err := dc.Cli.ContainerExecCreate(dc.Ctx, dc.Response.ID, execConfig)
	if err != nil {
		return "", fmt.Errorf("failed to create exec instance: %v", err)
	}

	execAttachResp, err := dc.Cli.ContainerExecAttach(dc.Ctx, execResp.ID, types.ExecStartCheck{})
	if err != nil {
		return "", fmt.Errorf("failed to attach to exec instance: %v", err)
	}
	defer execAttachResp.Close()

	outputBuffer := new(bytes.Buffer)
	_, err = io.Copy(outputBuffer, execAttachResp.Reader)
	if err != nil {
		return "", fmt.Errorf("failed to read exec output: %v", err)
	}

	// Use a regular expression to remove non-printable characters
    cleanOutput := regexp.MustCompile(`[\x00-\x1F\x7F-\x9F]`).ReplaceAllString(outputBuffer.String(), "")
    return cleanOutput, nil
}

//NewDockerClient is use to create a docker enviroment for the development of the vision app
//the Caller is responsible for closing Docker client gracefully by using the defer *DockerClient.Close()
func NewDockerClient(vision string)(*DockerClient, error){
	imageName := "golang"
	containerName := "my-go-app"
	portMapping := "7681:7681"

	dc, err := NewDockerClientEnv(imageName, containerName, portMapping)
	if err != nil {
		log.Fatalf("Failed to create DockerClient environment: %v", err)
	}
	
	// Example code to execute
	code := "echo '"+vision+"'"

	output, err := dc.ExecuteCode(code)
	if err != nil {
		log.Fatalf("Error executing code: %v", err)
	}

	fmt.Printf("Vision output: %s\n", output)
	return dc, nil
}