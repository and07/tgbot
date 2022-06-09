package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/Squirrel-Network/gobotapi"
	"github.com/Squirrel-Network/gobotapi/methods"
	"github.com/Squirrel-Network/gobotapi/types"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/uploader"
	"github.com/kelseyhightower/envconfig"
	"github.com/kkdai/youtube/v2"
	"go.uber.org/zap"
)

func helper(ctx context.Context, token, name string, f io.Reader) {
	appID := 17294691
	appHash := "cddcfe9c3d0d6d40cab8ed031e454df3"

	// Using custom session storage.
	// You can save session to database, e.g. Redis, MongoDB or postgres.
	// See memorySession for implementation details.
	sessionStorage := &memorySession{}

	client := telegram.NewClient(appID, appHash, telegram.Options{
		SessionStorage: sessionStorage,
	})

	if err := client.Run(context.Background(), func(ctx context.Context) error {

		// Checking auth status.
		status, err := client.Auth().Status(ctx)
		if err != nil {
			return err
		}
		// Can be already authenticated if we have valid session in
		// session storage.
		if !status.Authorized {
			// Otherwise, perform bot authentication.
			if _, err := client.Auth().Bot(ctx, token); err != nil {
				return err
			}
		}

		// All good, manually authenticated.
		log.Println("Done Auth")

		// It is only valid to use client while this function is not returned
		// and ctx is not cancelled.
		api := client.API()

		// Helper for uploading. Automatically uses big file upload when needed.
		u := uploader.NewUploader(api)

		// Helper for sending messages.
		sender := message.NewSender(api).WithUploader(u)

		// Uploading directly from path. Note that you can do it from
		// io.Reader or buffer, see From* methods of uploader.
		log.Println("Uploading file")

		log.Println("ReadAll")
		dat, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		log.Println("ReadAll Done")

		upload, err := u.FromBytes(ctx, name, dat)
		if err != nil {
			return fmt.Errorf("upload %q: %w", name, err)
		}

		// Resolving target. Can be telephone number or @nickname of user,
		// group or channel.
		target := sender.Resolve("@and07mbot")

		// Sending message with media.
		log.Println("Sending file")
		/*
			_, err = target.Video(ctx, upload)
			if err != nil {
				return fmt.Errorf("send: %w", err)
			}
		*/

		// Now we have uploaded file handle, sending it as styled message.
		// First, preparing message.
		document := message.UploadedDocument(upload)

		// You can set MIME type, send file as video or audio by using
		// document builder:
		document.
			MIME("video/mp4").
			Filename(name).
			Video()

		_, err = target.Media(ctx, document)
		if err != nil {
			return fmt.Errorf("send: %w", err)
		}

		// Return to close client connection and free up resources.
		return nil
	}); err != nil {
		panic(err)
	}
}

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

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "hello\n")
	})

	go func() { http.ListenAndServe(":"+serviceEnv.Port, nil) }()
	/*
		bot, err := tgbotapi.NewBotAPI(serviceEnv.Token)
		if err != nil {
			log.Panic(err)
		}
	*/
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
		/*
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
			/
			videoSend := tgbotapi.NewVideo(msg.Chat.ID, tgbotapi.FileReader{
				Name:   msg.Text + ".mp4",
				Reader: stream,
			})

			msgXX, err := bot.Send(videoSend)
			if err != nil {
				logger.Error(err)
				text.Text = err.Error()
				client.Invoke(text)
				return
			}


			logger.Info("Done ", msgXX)
		*/
		helper(context.Background(), serviceEnv.Token, msg.Text+".mp4", stream)

		logger.Info("Done")

	})
	client.Run()

}
