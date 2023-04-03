package poddeploy

// package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"k8s.io/client-go/tools/clientcmd"
)

// ! main and package also need same struct - can we use different thing here

type PodDetails struct {
	Podname          string
	Imagename        string
	PodIp            string
	HostIP           string
	StartTime        string //string //metav1.Time
	Labels           map[string]string
	PodUID           string
	ContainerRunning bool
	RestartCount     int
	OpenPort         int
	//todo Add cpu usage with this somehow
	EventInfo EventDetails
}

type EventDetails struct {
	EventType string
	Reason    string
	Message   string
}

func GetKubernetesClientSet() (*kubernetes.Clientset, error) {
	// Get the user's home directory
	//! whats the use of ctx here
	//Todo Need to set a client directory directly where we can use this connection directly to run kubernetes connection

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("check1")
		return nil, err
	}

	// Build the path to the kubeconfig file
	kubeConfigPath := filepath.Join(homeDir, ".kube", "config")

	// Load the kubeconfig file
	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		fmt.Println("check2")

		return nil, err
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		fmt.Println("check3")

		return nil, err
	}

	return clientset, nil
}

func GetPodLogs(podname string) {
	// ([]byte, error)
	// clientset, err := GetKubernetesClientSet()
	// if err != nil {
	// 	fmt.Println("K8s client not set", err)
	// }
	// ctx := context.Background()

	// podLogOpts := &v1.PodLogOptions{
	// 	// Container: containerName,	//defaults to only available container
	// 	Follow: true,
	// }

	// req := clientset.CoreV1().Pods("default").GetLogs(podname, podLogOpts)

	// stream, err := req.Stream(ctx)

	// defer stream.Close()

	// if err != nil {
	// 	// http.Error(w, fmt.Sprintf("Error opening stream: %v", err), http.StatusInternalServerError)
	// 	return nil, err
	// }

	// buf := make([]byte, 1024)
	// for {
	// 	n, err := stream.Read(buf)
	// 	if err != nil {
	// 		break
	// 	}
	// 	if err := ws.WriteMessage(websocket.TextMessage, buf[:n]); err != nil {
	// 		break
	// 	}
	// 	// Sleep for a short time to avoid overwhelming the client with data
	// 	time.Sleep(50 * time.Millisecond)
	// }

}

// ! main and package also need same struct - can we use different thing here

func PodListing() ([]PodDetails, error) {
	// func main() {

	clientset, err := GetKubernetesClientSet()
	if err != nil {
		fmt.Println("check4")

		fmt.Println("K8s client not set", err)
	}
	ctx := context.Background()
	// metav1list := metav1.ListOptions{}
	// PodListing()
	pods, err := clientset.CoreV1().Pods("default").List(context.Background(), metav1.ListOptions{})

	if err != nil {
		// panic(err.Error())
		return nil, err
	}

	// arrayofpods := []PodDetails{}
	// var finalArrayofPodDetails []PodDetails
	finalarray := []PodDetails{}

	// fmt.Printf("pods: %#v\n", pods)
	for _, pod := range pods.Items {

		var eventInfo EventDetails

		//!Events testing & pod metrics
		//todo

		events, err := clientset.CoreV1().Events("default").List(ctx, metav1.ListOptions{FieldSelector: fmt.Sprintf("involvedObject.name=%s", pod.Name)})
		if err != nil {
			panic(err.Error())
		}
		if len(events.Items) == 0 {
			fmt.Println("Pod running")

		} else {
			lastevent := events.Items[len(events.Items)-1] //retreiving last values
			fmt.Printf("Event Type: %s, Reason: %s, Message: %s,\n", lastevent.Type, lastevent.Reason, lastevent.Message)
			eventInfo = EventDetails{
				EventType: lastevent.Type,
				Reason:    lastevent.Reason,
				Message:   lastevent.Message,
			}
		}
		// tempvaluefortime := time.Time(pod.Status.ContainerStatuses[0].State.Running.StartedAt.Time).String()
		// fmt.Println(tempvaluefortime)

		tempPod := PodDetails{
			Podname:   pod.Name,
			Imagename: pod.Spec.Containers[0].Image,
			PodIp:     pod.Status.PodIP,
			HostIP:    pod.Status.HostIP,
			StartTime: "",
			// StartTime: startTimeAsTime,
			// StartTime:        pod.Status.ContainerStatuses[0].State.Running.StartedAt,
			Labels:           pod.Labels,
			PodUID:           string(pod.UID),
			ContainerRunning: pod.Status.ContainerStatuses[0].Ready,
			RestartCount:     int(pod.Status.ContainerStatuses[0].RestartCount),
			OpenPort:         80,
			EventInfo:        eventInfo,
		}
		// fmt.Println(pod.Name, pod.Spec.Containers[0].Image, pod.Status.PodIP, pod.Status.ContainerStatuses[0].Ready, pod.Status.ContainerStatuses[0].State.Running.StartedAt, pod.Status.ContainerStatuses[0].RestartCount, pod.UID, pod.Labels, pod.Status.HostIP, pod.Status.Phase, pod.Spec.NodeName)
		finalarray = append(finalarray, tempPod)

		// Iterate through the events and print their details
		// for _, event := range events.Items {
		// 	fmt.Printf("Event Type: %s, Reason: %s, Message: %s, Container Readiness: %v\n", event.Type, event.Reason, event.Message, pod.Status.ContainerStatuses[0].Ready)
		// }
		// -----------------------------------------------
		//todo
		// Get pod metrics
		// podName := "my-pod"
		// podMetrics, err := metricsClient.MetricsV1beta1().PodMetricses("default").Get(podName, v1beta1.GetOptions{})
		// if err != nil {
		// 	panic(err.Error())
		// }

		// // Get CPU usage
		// cpuUsage := int64(0)
		// for _, container := range podMetrics.Containers {
		// 	cpuUsage += container.Usage.Cpu().MilliValue()
		// }

		// fmt.Printf("CPU usage of pod %s: %d milliCPU\n", podName, cpuUsage)
		// ---------------------------------------------------------
		// pod.Spec.
		// fmt.Println(pod.Name, pod.Status.Phase, pod.CreationTimestamp.Time, pod.Labels, pod.ManagedFields, pod.Namespace, pod.Spec.Hostname, pod.Spec.HostUsers, pod.Spec.NodeName, pod.Spec.Volumes, pod.Status)
		// fmt.Printf("Name: %s, Status: %s\n", pod.Name, pod.Status.Phase, pod.Labels, pod.CreationTimestamp.Time)
	}

	// return finalArrayofPodDetails
	fmt.Println(len(finalarray))
	return finalarray, nil
	//
}

func PodCreation(podname string, podimage string) error {

	ctx := context.Background() //! whats the use of ctx here
	//Todo Need to set a client directory directly where we can use this connection directly to run kubernetes connection
	userHomeDir, err := os.UserHomeDir() //?Getting user home directory
	if err != nil {
		fmt.Printf("error getting user home dir: %v\n", err)
		os.Exit(1)
	}
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
	fmt.Printf("Using kubeconfig: %s\n", kubeConfigPath)

	kubeConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		fmt.Printf("Error getting kubernetes config: %v\n", err)
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(kubeConfig)

	if err != nil {
		log.Fatal(err)
	}

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			// Name: "example-pod",
			Name:   podname,
			Labels: map[string]string{"app": "my-app", "tier": "frontend"},
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  podname,
					Image: podimage,
					Ports: []v1.ContainerPort{
						{
							Name:          "http",
							ContainerPort: 80,
							Protocol:      v1.ProtocolTCP,
						},
					},
					Resources: v1.ResourceRequirements{
						Requests: v1.ResourceList{
							v1.ResourceCPU:    resource.MustParse("100m"),
							v1.ResourceMemory: resource.MustParse("128Mi"),
						},
						Limits: v1.ResourceList{
							v1.ResourceCPU:    resource.MustParse("200m"),
							v1.ResourceMemory: resource.MustParse("256Mi"),
						},
					},
					//todo for multiple ports on one single pod - need to change for the deployments
					//** mounting part

					// Name:       "example-container",
					// Image:      "nginx",
					// WorkingDir: "/usr/share/nginx/html",
					// volume mounts in a pod, liveness, readiness probes,
					// Command:         []string{"sh", "-c", "echo Hello Kubernetes!"},
					// Args:            []string{"arg1", "arg2"},
					// Env: []v1.EnvVar{
					// 	{
					// 		Name:  "ENV_VAR_1",
					// 		Value: "value1",
					// 	},
					// 	{
					// 		Name:  "ENV_VAR_2",
					// 		Value: "value2",
					// 		ValueFrom: &v1.EnvVarSource{
					// 			ConfigMapKeyRef: &v1.ConfigMapKeySelector{
					// 				LocalObjectReference: v1.LocalObjectReference{
					// 					Name: "config-map",
					// 				},
					// 				Key: "key",
					// 			},
					// 		},
					// 	},
					// },

				},
			},
		},
	}

	// metav1.
	// metapod := metav1.CreateOptions{}

	// clientset.CoreV1().Pods("default").Create(pod)

	podCreation, err := clientset.CoreV1().Pods("default").Create(ctx, pod, metav1.CreateOptions{})
	fmt.Printf("err: %v\n", err)
	if err != nil {
		fmt.Printf("error creating pod: %v\n", err)
		return err
		// os.Exit(1)
	}
	fmt.Printf("podCreation.Labels: %v\n", podCreation.Labels)

	return nil

}
