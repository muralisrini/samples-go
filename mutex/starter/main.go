package main

import (
	"context"
	"log"

	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"

	"github.com/temporalio/samples-go/mutex"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	// This workflow ID can be user business logic identifier as well.
	resourceID := uuid.New()
	for i := 1; i < 10; i++ {
		workflow1Options := client.StartWorkflowOptions{
			ID:        "SampleWorkflow1WithMutex_" + uuid.New(),
			TaskQueue: "mutex",
		}

		workflow2Options := client.StartWorkflowOptions{
			ID:        "SampleWorkflow2WithMutex_" + uuid.New(),
			TaskQueue: "mutex",
		}

		ctx1 := context.Background()
		we1, err := c.ExecuteWorkflow(ctx1, workflow1Options, mutex.SampleWorkflowWithMutex, resourceID)
		if err != nil {
			log.Fatalln("Unable to execute workflow1", err)
		} else {
			log.Println("Started workflow1", "WorkflowID", we1.GetID(), "RunID", we1.GetRunID())
		}

		ctx2 := context.Background()
		we2, err := c.ExecuteWorkflow(ctx2, workflow2Options, mutex.SampleWorkflowWithMutex, resourceID)
		if err != nil {
			log.Fatalln("Unable to execute workflow2", err)
		} else {
			log.Println("Started workflow2", "WorkflowID", we2.GetID(), "RunID", we2.GetRunID())
		}

		err = we2.Get(ctx1, nil)
		if err != nil {
			log.Fatalln("Unable to get workflow1", err)
		}

		err = we2.Get(ctx2, nil)
		if err != nil {
			log.Fatalln("Unable to get workflow2", err)
		}
	}
}
