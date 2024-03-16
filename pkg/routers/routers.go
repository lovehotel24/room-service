package routers

import (
	"github.com/labstack/echo/v4"

	"github.com/lovehotel24/room-service/pkg/controllers"
)

func RouterV1(e *echo.Echo, api *controllers.API) {
	wrapper := ServerInterfaceWrapper{Handler: api}
	e.GET("/v1/room", wrapper.GetAllRoom)
	e.POST("/v1/room", wrapper.CreateRoom)
	e.DELETE("/v1/room/:RoomId", wrapper.DeleteRoomById)
	e.GET("/v1/room/:RoomId", wrapper.GetRoomById)
	e.PUT("/v1/room/:RoomId", wrapper.UpdateRoomById)
	e.GET("/v1/roomtype", wrapper.GetAllRoomType)
	e.POST("/v1/roomtype", wrapper.CreateRoomType, api.OnlyAdmin())
	e.POST("/v1/roomtype/upload", wrapper.UploadRoomTypePhotos)
	e.DELETE("/v1/roomtype/:RoomTypeId", wrapper.DeleteRoomTypeById, api.OnlyAdmin())
	e.GET("/v1/roomtype/:RoomTypeId", wrapper.GetRoomTypeById)
	e.PUT("/v1/roomtype/:RoomTypeId", wrapper.UpdateRoomTypeById, api.OnlyAdmin())
}
