package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
)

// orm introduction: https://my.oschina.net/u/252343/blog/829912

type User struct {
	UserId    string     `json:"userId" orm:"column(userId);PK;unique;size(32)"`
	Channel   string     `json:"channel" orm:"column(channel);size(12);null"`
	WxSession *WxSession `orm:"column(openId);rel(one)"`
	UserProfile
	BaseModel
}

type UserProfile struct {
	Username  string `json:"username" orm:"column(username);size(18)"`
	Password  string `json:"password" orm:"column(password);size(24)"`
	Nickname  string `json:"nickname" orm:"column(nickname);size(16);"`
	Telephone string `json:"telephone" orm:"column(telephone);size(11)"`
	Male      bool   `json:"male" orm:"column(male);default(false)"`
	Signature string `json:"signature" orm:";default(This guy is lazy...)"`
}

type WxSession struct {
	SessionId         string `json:"sessionId" orm:"column(sessionId);pk;unique"`
	Skey              string `json:"skey" orm:"column(skey)"`
	SessionKey        string `json:"session_key" orm:"column(sessionKey)"`
	WechatUserProfile string `json:"wechatUserProfile" orm:"column(wechatUserProfile)"`
	OpenId            string `json:"openId" orm:"column(openId);index"`
	User              *User  `orm:"reverse(one)"`
}

// 自定义表名
func (ws *WxSession) TableName() string {
	return "wxsession"
}

type IUserOperation interface {
	Register() error
	CheckIsUserExist(userId string, telephone string) bool
	QueryByUserId(userId string) *User
	LoginByTelephone(telephone string, password string) (error, *User)
}

func (user *User) Register() error {
	o := orm.NewOrm()
	_, err := o.Insert(&user.UserProfile)
	if err != nil {
		return err
	}
	return nil
}

func (user *User) LoginByTelephone(telephone string, password string) (error, *User) {
	o := orm.NewOrm()
	result, err := o.Raw("select * from user where telephone = ? and password = ?;", telephone, password).Exec()
	if err != nil {
		return err, nil
	}
	if result == nil {
		return errors.New("telephone or password is invalid"), nil
	}
	return nil, user
}

func (user *User) CheckIsUserExist(userId string, telephone string) bool {
	panic("implement me")
}

func (user *User) QueryByUserId(userId string) *User {
	panic("implement me")
}
