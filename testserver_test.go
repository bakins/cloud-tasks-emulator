package emulator

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	cloudtasks "cloud.google.com/go/cloudtasks/apiv2"
	"github.com/stretchr/testify/require"
	"google.golang.org/api/option"
	taskspb "google.golang.org/genproto/googleapis/cloud/tasks/v2"
	"google.golang.org/grpc"
)

func TestTestServer(t *testing.T) {
	svr := NewTestServer()
	defer svr.Close()

	require.NotEmpty(t, svr.Address())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	conn, err := grpc.DialContext(ctx, svr.Address(), grpc.WithInsecure())
	require.NoError(t, err)

	client, err := cloudtasks.NewClient(ctx, option.WithGRPCConn(conn))
	require.NoError(t, err)

	createQueueRequest := taskspb.CreateQueueRequest{
		Parent: formattedParent,
		Queue:  newQueue(formattedParent, "test"),
	}

	queue, err := client.CreateQueue(ctx, &createQueueRequest)
	require.NoError(t, err)

	require.Equal(t, createQueueRequest.Queue.Name, queue.Name)

	done := make(chan struct{})
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		close(done)
	}

	target := httptest.NewServer(http.HandlerFunc(handler))
	defer target.Close()

	createTaskRequest := taskspb.CreateTaskRequest{
		Parent: queue.GetName(),
		Task: &taskspb.Task{
			MessageType: &taskspb.Task_HttpRequest{
				HttpRequest: &taskspb.HttpRequest{
					Url: target.URL,
				},
			},
		},
	}

	createdTask, err := client.CreateTask(ctx, &createTaskRequest)
	require.NoError(t, err)

	require.Equal(t, taskspb.HttpMethod_POST, createdTask.GetHttpRequest().GetHttpMethod())
	require.EqualValues(t, 0, createdTask.GetDispatchCount())

	select {
	case <-done:
	case <-ctx.Done():
		t.Errorf("failed to receive task")
	}
}
