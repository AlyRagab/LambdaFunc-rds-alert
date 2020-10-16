package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	Records []struct {
		SNS struct {
			Type       string `json:"Type"`
			Timestamp  string `json:"Timestamp"`
			SNSMessage string `json:"Message"`
		} `json:"Sns"`
	} `json:"Records"`
}

type SNSMessage struct {
	AlarmName      string `json:"AlarmName"`
	NewStateValue  string `json:"NewStateValue"`
	NewStateReason string `json:"NewStateReason"`
}

type SlackMessage struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Text  string `json:"text"`
	Color string `json:"color"`
	Title string `json:"title"`
}

// 1.Parsing the SNS data based on the AlarmName.
// 2.Sending it as argument to the slack function "buildMessage".
func handler(request Request) error {
	var snsMessage SNSMessage
	err := json.Unmarshal([]byte(request.Records[0].SNS.SNSMessage), &snsMessage)
	if err != nil {
		return err
	}

	slackMessage := buildMessage(snsMessage)
	sendToSlack(slackMessage)
	return nil
}

// Building/Formating Slack Message from SlackMessage Struct
func buildMessage(message SNSMessage) SlackMessage {
	return SlackMessage{
		Text: fmt.Sprintf("%s", message.AlarmName),
		Attachments: []Attachment{
			Attachment{
				Text:  message.NewStateReason,
				Color: "danger",
				Title: "Reason",
			},
		},
	}
}

func sendToSlack(message SlackMessage) error {
	client := &http.Client{}
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", os.Getenv("SLACK_WEBHOOK"), bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		return err
	}

	return nil
}
func main() {
	lambda.Start(handler)
}
