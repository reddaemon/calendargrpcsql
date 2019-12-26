package service

import (
	"context"
	"errors"

	//	"github.com/labstack/gommon/log"
	"github.com/reddaemon/calendargrpcsql/app"
	"github.com/reddaemon/calendargrpcsql/models/models"

	//	"github.com/reddaemon/calendargrpcsql/models/storage"
	//	"sync"
	"time"
)

type Repo interface {
	Create(ctx context.Context, event *models.Event) (*models.Event, error)
	Read(ctx context.Context, eventId uint64) (*models.Event, error)
	Update(ctx context.Context, event *models.Event, eventId uint64) (bool, error)
	Delete(ctx context.Context, eventId uint64) (bool, error)
}

type EventUseCase struct {
	App            *app.App
	eventRepo      Repo
	contextTimeout time.Duration
}

func NewEventUseCase(app *app.App, eventRepo Repo, contextTimeout time.Duration) *EventUseCase {
	return &EventUseCase{App: app, eventRepo: eventRepo, contextTimeout: contextTimeout}
}

func (e *EventUseCase) Create(ctx context.Context, event *models.Event) (*models.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()
	newEvent, err := e.eventRepo.Create(ctx, event)
	if err != nil {
		return nil, err
	}
	return newEvent, nil
}

func (e *EventUseCase) Read(ctx context.Context, eventId uint64) (*models.Event, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()
	res, err := e.eventRepo.Read(ctx, eventId)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (e *EventUseCase) Update(ctx context.Context, event *models.Event, eventId uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()
	return e.eventRepo.Update(ctx, event, eventId)
}

func (e *EventUseCase) Delete(ctx context.Context, eventId uint64) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, e.contextTimeout)
	defer cancel()
	existedEvent, err := e.eventRepo.Delete(ctx, eventId)
	if err != nil {
		return false, err
	}
	if !existedEvent {
		return false, errors.New("event not found")
	}
	return e.eventRepo.Delete(ctx, eventId)
}
