package models

import (
	. "ecommerce-sys/utils"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type User struct {
	UserId uint64 `json:"userId" gorm:"column:userId;primary_key;not null"`
	UserProfile
	Role    uint16 `json:"role" gorm:"column:role; default:10; not null;"`
	Status  string `json:"status" gorm:"column:status; type:varchar(10); default:'active'; not null;"`
	Channel string `json:"channel" gorm:"column:channel; type:varchar(12); not null;"`
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
	UserId            uint64 `json:"userId" gorm:"column:userId; not null;"`
	BaseModel
}

type Address struct {
	AddressId    uint64 `json:"addressId" gorm:"column:addressId; primary_key; not null;"`
	Contact      string `json:"contact" gorm:"column:contact; type:varchar(32); not null;"`
	Telephone    string `json:"telephone" gorm:"column:telephone; type:varchar(15); not null;"`
	IsDefault    bool   `json:"isDefault" gorm:"column:isDefault; default:false; not null;"`
	Country      string `json:"country" gorm:"column:country; not null;"`
	ProvinceCity string `json:"city" gorm:"column:city; not null;"`
	Status       string `json:"status" gorm:"column:status; type:varchar(10); default:'inactive'; not null;"`
	UserId       uint64 `json:"userId" gorm:"column:userId; not null;"`
	BaseModel
}

type Team struct {
	TeamId         uint64 `json:"teamId" gorm:"column:teamId; primary_key; not null;"`
	TopLevelAgent  uint64 `json:"topLevelAgent" gorm:"column:topLevelAgent; not null;"`
	SuperiorAgent  uint64 `json:"superiorAgent" gorm:"column:superiorAgent; not null;"`
	Status         string `json:"status" gorm:"column:status; default:'active'; not null;"`
	Channel        string `json:"channel" gorm:"column:channel; default:'Wechat'; not null;"`
	InvitationCode string `json:"invitationCode" gorm:"column:invitationCode; unique; type:varchar(6); not null;"`
	UserId         uint64 `json:"userId" gorm:"column:userId; not null;"`
	BaseModel
}

// 自定义表名
func (WxSession) TableName() string {
	return "wxsessions"
}

type IUserOperation interface {
	Register() error
	CheckIsUserExistByUserId(userId uint64) (bool, error)
	CheckIsUserExistByTelephone(telephone string) (bool, error)
	QueryByUserId(userId string) *User
	LoginByTelephone(telephone string, password string) error
	LoginByWechat(jsCode string, userInfo string, invitationCode string) (interface{}, error)
}

func (user *User) Register() error {
	// isExist, err := user.CheckIsUserExistByTelephone(user.Telephone)
	// if isExist == true {
	// 	return ErrCurrentUserIsExist
	// }
	//
	// user.UserId = GetWuid()
	// user.Status = "active"
	// o := orm.NewOrm()
	// _, err = o.Insert(user)
	// if err != nil {
	// 	logs.Error(err)
	// }
	// return err
	return nil
}

func (user *User) LoginByTelephone(telephone string, password string) error {
	o := orm.NewOrm()
	fmt.Println("telephone, password", telephone, password)
	err := o.Raw("SELECT * FROM user WHERE user.telephone = ? and user.password = ?;", telephone, password).QueryRow(user)
	if err == orm.ErrNoRows {
		return ErrTelOrPswInvalid
	}
	return nil
}

func (user *User) CheckIsUserExistByUserId(userId uint64) (bool, error) {
	if userId == 0 {
		logs.Warn(WarnParamsMissing.Error())
		return false, WarnParamsMissing
	}
	o := orm.NewOrm()
	queryUser := User{UserId: userId}
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
	o := orm.NewOrm()
	var total uint64
	err := o.Raw("SELECT COUNT(1) FROM user WHERE telephone = ?", telephone).QueryRow(&total)
	if err != nil {
		return false, err
	}
	if total > 0 {
		return true, nil
	}
	return false, nil
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
	team.TeamId = GetWuid()
	team.Status = "active"
	team.Channel = "Wechat"
	return nil, nil
}
