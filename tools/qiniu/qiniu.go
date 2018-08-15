package qiniu

import (
	"bytes"
	"context"
	"fmt"
	"github.com/qiniu/api.v7/auth/qbox"
	"github.com/qiniu/api.v7/storage"
	"github.com/tokenme/tmm/common"
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
