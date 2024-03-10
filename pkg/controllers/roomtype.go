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
		a.Log.WithError(err).Errorf("failed to get room type: %s", err)
		return sendError(ctx, http.StatusBadRequest, "failed to get room type")
	}

	return ctx.JSON(http.StatusOK, roomType)
}

func (a API) CreateRoomType(ctx echo.Context, params routers.CreateRoomTypeParams) error {
	var roomType routers.RoomType

	if err := ctx.Bind(&roomType); err != nil {
		a.Log.WithError(err).Errorf("failed to bind room type: %s", err)
		return sendError(ctx, http.StatusBadRequest, "invalid input data")
	}
	if err := a.DB.Create(&roomType).Error; err != nil {
		a.Log.WithError(err).Errorf("failed to create room type: %s", err)
		return sendError(ctx, http.StatusBadRequest, "failed to create room type")
	}

	a.Log.Infof("room type created with id: %s", roomType.Id)
	return ctx.JSON(http.StatusCreated, roomType.Id)
}

func (a API) DeleteRoomTypeById(ctx echo.Context, roomTypeId string) error {
	var roomType routers.RoomType
	if err := a.DB.Where("id = ?", roomTypeId).Delete(&roomType).Error; err != nil {
		a.Log.WithError(err).Errorf("failed to delete room type: %s", err)
		return sendError(ctx, http.StatusBadRequest, "failed to delete room type")
	}

	a.Log.Infof("room type deleted with id: %s", roomType.Id)
	return ctx.NoContent(http.StatusNoContent)
}

func (a API) GetRoomTypeById(ctx echo.Context, roomTypeId string) error {
	rt, err := a.getRoomTypeById(roomTypeId)
	if err != nil {
		return sendError(ctx, http.StatusNotFound, "room type not found")
	}

	return ctx.JSON(http.StatusOK, rt)
}

func (a API) UpdateRoomTypeById(ctx echo.Context, roomTypeId string) error {
	var requestRT routers.RoomType
	if err := ctx.Bind(&requestRT); err != nil {
		a.Log.WithError(err).Errorf("failed to bind room type: %s", err)
		return sendError(ctx, http.StatusBadRequest, "invalid input data")
	}
	rt, err := a.getRoomTypeById(roomTypeId)
	if err != nil {
		return sendError(ctx, http.StatusNotFound, "room type not found")
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
	if requestRT.ImageUrls != nil {
		rt.ImageUrls = requestRT.ImageUrls
	}

	if err := a.DB.Save(&rt).Error; err != nil {
		a.Log.WithError(err).Errorf("failed to update room type: %s", err)
		return sendError(ctx, http.StatusBadRequest, "failed to update room type")
	}

	a.Log.Infof("room type updated with id: %s", rt.Id)
	return ctx.JSON(http.StatusOK, rt.Id)
}
