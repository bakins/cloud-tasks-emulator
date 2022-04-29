package emulator

import (
	"context"
	"testing"

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

	conn, err := grpc.Dial(svr.Address(), grpc.WithInsecure())
	require.NoError(t, err)

	client, err := cloudtasks.NewClient(context.Background(), option.WithGRPCConn(conn))
	require.NoError(t, err)

	createQueueRequest := taskspb.CreateQueueRequest{
		Parent: formattedParent,
		Queue:  newQueue(formattedParent, "test"),
	}

	queue, err := client.CreateQueue(context.Background(), &createQueueRequest)
	require.NoError(t, err)

	require.Equal(t, createQueueRequest.Queue.Name, queue.Name)
}
