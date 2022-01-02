package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

var (
	APIKey   string
	UserPass string
)

const (
	GARAGE = 8
	WEST   = 7
	EAST   = 6
	POOL   = 5
	FRONT  = 4
	DOOR   = 2
)

func main() {
	if key, ok := os.LookupEnv("CFAPIKEY"); ok {
		APIKey = key
	} else {
		log.Fatal("CFAPIKEY environment variable must be set")
	}

	if userpass, ok := os.LookupEnv("USERPASS"); ok {
		UserPass = userpass
	} else {
		log.Fatal("USERPASS environment variable must be set")
	}

	for {
		err := Capture(GARAGE)
		if err != nil {
			log.Println(err)
		} else {
			err := Sync(GARAGE)
			if err != nil {
				log.Println(err)
			}
		}

		err = Capture(WEST)
		if err != nil {
			log.Println(err)
		} else {
			err := Sync(WEST)
			if err != nil {
				log.Println(err)
			}
		}
		err = Capture(FRONT)
		if err != nil {
			log.Println(err)
		} else {
			err := Sync(FRONT)
			if err != nil {
				log.Println(err)
			}
		}

		err = Capture(DOOR)
		if err != nil {
			log.Println(err)
		} else {
			err := Sync(DOOR)
			if err != nil {
				log.Println(err)
			}
		}

		err = Capture(EAST)
		if err != nil {
			log.Println(err)
		} else {
			err := Sync(EAST)
			if err != nil {
				log.Println(err)
			}
		}

		err = Sync(DOOR)
		if err != nil {
			log.Println(err)
		} else {
			err := Sync(DOOR)
			if err != nil {
				log.Println(err)
			}
		}
		time.Sleep(time.Second * 3)
	}
}

func Capture(camera int) error {
	var (
		buf bytes.Buffer
	)
	url := fmt.Sprintf("rtsp://%s@192.168.86.74:8554/Streaming/Channels/%d01", UserPass, camera)
	outfile := fmt.Sprintf("%d.jpg", camera)
	cmd := exec.Command("ffmpeg", "-y", "-i", url, "-frames:v", "1", outfile)
	cmd.Stderr = &buf
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("Could not run ffmpeg command for camera %d: %s, %s", camera, err, buf)
	}
	_, err = os.Stat(outfile)
	return err
}

func Sync(camera int) error {
	outfile := fmt.Sprintf("%d.jpg", camera)
	client := &http.Client{}
	buf, err := ioutil.ReadFile(outfile)
	if err != nil {
		return fmt.Errorf("Could not read file: %s", err)
	}
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/51a194e0a18776642ce8563cdaf9d3bd/storage/kv/namespaces/f214fc01830844c183d058e7cf335077/values/%s", outfile)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(buf))
	if err != nil {
		return fmt.Errorf("Could not create request %s", err)
	}
	req.Header.Add("Content-type", "plain")
	req.Header.Add("Authorization", "Bearer "+APIKey)
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Could not send request to cf: %s", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("CF API returned non-200 status code %d", resp.StatusCode)
	}
	return nil
}
