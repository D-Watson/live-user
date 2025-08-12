package mysql

import (
	"context"

	"live-user/entity"
	"live-user/utils/log"
)

func InsertUser(ctx context.Context, en entity.Users) error {
	res := UserEngineDB.WithContext(ctx).Create(en)
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

func QueryUser(ctx context.Context, id int64, email, phone string) (*entity.Users, error) {
	tx := UserEngineDB.WithContext(ctx)
	en := &entity.Users{}
	if id >= 0 {
		tx = tx.Where("id = ?", id)
	}
	if email != "" {
		tx = tx.Where("email = ?", email)
	}
	if phone != "" {
		tx = tx.Where("phone = ?", phone)
	}
	err := tx.Find(en).Error
	if err != nil {

		return &entity.Users{}, nil
	}
	return &entity.Users{}, nil
}
