package models

import (
	. "ecommerce-sys/models"
	"testing"
)

func TestUser_CheckIsUserExist(t *testing.T) {
	type fields struct {
		UserId      uint64
		UserProfile UserProfile
		Role        uint16
		Status      string
		Channel     string
		BaseModel   BaseModel
	}
	type args struct {
		userId    string
		telephone string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := &User{
				UserId:      tt.fields.UserId,
				UserProfile: tt.fields.UserProfile,
				Role:        tt.fields.Role,
				Status:      tt.fields.Status,
				Channel:     tt.fields.Channel,
				BaseModel:   tt.fields.BaseModel,
			}
			got, err := user.CheckIsUserExistByUserId(tt.args.userId, tt.args.telephone)
			if (err != nil) != tt.wantErr {
				t.Errorf("User.CheckIsUserExist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("User.CheckIsUserExist() = %v, want %v", got, tt.want)
			}
		})
	}
}
