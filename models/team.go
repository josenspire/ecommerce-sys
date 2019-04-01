package models

import "ecommerce-sys/db"

type Team struct {
	TeamId         uint64 `json:"teamId" gorm:"column:teamId; primary_key; not null;"`
	TopLevelAgent  uint64 `json:"topLevelAgent" gorm:"column:topLevelAgent; not null;"`
	SuperiorAgent  uint64 `json:"superiorAgent" gorm:"column:superiorAgent; not null;"`
	Status         string `json:"status" gorm:"column:status; default:'active'; not null;"`
	Channel        string `json:"channel" gorm:"column:channel; default:'Wechat'; not null;"`
	InvitationCode string `json:"invitationCode" gorm:"column:invitationCode; unique; type:varchar(6); not null;"`
	UserId         uint64 `json:"userId" gorm:"column:userId;not null;"`
	BaseModel
}

type ITeamOperation interface {
	QueryTeamByInvitationCode(invitationCode string) error
	QueryUserTeams(userId uint64) error
}

func (team *Team) QueryTeamByInvitationCode(invitationCode string) error {
	mysqlDB := db.GetMySqlConnection().GetMySqlDB()
	err := mysqlDB.Where("invitationCode = ? and status = 'active'", invitationCode).First(&team).Error
	return err
}

func (team *Team) QueryUserTeams(userId uint64) error {
	mysqlDB := db.GetMySqlConnection().GetMySqlDB()
	err := mysqlDB.Where("userId = ? and status = 'active'", userId).First(team).Error
	return err
}
