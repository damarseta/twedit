package twedit

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func Download(cleanURL string, filename string) (string, error) {
	res, err := http.Get(cleanURL)

	if err != nil {
		log.Printf("http.Get -> %v", err)
		return "", err
	}

	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	// You have to manually close the body, check docs
	// This is required if you want to use things like
	// Keep-Alive and other HTTP sorcery.
	res.Body.Close()

	folder := "medias"
	fullpath := folder + string(filepath.Separator) + filename

	os.MkdirAll(folder, 0777)

	// You can now save it to disk or whatever...
	if err = ioutil.WriteFile(fullpath, data, 0777); err != nil {
		log.Println("Error Saving:", filename, err)
	} else {
		log.Println("Saved:", filename)
	}

	return fullpath, nil
}
