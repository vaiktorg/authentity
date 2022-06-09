package handlers

import (
	"gorm.io/gorm"
	"sync"
)

type DataProvider struct {
	Mutex sync.Mutex
	DB    *gorm.DB
}
