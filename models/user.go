package models

import (
	. "ecommerce-sys/utils"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// orm introduction: https://my.oschina.net/u/252343/blog/829912

type User struct {
	UserId uint64 `json:"userId" gorm:"column:userId;primary_key;not null"`
	UserProfile
	Role    uint16 `json:"role" gorm:"column:role;default:10"`
	Status  string `json:"status" gorm:"column:status;type:varchar(10);default:'active';"`
	Channel string `json:"channel" gorm:"column:channel;type:varchar(12)"`
	// `json:"wxSession" gorm:"column:sessionId"`
	// WxSession WxSession
	// Address   []*Address `orm:"reverse(many)"`
	// Team      *Team      `json:"teamId" gorm:"column:teamId;reverse(one)"`
	BaseModel
}

type UserProfile struct {
	Telephone string `json:"telephone" gorm:"column:telephone; type:varchar(11);not null"`
	Username  string `json:"username" gorm:"column:username; type:varchar(18);not null"`
	Password  string `json:"password" gorm:"column:password; type:varchar(24);not null"`
	Nickname  string `json:"nickname" gorm:"column:nickname; type:varchar(16);not null;"`
	Male      bool   `json:"male" gorm:"column:male; default:false;"`
	Signature string `json:"signature" gorm:"not null; default:'This guy is lazy...'"`
}

type WxSession struct {
	SessionId         uint64 `json:"sessionId" gorm:"column:sessionId; primary_key;"`
	Skey              string `json:"skey" gorm:"column:skey"`
	SessionKey        string `json:"session_key" gorm:"column:sessionKey" `
	WechatUserProfile string `json:"wechatUserProfile" gorm:"column:wechatUserProfile"`
	OpenId            string `json:"openId" gorm:"column:openId; index"`
	User              User   `gorm:"foreignkey:userId"`
	BaseModel
}

type Address struct {
	AddressId    uint64 `json:"addressId" orm:"column(addressId);PK;unique;size(64)"`
	Contact      string `json:"contact" orm:"column(contact);size(32)"`
	Telephone    string `json:"telephone" orm:"column(telephone);size(15)"`
	IsDefault    bool   `json:"isDefault" orm:"column(isDefault);default(false)"`
	Country      string `json:"country" orm:"column(country);null"`
	ProvinceCity string `json:"city" orm:"column(city)"`
	Status       string `json:"status" orm:"column(status);size(10);default(inactive);on_delete(set_default)"`
	User         *User  `json:"user" orm:"column(userId);rel(fk)"`
	BaseModel
}

type Team struct {
	TeamId         uint64 `json:"teamId" orm:"column(teamId);PK;unique;size(64)"`
	TopLevelAgent  uint64 `json:"topLevelAgent" orm:"column(topLevelAgent)"`
	SuperiorAgent  uint64 `json:"superiorAgent" orm:"column(superiorAgent)"`
	Status         string `json:"status" orm:"column(status);default(inactive);on_delete(set_default);"`
	Channel        string `json:"channel" orm:"column(channel);default(Wechat);description(The channel which is user use)"`
	InvitationCode string `json:"invitationCode" orm:"column(invitationCode);unique;size(6);description(This is de unique code about team invitation code)"`
	User           *User  `json:"userId" orm:"column(userId);rel(one)"`
	BaseModel
}

// 自定义表名
func (WxSession) TableName() string {
	return "wxsession"
}

func init() {

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
