package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

func GetCommunityList() ([]models.Community, error){
	// 查与commnunity相关的数据
	return mysql.GetCommunityList()
} 

func GetCommunityDetail(cId int64) ( *models.CommunityDetail, error){
	return mysql.GetCommunityDetail(cId)	
}