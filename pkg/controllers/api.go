package controllers

import (
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/lovehotel24/room-service/pkg/routers"
)

type API struct {
	DB  *gorm.DB
	Log *logrus.Logger
}

func (a API) PostV1RoomtypeUpload(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}

func NewApp(db *gorm.DB, log *logrus.Logger) *API {
	return &API{DB: db, Log: log}
}

func sendError(ctx echo.Context, code int, message string) error {
	petErr := routers.Error{
		Code:    int32(code),
		Message: message,
	}
	err := ctx.JSON(code, petErr)
	return err
}

func (a API) getRoomTypeById(rtId interface{}) (routers.RoomType, error) {
	var rt routers.RoomType

	if err := a.DB.Where("id = ?", rtId).First(&rt).Error; err != nil {
		a.Log.WithError(err).Errorf("failed to get room type by id: %s", rtId)
		return routers.RoomType{}, err
	}

	return rt, nil
}

func (a API) getRoomById(rId interface{}) (routers.Room, error) {
	var rt routers.Room

	if err := a.DB.Where("id = ?", rId).First(&rt).Error; err != nil {
		a.Log.WithError(err).Errorf("failed to get room by id: %s", rId)
		return routers.Room{}, err
	}

	return rt, nil
}

//func OnlyAdmin(grpcClient userpb.UserServiceClient) echo.MiddlewareFunc {
//	return func(next echo.HandlerFunc) echo.HandlerFunc {
//		return func(c echo.Context) error {
//			userId, ok := c.Get(userKey).(string)
//			if !ok {
//				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Unauthorized"})
//			}
//
//			user, done := getBookUserById(userId, grpcClient)
//			if !done {
//				return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to fetch user info"})
//			}
//
//			if user.Role != "ADMIN" {
//				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Unauthorized"})
//			}
//
//			return next(c)
//		}
//	}
//}
