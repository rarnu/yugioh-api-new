package html

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func DownloadFile(url string, localFile string, callback func(progress int64, total int64)) error {
	var fsize int64
	var buf = make([]byte, 32*1024)
	var written int64
	client := new(http.Client)
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	fsize, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		log.Printf("%s", err)
	}
	//创建文件
	file, err := os.Create(localFile)
	if err != nil {
		return err
	}
	defer func(f *os.File) { _ = f.Close() }(file)
	if resp.Body == nil {
		return fmt.Errorf("body is null")
	}
	defer func(rc io.ReadCloser) { _ = rc.Close() }(resp.Body)

	for {
		nr, er := resp.Body.Read(buf)
		if nr > 0 {
			nw, ew := file.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er != nil {
			if er != io.EOF {
				err = er
			}
			break
		}
		callback(written, fsize)
	}
	return err
}
