package models

import (
	"ecommerce-sys/db"
	. "ecommerce-sys/utils"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"strings"
	"time"
)

type SMS struct {
	UserId        uint64    `json:"userId" gorm:"column:userId; not null;"`
	OperationMode string    `json:"operationMode" gorm:"column:operationMode; not null; type:varchar(10);"`
	SecurityCode  string    `json:"securityCode" gorm:"column:securityCode; type:varchar(6); not null;"`
	Telephone     string    `json:"telephone"`
	ExpiresAt     time.Time `json:"expiresAt" gorm:"column:expiresAt; type:dateTime; not null;"`
	CreatedAt     time.Time `json:"createdAt" gorm:"column:createdAt;type:dateTime"`
}

func (sms *SMS) ObtainSecurityCode(telephone string, userId uint64, operationMode string) (*SMS, error) {
	redisClient := db.GetRedisConnection().GetRedisClient()

	key, value, smsContent := buildSMSContent(telephone, userId, operationMode)
	err := redisClient.Set(key, value, 20*time.Minute).Err()
	if err != nil {
		beego.Error(err.Error())
		return nil, err
	}
	return smsContent, nil
}

func (sms *SMS) VerifySecurityCode(telephone string, userId uint64, securityCode string, operationMode string) (bool, error) {
	redisClient := db.GetRedisConnection().GetRedisClient()

	operationMode = strings.ToUpper(operationMode)
	key := fmt.Sprintf("%d-%s", userId, operationMode)
	resultBytes, err := redisClient.Get(key).Bytes()
	if err != nil {
		beego.Error(err.Error())
		return false, err
	}
	err = json.Unmarshal(resultBytes, &sms)
	if err != nil {
		beego.Error(err.Error())
		return false, err
	}
	if securityCode == sms.SecurityCode {
		return true, nil
	}
	return false, nil
}

func buildSMSContent(telephone string, userId uint64, operationMode string) (key string, value []byte, vo *SMS) {
	operationMode = strings.ToUpper(operationMode)
	key = fmt.Sprintf("%d-%s", userId, operationMode)
	sms := SMS{
		OperationMode: operationMode,
		SecurityCode:  GenerateRandString(6),
		Telephone:     telephone,
		UserId:        userId,
		ExpiresAt:     GenerateAndAddDurationFromNow(15 * time.Minute), // 15 mins
		CreatedAt:     time.Now(),
	}
	result, _ := json.Marshal(sms)
	return key, result, &sms
}
