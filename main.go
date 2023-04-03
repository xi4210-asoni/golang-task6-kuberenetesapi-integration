package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	// "github.com/akshat-soni-xebia/k8scode"

	"github.com/akshat-soni-xebia/kubernetes_api/poddeploy"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	v1 "k8s.io/api/core/v1"
)

// type server struct {
// 	router mux
// }

func main() {
	r := mux.NewRouter()
	// s := &server{}
	r.HandleFunc("/create", createPod).Methods("POST")
	r.HandleFunc("/list", getPodsList).Methods("GET")
	//!For LogStream Handler
	r.HandleFunc("/logs/{podname}", logPods).Methods("GET")
	fmt.Println("Server running on port 3004...")
	http.ListenAndServe("localhost:3004", r)
}

var upgrader = websocket.Upgrader{}

type CreatePodRequest struct {
	Name  string `json:"name" mapstructure:"name"`
	Image string `json:"image" mapstructure:"image"`
}

func createPod(w http.ResponseWriter, r *http.Request) {

	var podDetails CreatePodRequest

	//decoding the body response in the podDetails structure
	err := json.NewDecoder(r.Body).Decode(&podDetails)
	// err := mapstructure.Decode(r.Body, &podDetails)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("podDetails: %v\n", podDetails)
	// poddeploy.podCreation
	fmt.Println(podDetails.Name)
	fmt.Println(podDetails.Image)

	//? Call the package for creating a new pod
	err = poddeploy.PodCreation(podDetails.Name, podDetails.Image)
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	} else {

		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		emptyresp, _ := json.Marshal(struct{}{})
		w.Write(emptyresp)
	}

}

func getPodsList(w http.ResponseWriter, r *http.Request) {

	arrayofobjects, err := poddeploy.PodListing()

	if err != nil {
		//Give Error printing details
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		// json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
	} else {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		arrayresponse, _ := json.Marshal(arrayofobjects)
		w.Write(arrayresponse)

	}
}

func logPods(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	podName := vars["podname"]
	// getting the podname for searching directly
	clientset, err := poddeploy.GetKubernetesClientSet()
	if err != nil {
		fmt.Println("K8s client not set", err)
	}
	ctx := context.Background()

	podLogOpts := &v1.PodLogOptions{
		// Container: containerName,	//defaults to only available container
		Follow: true,
	}

	req := clientset.CoreV1().Pods("default").GetLogs(podName, podLogOpts)

	stream, err := req.Stream(ctx)

	defer stream.Close()
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error upgrading to websocket: %v", err), http.StatusInternalServerError)
		return
	}
	defer ws.Close()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error opening stream: %v", err), http.StatusInternalServerError)
		// return nil,err
	}

	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		if err := ws.WriteMessage(websocket.TextMessage, scanner.Bytes()); err != nil {
			break
		}
	}

	// buf := make([]byte, 1024)
	// for {
	// 	n, err := stream.Read(buf)
	// 	if err != nil {
	// 		break
	// 	}
	// 	chunk := buf[:n]
	// 	if err := ws.WriteMessage(websocket.TextMessage, chunk); err != nil {
	// 		break
	// 	}
	// 	// Sleep for a short time to avoid overwhelming the client with data
	// 	time.Sleep(50 * time.Millisecond)
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
