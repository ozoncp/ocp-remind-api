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

func (api *RemindAPIV1) CreateV1(ctx context.Context, req *pkg.CreateRemindRequest) *pkg.Remind {
	log.Printf("Create request: %v", req)
	return &pkg.Remind{} //nolint:exhaustivestruct
}

func (api *RemindAPIV1) DescribeV1(ctx context.Context, req *pkg.DescribeRemindRequest) *pkg.Remind {
	log.Printf("Descibe request: %v", req)
	return &pkg.Remind{}
}

func (api *RemindAPIV1) ListV1(ctx context.Context, req *pkg.ListRemindsResponse) *pkg.ListRemindsResponse {
	log.Printf("List request: %v", req)
	return &pkg.ListRemindsResponse{}
}

func (api *RemindAPIV1) RemoveV1(ctx context.Context, req *pkg.RemoveRemindRequest) *pkg.RemoveRemindResponse {
	log.Printf("Remove request: %v", req)
	return &pkg.RemoveRemindResponse{}
}
