package repositories

import (
	"github.com/holycann/cultour-backend/internal/cultural/models"
	"github.com/holycann/cultour-backend/pkg/repository"
)

type EventRepository interface {
	repository.BaseRepository[models.Event]
}

type LocalStoryRepository interface {
	repository.BaseRepository[models.LocalStory]
}
