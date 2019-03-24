package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2beta3"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2beta3"
)

var (
	projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	service   = os.Getenv("GAE_SERVICE")
	location  = os.Getenv("TASKS_LOCATION")
	queueID   = os.Getenv("TASKS_QUEUE_NAME")
)

func main() {
	http.HandleFunc("/cron/hours", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-Appengine-Cron") == "" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		ctx := r.Context()
		client, err := cloudtasks.NewClient(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer client.Close()

		var (
			tasks        = []string{"/tasks/task1", "/tasks/task2"}
			createdTasks = make([]string, 2)
		)

		for i, task := range tasks {
			createdTask, err := client.CreateTask(ctx, &taskspb.CreateTaskRequest{
				Parent: fmt.Sprintf("projects/%s/locations/%s/queues/%s", projectID, location, queueID),
				Task: &taskspb.Task{
					PayloadType: &taskspb.Task_AppEngineHttpRequest{
						AppEngineHttpRequest: &taskspb.AppEngineHttpRequest{
							HttpMethod: taskspb.HttpMethod_GET,
							AppEngineRouting: &taskspb.AppEngineRouting{
								Service: service,
							},
							RelativeUri: task,
						},
					},
				},
			})
			if err != nil {
				log.Printf("failed to create task: %v", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			createdTasks[i] = createdTask.Name
		}

		fmt.Fprintf(w, "Task created: %v", createdTasks)
	})

	http.HandleFunc("/tasks/task1", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Exec TASK1")
	})
	http.HandleFunc("/tasks/task2", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Exec TASK1")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}
