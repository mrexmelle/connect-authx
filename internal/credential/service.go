package credential

import (
	"github.com/mrexmelle/connect-authx/internal/config"
	"github.com/mrexmelle/connect-authx/internal/mapper"
	"github.com/mrexmelle/connect-authx/internal/profile"
)

type Service struct {
	ConfigService        *config.Service
	CredentialRepository Repository
	ProfileRepository    profile.Repository
}

func NewService(
	cfg *config.Service,
	cr Repository,
	pr profile.Repository,
) *Service {
	return &Service{
		ConfigService:        cfg,
		CredentialRepository: cr,
		ProfileRepository:    pr,
	}
}

func (s *Service) Create(req PostRequestDto) (*ResponseDto, error) {
	trx := s.ConfigService.Db.Begin()
	defer func() {
		if r := recover(); r != nil {
			trx.Rollback()
		}
	}()
	if trx.Error != nil {
		return nil, trx.Error
	}

	err := s.CredentialRepository.CreateWithDb(
		trx,
		req.EmployeeId,
		req.Password,
	)
	if err != nil {
		trx.Rollback()
		return nil, err
	}

	err = s.ProfileRepository.PrepareWithDb(
		trx,
		req.EmployeeId,
	)
	if err != nil {
		trx.Rollback()
		return nil, err
	}

	return &ResponseDto{
		Status: mapper.ToStatus(err),
	}, trx.Commit().Error
}

func (s *Service) UpdatePasswordByEmployeeId(
	employeeId string,
	req PatchPasswordRequestDto,
) (ResponseDto, error) {
	err := s.CredentialRepository.UpdatePasswordByEmployeeIdAndPassword(
		req.NewPassword, employeeId, req.CurrentPassword)
	return ResponseDto{
		Status: mapper.ToStatus(err),
	}, err
}

func (s *Service) ResetPasswordByEmployeeId(employeeId string) error {
	return s.CredentialRepository.ResetPasswordByEmployeeId(employeeId)
}

func (s *Service) DeleteByEmployeeId(employeeId string) error {
	return s.CredentialRepository.DeleteByEmployeeId(employeeId)
}
