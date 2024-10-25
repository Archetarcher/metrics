package compression

import (
	"bytes"
	"compress/gzip"
	"encoding/json"

	"github.com/go-resty/resty/v2"
)

func GzipMiddleware(c *resty.Client, req *resty.Request) error {
	if req.Header.Get("Content-Encoding") != "gzip" {

		buf := bytes.NewBuffer(nil)
		zb := gzip.NewWriter(buf)
		js, err := json.Marshal(req.Body)

		if err != nil {
			return err
		}

		_, err = zb.Write(js)

		if err != nil {
			return err
		}
		err = zb.Close()
		if err != nil {
			return err
		}
		req.Header.Set(
			"Content-Encoding", "gzip")
		req.SetBody(buf)
	}

	return nil
}
