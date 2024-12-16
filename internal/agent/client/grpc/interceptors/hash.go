package interceptors

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/Archetarcher/metrics.git/internal/agent/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func HashInterceptor(ctx context.Context, method string, req interface{},
	reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, config *config.AppConfig,
	opts ...grpc.CallOption) error {

	h := hmac.New(sha256.New, []byte(config.Key))

	js, err := json.Marshal(req)

	if err != nil {
		return err
	}

	h.Write(js)
	hash := h.Sum(nil)

	md := metadata.New(map[string]string{"HashSHA256g": hex.EncodeToString(hash)})
	ctx = metadata.NewOutgoingContext(context.Background(), md)

	return invoker(ctx, method, req, reply, cc, opts...)
}
