package server

import (
	"context"
	"fmt"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/omny-n/pi_task/models"
	pb "github.com/omny-n/pi_task/pb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userServer struct {
	pb.UnimplementedUserServiceServer
}

func newUserServer() pb.UserServiceServer {
	return new(userServer)
}

func (s *userServer) CreateUser(ctx context.Context, req *pb.CreateUserReq) (*pb.CreateUserRes, error) {
	db := MongCtx
	data := &models.UserStruct{
		ID:        primitive.NewObjectID(),
		FirstName: req.GetFirstname(),
		LastName:  req.GetLastname(),
		Age:       int(req.GetAge()),
		Email:     req.GetEmail(),
	}

	insertUser, err := db.InsertOne(ctx, data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}
	if oid, ok := insertUser.InsertedID.(primitive.ObjectID); ok {
		return &pb.CreateUserRes{Id: oid.Hex()}, nil
	}
	return nil, status.Errorf(
		codes.InvalidArgument, "Incorrect data",
	)
}

func (s *userServer) UpdateUser(ctx context.Context, req *pb.UpdateUserReq) (*pb.UpdateUserRes, error) {
	User := req.GetUser()

	objectId, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Incorrect ID: %v", err))
	}
	db := *MongCtx

	update := bson.M{
		"firstname": User.GetFirstname(),
		"lastname":  User.GetLastname(),
		"age":       User.GetAge(),
		"email":     User.GetEmail(),
	}

	filter := bson.M{"_id": objectId}

	result := db.FindOneAndUpdate(ctx, filter, bson.M{"$set": update}, options.FindOneAndUpdate().SetReturnDocument(1))
	decoded := &models.UserStruct{}
	err = result.Decode(&decoded)
	if err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Could not find user with supplied ID: %v", err),
		)
	}
	return &pb.UpdateUserRes{
		User: &pb.User{
			Id:        decoded.ID.Hex(),
			Firstname: decoded.FirstName,
			Lastname:  decoded.LastName,
			Age:       int32(decoded.Age),
			Email:     decoded.Email,
		},
	}, nil
}

func (s *userServer) ReadUser(ctx context.Context, req *pb.ReadUserReq) (*pb.ReadUserRes, error) {
	objectId, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Incorrect ID: %v", err))
	}
	db := *MongCtx
	result := db.FindOne(ctx, bson.M{"_id": objectId})
	data := &models.UserStruct{}
	if err := result.Decode(&data); err != nil {
		return nil, status.Errorf(
			codes.NotFound, fmt.Sprintf("Could not find User with ID %s: %v", req.GetId(), err),
		)
	}
	response := &pb.ReadUserRes{
		User: &pb.User{
			Id:        data.ID.Hex(),
			Firstname: data.FirstName,
			Lastname:  data.LastName,
			Age:       int32(data.Age),
			Email:     data.Email,
		},
	}
	return response, nil
}

func (s *userServer) ListUsers(ctx context.Context, req *empty.Empty) (*pb.ListUsersRes, error) {
	data := &models.UserStruct{}
	db := *MongCtx
	cursor, err := db.Find(context.Background(), bson.M{})
	if cursor == nil {
		status.New(codes.FailedPrecondition, "No users in db :(")
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Unknown internal error: %v", err))
	}
	res := &pb.ListUsersRes{}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		err := cursor.Decode(data)
		if err != nil {
			return nil, status.Errorf(codes.Unavailable, fmt.Sprintf("Could not decode data: %v", err))
		}
		res.Users = append(res.Users, &pb.User{
			Id:        data.ID.Hex(),
			Firstname: data.FirstName,
			Lastname:  data.LastName,
			Age:       int32(data.Age),
			Email:     data.Email,
		})
	}

	if err := cursor.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Unkown cursor error: %v", err))
	}
	return res, nil
}

func (s *userServer) DeleteUser(ctx context.Context, req *pb.DeleteUserReq) (*pb.DeleteUserRes, error) {
	objectId, err := primitive.ObjectIDFromHex(req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Incorrect ID: %v", err))
	}
	db := *MongCtx

	res, err := db.DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find/delete User with ID %s: %v", req.GetId(), err))
	}
	if res.DeletedCount == 0 {
		return &pb.DeleteUserRes{
			Success: false,
		}, nil
	}
	return &pb.DeleteUserRes{
		Success: true,
	}, nil
}
