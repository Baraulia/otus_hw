package grpcserver

//nolint:depguard
import (
	"context"
	"fmt"

	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/api"
	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/api/grpc/pb"
	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/models"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	service api.ApplicationInterface
	logger  app.Logger
	pb.UnimplementedEventServiceServer
}

func NewServer(service api.ApplicationInterface, logger app.Logger) *Server {
	return &Server{
		service: service,
		logger:  logger,
	}
}

func (s *Server) CreateEvent(ctx context.Context, req *pb.Event) (*pb.CreateEventResponse, error) {
	serviceEvent := models.Event{
		Header:      req.Header,
		Description: req.Description,
		UserID:      req.UserID,
	}

	if req.EventTime != nil {
		eventTime := req.EventTime.AsTime()
		valid := req.EventTime.IsValid()
		switch valid {
		case false:
			s.logger.Error("invalid event time", map[string]interface{}{"event time": eventTime})
			return nil, fmt.Errorf("invalid eventTime:%v", eventTime)
		default:
			serviceEvent.EventTime = eventTime
		}
	}

	if req.FinishEventTime != nil {
		finishTime := req.FinishEventTime.AsTime()
		valid := req.FinishEventTime.IsValid()
		switch valid {
		case false:
			s.logger.Error("invalid finish time", map[string]interface{}{"finish time": finishTime})
			return nil, fmt.Errorf("invalid finishTime:%v", finishTime)
		default:
			serviceEvent.FinishEventTime = &finishTime
		}
	}

	if req.NotificationTime != nil {
		notificationTime := req.NotificationTime.AsTime()
		valid := req.NotificationTime.IsValid()
		switch valid {
		case false:
			s.logger.Error("invalid notification time", map[string]interface{}{"notification time": notificationTime})
			return nil, fmt.Errorf("invalid finishTime:%v", notificationTime)
		default:
			serviceEvent.NotificationTime = &notificationTime
		}
	}

	result, err := s.service.CreateEvent(ctx, serviceEvent)
	if err != nil {
		return nil, fmt.Errorf("error while creating event:%w", err)
	}

	return &pb.CreateEventResponse{Id: result}, nil
}

func (s *Server) UpdateEvent(ctx context.Context, req *pb.Event) (*empty.Empty, error) {
	serviceEvent := models.Event{
		ID:          req.ID,
		Header:      req.Header,
		Description: req.Description,
		UserID:      req.UserID,
	}

	if req.EventTime != nil {
		eventTime := req.EventTime.AsTime()
		valid := req.EventTime.IsValid()
		switch valid {
		case false:
			s.logger.Error("invalid event time", map[string]interface{}{"event time": eventTime})
			return nil, fmt.Errorf("invalid eventTime:%v", eventTime)
		default:
			serviceEvent.EventTime = eventTime
		}
	}

	if req.FinishEventTime != nil {
		finishTime := req.FinishEventTime.AsTime()
		valid := req.FinishEventTime.IsValid()
		switch valid {
		case false:
			s.logger.Error("invalid finish time", map[string]interface{}{"finish time": finishTime})
			return nil, fmt.Errorf("invalid finishTime:%v", finishTime)
		default:
			serviceEvent.FinishEventTime = &finishTime
		}
	}

	if req.NotificationTime != nil {
		notificationTime := req.NotificationTime.AsTime()
		valid := req.NotificationTime.IsValid()
		switch valid {
		case false:
			s.logger.Error("invalid notification time", map[string]interface{}{"notification time": notificationTime})
			return nil, fmt.Errorf("invalid finishTime:%v", notificationTime)
		default:
			serviceEvent.NotificationTime = &notificationTime
		}
	}

	err := s.service.UpdateEvent(ctx, serviceEvent)
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*empty.Empty, error) {
	err := s.service.DeleteEvent(ctx, req.Id)
	if err != nil {
		return &empty.Empty{}, err
	}

	return &empty.Empty{}, nil
}

func (s *Server) GetListEvents(ctx context.Context, req *pb.GetListEventsRequest) (*pb.GetListEventsResponse, error) {
	if req.Start == nil {
		return nil, fmt.Errorf("not specified start date")
	}

	start := req.Start.AsTime()
	valid := req.Start.IsValid()
	switch valid {
	case false:
		s.logger.Error("invalid start time", map[string]interface{}{"start time": start})
		return nil, fmt.Errorf("invalid startTime:%v", start)
	default:
	}

	var events []models.Event
	var err error
	switch req.AmountDays {
	case 0:
		events, err = s.service.GetListEventsDuringDay(ctx, start)
		if err != nil {
			return nil, err
		}
	default:
		events, err = s.service.GetListEventsDuringFewDays(ctx, start, int(req.AmountDays))
		if err != nil {
			return nil, err
		}
	}

	pbEvents := make([]*pb.Event, 0, len(events))

	for _, event := range events {
		pbEvents = append(pbEvents, convert(event))
	}

	return &pb.GetListEventsResponse{Events: pbEvents}, nil
}

func convert(event models.Event) *pb.Event {
	pbEvent := &pb.Event{
		ID:          event.ID,
		Header:      event.Header,
		Description: event.Description,
		UserID:      event.UserID,
		EventTime:   timestamppb.New(event.EventTime),
	}

	if event.FinishEventTime != nil {
		pbEvent.FinishEventTime = timestamppb.New(*event.FinishEventTime)
	}

	if event.NotificationTime != nil {
		pbEvent.NotificationTime = timestamppb.New(*event.NotificationTime)
	}

	return pbEvent
}
