package services

import (
	"xhyovo.cn/community/server/model"
)

type LogServices struct {
}

func (*LogServices) GetPageOperLog(page, limit int, logSearch model.LogSearch) (logs []model.OperLogs, count int64) {
	db := model.OperLog()
	if logSearch.RequestMethod != "" {
		db.Where("request_method = ?", logSearch.RequestMethod)
	}
	if logSearch.RequestInfo != "" {
		db.Where("request_info like ?", "%s"+logSearch.RequestInfo+"%s")
	}
	if logSearch.Ip != "" {
		db.Where("ip like ?", "%s"+logSearch.Ip+"%s")
	}
	if logSearch.StartTime != "" {
		db.Where("created_at <= ? and ? >= created_at", logSearch.StartTime, logSearch.EndTime)
	}
	if logSearch.UserName != "" {
		var userS UserService
		ids := userS.SearchNameSelectId(logSearch.UserName)
		db.Where("user_id in ?", ids)
	}
	db.Limit(limit).Offset((page - 1) * limit).Order("created_at desc").Find(&logs)
	db.Count(&count)
	return
}

func (*LogServices) InsertOperLog(log model.OperLogs) {
	model.OperLog().Create(&log)
}

func (*LogServices) DeletesOperLogs(ids []int) {
	model.OperLog().Where("id ? in", ids).Delete(model.OperLogs{})
}
