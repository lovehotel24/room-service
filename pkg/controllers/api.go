package controllers

import (
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type API struct {
	DB  *gorm.DB
	Log *logrus.Logger
}

func NewApp(db *gorm.DB, log *logrus.Logger) *API {
	return &API{DB: db, Log: log}
}
