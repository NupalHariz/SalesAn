package supabase

import (
	"mime/multipart"

	"github.com/NupalHariz/SalesAn/src/utils/config"
	spbs "github.com/adityarizkyramadhan/supabase-storage-uploader"
)

type Interface interface {
	Upload(file *multipart.FileHeader) (string, error)
}

type supabase struct {
	client *spbs.Client
	cfg    config.SupabaseConfig
}

type InitParam struct {
	Cfg config.SupabaseConfig
}

func Init(param InitParam) Interface {
	client := spbs.New(param.Cfg.SupabaseUrl, param.Cfg.Token, param.Cfg.BucketName)
	return &supabase{
		client: client,
		cfg:    param.Cfg,
	}
}

func (s *supabase) Upload(file *multipart.FileHeader) (string, error) {
	url, err := s.client.Upload(file)

	return url, err
}
