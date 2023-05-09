package usecase

import (
	"biocad/internal/entity"
	"context"
)

type UserUseCase struct {
	repo UserRp
}

var _ UserContract = (*UserUseCase)(nil)

func NewUserUseCase(repo UserRp) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}

func (u UserUseCase) GetAllUniqueGuidList(ctx context.Context) ([]string, error) {
	list, err := u.repo.GetAllGuidList(ctx)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func (u UserUseCase) GetAllDataByGuid(ctx context.Context, guid string, limit uint64, page uint64) ([]entity.Data, error) {
	exists, err := u.repo.CheckGuidExists(ctx, guid)
	if err != nil {
		return nil, err
	}
	if exists {
		offset := limit * (page - 1)
		data, err := u.repo.GetDataByUnitGuid(ctx, guid, limit, offset)
		if err != nil {
			return nil, err
		}
		return data, nil
	} else {
		return nil, nil
	}
}
