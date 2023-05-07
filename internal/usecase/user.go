package usecase

type UserUseCase struct {
	repo UserRp
}

var _ UserContract = (*UserUseCase)(nil)

func NewUserUseCase(repo UserRp) *UserUseCase {
	return &UserUseCase{
		repo: repo,
	}
}
