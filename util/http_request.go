package util

import (
	"Pay-Server/log"
	"io/ioutil"
	"net/http"
	"time"
)

func GetRequest(url string, urlParam map[string]string) (bool, []byte) {
	defer func() {
		err := recover()
		if err != nil {
			log.PrintStackTrace(err)
		}
	}()
	if url == "" {
		return false, nil
	}
	request, _ := http.NewRequest("GET", url, nil)
	//加入get参数
	q := request.URL.Query()
	if urlParam != nil {
		for k, v := range urlParam {
			q.Add(k, v)
		}
	}
	request.URL.RawQuery = q.Encode()

	timeout := time.Duration(600 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Do(request)
	if err != nil {
		return false, nil
	}

	data, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return false, nil
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	return true, data
}
