package server

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/golang/glog"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/omny-n/pi_task/models"
	pb "github.com/omny-n/pi_task/pb"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

var MongCtx *mongo.Collection

func Run(ctx context.Context, network, address string) error {
	l, err := net.Listen(network, address)
	if err != nil {
		return err
	}
	defer func() {
		if err := l.Close(); err != nil {
			glog.Errorf("Failed to close %s %s: %v", network, address, err)
		}
	}()

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, newUserServer())
	db, mongoCtx, err := models.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	MongCtx = models.NewUserCollection(db)
	go func() {
		defer s.GracefulStop()
		defer db.Disconnect(mongoCtx)
		<-ctx.Done()
	}()
	return s.Serve(l)
}

func RunInProcessGateway(ctx context.Context, addr string, opts ...runtime.ServeMuxOption) error {
	mux := runtime.NewServeMux(opts...)

	pb.RegisterUserServiceHandlerServer(ctx, mux, newUserServer())

	s := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		glog.Infof("Shutting down the http gateway server")
		if err := s.Shutdown(context.Background()); err != nil {
			glog.Errorf("Failed to shutdown http gateway server: %v", err)
		}
	}()

	if err := s.ListenAndServe(); err != http.ErrServerClosed {
		glog.Errorf("Failed to listen and serve: %v", err)
		return err
	}
	return nil

}
