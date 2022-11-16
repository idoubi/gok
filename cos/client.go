package cos

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/idoubi/goutils"
	"github.com/spf13/viper"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type Config struct {
	Domain    string `json:"domain"`
	CdnDomain string `json:"cdnDomain"`
	SecretId  string `json:"secretId"`
	SecretKey string `json:"secretKey"`
	Dir       string `json:"dir"`
}

type Client struct {
	*cos.Client
	ctx  context.Context
	conf Config
}

type Response struct {
	*cos.Response
	file *CosFile
}

type CosFile struct {
	Url   string `json:"url"`
	Name  string `json:"name"`
	Cache bool   `json:"cache,omitempty"`
}

func New(name string) (*Client, error) {
	var conf Config
	sub := viper.Sub("cos." + name)
	if sub == nil {
		return nil, fmt.Errorf("invalid cos config under %s", name)
	}
	if err := sub.Unmarshal(&conf); err != nil {
		return nil, err
	}

	u, err := url.Parse(conf.Domain)
	if err != nil {
		return nil, err
	}

	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  conf.SecretId,
			SecretKey: conf.SecretKey,
		},
	})

	return &Client{c, context.Background(), conf}, nil
}

func (c *Client) UploadImage(filename string, r io.Reader) (*Response, error) {
	if filename == "" {
		filename = goutils.TimeSeq() + ".jpeg"
	}

	if c.conf.Dir != "" {
		filename = c.conf.Dir + "/" + filename
	}

	fileDomain := c.conf.Domain
	if c.conf.CdnDomain != "" {
		fileDomain = c.conf.CdnDomain
	}

	fileUrl := fmt.Sprintf("%s/%s", fileDomain, filename)

	cosFile := &CosFile{
		Name: filename,
		Url:  fileUrl,
	}

	if c.IsFileExists(filename) {
		cosFile.Cache = true

		return &Response{nil, cosFile}, nil
	}

	resp, err := c.Object.Put(c.ctx, filename, r, nil)
	if err != nil {
		return nil, err
	}

	return &Response{resp, cosFile}, nil
}

func (c *Client) IsFileExists(filename string) bool {
	if ok, err := c.Object.IsExist(c.ctx, filename); err == nil && ok {
		return true
	}

	return false
}

func (c *Client) DeleteFile(filename string) error {
	if c.conf.Dir != "" {
		filename = c.conf.Dir + "/" + filename
	}

	_, err := c.Object.Delete(c.ctx, filename)

	return err
}

func (r *Response) GetFile() *CosFile {
	return r.file
}
