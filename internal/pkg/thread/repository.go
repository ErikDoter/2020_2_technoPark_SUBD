package thread

import "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"

type Repository interface {
	Find(soi models.IdOrSlug) (*models.Thread, *models.Error)
}
