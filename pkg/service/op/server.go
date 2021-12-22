package op

import (
	"context"
	pb "github.com/x0y14/msnger/pkg/protobuf"
	"log"
)

type ServiceServer struct {
	pb.UnimplementedOpServiceServer
}

func (s *ServiceServer) FetchOps(req *pb.FetchOpsRequest, stream pb.OpService_FetchOpsServer) error {
	var err error
	base := req.LastRevisionId

	log.Printf("receive FetchOps Request: %v\n", base)

	for i := 0; i < 20; i++ {
		err = stream.Send(&pb.Operation{
			RevisionId: base + int32(i) + 1,
			Type:       pb.OperationType_NOOP,
		})
	}
	return err
}
func (s *ServiceServer) SendOp(ctx context.Context, req *pb.SendOpRequest) (*pb.SendOpResult, error) {
	return &pb.SendOpResult{}, nil
}
