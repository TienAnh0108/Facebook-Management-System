package tasks

import (
	"fb-management-system/data-worker/database"
	"fb-management-system/data-worker/facebook"
	"log"
	"time"

	"gorm.io/gorm"
)

// SyncPageData quét qua toàn bộ Page trong DB, gọi Facebook API lấy số liệu và lưu lại
func SyncPageData(db *gorm.DB) {
	// 1. Lấy danh sách tất cả các Page mà hệ thống đang quản lý
	var pages []database.Page
	if err := db.Find(&pages).Error; err != nil {
		log.Println("Lỗi không lấy được danh sách Page từ Database:", err)
		return
	}

	fbClient := facebook.NewClient()

	// 2. Lặp qua từng Page để thu thập dữ liệu
	for _, p := range pages {
		log.Printf("-> Đang quét dữ liệu cho Fanpage: %s (ID: %s)", p.Name, p.ID)

		// Lấy thông tin cơ bản
		fbInfo, err := fbClient.FetchPageInfo(p.ID, p.AccessToken)
		if err != nil {
			log.Printf("Lỗi lấy thông tin cơ bản của Page %s: %v", p.ID, err)
			continue
		}

		// Lấy chỉ số tương tác trong ngày
		reacts, comments, shares, err := fbClient.FetchDailyInteractions(p.ID, p.AccessToken)
		if err != nil {
			log.Printf("Lỗi lấy tương tác của Page %s: %v", p.ID, err)
			continue
		}

		// Định dạng ngày hôm nay (chỉ giữ lại Năm-Tháng-Ngày, đặt giờ về 00:00:00)
		today := time.Now()
		dateOnly := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

		// Chuẩn bị struct lưu vào bảng DailyStat
		dailyStat := database.DailyStat{
			PageID:         p.ID,
			Date:           dateOnly,
			FollowersCount: fbInfo.FollowersCount,
			ReactsCount:    reacts,
			CommentsCount:  comments,
			SharesCount:    shares,
		}

		// 3. Kiểm tra xem hôm nay Page này đã có dòng dữ liệu nào chưa để tránh ghi trùng (Idempotency)
		var existing database.DailyStat
		err = db.Where("page_id = ? AND date = ?", p.ID, dateOnly.Format("2006-01-02")).First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			// Nếu chưa có dữ liệu ngày hôm nay -> Tạo mới
			if err := db.Create(&dailyStat).Error; err != nil {
				log.Printf("Lỗi khi lưu dữ liệu mới cho Page %s: %v", p.ID, err)
			} else {
				log.Printf("Thành công: Đã tạo bản ghi thống kê mới cho ngày hôm nay!")
			}
		} else if err == nil {
			// Nếu hôm nay đã chạy rồi -> Cập nhật đè số liệu mới nhất
			existing.FollowersCount = dailyStat.FollowersCount
			existing.ReactsCount = dailyStat.ReactsCount
			existing.CommentsCount = dailyStat.CommentsCount
			existing.SharesCount = dailyStat.SharesCount
			db.Save(&existing)
			log.Printf("Thành công: Đã cập nhật số liệu mới nhất cho ngày hôm nay!")
		}
	}
}
