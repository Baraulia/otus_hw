package scripts

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/suite"
	"log"
	"os"
	"sync"
	"time"

	"github.com/cucumber/godog"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	amqpDSN            = os.Getenv("TESTS_AMQP_DSN")
	header             = "testHeader"
	headerForScheduler = "schedulerHeader"
	description        = "test description"
)
var createEventQuery = "INSERT INTO event (header, description, user_id, event_time, notification_time) VALUES($1,$2,$3,$4,$5) RETURNING id"

func init() {
	if amqpDSN == "" {
		amqpDSN = "amqp://guest:guest@rabbitmq:5672/"
	}
}

const (
	queueName                 = "ToNotification"
	notificationsExchangeName = "test-exchange"
	routingKey                = "test-route"
)

type notifyTest struct {
	suite.Suite
	db            *sql.DB
	conn          *amqp.Connection
	ch            *amqp.Channel
	messages      [][]byte
	messagesMutex *sync.RWMutex
	stopSignal    chan struct{}

	eventID string
}

func (test *notifyTest) setupTest() error {
	//start consuming
	var err error
	test.messages = make([][]byte, 0)
	test.messagesMutex = new(sync.RWMutex)
	test.stopSignal = make(chan struct{})

	test.conn, err = amqp.Dial(amqpDSN)
	if err != nil {
		return err
	}

	test.ch, err = test.conn.Channel()
	if err != nil {
		return err
	}

	// Consume
	_, err = test.ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = test.ch.QueueBind(queueName, routingKey, notificationsExchangeName, false, nil)
	if err != nil {
		return err
	}

	events, err := test.ch.Consume(queueName, "", true, true, false, false, nil)
	if err != nil {
		return err
	}

	go func(stop <-chan struct{}) {
		for {
			select {
			case <-stop:
				return
			case event := <-events:
				test.messagesMutex.Lock()
				test.messages = append(test.messages, event.Body)
				test.messagesMutex.Unlock()
			}
		}
	}(test.stopSignal)

	db, err := sql.Open("postgres", dataBaseConnectionString)
	if err != nil {
		return err
	}

	test.db = db

	log.Println("Creating new event for scheduler test...")
	var id string
	err = test.db.QueryRow(createEventQuery, headerForScheduler, description, userID, eventTime, notificationTime).Scan(&id)
	if err != nil {
		return err
	}

	test.eventID = id

	return nil
}

func (test *notifyTest) tearDownTest() error {
	close(test.stopSignal)

	log.Println("Clearing database from scheduler test...")
	query := `DELETE FROM event`
	_, err := test.db.Exec(query)
	if err != nil {
		return err
	}

	err = test.db.Close()
	if err != nil {
		return err
	}

	err = test.ch.Close()
	if err != nil {
		return err
	}

	err = test.conn.Close()
	if err != nil {
		return err
	}

	test.messages = nil

	return nil
}

func (test *notifyTest) iReceiveEventWithText() error {
	time.Sleep(3 * time.Second) // waiting for processing of event
	testMessage := fmt.Sprintf(`{"eventId":"%s","eventHeader":"%s","eventTime":"%s","userId":"%s"}`,
		test.eventID, headerForScheduler, eventTime, userID)

	test.messagesMutex.RLock()
	defer test.messagesMutex.RUnlock()

	for _, msg := range test.messages {
		if string(msg) == testMessage {
			return nil
		}
	}
	return fmt.Errorf("event with text '%s' was not found in %s", testMessage, test.messages)
}

func InitializeSchedulerScenario(ctx *godog.ScenarioContext) {
	test := new(notifyTest)

	ctx.Step(`Setup test for scheduler`, test.setupTest)
	ctx.Step(`Sender processed a new notification I want to receive notification about it`, test.iReceiveEventWithText)
	ctx.Step(`Teardown test for scheduler`, test.tearDownTest)
}
