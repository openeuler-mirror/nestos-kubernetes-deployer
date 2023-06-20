/*
Copyright 2023 KylinSoft  Co., Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
// connection between client and server
package connection

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"

	pb "housekeeper.io/pkg/connection/proto"
)

type Client struct {
	socketAddress string
	client        pb.UpgradeClusterClient
}

type PushInfo struct {
	OSImageURL   string
	OSVersion    string
	KubeVersion  string
	ControlPlane bool
}

// Create a grpc channel
func New(socketAddr string) (*Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	bc := backoff.DefaultConfig
	bc.MaxDelay = 5 * time.Second

	connection, err := grpc.DialContext(ctx, socketAddr, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithConnectParams(grpc.ConnectParams{Backoff: bc}))
	if err != nil {
		return nil, err
	}
	return &Client{socketAddress: socketAddr, client: pb.NewUpgradeClusterClient(connection)}, nil
}

// send update requests
func (c *Client) UpgradeKubeSpec(pushInfo *PushInfo) error {
	_, err := c.client.Upgrade(context.Background(),
		&pb.UpgradeRequest{
			KubeVersion:  pushInfo.KubeVersion,
			OsImageUrl:   pushInfo.OSImageURL,
			OsVersion:    pushInfo.OSVersion,
			ControlPlane: pushInfo.ControlPlane,
		})
	return err
}
