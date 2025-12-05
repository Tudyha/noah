package dao

import (
	"context"
	"noah/internal/model"

	"github.com/samber/lo"
	"gorm.io/gorm"
)

type workSpaceDao struct {
	db *gorm.DB
}

func newWorkSpaceDao(db *gorm.DB) WorkSpaceDao {
	return &workSpaceDao{db: db}
}

func (s *workSpaceDao) Create(ctx context.Context, name string, description string) (*model.WorkSpace, error) {
	space := &model.WorkSpace{
		Name:        name,
		Description: description,
		Status:      1,
	}
	err := s.db.Create(space).Error
	return space, err
}

func (s *workSpaceDao) CreateSpaceUser(ctx context.Context, spaceId uint64, userId uint64, role int) error {
	return s.db.Create(&model.WorkSpaceUser{
		SpaceID: spaceId,
		UserID:  userId,
		Role:    role,
	}).Error
}

func (s *workSpaceDao) CreateApp(ctx context.Context, spaceId uint64, name string, description string) error {
	return s.db.Create(&model.WorkSpaceApp{
		SpaceID:     spaceId,
		Name:        name,
		Description: description,
		Status:      1,
	}).Error
}

func (s *workSpaceDao) GetByUserID(ctx context.Context, userID uint64) ([]*model.WorkSpace, error) {
	var list []*model.WorkSpaceUser
	err := s.db.Where("user_id = ?", userID).Find(&list).Error
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	var spaces []*model.WorkSpace
	spaceIds := lo.Map(list, func(item *model.WorkSpaceUser, index int) uint64 {
		return item.SpaceID
	})
	return spaces, s.db.Where("id IN ?", spaceIds).Find(&spaces).Error
}

func (s *workSpaceDao) GetAppBySpaceIDs(ctx context.Context, spaceIDs []uint64) ([]*model.WorkSpaceApp, error) {
	var list []*model.WorkSpaceApp
	return list, s.db.Where("space_id IN ?", spaceIDs).Find(&list).Error
}
