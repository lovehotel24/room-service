package controllers

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/lovehotel24/room-service/pkg/routers"
)

func (a API) GetAllRoomType(ctx echo.Context, params routers.GetAllRoomTypeParams) error {
	var roomType []routers.RoomType

	limit := 10
	offSet := 1

	if params.Limit != nil {
		limit = *params.Limit
	}
	if params.Offset != nil {
		offSet = *params.Offset
	}

	if err := a.DB.Limit(limit).Offset(offSet).Find(&roomType).Error; err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, roomType)
}

func (a API) CreateRoomType(ctx echo.Context, params routers.CreateRoomTypeParams) error {
	var roomType routers.RoomType

	if err := ctx.Bind(&roomType); err != nil {
		return err
	}
	if err := a.DB.Create(&roomType).Error; err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, roomType.Id)
}

func (a API) DeleteRoomTypeById(ctx echo.Context, roomTypeId string) error {
	var roomType routers.RoomType
	if err := a.DB.Where("id = ?", roomTypeId).Delete(&roomType).Error; err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (a API) GetRoomTypeById(ctx echo.Context, roomTypeId string) error {
	rt, err := a.getRoomTypeById(roomTypeId)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, rt)
}

func (a API) UpdateRoomTypeById(ctx echo.Context, roomTypeId string) error {
	var requestRT routers.RoomType
	if err := ctx.Bind(&requestRT); err != nil {
		return err
	}
	rt, err := a.getRoomTypeById(roomTypeId)
	if err != nil {
		a.Log.WithError(err).Errorf("failed to get booking by id: %s", roomTypeId)
		return err
	}
	if requestRT.Name != "" {
		rt.Name = requestRT.Name
	}
	if requestRT.Description != "" {
		rt.Description = requestRT.Description
	}
	if requestRT.Price != 0 {
		rt.Price = requestRT.Price
	}
	if requestRT.Photo != nil {
		rt.Photo = requestRT.Photo
	}

	if err := a.DB.Save(&rt).Error; err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, rt.Id)
}

func (a API) getRoomTypeById(rtId interface{}) (routers.RoomType, error) {
	var rt routers.RoomType

	if err := a.DB.Where("id = ?", rtId).First(&rt).Error; err != nil {
		a.Log.WithError(err).Errorf("failed to get room type by id: %s", rtId)
		return routers.RoomType{}, err
	}

	return rt, nil
}
