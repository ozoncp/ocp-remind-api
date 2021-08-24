package ocpremindapi

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/ozoncp/ocp-remind-api/internal/models"
	repository "github.com/ozoncp/ocp-remind-api/internal/repo"
	"github.com/ozoncp/ocp-remind-api/pkg"
)

type RemindAPIV1 struct {
	pkg.UnimplementedRemindApiV1Server
	r repository.RemindsRepo
}

func NewRemindAPIV1() (*RemindAPIV1, error) {
	repo, err := repository.NewRemindDBRepository()
	if err != nil {
		return nil, err
	}
	return &RemindAPIV1{
		r: repo,
	}, nil
}

func (api *RemindAPIV1) CreateRemind(ctx context.Context, req *pkg.CreateRemindRequest) (*emptypb.Empty, error) {
	log.Printf("Create request: %v", req)

	reminds := make([]models.Remind, 1)
	reminds[0] = models.Remind{
		Id:       req.RemindId,
		UserId:   req.UserId,
		Deadline: time.Unix(req.Time.Seconds, 0),
		Text:     req.Text,
	}

	err := api.r.Add(ctx, reminds)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
func (api *RemindAPIV1) DescribeRemind(ctx context.Context, req *pkg.DescribeRemindRequest) (*pkg.Remind, error) {
	log.Printf("Descibe request: %v", req)
	describe, err := api.r.Describe(ctx, req.RemindId)
	if err != nil {
		return nil, err
	}

	return &pkg.Remind{
		Text:     describe.Text,
		UserId:   describe.UserId,
		UnixTime: describe.Deadline.Unix(),
		Id:       describe.Id,
	}, nil
}

func (api *RemindAPIV1) ListReminds(ctx context.Context, req *pkg.ListRemindsRequest) (*pkg.ListRemindsResponse, error) {
	log.Printf("List request: %v", req)
	reminds, err := api.r.List(ctx, req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	var resp pkg.ListRemindsResponse
	for _, remind := range reminds {
		resp_remind := &pkg.Remind{
			Text:     remind.Text,
			UserId:   remind.UserId,
			Id:       remind.Id,
			UnixTime: remind.Deadline.Unix(),
		}
		resp.Reminds = append(resp.Reminds, resp_remind)
	}

	return &resp, nil
}

func (api *RemindAPIV1) RemoveRemind(ctx context.Context, req *pkg.RemoveRemindRequest) (*emptypb.Empty, error) {
	log.Printf("Remove request: %v", req)
	err := api.r.Remove(ctx, req.Id, req.UserId)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
