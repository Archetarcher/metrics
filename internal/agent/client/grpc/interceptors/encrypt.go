package interceptors

import (
	"context"
	"encoding/json"
	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"github.com/Archetarcher/metrics.git/internal/agent/encryption"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func EncryptInterceptor(ctx context.Context, method string, req interface{},
	reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, config *config.AppConfig,
	opts ...grpc.CallOption) error {

	js, err := json.Marshal(req)

	if err != nil {
		return err
	}

	md := metadata.New(map[string]string{"Encrypted": "1"})
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	encrypted := encryption.EncryptSymmetric(js, config.Session.Key)

	req = encrypted

	return invoker(ctx, method, req, reply, cc, opts...)
}
