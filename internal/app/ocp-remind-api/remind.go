package ocpremindapi

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/ozoncp/ocp-remind-api/pkg"
)

type RemindAPIV1 struct {
	pkg.UnimplementedRemindApiV1Server
}

func NewRemindAPIV1() *RemindAPIV1 {
	return &RemindAPIV1{}
}

func (api *RemindAPIV1) CreateV1(_ context.Context, req *pkg.CreateRemindRequest) (*pkg.Remind, error) {
	log.Printf("Create request: %v", req)
	return &pkg.Remind{}, nil //nolint:exhaustivestruct
}

func (api *RemindAPIV1) DescribeV1(_ context.Context, req *pkg.DescribeRemindRequest) (*pkg.Remind, error) {
	log.Printf("Descibe request: %v", req)
	return &pkg.Remind{}, nil
}

func (api *RemindAPIV1) ListV1(_ context.Context, req *pkg.ListRemindsResponse) (*pkg.ListRemindsResponse, error) {
	log.Printf("List request: %v", req)
	return &pkg.ListRemindsResponse{}, nil
}

func (api *RemindAPIV1) RemoveV1(_ context.Context, req *pkg.RemoveRemindRequest) (*pkg.RemoveRemindResponse, error) {
	log.Printf("Remove request: %v", req)
	return &pkg.RemoveRemindResponse{}, nil
}
