package k8sclient

import (
	"fmt"
	"os"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

func T2() error {
	config := restclient.Config{

		Host: "https://114.67.244.228:6443",
		//Username: "kube-admin",
		//Password: "E5m*O9z8T3E5",
		BearerToken: "Basic a3ViZS1hZG1pbjpFNW0qTzl6OFQzRTU=",
		TLSClientConfig: restclient.TLSClientConfig{
			CertFile: "/Users/fengzhikui/tmp/k8s/client.crt",
			KeyFile:  "/Users/fengzhikui/tmp/k8s/client.key",
			CAFile:   "/Users/fengzhikui/tmp/k8s/ca.crt",
		},
	}
	podName := "ide-cpongo1-62f9b28804b1185df1be16af"
	namespace := "default"
	//containerName := podName
	//cmd := []string{"pwd"}

	cmd := []string{
		"sh",
		"-c",
		"pwd",
	}
	client, err := kubernetes.NewForConfig(&config)
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}
	req := client.CoreV1().RESTClient().Post().Resource("pods").Name(podName).
		Namespace(namespace).SubResource("exec")
	option := &v1.PodExecOptions{
		Command: cmd,
		Stdin:   false,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}

	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)
	exec, err := remotecommand.NewSPDYExecutor(&config, "POST", req.URL())
	if err != nil {
		return err
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}
	return nil
}

func T3() error {
	config := restclient.Config{
		Host: "https://119.3.190.86:6443",
		TLSClientConfig: restclient.TLSClientConfig{
			CertFile: "/Users/fengzhikui/tmp/key/client.crt",
			KeyFile:  "/Users/fengzhikui/tmp/key/client.key",
			CAFile:   "/Users/fengzhikui/tmp/key/ca.crt",
		},
	}
	podName := "ide-cpongo2-62f9b28804b1185df1be16af"
	namespace := "inscode"
	//containerName := podName
	//cmd := []string{"pwd"}

	cmd := []string{"sh", "-c", "pwd"}
	client, err := kubernetes.NewForConfig(&config)
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}
	req := client.CoreV1().RESTClient().Post().Resource("pods").Name(podName).
		Namespace(namespace).SubResource("exec")
	option := &v1.PodExecOptions{
		Command: cmd,
		Stdin:   false,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}

	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)
	exec, err := remotecommand.NewSPDYExecutor(&config, "POST", req.URL())
	if err != nil {
		return err
	}
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  nil,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	})
	if err != nil {
		fmt.Printf("%s", err)
		return err
	}
	return nil
}
