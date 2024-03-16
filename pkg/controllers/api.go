package controllers

import (
	"context"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/lovehotel24/booking-service/pkg/grpc/userpb"
)

const (
	userKey = "userId"
)

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	Password string `json:"password,omitempty"`
}

type API struct {
	DB   *gorm.DB
	Log  *logrus.Logger
	S3   *s3.Client
	Grpc userpb.UserServiceClient
}

func NewApp(db *gorm.DB, log *logrus.Logger, s3 *s3.Client, grpc userpb.UserServiceClient) *API {
	return &API{
		DB:   db,
		Log:  log,
		S3:   s3,
		Grpc: grpc,
	}
}

func sendError(ctx echo.Context, code int, message string) error {
	petErr := Error{
		Code:    int32(code),
		Message: message,
	}
	err := ctx.JSON(code, petErr)
	return err
}

func (a API) getRoomTypeById(rtId interface{}) (RoomType, error) {
	var rt RoomType

	if err := a.DB.Where("id = ?", rtId).First(&rt).Error; err != nil {
		a.Log.WithError(err).Errorf("failed to get room type by id: %s", rtId)
		return RoomType{}, err
	}

	return rt, nil
}

func (a API) getRoomById(rId interface{}) (Room, error) {
	var rt Room

	if err := a.DB.Where("id = ?", rId).First(&rt).Error; err != nil {
		a.Log.WithError(err).Errorf("failed to get room by id: %s", rId)
		return Room{}, err
	}

	return rt, nil
}

func getBookUserById(userId string, grpcClient userpb.UserServiceClient) (User, bool) {
	getUser, err := grpcClient.GetUser(context.Background(), &userpb.GetUserRequest{Id: &userpb.UUID{Value: userId}})
	if err != nil {
		return User{}, false
	}
	gUser := getUser.GetUser()
	user := User{
		Id:    gUser.GetId().GetValue(),
		Name:  gUser.GetName(),
		Phone: gUser.GetPhone(),
		Role:  gUser.GetRole(),
	}
	return user, true
}

func (a API) OnlyAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userId, ok := c.Get(userKey).(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Unauthorized"})
			}

			user, done := getBookUserById(userId, a.Grpc)
			if !done {
				return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": "Failed to fetch user info"})
			}

			if user.Role != "ADMIN" {
				return c.JSON(http.StatusUnauthorized, map[string]interface{}{"error": "Unauthorized"})
			}

			return next(c)
		}
	}
}
