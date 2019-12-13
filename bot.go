package go_messenger

import (
	"bytes"
	"encoding/json"
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
	defaultAPIVersion = 2.6

	NotificationTypeRegular    = NotificationType("REGULAR")
	NotificationTypeSilentPush = NotificationType("SILENT_PUSH")
	NotificationNoPush         = NotificationType("NO_PUSH")

	SenderActionMarkSeen  = SenderAction("mark_seen")
	SenderActionTypingOn  = SenderAction("typing_on")
	SenderActionTypingOff = SenderAction("typing_off")
)

type (
	Payload struct {
		Recipient        Recipient        `json:"recipient,omitempty"`
		GetStarted       GetStarted       `json:"get_started,omitempty"`
		PersistentMenu   []PersistentMenu `json:"persistent_menu,omitempty"`
		NotificationType NotificationType `json:"notification_type,omitempty"`
		SenderAction     SenderAction     `json:"sender_action,omitempty"`
		DeletedFields    []string         `json:"fields,omitempty"`
		Message          Message          `json:"message,omitempty"`
	}

	Recipient struct {
		Id string `json:"id"`
	}
	GetStarted struct {
		Payload string `json:"payload"`
	}
	PersistentMenu struct {
		Locale                string   `json:"locale"`
		ComposerInputDisabled bool     `json:"composer_input_disabled"`
		CallToActions         []Button `json:"call_to_actions"`
	}
)

type Bot struct {
	AccessToken string
	ApiVersion  int
	GraphUrl    string
}

// Create new bot with your page access token, and an api version.
// If you want to use
func NewBot(accessToken string, apiVersion int) Bot {
	if apiVersion == -1 {
		apiVersion = defaultAPIVersion
	}
	return Bot{
		AccessToken: accessToken,
		ApiVersion:  apiVersion,
		GraphUrl:    kGraphUrl + string(apiVersion),
	}
}

func (bot *Bot) sendRaw(requestSubPath string, method string, payload Payload) (*http.Response, error) {
	requestEndpoint := bot.GraphUrl + requestSubPath

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(payload); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	req, _ := http.NewRequest(method, requestEndpoint, body)
	req.Header.Add("Content-Type", "application/json")

	// Add access token to request params
	q := req.URL.Query()
	q.Add(kAccessToken, bot.AccessToken)
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error request:", err.Error())
		return resp, err
	}

	if err := json.NewEncoder(body).Encode(resp.Body); err != nil {
		log.Println(err.Error())
		return resp, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println("Error http response -> %v", resp)
		log.Println("Error http " + strconv.Itoa(resp.StatusCode) + " -> " + body.String())
	}

	return resp, err
}

func (bot *Bot) SendRawMessage(payload Payload) (*http.Response, error) {
	return bot.sendRaw("/me/messages", http.MethodPost, payload)
}

func (bot *Bot) SendRecipient(recipientId string, payload Payload, notificationType NotificationType) (*http.Response, error) {
	payload.Recipient = Recipient{Id: recipientId}
	payload.NotificationType = notificationType
	return bot.SendRawMessage(payload)
}

func (bot *Bot) SendAction(recipientId string, action SenderAction, notificationType NotificationType) (*http.Response, error) {
	payload := Payload{SenderAction: action}
	return bot.SendRecipient(recipientId, payload, notificationType)
}

func (bot *Bot) SendMessage(recipientId string, message Message) (*http.Response, error) {
	payload := Payload{
		Recipient: Recipient{Id: recipientId},
		Message:   message,
	}
	return bot.SendRawMessage(payload)
}

func (bot *Bot) SendTextMessage(recipientId string, text string) (*http.Response, error) {
	message := Message{
		Text: text,
	}
	return bot.SendMessage(recipientId, message)
}

func (bot *Bot) SendQuickReplies(recipientId string, text string, quickReplies []QuickReply) (*http.Response, error) {
	message := Message{
		Text:         text,
		QuickReplies: quickReplies,
	}
	return bot.SendMessage(recipientId, message)
}

func (bot *Bot) SendAttachmentMessage(recipientId string, attachment Attachment) (*http.Response, error) {
	message := Message{
		Attachment: attachment,
	}
	return bot.SendMessage(recipientId, message)
}

func (bot *Bot) SendAttachmentUrl(recipientId string, attachmentType AttachmentType, attachmentUrl string) (*http.Response, error) {
	attachment := Attachment{
		Type:    attachmentType,
		Payload: AttachmentPayload{URL: attachmentUrl},
	}
	return bot.SendAttachmentMessage(recipientId, attachment)
}

func (bot *Bot) SendGenericMessage(recipientId string, elements interface{}) (*http.Response, error) {
	attachment := Attachment{
		Type: AttachmentTypeTemplate,
		Payload: AttachmentPayload{
			TemplateType: TemplateTypeGeneric,
			Elements:     elements,
		},
	}
	return bot.SendAttachmentMessage(recipientId, attachment)
}

func (bot *Bot) SendButtonMessage(recipientId string, text string, buttons []Button) (*http.Response, error) {
	attachment := Attachment{
		Type: AttachmentTypeTemplate,
		Payload: AttachmentPayload{
			TemplateType: TemplateTypeButton,
			Text:         text,
			Buttons:      buttons,
		},
	}
	return bot.SendAttachmentMessage(recipientId, attachment)
}

func (bot *Bot) SendImageUrl(recipientId string, imageUrl string) (*http.Response, error) {
	return bot.SendAttachmentUrl(recipientId, AttachmentTypeImage, imageUrl)
}

func (bot *Bot) SendAudioUrl(recipientId string, audioUrl string) (*http.Response, error) {
	return bot.SendAttachmentUrl(recipientId, AttachmentTypeAudio, audioUrl)
}

func (bot *Bot) SendVideoUrl(recipientId string, videoUrl string) (*http.Response, error) {
	return bot.SendAttachmentUrl(recipientId, AttachmentTypeVideo, videoUrl)
}

func (bot *Bot) SendFileUrl(recipientId string, fileUrl string) (*http.Response, error) {
	return bot.SendAttachmentUrl(recipientId, AttachmentTypeFile, fileUrl)
}

func (bot *Bot) SetGetStarted(gsPayload Payload) (*http.Response, error) {
	return bot.sendRaw("/me/messenger_profile", http.MethodPost, gsPayload)
}

func (bot *Bot) RemoveGetStarted() (*http.Response, error) {
	payload := Payload{DeletedFields: []string{"get_started"}}
	return bot.sendRaw("/me/messenger_profile", http.MethodDelete, payload)
}

func (bot *Bot) SetPersistentMenu(pmPayload Payload) (*http.Response, error) {
	return bot.sendRaw("/me/messenger_profile", http.MethodPost, pmPayload)
}

func (bot *Bot) RemovePersistentMenu() (*http.Response, error) {
	payload := Payload{DeletedFields: []string{"persistent_menu"}}
	return bot.sendRaw("/me/messenger_profile", http.MethodDelete, payload)
}
