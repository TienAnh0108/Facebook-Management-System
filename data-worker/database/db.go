package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Page struct {
	ID          string `gorm:"primaryKey"`
	Name        string
	Category    string
	AccessToken string
	CreatedAt   time.Time
}

type DailyStat struct {
	ID             uint      `gorm:"primaryKey"`
	PageID         string    `gorm:"index"`
	Date           time.Time `gorm:"type:date;index"`
	FollowersCount int
	ReactsCount    int
	CommentsCount  int
	SharesCount    int

	Page Page `gorm:"foreignKey:PageID"`
}

func ConnectDB(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Không thể kết nối đến database:", err)
	}

	fmt.Println("Đã kết nối PostgreSQL thành công!")

	// Tính năng AutoMigrate: GORM sẽ tự động tạo bảng Page và DailyStat nếu chưa có
	err = db.AutoMigrate(&Page{}, &DailyStat{})
	if err != nil {
		log.Fatal("Lỗi khi tạo bảng:", err)
	}

	fmt.Println("Đã đồng bộ cấu trúc bảng xong!")
	return db
}
