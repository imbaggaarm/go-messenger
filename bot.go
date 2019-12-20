// go-messenger is a wrapper for Facebook Messenger API alternative to pymessenger
package messenger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

type (
	NotificationType string
	SenderAction     string
)

const (
	kGraphUrl         = "https://graph.facebook.com/v"
	kAccessToken      = "access_token"
	DefaultApiVersion = "2.6"

	NotificationTypeRegular    = NotificationType("REGULAR")
	NotificationTypeSilentPush = NotificationType("SILENT_PUSH")
	NotificationNoPush         = NotificationType("NO_PUSH")

	SenderActionMarkSeen  = SenderAction("mark_seen")
	SenderActionTypingOn  = SenderAction("typing_on")
	SenderActionTypingOff = SenderAction("typing_off")
)

type (
	Payload struct {
		Recipient        *Recipient       `json:"recipient,omitempty"`
		GetStarted       *GetStarted      `json:"get_started,omitempty"`
		PersistentMenu   []PersistentMenu `json:"persistent_menu,omitempty"`
		NotificationType NotificationType `json:"notification_type,omitempty"`
		SenderAction     SenderAction     `json:"sender_action,omitempty"`
		DeletedFields    []string         `json:"fields,omitempty"`
		Message          *Message         `json:"message,omitempty"`
	}

	Recipient struct {
		ID string `json:"id"`
	}

	GetStarted struct {
		Payload string `json:"payload,omitempty"`
	}

	PersistentMenu struct {
		Locale                string   `json:"locale"`
		ComposerInputDisabled bool     `json:"composer_input_disabled"`
		CallToActions         []Button `json:"call_to_actions"`
	}
)

type Bot struct {
	AccessToken string
	ApiVersion  string
	GraphUrl    string
}

// Create a new Bot instance with your page access token, and an api version.
//
// Input:
// 		accessToken: your page access token
//		apiVersion: specified api version, use DefaultAPIVersion if you want to use default api version
// Output:
// 		A Bot instance
func NewBot(accessToken string, apiVersion string) *Bot {
	if apiVersion == "" {
		apiVersion = DefaultApiVersion
	}
	return &Bot{
		AccessToken: accessToken,
		ApiVersion:  apiVersion,
		GraphUrl:    kGraphUrl + apiVersion,
	}
}

// Send raw message with a sub path, a httpMethod, and a payload object
// This method can not be used outside the package
//
// Input:
// 		requestSubPath: sub path of endpoint
// 		method: http method of this request
// 		payload: a Payload object to send
// Output:
// 		Response from API and an error if exists
func (bot *Bot) sendRaw(requestSubPath string, method string, payload Payload) (*http.Response, error) {

	fmt.Println("--------------------")
	defer fmt.Println("--------------------")
	// Create request endpoint with given sub path
	requestEndpoint := bot.GraphUrl + requestSubPath

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Encode the payload into request body
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(payload); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	jsonPayload, _ := json.MarshalIndent(payload, "", "\t")
	fmt.Println(string(jsonPayload))

	req, _ := http.NewRequest(method, requestEndpoint, body)
	req.Header.Add("Content-Type", "application/json")

	// Add access token to request params
	q := req.URL.Query()
	q.Add(kAccessToken, bot.AccessToken)
	req.URL.RawQuery = q.Encode()

	// Start the request
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error request:", err.Error())
		return resp, err
	}

	defer resp.Body.Close()

	if err := json.NewEncoder(body).Encode(resp.Body); err != nil {
		log.Println(err.Error())
		return resp, err
	}

	if resp.StatusCode != 200 {
		log.Println("Error http response -> %v", resp)
		log.Println("Error http " + strconv.Itoa(resp.StatusCode) + " -> " + body.String())
	}

	return resp, err
}

// Send raw message with a payload instance
// https://developers.facebook.com/docs/messenger-platform/reference/send-api/
//
// Input:
// 		payload: a Payload object to send
// Output:
// 		Response from API and an error if exists
func (bot *Bot) SendRawMessage(payload Payload) (*http.Response, error) {
	return bot.sendRaw("/me/messages", http.MethodPost, payload)
}

// Send message to a recipient with recipientID
// https://developers.facebook.com/docs/messenger-platform/reference/send-api/
//
// Input:
// 		recipientID: recipient id to send to
// 		payload: a Payload object to send
// 		notificationType: type of notification, see NotificationType
// Output:
// 		Response from API and an error if exists
func (bot *Bot) SendRecipient(recipientID string, payload Payload, notificationType NotificationType) (*http.Response, error) {
	payload.Recipient = &Recipient{ID: recipientID}
	payload.NotificationType = notificationType
	return bot.SendRawMessage(payload)
}

// Send typing indicators or send read receipts to the specified recipient.
// https://developers.facebook.com/docs/messenger-platform/send-api-reference/sender-actions
//
// Input:
// 		recipientID: recipient id to send to
// 		action: action type (mark_seen, typing_on, typing_off)
// Output:
// 		Response from API and an error if exists
func (bot *Bot) SendAction(recipientID string, action SenderAction, notificationType NotificationType) (*http.Response, error) {
	payload := Payload{SenderAction: action}
	return bot.SendRecipient(recipientID, payload, notificationType)
}

// Send message to the specified recipient.
// https://developers.facebook.com/docs/messenger-platform/send-messages
//
// Input:
// 		recipientID: recipient id to send to
// 		message: a Message object
// Output:
// 		Response from API and an error if exists
func (bot *Bot) SendMessage(recipientID string, message Message) (*http.Response, error) {
	payload := Payload{
		Recipient: &Recipient{ID: recipientID},
		Message:   &message,
	}
	return bot.SendRawMessage(payload)
}

// Send text message to the specified recipient.
// https://developers.facebook.com/docs/messenger-platform/send-messages#sending_text
//
// Input:
// 		recipientID: recipient id to send to
// 		text: a text message
// Output:
// 		Response from API and an error if exists
func (bot *Bot) SendTextMessage(recipientID string, text string) (*http.Response, error) {
	message := Message{
		Text: text,
	}
	return bot.SendMessage(recipientID, message)
}

// Send quick replies to the specified recipient.
// https://developers.facebook.com/docs/messenger-platform/send-messages/quick-replies
//
// Input:
// 		recipientID: recipient id to send to
// 		text: title of message
// 		quickReplies: an array of QuickReply objects, up to 13 elements
// Output:
// 		Response from API and an error if exists
func (bot *Bot) SendQuickReplies(recipientID string, text string, quickReplies []QuickReply) (*http.Response, error) {
	message := Message{
		Text:         text,
		QuickReplies: &quickReplies,
	}
	return bot.SendMessage(recipientID, message)
}

// Send attachment message to the specified recipient.
// https://developers.facebook.com/docs/messenger-platform/send-messages#sending_attachments
//
// Input:
// 		recipientID: recipient id to send to
// 		attachment: an attachment
// Output:
// 		Response from API and an error if exists
func (bot *Bot) SendAttachmentMessage(recipientID string, attachment Attachment) (*http.Response, error) {
	message := Message{
		Attachment: &attachment,
	}
	return bot.SendMessage(recipientID, message)
}

// Send attachment message to the specified recipient using URL.
// https://developers.facebook.com/docs/messenger-platform/send-messages#sending_attachments
//
// Input:
// 		recipientID: recipient id to send to
// 		attachmentType: type of the attachment
// 		attachmentUrl: url of the attachment
// Output:
// 		Response from API and an error if exists
func (bot *Bot) SendAttachmentUrl(recipientID string, attachmentType AttachmentType, attachmentUrl string) (*http.Response, error) {
	attachment := Attachment{
		Type:    attachmentType,
		Payload: AttachmentPayload{URL: attachmentUrl},
	}
	return bot.SendAttachmentMessage(recipientID, attachment)
}

// Send generic message to the specified recipient.
// https://developers.facebook.com/docs/messenger-platform/reference/template/generic
//
// Input:
// 		recipientID: recipient id to send to
// 		elements: an array of Element objects, can up to 10 elements
// Output:
// 		Response from API and an error if exists
func (bot *Bot) SendGenericMessage(recipientID string, elements []Element) (*http.Response, error) {
	attachment := Attachment{
		Type: AttachmentTypeTemplate,
		Payload: AttachmentPayload{
			TemplateType: TemplateTypeGeneric,
			Elements:     elements,
		},
	}
	return bot.SendAttachmentMessage(recipientID, attachment)
}

// Send button message to the specified recipient.
// https://developers.facebook.com/docs/messenger-platform/send-messages/buttons
//
// Input:
// 		recipientID: recipient id to send to
// 		text: text of message to send
// 		buttons: An array of Button objects, can up to 3 buttons
// Output:
// 		Response from API and an error if exists
func (bot *Bot) SendButtonMessage(recipientID string, text string, buttons []Button) (*http.Response, error) {
	attachment := Attachment{
		Type: AttachmentTypeTemplate,
		Payload: AttachmentPayload{
			TemplateType: TemplateTypeButton,
			Text:         text,
			Buttons:      buttons,
		},
	}
	return bot.SendAttachmentMessage(recipientID, attachment)
}

// Send an image message with image url
// https://developers.facebook.com/docs/messenger-platform/send-messages#sending_attachments
//
// Input:
// 		recipientID: recipient id to send to
// 		imageUrl: url of the image that we want to send
// Output:
// 		Response from API and an error if exists
func (bot *Bot) SendImageUrl(recipientID string, imageUrl string) (*http.Response, error) {
	return bot.SendAttachmentUrl(recipientID, AttachmentTypeImage, imageUrl)
}

// Send an audio message with audio url
// https://developers.facebook.com/docs/messenger-platform/send-messages#sending_attachments
//
// Input:
// 		recipientID: recipient id to send to
// 		imageUrl: url of the audio to send
// Output:
// 		Response from API and and error if exists
func (bot *Bot) SendAudioUrl(recipientID string, audioUrl string) (*http.Response, error) {
	return bot.SendAttachmentUrl(recipientID, AttachmentTypeAudio, audioUrl)
}

// Send a video message with video url
// https://developers.facebook.com/docs/messenger-platform/send-messages#sending_attachments
//
// Input:
// 		recipientID: recipient id to send to
// 		videoUrl: url of the video to send
// Output:
// 		Response from API and an error if exists
func (bot *Bot) SendVideoUrl(recipientID string, videoUrl string) (*http.Response, error) {
	return bot.SendAttachmentUrl(recipientID, AttachmentTypeVideo, videoUrl)
}

// Send file with file url
// https://developers.facebook.com/docs/messenger-platform/send-messages#sending_attachments
//
// Input:
// 		recipientID: recipient id to send to
// 		fileUrl: url of the file
// Output:
// 		Response from API and an error if exists
func (bot *Bot) SendFileUrl(recipientID string, fileUrl string) (*http.Response, error) {
	return bot.SendAttachmentUrl(recipientID, AttachmentTypeFile, fileUrl)
}

// Set a get started button for the page, this button will be shown on welcome screen for new users
// https://developers.facebook.com/docs/messenger-platform/reference/messenger-profile-api/get-started-button
//
// Input:
// 		gsPayload: a Payload object has GetStarted property as described by the API docs
// Output:
// 		Response from API and an error if exists
func (bot *Bot) SetGetStarted(gsPayload Payload) (*http.Response, error) {
	return bot.sendRaw("/me/messenger_profile", http.MethodPost, gsPayload)
}

// Remove get started button from the page
// https://developers.facebook.com/docs/messenger-platform/reference/messenger-profile-api/#delete
//
// Output:
// 		Response from API and an error if exists
func (bot *Bot) RemoveGetStarted() (*http.Response, error) {
	payload := Payload{DeletedFields: []string{"get_started"}}
	return bot.sendRaw("/me/messenger_profile", http.MethodDelete, payload)
}

// Set a persistent menu for the page. You have to set a get started button before use this
// https://developers.facebook.com/docs/messenger-platform/reference/messenger-profile-api/persistent-menu
//
// Input:
// 		pmPayload: a Payload object which has PersistentMenu property as described by the API docs
// Output:
// 		Response from API and an error if exists
func (bot *Bot) SetPersistentMenu(pmPayload Payload) (*http.Response, error) {
	return bot.sendRaw("/me/messenger_profile", http.MethodPost, pmPayload)
}

// Remove persistent menu from the page
// https://developers.facebook.com/docs/messenger-platform/reference/messenger-profile-api/#delete
//
// Output:
// 		Response from API and an error if exists
func (bot *Bot) RemovePersistentMenu() (*http.Response, error) {
	payload := Payload{DeletedFields: []string{"persistent_menu"}}
	return bot.sendRaw("/me/messenger_profile", http.MethodDelete, payload)
}
