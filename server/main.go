package main

import (
	"context"
	"grpc-services/db"
	"grpc-services/model"
	pb "grpc-services/proto"
	"net"
	"os"
	"os/signal"
	"time"

	formatter "github.com/antonfisher/nested-logrus-formatter"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc"
)

var (
	datab = db.MongodbDBInstance()
	//pd    = ProductServices(datab)
)

type server struct {
	pb.UnimplementedUserServiceServer
}

func init() {
	log.SetFormatter(&formatter.Formatter{
		HideKeys:      true,
		NoFieldsSpace: false,
		TrimMessages:  true,
	})

}
func main() {
	log.Info("Starting server on port :50051...")
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Unable to listen on port :50051: %v", err)
	}
	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	pb.RegisterUserServiceServer(s, &server{})
	db.CreateNewDBConnection()
	go func() {
		err := s.Serve(listener)
		if err != nil {
			log.Info("Error Starting server:", err)
			db.DisconnectToMongoDB()
			os.Exit(1)
		}
	}()
	log.Info("Server succesfully started on port :50051")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)
	sig := <-c
	log.Info("Got signal:", sig)
	log.Info("Stopping the server...")
	s.Stop()
	listener.Close()

}

func (s *server) GetUserById(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	id := req.UserId
	collection := datab.ConnectToMongoDB().Collection("user-service-db")
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result := collection.FindOne(ctx, bson.M{"id": id})
	data := model.User{}
	if err := result.Decode(&data); err != nil {
		log.Info("Error Found::", err)
	}
	response := &pb.UserResponse{
		User: &pb.User{
			Id:     req.UserId,
			Fname:  data.FName,
			City:   data.City,
			Phone:  data.Phone,
			Height: data.Height,
		},
	}
	return response, nil
}
func (s *server) GetUsersByIds(ctx context.Context, req *pb.UserIdsRequest) (*pb.UserListResponse, error) {
	data := &model.User{}
	var userList []*pb.User
	collection := datab.ConnectToMongoDB().Collection("user-service-db")
	_, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		log.Info("Error Found::")
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		err := cursor.Decode(data)
		if err != nil {
			log.Info("Error Found:::")
		}
		response := &pb.UserResponse{
			User: &pb.User{
				Id:     int32(data.ID),
				Fname:  data.FName,
				City:   data.City,
				Phone:  data.Phone,
				Height: data.Height,
			},
		}
		userList = append(userList, response.User)
	}
	if err := cursor.Err(); err != nil {
		log.Info("Error Found::", err)
	}
	return &pb.UserListResponse{Users: userList}, nil
}
