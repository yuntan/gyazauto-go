package main

import (
	"github.com/mattn/go-jsonpointer"
	"github.com/pkg/browser"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"

	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

var (
	chCloseServer = make(chan int, 1)
)

func authorize() {
	conf := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://api.gyazo.com/oauth/authorize",
			TokenURL: "https://api.gyazo.com/oauth/token",
		},
		RedirectURL: "http://localhost:9090/authorize",
	}

	authURL := conf.AuthCodeURL("state")
	log.Info("Opening authorize URL (%s)...", authURL)
	err := browser.OpenURL(authURL)
	if err != nil {
		log.Fatal(err)
	}

	s := &http.Server{
		Addr:    "localhost:9090",
		Handler: getHandler(conf),
	}
	go func() {
		<-chCloseServer
		s.Shutdown(context.Background())
	}()
	err = s.ListenAndServe()
	if err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func getHandler(conf *oauth2.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("handler: URL: %s", r.URL)

		if r.URL.Path != "/authorize" {
			return
		}
		if r.URL.Query().Get("state") != "state" {
			log.Error("Invalid token, must be CSRF attack!")
			return
		}
		code := r.URL.Query().Get("code")
		token, err := conf.Exchange(context.Background(), code)
		if err != nil {
			log.Fatal(err)
			return
		}

		config.AccessToken = token.AccessToken
		saveConfig()
		fmt.Fprintf(w, "Authorization success! You can close this page.")
		log.Info("Authorization success!")
		chCloseServer <- 0
	}
}

func upload(path string) string {
	name := filepath.Base(path)

	file, err := os.Open(path)
	if err != nil {
		log.Error("Failed to read file %s: %v", path, err)
		return ""
	}
	defer file.Close()

	reqBody := &bytes.Buffer{}
	writer := multipart.NewWriter(reqBody)

	w, _ := writer.CreateFormFile("imagedata", name)
	if _, err := io.Copy(w, file); err != nil {
		log.Error("Failed to read file %s: %v", path, err)
		return ""
	}

	writer.WriteField("access_token", config.AccessToken)
	writer.WriteField("title", filepath.Base(path))
	writer.Close()

	req, _ := http.NewRequest("POST", "https://upload.gyazo.com/api/upload", reqBody)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	log.Info("Uploading...")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("Failed to upload: %v", err)
		return ""
	} else if resp.StatusCode != 200 {
		log.Error("Failed to upload, code: %d", resp.StatusCode)
		return ""
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read body: %v", err)
		return ""
	}
	var obj interface{}
	if err = json.Unmarshal(respBody, &obj); err != nil {
		log.Error("Failed to parse JSON: %v", err)
	}
	obj, err = jsonpointer.Get(obj, "/permalink_url")
	if err != nil {
		log.Error("Failed to parse JSON: %v", err)
		return ""
	}

	url := obj.(string)
	log.Info("Uploaded. Gyazo URL: %s", url)
	return url
}
