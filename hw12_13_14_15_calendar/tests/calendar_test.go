package scripts

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cucumber/godog"
)

var eventTime = time.Now().Add(24 * time.Hour).Format("2006-01-02T15:04:05Z07:00")
var notificationTime = time.Now().Format("2006-01-02T15:04:05Z07:00")
var userID = "d9b64851-b955-4e29-aac2-0c0eea95d5fd"

var dataBaseConnectionString = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
	"postgres", "password", os.Getenv("POSTGRES_HOST"), 5432, "postgres")

var testEvent = fmt.Sprintf(`
	{
		"header":"testHeader",
		"description":"test description",
		"userId":"%s",
		"eventTime":"%s",
		"notificationTime":"%s" 
	}
`, userID, eventTime, notificationTime)

var updatingEvent = fmt.Sprintf(`
	{
		"header":"newHeader",
		"description":"new description",
		"userId":"%s",
		"eventTime":"%s",
		"notificationTime":"%s" 
	}
`, userID, eventTime, notificationTime)

type calendarTest struct {
	httpClient http.Client
	db         *sql.DB

	eventID            string
	responseStatusCode int
	responseBody       []byte
}

func (test *calendarTest) setupTest() error {
	db, err := sql.Open("postgres", dataBaseConnectionString)
	if err != nil {
		return err
	}

	test.db = db
	test.httpClient = http.Client{}

	return nil
}

func (test *calendarTest) tearDownTest() error {
	log.Println("Clearing database from calendar test...")
	query := `DELETE FROM event`
	_, err := test.db.Exec(query)
	if err != nil {
		return err
	}

	err = test.db.Close()
	if err != nil {
		return err
	}

	return nil
}

func (test *calendarTest) createNewEvent() error {
	var id string
	log.Println("Creating new event for calendar test...")
	err := test.db.QueryRow(createEventQuery, header, description, userID, eventTime, notificationTime).Scan(&id)
	if err != nil {
		return err
	}

	test.eventID = id

	return nil
}

func (test *calendarTest) iSendRequestTo(httpMethod, addr string) error {
	var response *http.Response
	var err error
	var request *http.Request
	switch httpMethod {
	case http.MethodGet:
		params := fmt.Sprintf("?start=%v&amount_days=3", time.Now().Format("2006-01-02"))
		addr += params
		request, err = http.NewRequest("GET", addr, nil)
		if err != nil {
			return err
		}

		response, err = test.httpClient.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		body, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		test.responseBody = body

	case http.MethodPost:
		request, err = http.NewRequest("POST", addr, bytes.NewBuffer([]byte(testEvent)))
		if err != nil {
			return err
		}

		response, err = test.httpClient.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()

		test.responseBody = nil
		test.eventID = response.Header.Get("id")

	case http.MethodPut:
		addr += test.eventID
		request, err = http.NewRequest("PUT", addr, bytes.NewBuffer([]byte(updatingEvent)))
		if err != nil {
			return err
		}

		response, err = test.httpClient.Do(request)
		if err != nil {
			return err
		}
		defer response.Body.Close()
		test.responseBody = nil

	case http.MethodDelete:
		addr += test.eventID
		request, err = http.NewRequest("DELETE", addr, nil)
		if err != nil {
			return err
		}

		response, err = test.httpClient.Do(request)
		defer response.Body.Close()
		test.responseBody = nil

	default:
		return fmt.Errorf("unknown method: %s", httpMethod)
	}

	test.responseStatusCode = response.StatusCode

	return nil
}

func (test *calendarTest) theResponseCodeShouldBe(code int) error {
	if test.responseStatusCode != code {
		return fmt.Errorf("unexpected status code: %d != %d", test.responseStatusCode, code)
	}
	return nil
}

func (test *calendarTest) theResponseShouldMatchText() error {
	list := fmt.Sprintf(`[{"id":"%s","header":"testHeader","description":"test description","userId":"%s","eventTime":"%s","notificationTime":"%s"}]`,
		test.eventID, userID, eventTime, notificationTime)

	if string(test.responseBody) != list {
		return fmt.Errorf("unexpected text: %s != %s", test.responseBody, list)
	}

	return nil
}

func InitializeCalendarScenario(ctx *godog.ScenarioContext) {
	test := &calendarTest{httpClient: http.Client{}}

	ctx.Step(`Setup test for calendar`, test.setupTest)
	ctx.Step(`Create new event`, test.createNewEvent)
	ctx.Step(`^I send "([^"]*)" request to "([^"]*)"$`, test.iSendRequestTo)
	ctx.Step(`^The response code should be (\d+)$`, test.theResponseCodeShouldBe)
	ctx.Step(`^The response should match text in variable listEvents`, test.theResponseShouldMatchText)
	ctx.Step(`Teardown test for calendar`, test.tearDownTest)
}
