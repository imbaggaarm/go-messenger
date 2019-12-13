package go_messenger

type (
	TemplateType     string
	AttachmentType   string
	QuickReplyType   string
)

const (
	AttachmentTypeAudio    = AttachmentType("audio")
	AttachmentTypeVideo    = AttachmentType("video")
	AttachmentTypeImage    = AttachmentType("image")
	AttachmentTypeFile     = AttachmentType("file")
	AttachmentTypeTemplate = AttachmentType("template")

	TemplateTypeGeneric = TemplateType("generic")
	TemplateTypeButton  = TemplateType("button")
	TemplateTypeReceipt = TemplateType("receipt")
	TemplateTypeAirline = TemplateType("airline_boardingpass")
	TemplateTypeMedia   = TemplateType("media")

	QuickReplyTypeText            = QuickReplyType("text")
	QuickReplyTypeUserPhoneNumber = QuickReplyType("user_phone_number")
	QuickReplyTypeUserEmail       = QuickReplyType("user_email")
)

type (
	Message struct {
		Text         string       `json:"text,omitempty"`
		Attachment   Attachment   `json:"attachment,omitempty"`
		QuickReplies []QuickReply `json:"quick_replies,omitempty"`
	}

	Attachment struct {
		Type    AttachmentType    `json:"type"`
		Payload AttachmentPayload `json:"payload"`
	}
	AttachmentPayload struct {
		TemplateType TemplateType `json:"template_type,omitempty"`
		Text         string       `json:"text,omitempty"`
		Elements     interface{}  `json:"elements,omitempty"`
		Buttons      []Button     `json:"buttons,omitempty"`
		URL          string       `json:"url,omitempty"`
	}
	Button struct {
		Type    string `json:"type"`
		Title   string `json:"title"`
		Payload string `json:"payload"`
	}
	QuickReply struct {
		ContentType QuickReplyType `json:"content_type"`
		Title       string         `json:"title,omitempty"`
		Payload     string         `json:"payload,omitempty"`
		ImageURL    string         `json:"image_url,omitempty"`
	}
)
