package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2beta3"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2beta3"
)

var (
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	location  = os.Getenv("TASKS_LOCATION")
	queueID   = os.Getenv("TASKS_QUEUE_NAME")
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		client, err := cloudtasks.NewClient(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		payload, _ := json.Marshal(map[string]string{"foo": "bar"})
		task, err := client.CreateTask(ctx, &taskspb.CreateTaskRequest{
			Parent: fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectID, location, queueID),
			Task: &taskspb.Task{
				PayloadType: &taskspb.Task_AppEngineHttpRequest{
					AppEngineHttpRequest: &taskspb.AppEngineHttpRequest{
						HttpMethod: taskspb.HttpMethod_POST,
						AppEngineRouting: &taskspb.AppEngineRouting{
							Service: "tasks",
						},
						RelativeUri: "/foo",
						Headers: map[string]string{
							"Content-Type": "application/json",
						},
						Body: payload,
					},
				},
			},
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "created job: %v", task.Name)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}
