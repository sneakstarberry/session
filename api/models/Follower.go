package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

//Follow when we liked post maked that post
type Follow struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UserAID   uint32    `gorm:"not null" json:"usera_id"`
	UserBID   uint32    `gorm:"not null" json:"userb_id"`
	CreatedAt time.Time `gorm:"default_CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default_CURRENT_TIMESTAMP" jsont:"updated_at"`
}

//SaveFollow Save like post
func (f *Follow) SaveFollow(db *gorm.DB) (*Follow, error) {
	// Check if the auth user has liked this post before:
	err := db.Debug().Model(&Follow{}).Where("userb_id =? AND usera_id = ?", f.UserBID, f.UserAID).Take(&f).Error
	if err != nil {
		if err.Error() == "record not found" {
			// the user has not liked this post before, so lets save incomming like:
			err = db.Debug().Model(&Follow{}).Create(&f).Error
			if err != nil {
				return &Follow{}, err
			}
		}
	} else {
		err = errors.New("double like")
		return &Follow{}, err
	}
	return f, nil
}

//DeleteFollow is delete like
func (f *Follow) DeleteFollow(db *gorm.DB) (*Follow, error) {
	var err error
	var deletedFollow *Follow

	err = db.Debug().Model(Follow{}).Where("id = ?", f.ID).Take(&f).Error
	if err != nil {
		return &Follow{}, err
	} else {
		//If the like exist, save it in deleted like and delete it
		deletedFollow = f
		db = db.Debug().Model(&Follow{}).Where("id = ?", f.ID).Take(&Follow{}).Delete(&Follow{})
		if db.Error != nil {
			fmt.Println("cant delete like: ", db.Error)
			return &Follow{}, db.Error
		}
	}
	return deletedFollow, nil
}

//GetFollowsInfo get what we like post information
func (f *Follow) GetFollowsInfo(db *gorm.DB, pid uint64) (*[]Follow, error) {
	follows := []Follow{}
	err := db.Debug().Model(&Follow{}).Where("usera_id=?", pid).Find(&follows).Error
	if err != nil {
		return &[]Follow{}, err
	}
	return &follows, err
}

//DeleteUserAFollows When a post is deleted, we also delete the likes that the post had
func (f *Follow) DeleteUserAFollows(db *gorm.DB, uid uint32) (int64, error) {
	follows := []Follow{}
	db = db.Debug().Model(&Follow{}).Where("usera_id =?", uid).Find(&follows).Delete(&follows)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

//DeleteUserBFollows When a post is deleted, we also delete the likes that the post had
func (f *Follow) DeleteUserBFollows(db *gorm.DB, uid uint64) (int64, error) {
	follows := []Follow{}
	db = db.Debug().Model(&Follow{}).Where("userb_id = ?", uid).Find(&follows).Delete(&follows)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
