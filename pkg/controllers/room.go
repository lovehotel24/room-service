package controllers

import (
	"fmt"
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
		a.Log.WithError(err).Errorf("failed to get rooms: %s", err)
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to get rooms. details: %s", err))
	}

	return ctx.JSON(http.StatusOK, rooms)
}

func (a API) CreateRoom(ctx echo.Context, params routers.CreateRoomParams) error {
	var room routers.Room

	if err := ctx.Bind(&room); err != nil {
		a.Log.WithError(err).Errorf("failed to bind room: %s", err)
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("invalid input data. details: %s", err))
	}

	if _, err := a.getRoomTypeById(room.RoomTypeId); err != nil {
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("room type with id %s does not exist. details: %s", room.RoomTypeId, err))
	}

	if err := a.DB.Create(&room).Error; err != nil {
		a.Log.WithError(err).Errorf("failed to create room: %s", err)
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to create room. details: %s", err))
	}

	a.Log.Infof("room created with id: %s", room.Id)
	return ctx.JSON(http.StatusCreated, room.Id)
}

func (a API) DeleteRoomById(ctx echo.Context, roomId string) error {
	var room routers.Room
	if err := a.DB.Where("id = ?", roomId).Delete(&room).Error; err != nil {
		a.Log.WithError(err).Errorf("failed to delete room: %s", err)
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to delete room. details: %s", err))
	}

	a.Log.Infof("room deleted with id: %s", roomId)
	return ctx.NoContent(http.StatusNoContent)
}

func (a API) GetRoomById(ctx echo.Context, roomId string) error {
	room, err := a.getRoomById(roomId)
	if err != nil {
		return sendError(ctx, http.StatusNotFound, fmt.Sprintf("room id: %s doesn't exist. details: %s", roomId, err))
	}

	return ctx.JSON(http.StatusOK, room)
}

func (a API) UpdateRoomById(ctx echo.Context, roomId string) error {
	var requestRoom routers.Room
	if err := ctx.Bind(&requestRoom); err != nil {
		a.Log.WithError(err).Errorf("failed to bind room: %s", err)
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("invalid input data. details: %s", err))
	}

	room, err := a.getRoomById(roomId)
	if err != nil {
		return sendError(ctx, http.StatusNotFound, fmt.Sprintf("room id: %s doesn't exist. details: %s", roomId, err))
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
		a.Log.WithError(err).Errorf("failed to update room: %s", err)
		return sendError(ctx, http.StatusBadRequest, fmt.Sprintf("failed to update room. details: %s", err))
	}

	a.Log.Infof("room updated with id: %s", room.Id)
	return ctx.JSON(http.StatusOK, room.Id)
}
