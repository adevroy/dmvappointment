package main

import (
	"fmt"
	"strconv"
)

type slackData struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type Attachment struct {
	Text           string   `json:"text"`
	Fallback       string   `json:"fallback"`
	CallbackID     string   `json:"callback_id"`
	Color          string   `json:"color"`
	AttachmentType string   `json:"attachment_type"`
	Actions        []Action `json:"actions"`
}

type Action struct {
	Name    string `json:"name"`
	Text    string `json:"text"`
	Type    string `json:"type"`
	Value   string `json:"value"`
	Style   string `json:"style,omitempty"`
	URL     string `json:"url"`
	Confirm struct {
		Title       string `json:"title"`
		Text        string `json:"text"`
		OkText      string `json:"ok_text"`
		DismissText string `json:"dismiss_text"`
	} `json:"confirm,omitempty"`
}

func getSlackFormatedMessage(apptString string, v LocationData) slackData {
	sData := slackData{}
	sData.Text = slackUserID + " " + v.Name
	sData.Attachments = []Attachment{
		{
			Text: fmt.Sprintf("%v - %v, %v, %v", v.Street1, v.Street2, v.City, v.Zip),
		},
		{
			Text: apptString,
			Actions: []Action{{
				Name:  "Book Now",
				Text:  "Book Now",
				Type:  "button",
				Style: "primary",
				URL:   appointmentLinkURL + strconv.Itoa(v.LocAppointments[0].LocationID),
			}},
		},
	}
	return sData
}
