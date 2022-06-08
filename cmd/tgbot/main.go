package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Squirrel-Network/gobotapi"
	"github.com/Squirrel-Network/gobotapi/methods"
	"github.com/Squirrel-Network/gobotapi/types"
	"github.com/kelseyhightower/envconfig"
	"github.com/kkdai/youtube/v2"
	"go.uber.org/zap"
)

func main() {

	// logger
	simpleLogger, err := zap.NewDevelopment()
	if err != nil {
		os.Exit(1)
	}
	// flushes buffer, if any
	defer func() {
		if err := simpleLogger.Sync(); err != nil {
			fmt.Println("OOOPS Logger sync failed")
		}
	}()

	logger := simpleLogger.Sugar()

	var serviceEnv Configuration
	err = envconfig.Process("", &serviceEnv)
	if err != nil {
		logger.Error("msg", "failed to parse service env", "error", err)
		os.Exit(1)
	}

	clientYouTube := youtube.Client{}

	client := gobotapi.NewClient(serviceEnv.Token)
	client.OnMessage(func(client gobotapi.Client, msg types.Message) {

		text := &methods.SendMessage{
			ChatID: msg.Chat.ID,
			Text:   "Hello World!" + msg.Text,
		}
		logger.Info("GetVideo")
		video, err := clientYouTube.GetVideo(msg.Text)
		if err != nil {
			text.Text = err.Error()
			client.Invoke(text)
			return
		}
		logger.Info("GetVideo Done")

		logger.Info("GetStream")
		formats := video.Formats.WithAudioChannels() // only get videos with audio
		stream, _, err := clientYouTube.GetStream(video, &formats[0])
		if err != nil {
			text.Text = err.Error()
			client.Invoke(text)
			return
		}
		logger.Info("GetStream Done")

		logger.Info("ReadAll")
		dat, err := ioutil.ReadAll(stream)
		if err != nil {
			text.Text = err.Error()
			client.Invoke(text)
			return
		}
		logger.Info("ReadAll Done")

		client.Invoke(&methods.SendVideo{
			ChatID: msg.Chat.ID,
			Video: types.InputFile{
				Name:  msg.Text + ".mp4",
				Bytes: dat,
			},
		})

		logger.Info("Done")

	})
	client.Run()

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) { fmt.Fprintf(w, "hello\n") })
	http.ListenAndServe(serviceEnv.Port, nil)
}
