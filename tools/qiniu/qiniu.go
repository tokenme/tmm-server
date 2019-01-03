package qiniu

import (
	"bytes"
	"context"
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"github.com/tokenme/tmm/common"
	"strings"
)

func Upload(ctx context.Context, config common.QiniuConfig, path string, filename string, data []byte) (string, storage.PutRet, error) {
	key := fmt.Sprintf("%s/%s", path, filename)
	putPolicy := storage.PutPolicy{
		Scope:      fmt.Sprintf("%s:%s", config.Bucket, key),
		DetectMime: 1,
	}
	mac := qbox.NewMac(config.AK, config.Secret)
	upToken := putPolicy.UploadToken(mac)
	cfg := storage.Config{}
	cfg.Zone = &storage.ZoneHuadong
	cfg.UseHTTPS = true
	cfg.UseCdnDomains = true
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:path": path,
			"x:file": filename,
		},
	}
	dataLen := int64(len(data))
	err := formUploader.Put(ctx, &ret, upToken, key, bytes.NewReader(data), dataLen, &putExtra)
	if err != nil {
		return "", ret, err
	}
	publicURL := storage.MakePublicURL(config.Domain, ret.Key)
	return publicURL, ret, nil
}

func UpToken(config common.QiniuConfig, path string, filename string) (string, string, string) {
	key := fmt.Sprintf("%s/%s", path, filename)
	putPolicy := storage.PutPolicy{
		Scope:      fmt.Sprintf("%s:%s", config.Bucket, key),
		DetectMime: 1,
	}
	mac := qbox.NewMac(config.AK, config.Secret)
	upToken := putPolicy.UploadToken(mac)
	return upToken, key, storage.MakePublicURL(config.Domain, key)
}

func ConvertImage(url string, extension string, config common.QiniuConfig) (newUrl string, persistentId string, err error) {
	pathArr := strings.Split(url, "/")
	if len(pathArr) > 0 {
		filename := pathArr[len(pathArr)-1]
		newFilename := fmt.Sprintf("%s/c%s", config.ImagePath, filename)
		fopBaseConvert := fmt.Sprintf("imageView2/0/w/500/h/500/format/%s|saveas/%s", extension, storage.EncodedEntry(config.Bucket, newFilename))
		fopBatch := []string{fopBaseConvert}
		persistentId, err = Pfop(config, config.ImagePath, filename, fopBatch, false)
		if err != nil {
			return "", "", err
		}
		newUrl = fmt.Sprintf("%s/%s", config.Domain, newFilename)
	}
	return
}

func Pfop(config common.QiniuConfig, path string, filename string, fopBatch []string, force bool) (string, error) {
	key := fmt.Sprintf("%s/%s", path, filename)
	mac := qbox.NewMac(config.AK, config.Secret)
	cfg := storage.Config{
		Zone:          &storage.ZoneHuadong,
		UseCdnDomains: true,
		UseHTTPS:      false,
	}
	operationManager := storage.NewOperationManager(mac, &cfg)
	fops := strings.Join(fopBatch, ";")
	return operationManager.Pfop(config.Bucket, key, fops, config.Pipeline, config.NotifyURL, force)
}
