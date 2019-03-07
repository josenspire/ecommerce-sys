package models

import (
	. "ecommerce-sys/utils"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// orm introduction: https://my.oschina.net/u/252343/blog/829912

type User struct {
	UserId uint64 `json:"userId" orm:"column(userId);PK;unique;size(64)"`
	UserProfile
	Role    uint16 `json:"role" orm:"column(role);default(10)"`
	Status  string `json:"status" orm:"column(status);size(10);default(active)"`
	Channel string `json:"channel" orm:"column(channel);size(12);null"`
	BaseModel
}

type UserProfile struct {
	Telephone string `json:"telephone" orm:"column(telephone);size(11)"`
	Username  string `json:"username" orm:"column(username);size(18)"`
	Password  string `json:"password" orm:"column(password);size(24)"`
	Nickname  string `json:"nickname" orm:"column(nickname);size(16);"`
	Male      bool   `json:"male" orm:"column(male);default(false)"`
	Signature string `json:"signature" orm:";default(This guy is lazy...)"`
}

type WxSession struct {
	SessionId         string `json:"sessionId" orm:"column(sessionId);pk;unique"`
	Skey              string `json:"skey" orm:"column(skey)"`
	SessionKey        string `json:"session_key" orm:"column(sessionKey)"`
	WechatUserProfile string `json:"wechatUserProfile" orm:"column(wechatUserProfile)"`
	OpenId            string `json:"openId" orm:"column(openId);index"`
	User              *User  `orm:"column(userId);rel(one)"`
}

// 自定义表名
func (ws *WxSession) TableName() string {
	return "wxsession"
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
	isExist, err := user.CheckIsUserExistByTelephone(user.Telephone)
	if isExist == true {
		return ErrCurrentUserIsExist
	}

	user.UserId = GetWuid()
	user.Status = "active"
	o := orm.NewOrm()
	_, err = o.Insert(user)
	return err
}

func (user *User) LoginByTelephone(telephone string, password string) error {
	o := orm.NewOrm()
	fmt.Println("telephone, password", telephone, password)
	err := o.Raw("SELECT * FROM user WHERE telephone = ? and password = ?;", telephone, password).QueryRow(user)
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
	err := o.Raw("SELECT COUNT(*) FROM user").QueryRow(&total)
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

func (user *User) LoginByWechat(jsCode string, userInfo string, invitationCode string) (interface{}, error) {
	panic("implement me")
}
