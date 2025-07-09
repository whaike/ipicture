package api

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"ipicture-backend/models"
	"ipicture-backend/service"
	"ipicture-backend/utils"
	"github.com/gin-gonic/gin"
	"image"
	"image/jpeg"
	"github.com/nfnt/resize"
	"log"
	_ "image/png"
	"crypto/md5"
	"io"
	"encoding/hex"
)

func RegisterFileRoutes(r *gin.Engine) {
	api := r.Group("/api")
	api.POST("/upload", utils.JWTAuthMiddleware(), UploadFile)
	api.GET("/files", utils.JWTAuthMiddleware(), ListFiles)
	api.GET("/file/download/:id", utils.JWTAuthMiddleware(), DownloadFile)
	api.POST("/file/ai_tag/:id", utils.JWTAuthMiddleware(), AITagFile)
	api.POST("/files/check_md5", utils.JWTAuthMiddleware(), CheckFilesMD5)
	api.DELETE("/file/:id", utils.JWTAuthMiddleware(), DeleteFile)
}

// 批量MD5查重接口
func CheckFilesMD5(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	var req struct {
		MD5List []string `json:"md5_list"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || len(req.MD5List) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	files, err := service.GetFilesByMD5List(userID, req.MD5List)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查重失败"})
		return
	}
	existMD5 := make([]string, 0, len(files))
	for _, f := range files {
		existMD5 = append(existMD5, f.MD5)
	}
	c.JSON(http.StatusOK, gin.H{"exist_md5": existMD5})
}

func UploadFile(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未选择文件"})
		return
	}
	log.Printf("收到上传文件: %s, Content-Type: %s", file.Filename, file.Header.Get("Content-Type"))
	md5 := c.PostForm("md5")
	if md5 == "" {
		// App端未传md5，自动计算
		f, ferr := file.Open()
		if ferr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "文件读取失败"})
			return
		}
		defer f.Close()
		h := md5sum(f)
		md5 = h
		log.Printf("后端自动计算MD5: %s", md5)
		// 查重
		if exist, _ := service.GetFileByMD5(userID, md5); exist != nil {
			c.JSON(http.StatusOK, gin.H{"message": "文件已存在", "file": exist})
			return
		}
	} else if md5 != "" {
		if exist, _ := service.GetFileByMD5(userID, md5); exist != nil {
			c.JSON(http.StatusOK, gin.H{"message": "文件已存在", "file": exist})
			return
		}
	}
	// 保存文件到本地
	dstDir := "uploads/" + strconv.Itoa(int(userID))
	os.MkdirAll(dstDir, os.ModePerm)
	filename := filepath.Base(file.Filename)
	dstPath := filepath.Join(dstDir, filename)
	if err := c.SaveUploadedFile(file, dstPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存文件失败"})
		return
	}
	log.Printf("文件已保存: %s", dstPath)
	// 生成缩略图
	thumbPath := ""
	if strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		thumbDir := filepath.Join("thumbnails", strconv.Itoa(int(userID)))
		os.MkdirAll(thumbDir, os.ModePerm)
		thumbPath = filepath.Join(thumbDir, "thumb_"+filename+".jpg")
		log.Printf("准备生成缩略图: %s", thumbPath)
		err := generateImageThumbnail(dstPath, thumbPath, 200, 200)
		if err != nil {
			thumbPath = ""
			log.Printf("生成缩略图失败: %v, src=%s, dst=%s", err, dstPath, thumbPath)
		} else {
			log.Printf("缩略图生成成功: %s", thumbPath)
		}
	} else {
		log.Printf("Content-Type不是image，未生成缩略图: %s", file.Header.Get("Content-Type"))
	}
	// 记录元数据
	fileType := "other"
	if strings.HasPrefix(file.Header.Get("Content-Type"), "image/") {
		fileType = "image"
	} else if strings.HasPrefix(file.Header.Get("Content-Type"), "video/") {
		fileType = "video"
	}
	log.Printf("最终写入thumb_path: %s", thumbPath)
	meta := models.File{
		UserID:   userID,
		Filename: filename,
		Path:     dstPath,
		Type:     fileType,
		Size:     file.Size,
		MD5:      md5,
		ThumbPath: thumbPath,
	}
	if err := service.SaveFileMeta(&meta); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "保存元数据失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "上传成功", "file": meta})
}

// 生成图片缩略图
func generateImageThumbnail(srcPath, dstPath string, width, height uint) error {
	f, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return err
	}
	thumb := resize.Thumbnail(width, height, img, resize.Lanczos3)
	out, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer out.Close()
	return jpeg.Encode(out, thumb, &jpeg.Options{Quality: 80})
}

// 计算文件MD5
func md5sum(r io.Reader) string {
	h := md5.New()
	io.Copy(h, r)
	return hex.EncodeToString(h.Sum(nil))
}

func ListFiles(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	fileType := c.Query("type")
	date := c.Query("date")
	tag := c.Query("tag")
	files, total, err := service.QueryFiles(userID, page, pageSize, fileType, date, tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "查询失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"files": files, "total": total, "page": page, "page_size": pageSize})
}

func DownloadFile(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	id, _ := strconv.Atoi(c.Param("id"))
	file, err := service.GetFileByID(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}

	// 自动判断文件类型
	f, openErr := os.Open(file.Path)
	if openErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "文件读取失败"})
		return
	}
	defer f.Close()
	buf := make([]byte, 512)
	_, _ = f.Read(buf)
	contentType := http.DetectContentType(buf)
	// 保证 .graffle 等特殊后缀也能原样下载
	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", "attachment; filename=\""+file.Filename+"\"")
	c.File(file.Path)
}

func AITagFile(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	id, _ := strconv.Atoi(c.Param("id"))
	tags, err := service.AITagFile(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI标签识别失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tags": tags})
}

// 删除文件接口
func DeleteFile(c *gin.Context) {
	userID := utils.GetUserIDFromContext(c)
	id, _ := strconv.Atoi(c.Param("id"))
	file, err := service.GetFileByID(uint(id), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文件不存在"})
		return
	}
	// 物理删除原文件
	if file.Path != "" {
		os.Remove(file.Path)
	}
	// 物理删除缩略图
	if file.ThumbPath != "" {
		os.Remove(file.ThumbPath)
	}
	// 软删除数据库记录
	if err := service.DeleteFileByID(uint(id), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
} 