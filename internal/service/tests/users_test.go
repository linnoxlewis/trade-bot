package tests

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/linnoxlewis/trade-bot/config"
	"github.com/linnoxlewis/trade-bot/internal/domain/consts"
	srvErr "github.com/linnoxlewis/trade-bot/internal/errors"
	"github.com/linnoxlewis/trade-bot/internal/helper"
	"github.com/linnoxlewis/trade-bot/internal/service"
	mock_service "github.com/linnoxlewis/trade-bot/internal/service/tests/repo_mocks"
	"github.com/linnoxlewis/trade-bot/pkg/i18n"
	"github.com/linnoxlewis/trade-bot/pkg/log"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserService_CreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_service.NewMockUserRepo(ctrl)
	mockLogger := log.NewLogger()
	cfg := &config.Config{}
	i18nSrv := i18n.NewI18n("../../../"+consts.LocaleDataPath, helper.GetLanguageList())
	userService := service.NewUserService(cfg, mockUserRepo, i18nSrv, mockLogger)

	testCases := []struct {
		name              string
		existUser         bool
		createUserError   error
		expectedError     error
		expectedErrorCode int
	}{
		{
			name:              "Create user success",
			existUser:         false,
			createUserError:   nil,
			expectedError:     nil,
			expectedErrorCode: 200,
		},
		{
			name:              "User already exists",
			existUser:         true,
			createUserError:   nil,
			expectedError:     srvErr.BadRequestError(i18nSrv.T("userAlreadyExist", nil, helper.GetDefaultLg())),
			expectedErrorCode: 400,
		},
		{
			name:              "Error creating user",
			existUser:         false,
			createUserError:   errors.New("error creating user"),
			expectedError:     srvErr.InternalServerError(errors.New("error creating user")),
			expectedErrorCode: 500,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println("start test:", tc.name)
			mockUserRepo.EXPECT().ExistUser(gomock.Any(), gomock.Any()).Return(tc.existUser)
			mockUserRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(tc.createUserError)

			err := userService.CreateUser(context.Background(), "username", 123)
			if err != nil {
				resultErr := err.(srvErr.Error)
				assert.Equal(t, tc.expectedError, resultErr)
				assert.Equal(t, tc.expectedErrorCode, resultErr.GetCode())
			}
		})
	}
}

/*
func TestUserService_IsAdmin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_service.NewMockUserRepo(ctrl)
	mockLogger := log.NewLogger()
	cfg := &config.Config{}
	i18nSrv := i18n.NewI18n("../../../"+consts.LocaleDataPath, helper.GetLanguageList())
	userService := service.NewUserService(cfg, mockUserRepo, i18nSrv, mockLogger)

	testCases := []struct {
		name              string
		existUser         bool
		isAdmin           bool
		isAdminError      error
		expectedIsAdmin   bool
		expectedError     bool
		expectedErrorCode string
	}{
		{
			name:              "Is admin success",
			existUser:         true,
			isAdmin:           true,
			isAdminError:      nil,
			expectedIsAdmin:   true,
			expectedError:     false,
			expectedErrorCode: "",
		},
		{
			name:              "User doesn't exist",
			existUser:         false,
			isAdmin:           false,
			isAdminError:      nil,
			expectedIsAdmin:   false,
			expectedError:     false,
			expectedErrorCode: "",
		},
		{
			name:              "Error getting admin status",
			existUser:         true,
			isAdmin:           false,
			isAdminError:      errors.New("error getting admin status"),
			expectedIsAdmin:   false,
			expectedError:     true,
			expectedErrorCode: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepo.EXPECT().ExistUser(gomock.Any(), gomock.Any()).Return(tc.existUser)
			if tc.existUser {
				mockUserRepo.EXPECT().IsAdmin(gomock.Any(), gomock.Any()).Return(tc.isAdmin, tc.isAdminError)
			}
			isAdmin := userService.IsAdmin(context.Background(), 123)
			assert.Equal(t, tc.expectedIsAdmin, isAdmin)
		})
	}
}

func TestUserService_GetAdmins(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_service.NewMockUserRepo(ctrl)
	mockLogger := log.NewLogger()
	i18nSrv := i18n.NewI18n("../../../"+consts.LocaleDataPath, helper.GetLanguageList())
	userService := service.NewUserService(&config.Config{}, mockUserRepo, i18nSrv, mockLogger)

	testCases := []struct {
		name              string
		adminIds          []int
		getAdminsError    error
		expectedAdminIds  []int
		expectedError     bool
		expectedErrorCode string
	}{
		{
			name:              "Get admins success",
			adminIds:          []int{1, 2, 3},
			getAdminsError:    nil,
			expectedAdminIds:  []int{1, 2, 3},
			expectedError:     false,
			expectedErrorCode: "",
		},
		{
			name:              "Error getting admin IDs",
			adminIds:          nil,
			getAdminsError:    errors.New("error getting admin IDs"),
			expectedAdminIds:  nil,
			expectedError:     true,
			expectedErrorCode: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockUserRepo.EXPECT().GetAdminIds(gomock.Any()).Return(tc.adminIds, tc.getAdminsError)
			admins, err := userService.GetAdmins(context.Background())
			assert.Equal(t, tc.expectedAdminIds, admins)
			if tc.expectedError {
				assert.Error(t, err)
				if tc.expectedErrorCode != "" {
					assert.Equal(t, tc.expectedErrorCode, srvErr.InternalServerError(err).GetCode())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}*/
