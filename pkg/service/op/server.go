package op

import (
	"context"
	"database/sql"
	"github.com/x0y14/msnger/pkg/db"
	pb "github.com/x0y14/msnger/pkg/protobuf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"time"
)

type ServiceServer struct {
	pb.UnimplementedOpServiceServer
}

func (s *ServiceServer) FetchOps(req *pb.FetchOpsRequest, stream pb.OpService_FetchOpsServer) error {

	// LastRevisionId -> LRId

	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return status.Errorf(codes.Internal, "failed to get header.")
	}

	LRIdHoldByUser := req.LastRevisionId
	log.Printf("[Op] FetchOps Will: Get userId from stream context.")
	userId := md.Get("userId")[0]
	log.Printf("%v", userId)

	for {
		// DBから、ユーザーの最終RevisionIdを取得。
		LRIdStoredOnDB, err := db.GetLastOpRevision(userId)
		if err != nil {
			log.Printf("[Op] FetchOps Err: %v", err)
			return status.Errorf(codes.Internal, "failed to get lastRevisionId from db.")
		}
		// 一致してたら、新しいデータはない。
		if LRIdHoldByUser == LRIdStoredOnDB {
			time.Sleep(time.Millisecond * 500)
			continue
		}
		// 一致しなかった
		operations, err := db.GetOpsBiggerThan(userId, LRIdHoldByUser)
		if err == sql.ErrNoRows {
			log.Printf("LRIdHoldByUser = %v, LRIdStoredOnDB = %v, But No Matche Ops\n", LRIdHoldByUser, LRIdStoredOnDB)
			continue
		}
		if err != nil {
			log.Printf("[Op] FetchOps Err: %v", err)
			return status.Errorf(codes.Internal, "failed to get ops.")
		}
		for _, operation := range operations {
			err = stream.Send(operation)
			if err != nil {
				log.Printf("[Op] FetchOps Err: %v", err)
				return status.Errorf(codes.Internal, "failed to stream ops.")
			}
			LRIdHoldByUser = operation.RevisionId
		}
	}
}
func (s *ServiceServer) SendOp(ctx context.Context, req *pb.SendOpRequest) (*pb.Operation, error) {
	userId := ctx.Value("userId").(string)
	log.Printf("[Op] SendOp UserId: %v", userId)
	ac, err := db.SelectAccountWithId(&db.SelectAccountReq{Id: userId})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user data.")
	}
	if !ac.IsAdmin {
		return nil, status.Errorf(codes.PermissionDenied, "failed to send op by normal user")
	}

	err = db.StoreOp(req.Op)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to store op by admin.")
	}

	log.Printf("[Op] SendOp Success")

	// todo : return real op data
	return &pb.Operation{}, nil
}
