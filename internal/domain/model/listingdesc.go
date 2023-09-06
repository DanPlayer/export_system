package model

import "export_system/internal/db"

type ListingDesc struct {
	ID           int64
	Title        string
	Keywords     string
	Description  string
	BulletPoint1 string
	BulletPoint2 string
	BulletPoint3 string
	BulletPoint4 string
	BulletPoint5 string
}

func (m *ListingDesc) ListRangeByID(id, size int) (list []ListingDesc, err error) {
	table := db.MasterClient.Model(&m)
	err = table.Where("id > ?", id).Limit(size).Find(&list).Error
	return
}
