package models

import (
	"ecommerce-sys/db"
	"github.com/astaxie/beego"
)

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

type TeamVO struct {
	Team              Team `json:"team"`
	SecondAgentsCount uint `json:"secondAgentsCount"`
	ThirdAgentsCount  uint `json:"thirdAgentsCount"`
}

type ITeamOperation interface {
	QueryTeamByInvitationCode(invitationCode string) error
	QueryUserTeams(userId uint64) (*TeamVO, error)
}

func (team *Team) QueryTeamByInvitationCode(invitationCode string) error {
	mysqlDB := db.GetMySqlConnection().GetMySqlDB()
	err := mysqlDB.Where("invitationCode = ? and status = 'active'", invitationCode).First(&team).Error
	return err
}

func (team *Team) QueryUserTeams(userId uint64) (*TeamVO, error) {
	mysqlDB := db.GetMySqlConnection().GetMySqlDB()
	var userTeam Team
	err := mysqlDB.Where("userId = ? and status = 'active'", userId).First(&userTeam).Error
	// query top agent
	secondAgentsCount, err := countUsersByUserId("topLevelAgent", userId)
	if err != nil {
		beego.Error(err.Error())
		return nil, err
	}
	// query supervisor agent
	thirdAgentsCount, err := countUsersByUserId("superiorAgent", userId)
	if err != nil {
		beego.Error(err.Error())
		return nil, err
	}
	vo := TeamVO{userTeam, secondAgentsCount, thirdAgentsCount}
	return &vo, err
}

func countUsersByUserId(agentLevel string, userId uint64) (uint, error) {
	var count uint
	mysqlDB := db.GetMySqlConnection().GetMySqlDB()
	var err error
	if agentLevel == "topLevelAgent" {
		err = mysqlDB.Model(&Team{}).Where("topLevelAgent = ?", userId).Count(&count).Error
	} else {
		err = mysqlDB.Model(&Team{}).Where("superiorAgent = ?", userId).Count(&count).Error
	}
	return count, err
}
