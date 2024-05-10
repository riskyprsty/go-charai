package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/harmony-ai-solutions/CharacterAI-Golang/cai"
)

var (
	token     string = "your_charai_token"
	isPlus    bool   = false
	caiClient *cai.GoCAI
)

func init() {
	// Create AI client
	var errClient error
	caiClient, errClient = cai.NewGoCAI(token, isPlus)
	if errClient != nil {
		fmt.Println(fmt.Errorf("unable to create client, error: %q", errClient))
		os.Exit(1)
	}
}

type Message struct {
	Status string `json:"status"`
	Body   string `json:"body"`
}

func getTestMessage(context *gin.Context) {
	var messages = []Message{
		{Status: "success", Body: "Hello! This is test message"},
	}

	context.IndentedJSON(http.StatusOK, messages)
}

func getAiParticipant(context *gin.Context) {
	character := context.Query("charid")

	// Get AI chat data
	chatData, errChat := caiClient.Chat.GetChat(character)
	if errChat != nil {
		if strings.Contains(errChat.Error(), "404") {
			// if chat with ai not initialized, initialized first
			chatData, errChat = caiClient.Chat.NewChat(character)
			if errChat != nil {
				fmt.Println(fmt.Errorf("[x] Unable to create chat, error: %q", errChat))
				var messages = []Message{
					{Status: "false", Body: errChat.Error()},
				}

				context.IndentedJSON(http.StatusOK, messages)
			}
		} else {
			fmt.Println(fmt.Errorf("[X] Unable to fetch chat data, error: %q", errChat))
			var messages = []Message{
				{Status: "false", Body: errChat.Error()},
			}

			context.IndentedJSON(http.StatusOK, messages)
		}
	}
	// looking for AI participant
	var aiParticipant *cai.ChatParticipant
	for _, participant := range chatData.Participants {
		if !participant.IsHuman {
			aiParticipant = participant
			break
		}
	}

	var messages = []Message{
		{Status: "success", Body: aiParticipant.Name},
	}

	context.IndentedJSON(http.StatusOK, messages)

}

func getAiProfileImage(context *gin.Context) {
	character := context.Query("charid")

	charInfo, errInfo := caiClient.Character.Info(character, token)

	if errInfo != nil {
		fmt.Println("[x] Unable to get character profile image info")
		var messages = []Message{
			{Status: "false", Body: errInfo.Error()},
		}

		context.IndentedJSON(http.StatusOK, messages)
	}

	characterMap := charInfo["character"].(map[string]interface{})
	avatarImage := characterMap["avatar_file_name"].(string)

	var messages = []Message{
		{Status: "success", Body: avatarImage},
	}

	context.IndentedJSON(http.StatusOK, messages)

}

func getMessage(context *gin.Context) {
	character := context.Query("charid")
	message := context.Query("body")

	// Get AI chat data
	chatData, errChat := caiClient.Chat.GetChat(character)
	if errChat != nil {
		if strings.Contains(errChat.Error(), "404") {
			// if chat with ai not initialized, initialized first
			chatData, errChat = caiClient.Chat.NewChat(character)
			if errChat != nil {
				fmt.Println(fmt.Errorf("[x] Unable to create chat, error: %q", errChat))
				os.Exit(3)
			}
		} else {
			fmt.Println(fmt.Errorf("[X] Unable to fetch chat data, error: %q", errChat))
			os.Exit(2)
		}
	}
	// looking for AI participant
	var aiParticipant *cai.ChatParticipant
	for _, participant := range chatData.Participants {
		if !participant.IsHuman {
			aiParticipant = participant
			break
		}
	}

	// send message to AI
	messageResult, errMessage := caiClient.Chat.SendMessage(chatData.ExternalID, aiParticipant.User.Username, message, nil)
	if errMessage != nil {
		fmt.Println(fmt.Errorf("[x] Unable to send message. Error: %v", errMessage))
	}

	// Handle result
	if len(messageResult.Replies) > 0 {
		firstReply := messageResult.Replies[0]
		log := fmt.Sprintf("%v: %v", aiParticipant.Name, firstReply.Text)
		fmt.Println(log)

		var messages = []Message{
			{Status: "success", Body: firstReply.Text},
		}

		context.IndentedJSON(http.StatusOK, messages)
	}

}

func main() {
	router := gin.Default()

	router.GET("/test", getTestMessage)
	router.GET("/ai", getAiParticipant)
	router.GET("/image", getAiProfileImage)
	router.GET("/message", getMessage)
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.Run()
}
