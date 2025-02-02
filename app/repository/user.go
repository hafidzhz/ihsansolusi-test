package repository

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/hrshadhin/fiber-go-boilerplate/app/dto"
	"github.com/hrshadhin/fiber-go-boilerplate/app/entity"
	"github.com/hrshadhin/fiber-go-boilerplate/platform/logger"
	"github.com/hrshadhin/fiber-go-boilerplate/shared"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{db}
}

func (repo *UserRepositoryImpl) CreateUser(request *dto.CreateUserRequest) (*entity.User, error) {
	tx := repo.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			logger.GetLogger().Errorf("[CreateUser-01] unexpected error: %s. User Info: Name=%s, IdentityNumber=%s, PhoneNumber=%s",
				r,
				request.Name,
				request.IdentityNumber,
				request.PhoneNumber)
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		logger.GetLogger().Errorf("[CreateUser-02] transaction error: %s. User Info: Name=%s, IdentityNumber=%s, PhoneNumber=%s",
			err,
			request.Name,
			request.IdentityNumber,
			request.PhoneNumber)
		return nil, err
	}

	user := entity.User{
		Name:           request.IdentityNumber,
		PhoneNumber:    request.PhoneNumber,
		IdentityNumber: request.IdentityNumber,
		Balance:        0,
		AccountNumber:  shared.GenerateAccountNumber(),
	}

	err := tx.Create(&user).Error

	if err != nil {
		tx.Rollback()
		var pqErr *pq.Error

		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			re := regexp.MustCompile(`\((.*?)\)`)
			matches := re.FindStringSubmatch(pqErr.Message)

			if len(matches) > 1 {
				logger.GetLogger().Warnf(
					"[CreateUser-03] Unique constraint violation detected. Duplicate entry for field: %s. Value: %s.",
					matches[1],
					getFieldValue(user, matches[1]),
				)
				return nil, fmt.Errorf("duplicate entry for field: %s", matches[1])
			}

			logger.GetLogger().Errorf(
				"[CreateUser-04] Unique constraint violation. Message: %s. User Info: Name=%s, IdentityNumber=%s, PhoneNumber=%s",
				pqErr.Message,
				request.Name,
				request.IdentityNumber,
				request.PhoneNumber,
			)
			return nil, fmt.Errorf("unique constraint violation: %v", pqErr.Message)
		}

		logger.GetLogger().Errorf(
			"[CreateUser-05] Error creating user: %s. User Info: Name=%s, IdentityNumber=%s, PhoneNumber=%s",
			err,
			request.Name,
			request.IdentityNumber,
			request.PhoneNumber,
		)
		return nil, err
	}

	err = tx.Commit().Error

	if err != nil {
		logger.GetLogger().Errorf("[CreateUser-03] error create user: %s. User Info: Name=%s, IdentityNumber=%s, PhoneNumber=%s",
			err,
			request.Name,
			request.IdentityNumber,
			request.PhoneNumber)
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepositoryImpl) DepositToUserAccount(request *dto.DepositRequest) (*entity.User, error) {
	tx := repo.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	result := tx.Model(&entity.User{}).
		Where("account_number = ?", request.AccountNumber).
		UpdateColumn("balance", gorm.Expr("balance + ?", request.Amount))

	if result.Error != nil {
		tx.Rollback()
		logger.GetLogger().Errorf("[DepositToUserAccount-01] error updating balance for account number: %s", request.AccountNumber)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		logger.GetLogger().Warnf("[DepositToUserAccount-02] user not found. account number: %s", request.AccountNumber)
		return nil, fmt.Errorf("user not found with account number: %s", request.AccountNumber)
	}

	var user entity.User
	err := tx.Where("account_number = ?", request.AccountNumber).First(&user).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		logger.GetLogger().Errorf("[DepositToUserAccount-03] error updating balance for account number: %s", request.AccountNumber)
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepositoryImpl) WithdrawFromUserAccount(request *dto.DepositRequest) (*entity.User, error) {
	tx := repo.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	result := tx.Model(&entity.User{}).
		Where("account_number = ?", request.AccountNumber).
		UpdateColumn("balance", gorm.Expr("balance - ?", request.Amount))

	if result.Error != nil {
		tx.Rollback()
		logger.GetLogger().Errorf("[WithdrawFromUserAccount-01] error updating balance for account number: %s", request.AccountNumber)
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		tx.Rollback()
		logger.GetLogger().Warnf("[WithdrawFromUserAccount-02] user not found. account number: %s", request.AccountNumber)
		return nil, fmt.Errorf("user not found with account number: %s", request.AccountNumber)
	}

	var user entity.User
	err := tx.Where("account_number = ?", request.AccountNumber).First(&user).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		logger.GetLogger().Errorf("[WithdrawFromUserAccount-03] error updating balance for account number: %s", request.AccountNumber)
		return nil, err
	}

	return &user, nil
}

func (repo *UserRepositoryImpl) ResolveByAccountNumber(accountNumber string) (user entity.User, err error) {
	err = repo.db.Model(&entity.User{}).
		Where("account_number = ?", accountNumber).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.GetLogger().Warnf("[ResolveByAccountNumber-01] user not found. account number: %s", accountNumber)
			return user, errors.New("user not found")
		}

		logger.GetLogger().Errorf("[ResolveByAccountNumber-02] unexpected error: %s. account number", err, accountNumber)
		return user, err
	}

	return user, nil
}

func getFieldValue(user entity.User, fieldName string) string {
	switch fieldName {
	case "identity_number":
		return user.IdentityNumber
	case "phone_number":
		return user.PhoneNumber
	case "account_number":
		return user.AccountNumber
	default:
		return "Unknown Field"
	}
}
