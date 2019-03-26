package models

import (
	"ecommerce-sys/db"
	. "ecommerce-sys/utils"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/jinzhu/gorm"
)

type User struct {
	// ID uint64 `json:"id" gorm:"column:id;NOT NULL;PRIMARY_KEY;"`
	UserId uint64 `json:"userId" gorm:"column:userId;primary_key;not null"`
	UserProfile
	Role      uint16    `json:"role" gorm:"column:role; default:10; not null;"`
	Status    string    `json:"status" gorm:"column:status; type:varchar(10); default:'active'; not null;"`
	Channel   string    `json:"channel" gorm:"column:channel; type:varchar(12); not null;"`
	Addresses []Address `json:"-"`
	Team      *Team     `json:"-"`
	BaseModel
}

type UserProfile struct {
	Telephone string `json:"telephone" gorm:"column:telephone; type:varchar(11);not null"`
	Username  string `json:"username" gorm:"column:username; type:varchar(18);not null"`
	Password  string `json:"password" gorm:"column:password; type:varchar(24);not null"`
	Nickname  string `json:"nickname" gorm:"column:nickname; type:varchar(16);not null;"`
	Male      bool   `json:"male" gorm:"column:male; not null; default:false;"`
	Signature string `json:"signature" gorm:"not null; default:'This guy is lazy...'"`
}

type WxSession struct {
	SessionId         uint64 `json:"sessionId" gorm:"column:sessionId; not null; primary_key;"`
	Skey              string `json:"skey" gorm:"column:skey; not null;"`
	SessionKey        string `json:"session_key" gorm:"column:sessionKey; not null;" `
	WechatUserProfile string `json:"wechatUserProfile" gorm:"column:wechatUserProfile; not null;"`
	OpenId            string `json:"openId" gorm:"column:openId; index; not null;"`
	BaseModel
}

type UserWechat struct {
	ID        uint64 `json:"id" gorm:"column:id;NOT NULL;PRIMARY_KEY;"`
	UserId    uint64 `json:"userId" gorm:"column:userId; not null;"`
	SessionId uint64 `json:"sessionId" gorm:"column:sessionId; not null; unique;"`
	BaseModel
}

type UserTeam struct {
	ID     uint64 `json:"id" gorm:"column:id;NOT NULL;PRIMARY_KEY;"`
	UserId uint64 `json:"userId" gorm:"column:userId; not null;"`
	TeamId uint64 `json:"teamId" gorm:"column:teamId; not null; unique;"`
	BaseModel
}

type UserRegisterDTO struct {
	Telephone      string `json:"telephone"`
	Password       string `json:"password"`
	Nickname       string `json:"nickname"`
	InvitationCode string `json:"invitationCode"`
}

func (WxSession) TableName() string {
	return "wxsessions"
}

func (UserWechat) TableName() string {
	return "userwechat"
}

func (UserTeam) TableName() string {
	return "userteams"
}

type IUserOperation interface {
	Register(dto UserRegisterDTO) error
	CheckIsUserExistByUserId(userId uint64) (bool, error)
	CheckIsUserExistByTelephone(telephone string) (bool, error)
	QueryByUserId(userId string) *User
	LoginByTelephone(telephone string, password string) error
	LoginByWechat(jsCode string, userInfo string, invitationCode string) (interface{}, error)
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

	var agentTeam Team
	team := Team{}
	team.UserId = user.UserId
	team.TeamId = GetWuid()
	team.InvitationCode = GenerateRandString(6)

	err = agentTeam.QueryTeamByInvitationCode(dto.InvitationCode)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// TOP ADMIN
			team.TopLevelAgent = 8888888888
			team.SuperiorAgent = 8888888888
		} else {
			logs.Error(err)
			return err
		}
	} else {
		team.TopLevelAgent = agentTeam.SuperiorAgent
		team.SuperiorAgent = agentTeam.UserId
	}
	userTeam := UserTeam{
		ID:     GetWuid(),
		UserId: user.UserId,
		TeamId: team.TeamId,
	}

	mysqlDB := db.GetMySqlConnection().GetMySqlDB()
	ts := mysqlDB.Begin()
	err = ts.Create(&user).Error
	err = ts.Create(&team).Error
	err = ts.Create(&userTeam).Error
	if err != nil {
		logs.Error(err)
		ts.Rollback()
	} else {
		ts.Commit()
	}
	return err
}

func (user *User) LoginByTelephone(telephone string, password string) error {
	mysqlDB := db.GetMySqlConnection().GetMySqlDB()
	fmt.Println("telephone, password", telephone, password)
	err := mysqlDB.Where("telephone = ? and password = ?", telephone, password).First(&user).Error
	// err = mysqlDB.Where("userId = ?", user.UserId).Find(&user.Addresses).Error
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

func (user *User) LoginByWechat(jsCode string, wechatUserProfile string, invitationCode string) (*WxSession, error) {
	wxSession, err := authorization(jsCode, wechatUserProfile)
	if err != nil {
		logs.Error(err)
		return nil, err
	}
	// userId := wxSession.SessionId
	// if invitationCode
	return wxSession, nil
}

func authorization(jsCode string, wechatUserProfile string) (*WxSession, error) {
	// openId, sessionKey, err := JsCode2Session(jsCode)
	// if err != nil {
	// 	return nil, err
	// }
	// fmt.Printf("User openId = %s, sessionKey = %s", openId, sessionKey)
	//
	// h := sha1.New()
	// _, err = io.WriteString(h, sessionKey)
	// if err != nil {
	// 	logs.Error(err)
	// 	return nil, err
	// }
	// sKey := h.Sum(nil)
	//
	// var wxSession *WxSession
	// wxSession.SessionId = GetWuid()
	// wxSession.SessionKey = sessionKey
	// wxSession.OpenId = openId
	// wxSession.Skey = string(sKey)
	// wxSession.WechatUserProfile = wechatUserProfile
	//
	// o := orm.NewOrm()
	// // transaction begin
	// err = o.Begin()
	// _, err = o.InsertOrUpdate(wxSession, "sessionId,openId")
	// if err != nil {
	// 	logs.Error(err)
	// 	_ = o.Rollback()
	// 	return nil, err
	// }
	// user := isAssociated(openId)
	// if user == nil {
	// 	// create user
	// 	var newUser *User
	// 	newUser.UserId = GetWuid()
	// 	newUser.Channel = "Wechat"
	// 	newUser.WxSession = wxSession
	// 	newUser.Status = "active"
	// 	newUser.Role = 10
	//
	// 	_, err = o.Insert(newUser)
	// 	if err != nil {
	// 		logs.Error(err)
	// 		// error rollback
	// 		_ = o.Rollback()
	// 	} else {
	// 		_ = o.Commit()
	// 	}
	// }
	// // query user with wechat information
	// var userSession *WxSession
	// err = o.Raw("SELECT * FROM wxsession t1, user t2 WHERE t1.sessionId = t2.sessionId and t1.sessionId = %s", wxSession.SessionId).QueryRow(&userSession)
	// if err != nil {
	// 	return nil, err
	// }
	// return userSession, nil
	return nil, nil
}

func isAssociated(openId string) *User {
	var user *User
	o := orm.NewOrm()
	err := o.Raw("SELECT * FROM user WHERE openId = %;", openId).QueryRow(user)
	if err != nil {
		return nil
	}
	return user
}

func initTeamForm() (interface{}, error) {
	var team *Team
	team.UserId = GetWuid()
	team.Status = "active"
	team.Channel = "Wechat"
	return nil, nil
}
