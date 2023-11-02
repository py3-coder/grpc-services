package main

import (
	"flag"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"

	pb "grpc-services/proto"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Info("did not connect: %v", err)
	}

	defer conn.Close()
	client := pb.NewUserServiceClient(conn)

	r := gin.Default()
	r.GET("/users", func(ctx *gin.Context) {
		res, err := client.GetUsersByIds(ctx, &pb.UserIdsRequest{})
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"users": res.Users,
		})
	})
	r.GET("/user/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		i, err := strconv.ParseInt(id, 10, 32)
		if err != nil {
			log.Info("Error", err)
		}
		res, err := client.GetUserById(ctx, &pb.UserRequest{UserId: int32(i)})
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"user": res.User,
		})
	})

	r.Run(":5000")

}
