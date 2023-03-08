package services

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gagliardetto/solana-go"
	"github.com/gin-gonic/gin"
	"github.com/mr-tron/base58"
	"github.com/nathanramli/solcare-backend/common"
	"github.com/nathanramli/solcare-backend/config"
	"github.com/nathanramli/solcare-backend/httpserver/controllers/params"
	"github.com/nathanramli/solcare-backend/httpserver/controllers/views"
	"github.com/nathanramli/solcare-backend/httpserver/repositories"
	"github.com/nathanramli/solcare-backend/httpserver/repositories/models"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

type userSvc struct {
	repo         repositories.UserRepo
	kycQueueRepo repositories.KycQueueRepo
	adminRepo    repositories.AdminRepo
}

func NewUserSvc(repo repositories.UserRepo, kycQueueRepo repositories.KycQueueRepo, adminRepo repositories.AdminRepo) UserSvc {
	return &userSvc{
		repo:         repo,
		kycQueueRepo: kycQueueRepo,
		adminRepo:    adminRepo,
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

	admin, err := svc.adminRepo.FindAdminByAddress(ctx, user.Address)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
		}
		admin = nil
	}

	claims := &common.CustomClaims{
		Address: user.Address,
		IsAdmin: admin != nil,
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
	user.Gender = params.Gender

	err = svc.repo.UpdateUser(ctx, user)
	if err != nil {
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	return views.SuccessResponse(http.StatusOK, views.M_OK, views.FindUser{
		Address:        user.WalletAddress,
		CreatedAt:      user.CreatedAt.Unix(),
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Gender:         *user.Gender,
		IsVerified:     *user.IsVerified,
		IsWarned:       *user.IsWarned,
		ProfilePicture: user.ProfilePicture,
	})
}

func (svc *userSvc) VerifyKyc(ctx context.Context, params *params.VerifyKyc) *views.Response {
	kycQueue, err := svc.kycQueueRepo.FindKycRequestByUser(ctx, params.Address)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, err)
		} else {
			return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
		}
	}

	if kycQueue.Status != models.KYC_STATUS_REQUESTED {
		return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, errors.New("kyc request is not in pending state"))
	}

	if *params.IsAccepted {
		kycQueue.Status = models.KYC_STATUS_APPROVED

		user, err := svc.repo.FindUserByAddress(ctx, params.Address)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, err)
			} else {
				return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
			}
		}
		user.IsVerified = params.IsAccepted
		err = svc.repo.UpdateUser(ctx, user)
		if err != nil {
			return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
		}
	} else {
		kycQueue.Status = models.KYC_STATUS_DECLINED
	}

	err = svc.kycQueueRepo.SaveKycQueue(ctx, kycQueue)
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
		Gender:         *user.Gender,
		IsVerified:     *user.IsVerified,
		IsWarned:       *user.IsWarned,
		ProfilePicture: user.ProfilePicture,
	})
}

func (svc *userSvc) UpdateAvatar(ctx context.Context, address string, params *params.UpdateUserAvatar) *views.Response {
	user, err := svc.repo.FindUserByAddress(ctx, address)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, err)
		} else {
			return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
		}
	}

	fileNameSplits := strings.Split(params.Picture.Filename, ".")
	ext := fileNameSplits[len(fileNameSplits)-1]

	err = ctx.(*gin.Context).SaveUploadedFile(params.Picture, "./resources/avatar_"+address+"."+ext)
	if err != nil {
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	user.ProfilePicture = "avatar_" + address + "." + ext
	err = svc.repo.UpdateUser(ctx, user)
	if err != nil {
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	return views.SuccessResponse(http.StatusOK, views.M_OK, views.FindUser{
		Address:        user.WalletAddress,
		CreatedAt:      user.CreatedAt.Unix(),
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		Gender:         *user.Gender,
		IsVerified:     *user.IsVerified,
		IsWarned:       *user.IsWarned,
		ProfilePicture: user.ProfilePicture,
	})
}

func (svc *userSvc) FindKycRequestByUser(ctx context.Context, address string) *views.Response {
	kycQueue, err := svc.kycQueueRepo.FindKycRequestByUser(ctx, address)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return views.SuccessResponse(http.StatusOK, views.M_NOT_FOUND, nil)
		}

		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	return views.SuccessResponse(http.StatusOK, views.M_OK, views.FindKycRequest{
		Id:                      kycQueue.Id,
		Nik:                     kycQueue.Nik,
		Name:                    kycQueue.Users.FirstName + kycQueue.Users.LastName,
		UsersWalletAddress:      kycQueue.UsersWalletAddress,
		RequestedAt:             kycQueue.RequestedAt.Unix(),
		IdCardPicture:           kycQueue.IdCardPicture,
		FacePicture:             kycQueue.FacePicture,
		SelfieWithIdCardPicture: kycQueue.SelfieWithIdCardPicture,
		Status:                  kycQueue.Status,
	})
}

func (svc *userSvc) FindAllUsers(ctx context.Context) *views.Response {
	users, err := svc.repo.FindAllUsers(ctx)
	if err != nil {
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	resp := make([]views.FindUser, len(users))
	for i, user := range users {
		r := views.FindUser{
			Address:        user.WalletAddress,
			CreatedAt:      user.CreatedAt.Unix(),
			Email:          user.Email,
			FirstName:      user.FirstName,
			LastName:       user.LastName,
			Gender:         *user.Gender,
			IsVerified:     *user.IsVerified,
			IsWarned:       *user.IsWarned,
			ProfilePicture: user.ProfilePicture,
		}
		resp[i] = r
	}
	return views.SuccessResponse(http.StatusOK, views.M_OK, resp)
}

func (svc *userSvc) FindAllKycRequest(ctx context.Context, status int) *views.Response {
	kycQueues, err := svc.kycQueueRepo.FindAllKycRequest(ctx, status)
	if err != nil {
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	resp := make([]views.FindKycRequest, len(kycQueues))
	for i, kycQueue := range kycQueues {
		r := views.FindKycRequest{
			Id:                      kycQueue.Id,
			Nik:                     kycQueue.Nik,
			Name:                    kycQueue.Users.FirstName + kycQueue.Users.LastName,
			UsersWalletAddress:      kycQueue.UsersWalletAddress,
			RequestedAt:             kycQueue.RequestedAt.Unix(),
			IdCardPicture:           kycQueue.IdCardPicture,
			FacePicture:             kycQueue.FacePicture,
			SelfieWithIdCardPicture: kycQueue.SelfieWithIdCardPicture,
			Status:                  kycQueue.Status,
		}
		resp[i] = r
	}
	return views.SuccessResponse(http.StatusOK, views.M_OK, resp)
}

func (svc *userSvc) RequestKyc(ctx context.Context, address string, params *params.RequestKyc) *views.Response {
	user, err := svc.repo.FindUserByAddress(ctx, address)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, err)
		} else {
			return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
		}
	}

	recentKycQueue, err := svc.kycQueueRepo.FindKycRequestByUser(ctx, address)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			recentKycQueue = nil
		} else {
			return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
		}
	}

	if recentKycQueue != nil && (recentKycQueue.Status == models.KYC_STATUS_REQUESTED || recentKycQueue.Status == models.KYC_STATUS_APPROVED) {
		return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, errors.New("you are not allowed to request another kyc in this phase"))
	}

	kycQueue := &models.KycQueues{
		RequestedAt:        time.Now(),
		UsersWalletAddress: user.WalletAddress,
		Status:             models.KYC_STATUS_REQUESTED,
		Nik:                params.Nik,
	}

	if recentKycQueue != nil {
		kycQueue.Id = recentKycQueue.Id
	}

	// save face picture
	fileNameSplits := strings.Split(params.Face.Filename, ".")
	ext := fileNameSplits[len(fileNameSplits)-1]

	name := "kyc_face_" + address + "." + ext
	err = ctx.(*gin.Context).SaveUploadedFile(params.Face, "./resources/"+name)
	if err != nil {
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}
	kycQueue.FacePicture = name

	// save id card picture
	fileNameSplits = strings.Split(params.IdCard.Filename, ".")
	ext = fileNameSplits[len(fileNameSplits)-1]

	name = "kyc_id_" + address + "." + ext
	err = ctx.(*gin.Context).SaveUploadedFile(params.IdCard, "./resources/"+name)
	if err != nil {
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}
	kycQueue.IdCardPicture = name

	// save selfie with id card picture
	fileNameSplits = strings.Split(params.FaceWithIdCard.Filename, ".")
	ext = fileNameSplits[len(fileNameSplits)-1]

	name = "kyc_selfie_with_id_" + address + "." + ext
	err = ctx.(*gin.Context).SaveUploadedFile(params.FaceWithIdCard, "./resources/"+name)
	if err != nil {
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}
	kycQueue.SelfieWithIdCardPicture = name

	// update to the database
	err = svc.kycQueueRepo.SaveKycQueue(ctx, kycQueue)
	if err != nil {
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	return views.SuccessResponse(http.StatusOK, views.M_OK, nil)
}

func (svc *userSvc) RemoveKyc(ctx context.Context, address string) *views.Response {
	user, err := svc.repo.FindUserByAddress(ctx, address)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, err)
		} else {
			return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
		}
	}

	if *user.IsVerified == false {
		return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, errors.New("user is not verified"))
	}

	recentKycQueue, err := svc.kycQueueRepo.FindKycRequestByUser(ctx, address)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return views.ErrorResponse(http.StatusBadRequest, views.M_BAD_REQUEST, errors.New("user never sent kyc request"))
		} else {
			return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
		}
	}

	recentKycQueue.Status = models.KYC_STATUS_REMOVED
	off := false
	user.IsVerified = &off

	err = svc.kycQueueRepo.SaveKycQueue(ctx, recentKycQueue)
	if err != nil {
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	err = svc.repo.UpdateUser(ctx, user)
	if err != nil {
		return views.ErrorResponse(http.StatusInternalServerError, views.M_INTERNAL_SERVER_ERROR, err)
	}

	return views.SuccessResponse(http.StatusOK, views.M_OK, nil)
}
