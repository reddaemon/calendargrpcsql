package storage

import (
	"context"
	"fmt"

	"github.com/reddaemon/calendargrpcsql/app"
	"github.com/reddaemon/calendargrpcsql/models/models"
)

type PsqlEventStorage struct {
	app *app.App
}

func NewPsqlRepository(app *app.App) *PsqlEventStorage {
	return &PsqlEventStorage{app: app}
}

func (p PsqlEventStorage) Create(ctx context.Context, event *models.Event) (*models.Event, error) {
	query := `
			INSERT INTO events(title, description, date)
			VALUES (:title, :description, :date)
			RETURNING id
	`

	stmt, err := p.app.Db.PrepareNamedContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var id int64

	err = stmt.GetContext(ctx, &id, map[string]interface{}{
		"description": event.Description,
		"title":       event.Title,
		"date":        event.Date,
	})
	if err != nil {
		return nil, err
	}

	event.Id = uint64(id)
	return event, nil
}

func (p PsqlEventStorage) Read(ctx context.Context, eventID uint64) (*models.Event, error) {
	query := `SELECT id, title, description, date FROM events WHERE ID = $1`
	row := p.app.Db.QueryRowContext(ctx, query, eventID)

	event := models.Event{}
	err := row.Scan(&event.Id, &event.Title, &event.Description, &event.Date)

	if err != nil {
		return nil, err
	}

	return &event, nil

}

func (p PsqlEventStorage) Update(ctx context.Context, event *models.Event, eventId uint64) (bool, error) {
	query := `UPDATE events set title=$1, description=$2, date=$3 WHERE ID = $4`
	stmt, err := p.app.Db.PrepareContext(ctx, query)
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, event.Title, event.Description, event.Date, eventId)
	if err != nil {
		return false, err
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if affect != 1 {
		err = fmt.Errorf("Strange behaviour. total affected: %d", affect)
		return false, err
	}

	return true, nil
}

func (p PsqlEventStorage) Delete(ctx context.Context, eventId uint64) (bool, error) {
	query := "DELETE FROM events WHERE id = $1"

	stmt, err := p.app.Db.PrepareContext(ctx, query)
	if err != nil {
		return false, err
	}

	_, err = stmt.ExecContext(ctx, eventId)
	if err != nil {
		return false, err
	}

	return true, nil
}
