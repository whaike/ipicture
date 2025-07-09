package service

import (
	"ipicture-backend/models"
	"ipicture-backend/config"
	"time"
	"encoding/json"
	"strings"
)

func SaveFileMeta(file *models.File) error {
	return config.DB.Create(file).Error
}

func QueryFiles(userID uint, page, pageSize int, fileType, date string, tag ...string) ([]models.File, int64, error) {
	db := config.DB.Model(&models.File{}).Where("user_id = ?", userID)
	if fileType == "image" || fileType == "video" {
		db = db.Where("type = ?", fileType)
	}
	if date != "" {
		if t, err := time.Parse("2006-01-02", date); err == nil {
			start := t
			end := t.Add(24 * time.Hour)
			db = db.Where("created_at >= ? AND created_at < ?", start, end)
		}
	}
	if len(tag) > 0 && tag[0] != "" {
		db = db.Where("tags LIKE ?", "%"+tag[0]+"%")
	}
	var files []models.File
	var total int64
	db.Count(&total)
	db = db.Order("created_at desc").Offset((page-1)*pageSize).Limit(pageSize)
	err := db.Find(&files).Error
	// 反序列化tags
	for i := range files {
		var tags []string
		json.Unmarshal([]byte(files[i].Tags), &tags)
		files[i].Tags = strings.Join(tags, ",")
	}
	return files, total, err
}

func GetFileByID(id, userID uint) (*models.File, error) {
	var file models.File
	err := config.DB.Where("id = ? AND user_id = ?", id, userID).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func AITagFile(id, userID uint) ([]string, error) {
	_, err := GetFileByID(id, userID)
	if err != nil {
		return nil, err
	}
	// 模拟AI识别，实际可调用第三方API或本地模型
	tags := []string{"风景", "人物"}
	tagsJson, _ := json.Marshal(tags)
	err = config.DB.Model(&models.File{}).Where("id = ? AND user_id = ?", id, userID).Updates(map[string]interface{}{"tags": string(tagsJson)}).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}

func GetFilesByMD5List(userID uint, md5List []string) ([]models.File, error) {
	var files []models.File
	err := config.DB.Where("user_id = ? AND md5 IN ?", userID, md5List).Find(&files).Error
	return files, err
}

func GetFileByMD5(userID uint, md5 string) (*models.File, error) {
	var file models.File
	err := config.DB.Where("user_id = ? AND md5 = ?", userID, md5).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func DeleteFileByID(id, userID uint) error {
	return config.DB.Where("id = ? AND user_id = ?", id, userID).Delete(&models.File{}).Error
} 