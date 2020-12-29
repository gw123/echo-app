package gather

import (
	"context"
	"fmt"

	"github.com/gw123/echo-app/app"
	"github.com/gw123/echo-app/gather/internal/domain/gather/entity"
	"github.com/pkg/errors"
)

func GatherSalesVolumes(ctx context.Context, targetID int64, startTime, endTime string) (num int64, err error) {
	return 0, nil
}
func GatherCommentsNumber(ctx context.Context, targetID int64, startTime, endTime string) (num int64, err error) {
	return 0, nil
}
func GatherViews(ctx context.Context, targetID int64, startTime, endTime string) (num int64, err error) {
	db, err := app.GetDb("user")
	if err != nil {
		return -1, errors.Wrap(err, "db user")
	}
	if startTime > endTime {
		startTime, endTime = endTime, startTime
	}
	result := []*entity.Statistics{}
	err = db.Table("user_history").
		Where("created_at>=? and created_at<=? ", startTime, endTime). //统计当天0时到现在
		Select("count(*) as total,target_id,type").
		Group("target_id,type").
		Find(&result).Error
	if err != nil {
		return 0, errors.Wrap(err, "db user_his groupby")
	}
	//db.AutoMigrate(&entity.Statistics{})
	for _, val := range result {
		val.Date = fmt.Sprintf("%s--%s", startTime, endTime)
		//log.Printf(val)
		if err := db.Save(val).Error; err != nil {
			return -1, errors.Wrap(err, "db create")
		}

	}
	return 0, nil
}
