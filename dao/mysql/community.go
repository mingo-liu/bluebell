package mysql

import (
	"bluebell/models"
	"database/sql"

	"go.uber.org/zap"
)

func GetCommunityList() (data []models.Community, err error){
	sqlStr := `select community_id, community_name from community`	
	if err = db.Select(&data, sqlStr); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no community in db")
			err = nil
		}
		return
	}
	return
}

func GetCommunityDetail(cId int64) (communityDetail *models.CommunityDetail, err error) {
	sqlStr := `select community_id, community_name, introduction, create_time 
			   from community where community_id = ?` 
	communityDetail = new(models.CommunityDetail)
	if err = db.Get(communityDetail, sqlStr, cId); err != nil {
		if err == sql.ErrNoRows {
			err = ErrorInvalidID 
		}
		return
	}
	return
}