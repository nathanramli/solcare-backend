package services

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
	"github.com/nathanramli/solcare-backend/common"
	"github.com/nathanramli/solcare-backend/config"
	"github.com/nathanramli/solcare-backend/httpserver/controllers/params"
	"github.com/nathanramli/solcare-backend/httpserver/controllers/views"
	"github.com/nathanramli/solcare-backend/httpserver/repositories"
	"github.com/nathanramli/solcare-backend/httpserver/repositories/models"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type userSvc struct {
	repo repositories.UserRepo
}

func NewUserSvc(repo repositories.UserRepo) UserSvc {
	return &userSvc{
		repo: repo,
	}
}

func (svc *userSvc) Login(ctx context.Context, user *params.Login) *views.Response {
	pubkey, err := solana.PublicKeyFromBase58(user.Address)
	if err != nil {
		return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, err)
	}

	signatureBytes, err := base58.Decode(user.SignedMessage)
	if err != nil {
		return views.ErrorResponse(http.StatusBadRequest, views.M_SIGNATURE_INVALID, err)
	}

	if !solana.SignatureFromBytes(signatureBytes).Verify(pubkey, common.LoginMessage) {
		return views.ErrorResponse(http.StatusBadRequest, views.M_SIGNATURE_INVALID, errors.New("signature is invalid"))
	}

	_, err = svc.repo.FindUserByAddress(ctx, user.Address)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = svc.repo.CreateUser(ctx, &models.Users{
				WalletAddress: user.Address,
			})

			if err != nil {
				return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
			}
		} else {
			return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
		}
	}

	claims := &common.CustomClaims{
		Address: user.Address,
	}

	claims.ExpiresAt = time.Now().Add(time.Minute * time.Duration(config.GetJwtExpiredTime())).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(config.GetJwtSignature())

	return views.SuccessResponse(http.StatusOK, views.M_OK, views.Login{
		Token: ss,
	})
}

func (svc *userSvc) UpdateUser(ctx context.Context, address string, params *params.UpdateUser) *views.Response {
	user, err := svc.repo.FindUserByAddress(ctx, address)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, err)
		} else {
			return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
		}
	}

	user.FirstName = params.FirstName
	user.LastName = params.LastName
	user.Email = params.Email
	user.Gender = *params.Gender

	err = svc.repo.UpdateUser(ctx, user)
	if err != nil {
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	return views.SuccessResponse(http.StatusOK, views.M_OK, nil)
}

func (svc *userSvc) FindUserByAddress(ctx context.Context, address string) *views.Response {
	user, err := svc.repo.FindUserByAddress(ctx, address)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, err)
		}
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	return views.SuccessResponse(http.StatusOK, views.M_OK, views.FindUser{
		Address:        user.WalletAddress,
		CreatedAt:      user.CreatedAt.Unix(),
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Gender:         user.Gender,
		IsVerified:     user.IsVerified,
		IsWarned:       user.IsWarned,
		ProfilePicture: user.ProfilePicture,
	})
}
