package ocpremindapi

import (
	"context"
	"fmt"
	"time"
	"unsafe"

	"github.com/jackc/pgx/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/ozoncp/ocp-remind-api/internal/metrics"
	"github.com/ozoncp/ocp-remind-api/internal/models"
	"github.com/ozoncp/ocp-remind-api/internal/producer"
	repository "github.com/ozoncp/ocp-remind-api/internal/repo"
	"github.com/ozoncp/ocp-remind-api/internal/utils"
	"github.com/ozoncp/ocp-remind-api/pkg"
)

type RemindAPIV1 struct {
	pkg.UnimplementedRemindApiV1Server
	r repository.RemindsRepo
	p producer.Producer
}

const chunkSize = 5

func NewRemindAPIV1(conn *pgx.Conn) (*RemindAPIV1, error) {
	repo, err := repository.NewRemindDBRepository(conn)
	if err != nil {
		return nil, err
	}

	return &RemindAPIV1{
		r: repo,
		p: producer.NewProducer(),
	}, nil
}

func (api RemindAPIV1) CreateRemind(_ context.Context, req *pkg.CreateRemindRequest) (*emptypb.Empty, error) {
	log.Printf("Create request: %v", req)

	reminds := make([]models.Remind, 1)
	reminds[0] = models.Remind{
		Id:       req.RemindId,
		UserId:   req.UserId,
		Deadline: req.Time.AsTime(),
		Text:     req.Text,
	}

	err := api.r.Add(reminds)
	if err != nil {
		return nil, err
	}
	err = api.p.Send(producer.CreateMessage(producer.Create, producer.EventMessage{
		ID:        reminds[0].Id,
		Action:    producer.Create.String(),
		Timestamp: time.Now().Unix(),
	}))
	if err != nil {
		log.Error().Msgf("failed send to kafka: %v", err)
	}
	metrics.CreateCounterUp()
	return &emptypb.Empty{}, nil
}
func (api RemindAPIV1) DescribeRemind(_ context.Context, req *pkg.DescribeRemindRequest) (*pkg.Remind, error) {
	log.Printf("Descibe request: %v", req)
	describe, err := api.r.Describe(req.RemindId)
	if err != nil {
		return nil, err
	}

	remind := pkg.Remind{
		Text:   describe.Text,
		UserId: describe.UserId,
		Id:     describe.Id,
		Time:   timestamppb.New(describe.Deadline),
	}

	return &remind, nil

}
func (api RemindAPIV1) ListReminds(_ context.Context, req *pkg.ListRemindsRequest) (*pkg.ListRemindsResponse, error) {
	log.Printf("List request: %v", req)
	reminds, err := api.r.List(req.Limit, req.Offset)
	if err != nil {
		return nil, err
	}

	var resp pkg.ListRemindsResponse
	for _, remind := range reminds {
		resp_remind := &pkg.Remind{
			Text:   remind.Text,
			UserId: remind.UserId,
			Id:     remind.Id,
			Time:   timestamppb.New(remind.Deadline),
		}
		resp.Reminds = append(resp.Reminds, resp_remind)
	}

	return &resp, nil
}

func (api RemindAPIV1) RemoveRemind(_ context.Context, req *pkg.RemoveRemindRequest) (*emptypb.Empty, error) {
	log.Printf("Remove request: %v", req)

	err := api.r.Remove(req.Id, req.UserId)
	if err != nil {
		return nil, err
	}
	err = api.p.Send(producer.CreateMessage(producer.Remove,
		producer.EventMessage{
			ID:        req.Id,
			Action:    producer.Create.String(),
			Timestamp: time.Now().Unix(),
		}))
	if err != nil {
		log.Error().Msgf("failed send to kafka: %v", err)
	}
	metrics.RemoveCounterUp()
	return &emptypb.Empty{}, nil
}

func (api RemindAPIV1) MultiCreateRemind(_ context.Context, req *pkg.MultiCreateRemindsRequest) (*emptypb.Empty, error) {
	log.Printf("Multi create request: %v", req)

	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("MultiCreateReminds")
	defer span.Finish()

	reminds := make([]models.Remind, len(req.Reminds))
	for i, remind := range req.Reminds {
		reminds[i].Text = remind.Text
		reminds[i].Id = remind.RemindId
		reminds[i].UserId = remind.UserId
		reminds[i].Deadline = remind.Time.AsTime()
	}

	batched := utils.BatchReminds(reminds, chunkSize)
	for _, v := range batched {
		err := api.r.Add(v)
		if err != nil {
			childSpan := tracer.StartSpan(
				"Size is 0 bytes", opentracing.ChildOf(span.Context()))
			childSpan.Finish()
			return &emptypb.Empty{}, err
		}
		childSpan := tracer.StartSpan(
			fmt.Sprintf("Size is %b bytes", unsafe.Sizeof(v)), opentracing.ChildOf(span.Context()),
		)
		childSpan.Finish()
	}

	return &emptypb.Empty{}, nil
}

func (api RemindAPIV1) UpdateRemind(_ context.Context, req *pkg.Remind) (*emptypb.Empty, error) {
	log.Printf("Update request: %v", req)
	metrics.UpdateCounterUp()
	remind := models.Remind{
		Text:     req.Text,
		Id:       req.Id,
		UserId:   req.UserId,
		Deadline: req.Time.AsTime(),
	}
	err := api.r.Update(remind)
	if err != nil {
		return nil, err
	}

	err = api.p.Send(producer.CreateMessage(producer.Update,
		producer.EventMessage{
			ID:        remind.Id,
			Action:    producer.Create.String(),
			Timestamp: time.Now().Unix(),
		}))
	if err != nil {
		log.Error().Msgf("failed send to kafka :%v", err)
	}
	return &emptypb.Empty{}, nil
}
