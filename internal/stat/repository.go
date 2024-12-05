package stat

import (
	"link-shortner-api/pkg/db"
	"time"

	"gorm.io/datatypes"
)

type StatRepository struct {
	*db.Db
}

func NewStatRepository(db *db.Db) *StatRepository {
	return &StatRepository{
		Db: db,
	}
}

func (repo *StatRepository) AddClick(linkId uint) {
	// если нет записи (за сегодня) статистики для ссылки, то создаем новую
	// если есть - инкремент

	var stat Stat

	currentDate := datatypes.Date(time.Now())
	repo.Db.Find(&stat, "link_id = ? AND date = ?", linkId, currentDate)
	if stat.ID == 0 {
		stat = Stat{
			LinkID: linkId,
			Clicks: 1,
			Date:   currentDate,
		}
		repo.Db.Create(&stat)
	} else {
		stat.Clicks += 1
		repo.DB.Save(&stat)
	}
}
