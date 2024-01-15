package sqlstorage

//nolint:depguard
import (
	"context"
	"fmt"
	"time"

	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/Baraulia/otus_hw/hw12_13_14_15_calendar/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	// Empty import to ensure execution of code in the package's init function.
	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

const (
	MaxConnections = 10
	EventTable     = "event"
)

type PostgresStorage struct {
	databaseURL    string
	migrationsPath string
	db             *pgxpool.Pool
	logger         app.Logger
}

type PgConfig struct {
	Host           string
	Username       string
	Password       string
	Port           string
	Database       string
	MigrationsPath string
}

func NewPostgresStorage(conf PgConfig, logger app.Logger) *PostgresStorage {
	url := fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		conf.Username, conf.Password,
		conf.Host, conf.Port, conf.Database,
	)
	storage := &PostgresStorage{databaseURL: url, logger: logger, migrationsPath: conf.MigrationsPath}

	storage.Connect()

	return storage
}

func (s *PostgresStorage) Connect() {
	poolConfig, err := pgxpool.ParseConfig(fmt.Sprintf("%s&pool_max_conns=%d", s.databaseURL, MaxConnections))
	if err != nil {
		s.logger.Fatal("Unable to parse databaseURL", map[string]interface{}{"error": err, "databaseURL": s.databaseURL})
	}

	db, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		s.logger.Fatal("Unable to create connection pool", map[string]interface{}{"error": err})
	}

	s.db = db

	gooseDB, err := goose.OpenDBWithDriver("postgres", s.databaseURL)
	if err != nil {
		s.logger.Fatal("error while connecting to postgres by goose", map[string]interface{}{"error": err})
	}

	s.logger.Info("executing migrations...", nil)
	if err := goose.Up(gooseDB, s.migrationsPath); err != nil {
		s.logger.Fatal("error while executing migrations", map[string]interface{}{"error": err})
	}
}

func (s *PostgresStorage) Close() {
	s.db.Close()
}

func (s *PostgresStorage) CreateEvent(ctx context.Context, eventDTO models.Event) (string, error) {
	var id string
	sql := fmt.Sprintf(
		"INSERT INTO %s ('header','description','user_id','event_time','finish_event_time','notification_time') "+
			"VALUES($1, $2, $3, $4, $5, $6) RETURNING id", EventTable)
	err := s.db.QueryRow(
		ctx,
		sql,
		eventDTO.Header, eventDTO.Description, eventDTO.UserID,
		eventDTO.EventTime, eventDTO.FinishEventTime, eventDTO.NotificationTime,
	).Scan(&id)
	if err != nil {
		s.logger.Error("error while creating new event", map[string]interface{}{"error": err})
		return "", fmt.Errorf("error while creating new event: %w", err)
	}

	return id, nil
}

func (s *PostgresStorage) UpdateEvent(ctx context.Context, eventDTO models.Event) error {
	sql := fmt.Sprintf(
		"UPDATE %s SET ("+
			"header = $1,description = $2, user_id = $3, event_time = $4,"+
			" finish_event_time = $5, notification_time = $6) "+
			"WHERE id = $7", EventTable)
	result, err := s.db.Exec(
		ctx,
		sql,
		eventDTO.Header, eventDTO.Description, eventDTO.UserID, eventDTO.EventTime, eventDTO.FinishEventTime,
		eventDTO.NotificationTime, eventDTO.ID,
	)
	if err != nil {
		s.logger.Error("error while updating event", map[string]interface{}{"error": err})
		return fmt.Errorf("error while updating event: %w", err)
	}

	if result.RowsAffected() == 0 {
		s.logger.Error("no objects have been modified", nil)
		return fmt.Errorf("no objects have been modified")
	}

	return nil
}

func (s *PostgresStorage) DeleteEvent(ctx context.Context, id string) error {
	sql := fmt.Sprintf(
		"DELETE FROM %s WHERE id = $1", EventTable)
	result, err := s.db.Exec(ctx, sql, id)
	if err != nil {
		s.logger.Error("error while deleting event", map[string]interface{}{"error": err, "eventID": id})
		return fmt.Errorf("error while deleting event: %w", err)
	}

	if result.RowsAffected() == 0 {
		s.logger.Error("no objects have been deleted", nil)
		return fmt.Errorf("no objects have been deleted")
	}

	return nil
}

func (s *PostgresStorage) GetListEventsDuringDay(ctx context.Context, targetDay time.Time) ([]models.Event, error) {
	date := time.Date(targetDay.Year(), targetDay.Month(), targetDay.Day(), 0, 0, 0, 0, targetDay.Location())
	sql := fmt.Sprintf(
		"SELECT id, header, description, user_id, event_time, finish_event_time, notification_time "+
			"FROM %s WHERE DATE(event_time) = $1", EventTable)
	rows, err := s.db.Query(ctx, sql, date)
	if err != nil {
		s.logger.Error(
			"error while getting list events per day", map[string]interface{}{"error": err, "targetDay": targetDay})
		return nil, fmt.Errorf("error while getting list events per day: %w", err)
	}
	defer rows.Close()

	events := make([]models.Event, 0)
	for rows.Next() {
		var event models.Event
		if err = rows.Scan(&event.ID, &event.Header, &event.Description, &event.UserID, &event.EventTime,
			&event.FinishEventTime, &event.NotificationTime); err != nil {
			s.logger.Error("error while scanning event", map[string]interface{}{"error": err})
			return nil, fmt.Errorf("error while scanning event: %w", err)
		}

		events = append(events, event)
	}

	return events, nil
}

func (s *PostgresStorage) GetListEventsDuringFewDays(
	ctx context.Context, start time.Time, amountDays int,
) ([]models.Event, error) {
	startDate := time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, start.Location())
	finish := start.AddDate(0, 0, amountDays)
	finishDate := time.Date(finish.Year(), finish.Month(), finish.Day(), 0, 0, 0, 0, finish.Location())
	sql := fmt.Sprintf(
		"SELECT id, header, description, user_id, event_time, finish_event_time, notification_time FROM %s "+
			"WHERE DATE(event_time) >= $1 AND DATE(event_time) < $2", EventTable)
	rows, err := s.db.Query(ctx, sql, startDate, finishDate)
	if err != nil {
		s.logger.Error("error while getting list events per week", map[string]interface{}{"error": err, "startDay": start})
		return nil, fmt.Errorf("error while getting list events per week: %w", err)
	}
	defer rows.Close()

	events := make([]models.Event, 0)
	for rows.Next() {
		var event models.Event
		if err = rows.Scan(&event.ID, &event.Header, &event.Description, &event.UserID, &event.EventTime,
			&event.FinishEventTime, &event.NotificationTime); err != nil {
			s.logger.Error("error while scanning event", map[string]interface{}{"error": err})
			return nil, fmt.Errorf("error while scanning event: %w", err)
		}

		events = append(events, event)
	}

	return events, nil
}
