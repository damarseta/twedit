package twedit

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
)

func GetBody(filepath string) (body io.ReadWriter, header string, err error) {
	var (
		mp     *multipart.Writer
		media  []byte
		writer io.Writer
	)
	body = bytes.NewBufferString("")
	mp = multipart.NewWriter(body)
	media, err = ioutil.ReadFile(filepath)
	if err != nil {
		return
	}
	mp.WriteField("status", "")
	writer, err = mp.CreateFormFile("media[]", "media.png")
	if err != nil {
		return
	}
	writer.Write(media)
	header = fmt.Sprintf("multipart/form-data;boundary=%v", mp.Boundary())
	mp.Close()
	return
}
