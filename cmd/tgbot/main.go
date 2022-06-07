package main

import "github.com/Squirrel-Network/gobotapi"
import "github.com/Squirrel-Network/gobotapi/types"
import "github.com/Squirrel-Network/gobotapi/methods"

func main() {
    client := gobotapi.NewClient("YOUR_TOKEN")
    client.OnMessage(func(message types.Message) {
        client.Invoke(&methods.SendMessage{
            ChatID: message.Chat.ID,
            Text:   "Hello World!",
        })
    })
    client.Run()
}
