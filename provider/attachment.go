package provider

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"irisblog/config"
	"irisblog/library"
	"irisblog/model"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"

	"github.com/nfnt/resize"
)

func AttachmentUpload(file multipart.File, info *multipart.FileHeader) (*model.Attachment, error) {
	db := config.DB
	bufFile := bufio.NewReader(file)
	img, imgType, err := image.Decode(bufFile)
	if err != nil {
		fmt.Println("无法获取图片尺寸")
		return nil, err
	}
	imgType = strings.ToLower(imgType)
	width := uint(img.Bounds().Dx())
	height := uint(img.Bounds().Dy())
	fmt.Println("width = ", width, " height = ", height)
	if imgType != "jpg" && imgType != "jpeg" && imgType != "gif" && imgType != "png" {
		return nil, errors.New("不支持的图片格式:" + imgType)
	}
	if imgType == "jpeg" {
		imgType = "jpg"
	}

	fileName := strings.TrimSuffix(info.Filename, path.Ext(info.Filename))
	log.Printf(fileName)

	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	md5hash := md5.New()
	bufFile = bufio.NewReader(file)
	_, err = io.Copy(md5hash, bufFile)
	if err != nil {
		return nil, err
	}
	md5Str := hex.EncodeToString(md5hash.Sum(nil))
	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, err
	}

	attachment, err := GetAttachmentByMd5(md5Str)
	if err == nil {
		if attachment.Status != 1 {
			attachment.Status = 1
			err = attachment.Save(db)
			if err != nil {
				return nil, err
			}
		}
		return attachment, nil
	}

	buff := &bytes.Buffer{}
	if width > 750 && imgType != "gif" {
		newImg := library.Resize(750, 0, img, resize.Lanczos3)
		width = uint(newImg.Bounds().Dx())
		height = uint(newImg.Bounds().Dy())
		if imgType == "jpg" {
			jpeg.Encode(buff, newImg, nil)
		} else if imgType == "png" {
			png.Encode(buff, newImg)
		}
	} else {
		io.Copy(buff, file)
	}
	tmpName := md5Str[8:24] + "." + imgType
	filePath := time.Now().Format("200601") + "/" + time.Now().Format("02") + "/"

	basePath := config.ExecPath + "/static/uploads//"

	if _, err = os.Stat(basePath + filePath); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(basePath+filePath, os.ModePerm)
			if err != nil {
				return nil, err
			}
		}
	}

	originFile, err := os.OpenFile(basePath+filePath+tmpName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	defer originFile.Close()
	_, err = io.Copy(originFile, buff)
	if err != nil {
		return nil, err
	}

	thumbName := "thumb_" + tmpName

	newThumbImg := library.ThumbnailCrop(250, 250, img)
	if imgType == "jpg" {
		jpeg.Encode(buff, newThumbImg, nil)
	} else if imgType == "png" {
		png.Encode(buff, newThumbImg)
	} else if imgType == "gif" {
		gif.Encode(buff, newThumbImg, nil)
	}

	thumbFile, err := os.OpenFile(basePath+filePath+thumbName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	defer thumbFile.Close()

	_, err = io.Copy(thumbFile, buff)
	if err != nil {
		return nil, err
	}

	attachment = &model.Attachment{
		FileName:     fileName,
		FileLocation: filePath + tmpName,
		FileSize:     info.Size,
		FileMd5:      md5Str,
		Width:        width,
		Height:       height,
		Status:       1,
	}
	attachment.GetThumb()
	err = attachment.Save(db)
	if err != nil {
		return nil, err
	}
	return attachment, nil

}

func GetAttachmentByMd5(md5 string) (*model.Attachment, error) {
	db := config.DB
	var attach model.Attachment
	if err := db.Where("`status` != 99").Where("`file_md5`=?", md5).First(&attach).Error; err != nil {
		return nil, err
	}
	attach.GetThumb()
	return &attach, nil
}
