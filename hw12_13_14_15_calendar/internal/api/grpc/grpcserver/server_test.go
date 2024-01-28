package grpcserver

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Baraulia/X-Labs_Test/pkg/logger"
	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/api/grpc/pb"
	mockservice "github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/api/mocks"
	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreateEvent(t *testing.T) {
	type mockBehavior func(s *mockservice.MockApplicationInterface, dto models.Event)
	logg, err := logger.GetLogger("INFO")
	require.NoError(t, err)
	ctx := context.Background()
	newUserUUID := uuid.New().String()
	newEventUUID := uuid.New().String()

	testTable := []struct {
		name           string
		inputData      *pb.Event
		convertData    models.Event
		expectedResult *pb.CreateEventResponse
		mockBehavior   mockBehavior
		expectedError  bool
	}{
		{
			name: "successful",
			inputData: &pb.Event{
				Header:      "header",
				Description: "desc",
				UserID:      newUserUUID,
			},
			convertData: models.Event{
				Header:      "header",
				Description: "desc",
				UserID:      newUserUUID,
			},
			expectedResult: &pb.CreateEventResponse{Id: newEventUUID},
			mockBehavior: func(s *mockservice.MockApplicationInterface, dto models.Event) {
				s.EXPECT().CreateEvent(ctx, dto).Return(newEventUUID, nil)
			},
			expectedError: false,
		},
		{
			name: "error from service",
			inputData: &pb.Event{
				Header:      "header",
				Description: "desc",
				UserID:      newUserUUID,
			},
			convertData: models.Event{
				Header:      "header",
				Description: "desc",
				UserID:      newUserUUID,
			},
			mockBehavior: func(s *mockservice.MockApplicationInterface, dto models.Event) {
				s.EXPECT().CreateEvent(ctx, dto).Return("", errors.New("service error"))
			},
			expectedError: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			service := mockservice.NewMockApplicationInterface(c)
			testCase.mockBehavior(service, testCase.convertData)
			server := NewServer(service, logg)

			response, err := server.CreateEvent(ctx, testCase.inputData)
			if testCase.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, testCase.expectedResult, response)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	type mockBehavior func(s *mockservice.MockApplicationInterface, dto models.Event)
	logg, err := logger.GetLogger("INFO")
	require.NoError(t, err)
	ctx := context.Background()
	newUserUUID := uuid.New().String()
	newEventUUID := uuid.New().String()

	testTable := []struct {
		name          string
		inputData     *pb.Event
		convertData   models.Event
		mockBehavior  mockBehavior
		expectedError bool
	}{
		{
			name: "successful",
			inputData: &pb.Event{
				ID:          newEventUUID,
				Header:      "header",
				Description: "desc",
				UserID:      newUserUUID,
			},
			convertData: models.Event{
				ID:          newEventUUID,
				Header:      "header",
				Description: "desc",
				UserID:      newUserUUID,
			},
			mockBehavior: func(s *mockservice.MockApplicationInterface, dto models.Event) {
				s.EXPECT().UpdateEvent(ctx, dto).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "server error",
			inputData: &pb.Event{
				ID:          newEventUUID,
				Header:      "header",
				Description: "desc",
				UserID:      newUserUUID,
			},
			convertData: models.Event{
				ID:          newEventUUID,
				Header:      "header",
				Description: "desc",
				UserID:      newUserUUID,
			},
			mockBehavior: func(s *mockservice.MockApplicationInterface, dto models.Event) {
				s.EXPECT().UpdateEvent(ctx, dto).Return(errors.New("server error"))
			},

			expectedError: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			service := mockservice.NewMockApplicationInterface(c)
			testCase.mockBehavior(service, testCase.convertData)
			server := NewServer(service, logg)

			_, err = server.UpdateEvent(ctx, testCase.inputData)
			if testCase.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	type mockBehavior func(s *mockservice.MockApplicationInterface, id string)
	logg, err := logger.GetLogger("INFO")
	require.NoError(t, err)
	ctx := context.Background()
	newUUID := uuid.New().String()

	testTable := []struct {
		name          string
		inputData     *pb.DeleteEventRequest
		id            string
		mockBehavior  mockBehavior
		expectedError bool
	}{
		{
			name: "successful",
			inputData: &pb.DeleteEventRequest{
				Id: newUUID,
			},
			id: newUUID,
			mockBehavior: func(s *mockservice.MockApplicationInterface, id string) {
				s.EXPECT().DeleteEvent(ctx, id).Return(nil)
			},
			expectedError: false,
		},
		{
			name: "error from service",
			inputData: &pb.DeleteEventRequest{
				Id: newUUID,
			},
			id: newUUID,
			mockBehavior: func(s *mockservice.MockApplicationInterface, id string) {
				s.EXPECT().DeleteEvent(ctx, id).Return(errors.New("service error"))
			},
			expectedError: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			service := mockservice.NewMockApplicationInterface(c)
			testCase.mockBehavior(service, testCase.id)
			server := NewServer(service, logg)

			_, err = server.DeleteEvent(ctx, testCase.inputData)
			if testCase.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetUsers(t *testing.T) {
	type mockBehavior func(s *mockservice.MockApplicationInterface, start time.Time, amountDays int)
	logg, err := logger.GetLogger("INFO")
	require.NoError(t, err)
	ctx := context.Background()
	testTimeForJSON := time.Now().In(time.UTC).Format(time.RFC3339Nano)
	testTime, _ := time.Parse(time.RFC3339Nano, testTimeForJSON)
	testTimePB := timestamppb.New(testTime)
	newUUID1 := uuid.New().String()
	newUUID2 := uuid.New().String()
	newUserUUID := uuid.New().String()

	testTable := []struct {
		name           string
		inputData      *pb.GetListEventsRequest
		start          time.Time
		amountDays     int
		expectedResult *pb.GetListEventsResponse
		mockBehavior   mockBehavior
		expectedError  bool
	}{
		{
			name: "successful",
			inputData: &pb.GetListEventsRequest{
				Start:      testTimePB,
				AmountDays: 2,
			},
			start:      testTime,
			amountDays: 2,
			expectedResult: &pb.GetListEventsResponse{
				Events: []*pb.Event{
					{
						ID:          newUUID1,
						Header:      "header",
						Description: "desc",
						UserID:      newUserUUID,
						EventTime:   testTimePB,
					},
					{
						ID:          newUUID2,
						Header:      "header",
						Description: "desc",
						UserID:      newUserUUID,
						EventTime:   testTimePB,
					},
				},
			},
			mockBehavior: func(s *mockservice.MockApplicationInterface, start time.Time, amountDays int) {
				s.EXPECT().GetListEventsDuringFewDays(ctx, start, amountDays).Return([]models.Event{
					{
						ID:          newUUID1,
						Header:      "header",
						Description: "desc",
						UserID:      newUserUUID,
						EventTime:   testTime,
					},
					{
						ID:          newUUID2,
						Header:      "header",
						Description: "desc",
						UserID:      newUserUUID,
						EventTime:   testTime,
					},
				}, nil)
			},
			expectedError: false,
		},
		{
			name: "error from service",
			inputData: &pb.GetListEventsRequest{
				Start:      testTimePB,
				AmountDays: 2,
			},
			start:      testTime,
			amountDays: 2,
			mockBehavior: func(s *mockservice.MockApplicationInterface, start time.Time, amountDays int) {
				s.EXPECT().GetListEventsDuringFewDays(ctx, start, amountDays).Return(nil, errors.New("service error"))
			},
			expectedError: true,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			service := mockservice.NewMockApplicationInterface(c)
			testCase.mockBehavior(service, testCase.start, testCase.amountDays)
			server := NewServer(service, logg)

			response, err := server.GetListEvents(ctx, testCase.inputData)
			if testCase.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				require.Equal(t, testCase.expectedResult, response)
			}
		})
	}
}
