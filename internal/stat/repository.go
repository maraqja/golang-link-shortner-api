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

func (repo *StatRepository) GetStats(by string, from, to time.Time) []GetStatResponse {
	var stats []GetStatResponse

	/*
		SELECT to_char(date, format) AS period, sum(clicks) FROM stats
		WHERE date BETWEEN '01/01/2000' AND '01/01/2026'
		GROUP BY period
		ORDER BY period


		в зависимости от by:
			by = month => format = 'YYYY-MM'
			by = day ==> format = 'YYYY-MM-DD'
	*/

	var selectQuery string
	switch by {
	case GroupByDay:
		selectQuery = "to_char(date, 'YYYY-MM-DD') AS period, sum(clicks)"
	case GroupByMonth:
		selectQuery = "to_char(date, 'YYYY-MM') AS period, sum(clicks)"
	}

	repo.Db.
		Table("stats").
		Select(selectQuery).
		Where("date BETWEEN ? AND ?", from, to).
		Group("period").
		Order("period").
		Scan(&stats)

	return stats

}
