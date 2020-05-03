package util

import "testing"

func TestRequestString(t *testing.T) {
	hg := NewHttpGet()
	urlstr := "https://www.mozilla.org/en-US/firefox/75.0/whatsnew/all/?oldversion=73.0.1"
	s, err := hg.RequestString(urlstr)
	if err != nil {
		t.Errorf("error: %v\n", err)
	}
	t.Logf("Request String Length: %d\n", len(s))
}

func TestDownloadFile(t *testing.T) {
	hg := NewHttpGet()
	urlstr := "https://www.mozilla.org/en-US/firefox/75.0/whatsnew/all/?oldversion=73.0.1"
	size, err := hg.DownloadFile(urlstr, "/tmp/a.html")
	if err != nil {
		t.Errorf("error: %v\n", err)
	}
	t.Logf("Download File Size: %d\n", size)
}
