package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type RoomType struct {
	gorm.Model
	Id          uuid.UUID      `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Name        string         `json:"name" form:"name"`
	Description string         `json:"description" form:"description"`
	ImageUrls   pq.StringArray `json:"images" form:"images" gorm:"type:text[]"`
	Price       float64        `json:"price" form:"price"`
}

func (rt *RoomType) BeforeCreate(tx *gorm.DB) (err error) {
	rt.Id, err = uuid.NewUUID()
	if err != nil {
		return err
	}
	return nil
}

type Room struct {
	gorm.Model
	Id         uuid.UUID `json:"id" gorm:"primary_key;type:uuid;default:uuid_generate_v4()"`
	Number     string    `json:"number"`
	RoomTypeId uuid.UUID `json:"roomTypeId" gorm:"type:uuid;"`
	Status     string    `json:"status"`
}

func (r *Room) BeforeCreate(tx *gorm.DB) (err error) {
	r.Id, err = uuid.NewUUID()
	if err != nil {
		return err
	}
	return nil
}

//
//func (Photo) GormDataType() string {
//	return "text"
//}
//func (Photo) GormDBDataType(db *gorm.DB, field *schema.Field) string {
//
//	// returns different database type based on driver name
//	switch db.Dialector.Name() {
//	case "mysql", "sqlite":
//		return "text"
//	}
//	return ""
//}

//func (o *Photo) Scan(src any) error {
//	bytes, ok := src.([]byte)
//	if !ok {
//		return errors.New("src value cannot cast to []byte")
//	}
//	*o = strings.Split(string(bytes), ",")
//	return nil
//}
//func (o Photo) Value() (driver.Value, error) {
//	if len(o) == 0 {
//		return nil, nil
//	}
//	return strings.Join(o, ","), nil
//}
