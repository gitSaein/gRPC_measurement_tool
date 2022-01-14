package interceptors

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"google.golang.org/grpc"
)

func CheckHttpHeader(ctx context.Context) {

	log.Println("================== HEADER start ==================")

	method, _ := grpc.Method(ctx)
	log.Printf("method: %s", method)
	md, _ := metadata.FromIncomingContext(ctx)
	p, ok := peer.FromContext(ctx)
	if !ok {
		log.Print("error")
	}
	log.Printf("%v", p)

	log.Printf("metadata: %v", md)
	log.Println("==================  HEADER end  ==================")

}

func UnaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()

	CheckHttpHeader(ctx)

	m, err := handler(ctx, req)
	if err != nil {
		log.Printf(" [error] server interceptor handler: %v", err)
	}

	log.Printf("Post Proc Message: %s", m)

	elapsed := time.Since(start)
	log.Printf("take time - %s", elapsed)
	return m, err
}
