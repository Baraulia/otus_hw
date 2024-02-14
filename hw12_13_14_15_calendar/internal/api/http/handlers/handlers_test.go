package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	mockservice "github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/api/mocks"
	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/models"
	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/pkg/logger"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetListEvents(t *testing.T) {
	type mockBehavior func(s *mockservice.MockApplicationInterface)
	testEventID1 := uuid.New().String()
	testEventID2 := uuid.New().String()
	testUserID := uuid.New().String()
	testTime := time.Now()
	testParamTime := testTime.Format(time.DateOnly)
	testDateOnly, _ := time.Parse(time.DateOnly, testParamTime)

	testTable := []struct {
		name                string
		mockBehavior        mockBehavior
		getParams           string
		expectedStatusCode  int
		expectedRequestBody string
	}{
		{
			name: "OK by one day",
			mockBehavior: func(s *mockservice.MockApplicationInterface) {
				s.EXPECT().GetListEventsDuringDay(gomock.Any(), testDateOnly).Return([]models.Event{
					{
						ID:          testEventID1,
						Header:      "test1",
						Description: "testDescription1",
						UserID:      testUserID,
						EventTime:   testTime,
					},
					{
						ID:          testEventID2,
						Header:      "test2",
						Description: "testDescription2",
						UserID:      testUserID,
						EventTime:   testTime,
					},
				}, nil)
			},
			getParams:          fmt.Sprintf("?start=%s", testParamTime),
			expectedStatusCode: 200,
			expectedRequestBody: fmt.Sprintf(
				//nolint: lll
				`[{"id":"%s","header":"test1","description":"testDescription1","userId":"%s","eventTime":"%s"},{"id":"%s","header":"test2","description":"testDescription2","userId":"%s","eventTime":"%s"}]`,
				testEventID1, testUserID, testTime.Format(time.RFC3339Nano),
				testEventID2, testUserID, testTime.Format(time.RFC3339Nano)),
		},
		{
			name: "OK by a few day",
			mockBehavior: func(s *mockservice.MockApplicationInterface) {
				s.EXPECT().GetListEventsDuringFewDays(gomock.Any(), testDateOnly, 3).Return([]models.Event{
					{
						ID:          testEventID1,
						Header:      "test1",
						Description: "testDescription1",
						UserID:      testUserID,
						EventTime:   testTime,
					},
					{
						ID:          testEventID2,
						Header:      "test2",
						Description: "testDescription2",
						UserID:      testUserID,
						EventTime:   testTime,
					},
				}, nil)
			},
			getParams:          fmt.Sprintf("?start=%s&amount_days=3", testParamTime),
			expectedStatusCode: 200,
			expectedRequestBody: fmt.Sprintf(
				//nolint: lll
				`[{"id":"%s","header":"test1","description":"testDescription1","userId":"%s","eventTime":"%s"},{"id":"%s","header":"test2","description":"testDescription2","userId":"%s","eventTime":"%s"}]`,
				testEventID1, testUserID, testTime.Format(time.RFC3339Nano),
				testEventID2, testUserID, testTime.Format(time.RFC3339Nano)),
		},
		{
			name: "Server error",
			mockBehavior: func(s *mockservice.MockApplicationInterface) {
				s.EXPECT().GetListEventsDuringFewDays(gomock.Any(), testDateOnly, 3).Return(nil, errors.New("server error"))
			},
			getParams:           fmt.Sprintf("?start=%s&amount_days=3", testParamTime),
			expectedStatusCode:  500,
			expectedRequestBody: "server error\n",
		},
		{
			name:                "No start time in input",
			mockBehavior:        func(_ *mockservice.MockApplicationInterface) {},
			getParams:           "",
			expectedStatusCode:  400,
			expectedRequestBody: "start is required parameter\n",
		},
		{
			name:               "invalid start time in input",
			mockBehavior:       func(_ *mockservice.MockApplicationInterface) {},
			getParams:          "?start=invalidTime",
			expectedStatusCode: 400,
			//nolint: lll
			expectedRequestBody: "invalid start parameter: parsing time \"invalidTime\" as \"2006-01-02\": cannot parse \"invalidTime\" as \"2006\"\n",
		},
		{
			name:                "invalid amount days in input",
			mockBehavior:        func(_ *mockservice.MockApplicationInterface) {},
			getParams:           fmt.Sprintf("?start=%s&amount_days=a", testParamTime),
			expectedStatusCode:  400,
			expectedRequestBody: "Invalid amount_days parameter\n",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			appInterface := mockservice.NewMockApplicationInterface(c)
			testCase.mockBehavior(appInterface)

			logg, err := logger.GetLogger("INFO")
			require.NoError(t, err)

			handler := NewHandler(logg, appInterface)
			r := handler.InitRoutes()
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/event/list%s", testCase.getParams), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedRequestBody, w.Body.String())
		})
	}
}

func TestCreateEvent(t *testing.T) {
	type mockBehavior func(s *mockservice.MockApplicationInterface, eventDTO models.Event)
	testUserID := uuid.New().String()
	testTimeForJSON := time.Now().Format(time.RFC3339Nano)
	testTime, _ := time.Parse(time.RFC3339Nano, testTimeForJSON)
	testEventID := uuid.New().String()

	testTable := []struct {
		name               string
		inputBody          string
		inputEvent         models.Event
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name: "OK",
			inputBody: fmt.Sprintf(`{"header":"test1","description":"testDescription1","userId":"%s","eventTime":"%s"}`,
				testUserID, testTimeForJSON),
			inputEvent: models.Event{
				Header:           "test1",
				Description:      "testDescription1",
				UserID:           testUserID,
				EventTime:        testTime,
				FinishEventTime:  nil,
				NotificationTime: nil,
			},
			mockBehavior: func(s *mockservice.MockApplicationInterface, eventDTO models.Event) {
				s.EXPECT().CreateEvent(gomock.Any(), eventDTO).Return(testEventID, nil)
			},

			expectedStatusCode: 201,
		},
		{
			name: "server error",
			inputBody: fmt.Sprintf(`{"header":"test1","description":"testDescription1","userId":"%s","eventTime":"%s"}`,
				testUserID, testTimeForJSON),
			inputEvent: models.Event{
				Header:           "test1",
				Description:      "testDescription1",
				UserID:           testUserID,
				EventTime:        testTime,
				FinishEventTime:  nil,
				NotificationTime: nil,
			},
			mockBehavior: func(s *mockservice.MockApplicationInterface, eventDTO models.Event) {
				s.EXPECT().CreateEvent(gomock.Any(), eventDTO).Return("", errors.New("server error"))
			},

			expectedStatusCode: 500,
		},
		{
			name: "invalid input",
			inputBody: fmt.Sprintf(`{"header":"test1","description":"testDescription1","userId":2,"eventTime":"%s"}`,
				testTimeForJSON),
			inputEvent:   models.Event{},
			mockBehavior: func(_ *mockservice.MockApplicationInterface, _ models.Event) {},

			expectedStatusCode: 400,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			appInterface := mockservice.NewMockApplicationInterface(c)
			testCase.mockBehavior(appInterface, testCase.inputEvent)
			logg, err := logger.GetLogger("INFO")
			require.NoError(t, err)
			handler := NewHandler(logg, appInterface)
			r := handler.InitRoutes()
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/event", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}

func TestUpdateEvent(t *testing.T) {
	type mockBehavior func(s *mockservice.MockApplicationInterface, eventDTO models.Event)
	testUserID := uuid.New().String()
	testTimeForJSON := time.Now().Format(time.RFC3339Nano)
	testTime, _ := time.Parse(time.RFC3339Nano, testTimeForJSON)
	testEventID := uuid.New().String()

	testTable := []struct {
		name               string
		inputBody          string
		pathID             interface{}
		inputEvent         models.Event
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name: "OK",
			inputBody: fmt.Sprintf(`{"header":"test1","description":"testDescription1","userId":"%s","eventTime":"%s"}`,
				testUserID, testTimeForJSON),
			inputEvent: models.Event{
				ID:               testEventID,
				Header:           "test1",
				Description:      "testDescription1",
				UserID:           testUserID,
				EventTime:        testTime,
				FinishEventTime:  nil,
				NotificationTime: nil,
			},
			pathID: testEventID,
			mockBehavior: func(s *mockservice.MockApplicationInterface, eventDTO models.Event) {
				s.EXPECT().UpdateEvent(gomock.Any(), eventDTO).Return(nil)
			},

			expectedStatusCode: 204,
		},
		{
			name: "server error",
			inputBody: fmt.Sprintf(`{"header":"test1","description":"testDescription1","userId":"%s","eventTime":"%s"}`,
				testUserID, testTimeForJSON),
			inputEvent: models.Event{
				ID:               testEventID,
				Header:           "test1",
				Description:      "testDescription1",
				UserID:           testUserID,
				EventTime:        testTime,
				FinishEventTime:  nil,
				NotificationTime: nil,
			},
			pathID: testEventID,
			mockBehavior: func(s *mockservice.MockApplicationInterface, eventDTO models.Event) {
				s.EXPECT().UpdateEvent(gomock.Any(), eventDTO).Return(errors.New("server error"))
			},

			expectedStatusCode: 500,
		},
		{
			name: "invalid input",
			inputBody: fmt.Sprintf(`{"header":"test1","description":"testDescription1","userId":"%s","eventTime":"%s"}`,
				testUserID, testTimeForJSON),
			inputEvent:   models.Event{},
			pathID:       1,
			mockBehavior: func(_ *mockservice.MockApplicationInterface, _ models.Event) {},

			expectedStatusCode: 400,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			appInterface := mockservice.NewMockApplicationInterface(c)
			testCase.mockBehavior(appInterface, testCase.inputEvent)
			logg, err := logger.GetLogger("INFO")
			require.NoError(t, err)
			handler := NewHandler(logg, appInterface)
			r := handler.InitRoutes()
			w := httptest.NewRecorder()
			req := httptest.NewRequest(
				"PUT", fmt.Sprintf("/event/%v", testCase.pathID), bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}

func TestDeleteEvent(t *testing.T) {
	type mockBehavior func(s *mockservice.MockApplicationInterface, id string)
	testEventID := uuid.New().String()

	testTable := []struct {
		name               string
		pathID             interface{}
		id                 string
		mockBehavior       mockBehavior
		expectedStatusCode int
	}{
		{
			name:   "OK",
			pathID: testEventID,
			id:     testEventID,
			mockBehavior: func(s *mockservice.MockApplicationInterface, id string) {
				s.EXPECT().DeleteEvent(gomock.Any(), id).Return(nil)
			},

			expectedStatusCode: 204,
		},
		{
			name:   "server error",
			pathID: testEventID,
			id:     testEventID,
			mockBehavior: func(s *mockservice.MockApplicationInterface, id string) {
				s.EXPECT().DeleteEvent(gomock.Any(), id).Return(errors.New("server error"))
			},

			expectedStatusCode: 500,
		},
		{
			name:               "invalid input",
			pathID:             1,
			mockBehavior:       func(_ *mockservice.MockApplicationInterface, _ string) {},
			expectedStatusCode: 400,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			appInterface := mockservice.NewMockApplicationInterface(c)
			testCase.mockBehavior(appInterface, testCase.id)
			logg, err := logger.GetLogger("INFO")
			require.NoError(t, err)
			handler := NewHandler(logg, appInterface)
			r := handler.InitRoutes()
			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/event/%v", testCase.pathID), nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
		})
	}
}
