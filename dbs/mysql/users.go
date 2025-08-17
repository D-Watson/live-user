package mysql

import (
	"context"

	"github.com/D-Watson/live-safety/log"
	"live-user/dbs"
	"live-user/entity"
)

func InsertUser(ctx context.Context, en entity.Users) error {
	res := dbs.MysqlEngine.WithContext(ctx).Create(en)
	if err := res.Error; err != nil {
		log.Errorf(ctx, "[DB] create error=%v", err)
		return err
	}
	return nil
}

func UpdateUser(ctx context.Context) error {
	return nil
}

func DeleteUser(ctx context.Context) error {
	return nil
}

func QueryUser(ctx context.Context, en *entity.Users) (*entity.Users, error) {
	tx := dbs.MysqlEngine.WithContext(ctx)
	if en.Id >= 0 {
		tx = tx.Where("id = ?", en.Id)
	}
	if en.Email != "" {
		tx = tx.Where("email = ?", en.Email)
	}
	if en.Phone != "" {
		tx = tx.Where("phone = ?", en.Phone)
	}
	if en.PasswordHash != "" {
		tx = tx.Where("password_hash = ?", en.PasswordHash)
	}
	err := tx.Find(&en).Error
	if err != nil {
		log.Errorf(ctx, "[mysql] error,err=", err)
		return &entity.Users{}, nil
	}
	return en, nil
}
