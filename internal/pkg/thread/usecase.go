package thread

import "github.com/ErikDoter/2020_2_technoPark_SUBD/internal/pkg/models"

type Usecase interface {
	Find(slugOrId string) (*models.Thread, *models.Error)
}
