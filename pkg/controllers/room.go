package controllers

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/lovehotel24/room-service/pkg/routers"
)

func (a API) GetAllRoom(ctx echo.Context, params routers.GetAllRoomParams) error {
	var rooms []routers.Room

	limit := 10
	offSet := 1

	if params.Limit != nil {
		limit = *params.Limit
	}
	if params.Offset != nil {
		offSet = *params.Offset
	}

	if err := a.DB.Limit(limit).Offset(offSet).Find(&rooms).Error; err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, rooms)
}

func (a API) CreateRoom(ctx echo.Context, params routers.CreateRoomParams) error {
	var room routers.Room

	if err := ctx.Bind(&room); err != nil {
		return err
	}
	if err := a.DB.Create(&room).Error; err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, room.Id)
}

func (a API) DeleteRoomById(ctx echo.Context, roomId string) error {
	var room routers.Room
	if err := a.DB.Where("id = ?", roomId).Delete(&room).Error; err != nil {
		return err
	}
	return ctx.NoContent(http.StatusNoContent)
}

func (a API) GetRoomById(ctx echo.Context, roomId string) error {
	room, err := a.getRoomById(roomId)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, room)
}

func (a API) UpdateRoomById(ctx echo.Context, roomId string) error {
	var requestRoom routers.Room
	if err := ctx.Bind(&requestRoom); err != nil {
		return err
	}

	room, err := a.getRoomById(roomId)
	if err != nil {
		a.Log.WithError(err).Errorf("failed to get room by id: %s", roomId)
		return err
	}

	if requestRoom.Number != "" {
		room.Number = requestRoom.Number
	}

	if requestRoom.RoomTypeId != uuid.Nil {
		room.RoomTypeId = requestRoom.RoomTypeId
	}

	if requestRoom.Status != "" {
		room.Status = requestRoom.Status
	}

	if err := a.DB.Save(&room).Error; err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, room.Id)
}

func (a API) getRoomById(rId interface{}) (routers.Room, error) {
	var rt routers.Room

	if err := a.DB.Where("id = ?", rId).First(&rt).Error; err != nil {
		a.Log.WithError(err).Errorf("failed to get room by id: %s", rId)
		return routers.Room{}, err
	}

	return rt, nil
}
