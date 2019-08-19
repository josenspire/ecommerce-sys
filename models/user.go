package models

import (
	"ecommerce-sys/db"
	. "ecommerce-sys/utils"
	"encoding/hex"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/jinzhu/gorm"
)

type User struct {
	UserId uint64 `json:"userId" gorm:"column:userId;primary_key;not null"`
	UserProfile
	Role      uint16    `json:"role" gorm:"column:role; default:10; not null;"`
	Status    string    `json:"status" gorm:"column:status; type:varchar(10); default:'active'; not null;"`
	Channel   string    `json:"channel" gorm:"column:channel; type:varchar(12); not null;"`
	Addresses []Address `json:"addresses" gorm:"column:addresses;"`
	BaseModel
}

type UserProfile struct {
	Telephone string `json:"telephone" gorm:"column:telephone; type:varchar(11);not null"`
	Username  string `json:"username" gorm:"column:username; type:varchar(18);not null"`
	Password  string `json:"-" gorm:"column:password; type:varchar(24);null"`
	Nickname  string `json:"nickname" gorm:"column:nickname; type:varchar(16);not null;"`
	Male      bool   `json:"male" gorm:"column:male; not null; default:false;"`
	Signature string `json:"signature" gorm:"not null; default:'This guy is lazy...'"`
}

type WxSession struct {
	SessionId         uint64 `json:"sessionId" gorm:"column:sessionId; not null; primary_key;"`
	Skey              string `json:"skey" gorm:"column:skey; not null;"`
	SessionKey        string `json:"sessionKey" gorm:"column:sessionKey; not null;" `
	WechatUserProfile string `json:"wechatUserProfile" gorm:"column:wechatUserProfile; not null;"`
	OpenId            string `json:"openId" gorm:"column:openId; index; not null;"`
	UserId            uint64 `json:"userId" gorm:"column:userId; not null;"`
	BaseModel
}

type UserRegisterDTO struct {
	Telephone      string `json:"telephone"`
	Password       string `json:"password"`
	Nickname       string `json:"nickname"`
	InvitationCode string `json:"invitationCode"`
}

type UserWechatVO struct {
	User          *User      `json:"user"`
	WechatSession *WxSession `json:"wechatSession"`
}

func (WxSession) TableName() string {
	return "wxsessions"
}

// Get global config
var AESSecretKey = beego.AppConfig.String("AESSecretKey")

// callbacks hock -- before create, encrypt password
func (user *User) BeforeCreate(scope *gorm.Scope) error {
	if IsEmptyString(user.Password) {
		return nil
	}
	encryptPassword, err := AESEncrypt(user.Password, AESSecretKey)
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	err = scope.SetColumn("password", encryptPassword)
	if err != nil {
		logs.Error(err.Error())
		return err
	}
	return nil
}

type IUserOperation interface {
	Register(dto UserRegisterDTO) error
	CheckIsUserExistByUserId(userId uint64) (bool, error)
	CheckIsUserExistByTelephone(telephone string) (bool, error)
	QueryByUserId(userId string) *User
	LoginByTelephone(telephone string, password string) error
	LoginByWechat(jsCode string, wechatUserProfile string, invitationCode string) (*UserWechatVO, error)
}

func (user *User) Register(dto UserRegisterDTO) error {
	isExist, err := user.CheckIsUserExistByTelephone(dto.Telephone)
	if err != nil {
		return err
	}
	if isExist == true {
		return ErrCurrentUserIsExist
	}
	user.UserId = GetWuid()
	user.Telephone = dto.Telephone
	user.Password = dto.Password
	user.Nickname = dto.Nickname

	team, err := initialUserAgentTeamsContent(user.UserId, dto.InvitationCode)

	mysqlDB := db.GetMySqlConnection().GetMySqlDB()
	ts := mysqlDB.Begin()
	err = ts.Create(&user).Error
	err = ts.Create(&team).Error
	if err != nil {
		logs.Error(err.Error())
		ts.Rollback()
	} else {
		ts.Commit()
	}
	return err
}

func (user *User) LoginByTelephone(telephone string, password string) error {
	mysqlDB := db.GetMySqlConnection().GetMySqlDB()
	var userProfile User
	err := mysqlDB.Where("telephone = ?", telephone).First(&userProfile).Error
	if err == gorm.ErrRecordNotFound {
		return ErrTelOrPswInvalid
	}
	if err != nil {
		logs.Error(err)
		return err
	}
	if IsEmptyString(userProfile.Password) {
		// 	TODO: should verify telephone, send SMS code
		return WarnAccountNeedVerify
	}
	encryptPassword, err := AESEncrypt(password, AESSecretKey)
	if err != nil {
		logs.Error(err.Error())
		return ErrDecrypt
	}
	err = mysqlDB.Where("telephone = ? and password = ?", telephone, encryptPassword).First(&user).Error
	return err
}

func (user *User) CheckIsUserExistByUserId(userId uint64) (bool, error) {
	if userId == 0 {
		logs.Warn(WarnParamsMissing.Error())
		return false, WarnParamsMissing
	}
	o := orm.NewOrm()
	queryUser := new(User)
	queryUser.UserId = userId
	err := o.Read(&queryUser)
	if err == orm.ErrNoRows {
		return false, nil
	} else if err == orm.ErrMissPK {
		return false, err
	}
	return true, nil
}

func (user *User) CheckIsUserExistByTelephone(telephone string) (bool, error) {
	if telephone == "" {
		logs.Warn(WarnParamsMissing.Error())
		return false, WarnParamsMissing
	}
	mysqlDB := db.GetMySqlConnection().GetMySqlDB()

	err := mysqlDB.Where("telephone = ?", telephone).Find(&user).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (user *User) QueryByUserId(userId string) *User {
	panic("implement me")
}

func (user *User) LoginByWechat(jsCode string, wechatUserProfile string, invitationCode string) (*UserWechatVO, error) {
	user, wxSession, err := authorizationWithInitialUser(jsCode, wechatUserProfile, invitationCode)
	if err != nil {
		logs.Error(err.Error())
		return nil, err
	}
	var vo = UserWechatVO{
		User:          user,
		WechatSession: wxSession,
	}
	return &vo, nil
}

func authorizationWithInitialUser(jsCode string, wechatUserProfile string, invitationCode string) (*User, *WxSession, error) {
	openId, sessionKey, err := JsCode2Session(jsCode)
	// openId, sessionKey := GenerateNowDateString(), GenerateRandString(8)

	if err != nil {
		logs.Error(err.Error())
		return nil, nil, err
	}
	fmt.Printf("User openId = %s, sessionKey = %s", openId, sessionKey)

	skeyBytes, err := SHA1Encrypt(sessionKey)
	if err != nil {
		logs.Error(err)
		return nil, nil, err
	}
	skeyHexStr := hex.EncodeToString(skeyBytes)

	mysqlDB := db.GetMySqlConnection().GetMySqlDB()
	tx := mysqlDB.Begin()

	wechatSession := WxSession{}
	err = tx.Where("openId = ?", openId).First(&wechatSession).Error
	if err == gorm.ErrRecordNotFound {
		var userId = GetWuid()
		user := User{
			UserId:  userId,
			Channel: "Wechat",
		}
		err = tx.Create(&user).Error
		if err != nil {
			logs.Error(err.Error())
			tx.Rollback()
			return nil, nil, err
		}
		wxSession := WxSession{
			SessionId:         GetWuid(),
			SessionKey:        sessionKey,
			OpenId:            openId,
			UserId:            userId,
			Skey:              skeyHexStr,
			WechatUserProfile: wechatUserProfile,
		}
		err = tx.Create(&wxSession).Error
		if err != nil {
			logs.Error(err.Error())
			tx.Rollback()
			return nil, nil, err
		}
		wechatSession = wxSession

		team, err := initialUserAgentTeamsContent(userId, invitationCode)
		if err != nil {
			logs.Error(err.Error())
			tx.Rollback()
			return nil, nil, err
		}
		err = tx.Create(&team).Error
		if err != nil {
			logs.Error(err.Error())
			tx.Rollback()
			return nil, nil, err
		}
	} else if err != nil {
		logs.Error(err.Error())
		tx.Rollback()
		return nil, nil, err
	}
	err = tx.Model(&wechatSession).Updates(map[string]interface{}{"sessionKey": sessionKey, "skey": skeyHexStr}).Error
	if err != nil {
		logs.Error(err.Error())
		return nil, nil, err
	}
	tx.Commit()

	var user = User{}
	var wxSession = WxSession{}
	err = mysqlDB.Where("userId = ?", wechatSession.UserId).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		logs.Error(err.Error())
		return nil, nil, err
	}
	err = mysqlDB.Model(&user).Where("userId = ?", user.UserId).First(&wxSession).Error
	return &user, &wxSession, nil
}

func initialUserAgentTeamsContent(userId uint64, invitationCode string) (*Team, error) {
	var team = Team{
		UserId:         userId,
		TeamId:         GetWuid(),
		InvitationCode: GenerateRandString(6),
	}
	var agentTeam Team
	err := agentTeam.QueryTeamByInvitationCode(invitationCode)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// TOP ADMIN
			team.TopLevelAgent = TOP_AGENT
			team.SuperiorAgent = SUPERIOR_AGNET
		} else {
			logs.Error(err.Error())
			return nil, err
		}
	} else {
		team.TopLevelAgent = agentTeam.SuperiorAgent
		team.SuperiorAgent = agentTeam.UserId
	}
	return &team, nil
}
