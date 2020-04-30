package util

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	USER_AGENT = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:65.0) Gecko/20100101 Firefox/65.0"
)

type HttpGet struct {
	userAgent string
}

func NewHttpGet() *HttpGet {
	return &HttpGet{userAgent: USER_AGENT}
}

func (hg *HttpGet) SetUserAgent(ua string) {
	hg.userAgent = ua
}

func (hg *HttpGet) Do(urlstr string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", urlstr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", hg.userAgent)
	return client.Do(req)
}

func (hg *HttpGet) RequestBytes(urlstr string) ([]byte, error) {
	resp, err := hg.Do(urlstr)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (hg *HttpGet) RequestString(urlstr string) (string, error) {
	buf, err := hg.RequestBytes(urlstr)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func (hg *HttpGet) DownloadFile(urlstr string, fpath string) (int64, error) {
	resp, err := hg.Do(urlstr)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	f, err := os.Create(fpath)
	if err != nil {
		return 0, err
	}
	size, err := io.Copy(f, resp.Body)
	if err != nil {
		return 0, err
	}
	return size, nil
}
