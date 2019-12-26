package server

import (
	"context"
	"github.com/golang/protobuf/ptypes"
	"github.com/labstack/gommon/log"
	"github.com/reddaemon/calendargrpcsql/internal/domain/grpc/service"
	"github.com/reddaemon/calendargrpcsql/models/models"
	eventpb "github.com/reddaemon/calendargrpcsql/protofiles"
	"time"
)

type Server struct {
	*service.EventUseCase
}

func NewServer(usecase *service.EventUseCase) *Server {
	return &Server{
		EventUseCase: usecase,
	}
}

func (s *Server) Create(ctx context.Context, req *eventpb.CreateRequest) (*eventpb.CreateResponse, error) {
	success := false
	event := s.unmarshalPbEvent(req.Event)
	resEvent, err := s.EventUseCase.Create(ctx, &event)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	} else {
		success = true
	}

	res := eventpb.CreateResponse{
		Id:      uint32(resEvent.Id),
		Success: success,
	}
	return &res, nil
}

func (s *Server) Read(ctx context.Context, req *eventpb.ReadRequest) (*eventpb.ReadResponse, error) {
	event, err := s.EventUseCase.Read(ctx, uint64(req.Id))

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	res := eventpb.ReadResponse{
		Event: s.marshalEvent(event),
	}

	return &res, nil
}

func (s *Server) Update(ctx context.Context, req *eventpb.UpdateRequest) (*eventpb.UpdateResponse, error) {
	event := s.unmarshalPbEvent(req.Event)
	success, err := s.EventUseCase.Update(ctx, &event, uint64(req.Id))

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return &eventpb.UpdateResponse{
		Success: success,
	}, nil
}

func (s *Server) Delete(ctx context.Context, req *eventpb.DeleteRequest) (*eventpb.DeleteResponse, error) {
	success, err := s.EventUseCase.Delete(ctx, uint64(req.GetId()))
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &eventpb.DeleteResponse{
		Success: success,
	}, nil
}

func (s *Server) unmarshalPbEvent(e *eventpb.Event) models.Event {
	eventStruct := models.Event{
		Id:          uint64(e.GetId()),
		Title:       e.GetTitle(),
		Description: e.GetDescription(),
		Date:        ptypes.TimestampString(e.GetDate()),
	}
	return eventStruct
}

func (s *Server) marshalEvent(e *models.Event) *eventpb.Event {
	t, err := time.Parse(time.RFC3339, e.Date)
	if err != nil {
		log.Error(err.Error())
	}

	ts, err := ptypes.TimestampProto(t)
	if err != nil {
		log.Error(err.Error())
	}
	pbEvent := eventpb.Event{
		Id:          uint32(e.Id),
		Title:       e.Title,
		Description: e.Description,
		Date:        ts,
	}
	return &pbEvent
}
