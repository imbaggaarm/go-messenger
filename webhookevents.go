package messenger

type AccountLinkingStatus string
type PolicyEnforcementAction string
type ReferralSource string
type ReactionAction string

const (
	AccountLinkingStatusLinked   = AccountLinkingStatus("linked")
	AccountLinkingStatusUnlinked = AccountLinkingStatus("unlinked")

	PolicyEnforcementActionWarning = PolicyEnforcementAction("warning")
	PolicyEnforcementActionBlock   = PolicyEnforcementAction("block")
	PolicyEnforcementActionUnblock = PolicyEnforcementAction("unblock")

	ReferralSourceMessengerCode      = ReferralSource("MESSENGER_CODE")
	ReferralSourceDiscoverTab        = ReferralSource("DISCOVER_TAB")
	ReferralSourceAds                = ReferralSource("ADS")
	ReferralSourceShortlink          = ReferralSource("SHORTLINK")
	ReferralSourceCustomerChatPlugin = ReferralSource("CUSTOMER_CHAT_PLUGIN")

	ReactionActionReact   = ReactionAction("react")
	ReactionActionUnreact = ReactionAction("unreact")
)

type WebhookEvent struct {
	Object string  `json:"object"`
	Entry  []Entry `json:"entry"`
}

type Entry struct {
	ID        string          `json:"id"`
	Time      int             `json:"time"`
	Messaging *[]EntryMessage `json:"messaging,omitempty"`
	Standby   *[]EntryMessage `json:"standby,omitempty"`
}

type EntryMessage struct {
	Sender    Sender    `json:"sender"`
	Recipient Recipient `json:"recipient"`
	Timestamp int64     `json:"timestamp,omitempty"`

	Message              *WebhookMessage    `json:"message,omitempty"` // represent for messages event and message_echoes event
	AccountLinking       *AccountLinking    `json:"account_linking,omitempty"`
	MessageDelivery      *MessageDelivery   `json:"message_delivery,omitempty"`
	GamePlay             *GamePlay          `json:"game_play,omitempty"`
	PassThreadControl    *Handover          `json:"pass_thread_control,omitempty"`
	TakeThreadControl    *Handover          `json:"take_thread_control,omitempty"`
	AppRoles             interface{}        `json:"app_roles,omitempty"`
	RequestThreadControl *Handover          `json:"request_thread_control,omitempty"`
	Optin                *Optin             `json:"optin,omitempty"`
	PolicyEnforcement    *PolicyEnforcement `json:"policy_enforcement,omitempty"`
	Postback             *Postback          `json:"postback,omitempty"`
	Reaction             *Reaction          `json:"reaction,omitempty"`
	MessageRead          *MessageRead       `json:"message_read,omitempty"`
	Referral             *Referral          `json:"referral,omitempty"`
}

type Sender struct {
	ID string `json:"id"`
}

type MessageEcho struct {
	IsEcho   bool   `json:"is_echo,omitempty"`
	AppID    string `json:"app_id,omitempty"`
	Metadata string `json:"metadata,omitempty"`
	Mid      string `json:"mid,omitempty"`
}

type WebhookMessage struct { //message and message_echoes
	MessageEcho
	Text        *string       `json:"text,omitempty"`
	Attachments *[]Attachment `json:"attachments,omitempty"`
	QuickReply  *QuickReply   `json:"quick_reply,omitempty"`
	ReplyTo     *ReplyTo      `json:"reply_to,omitempty"`
}

type ReplyTo struct {
	Mid string `json:"mid,omitempty"`
}

type AccountLinking struct {
	Status            AccountLinkingStatus `json:"status"` // linked or unlinked
	AuthorizationCode string               `json:"authorization_code"`
}

type MessageDelivery struct {
	Mids      []string `json:"mids"`
	Watermark int      `json:"watermark"`
}

type GamePlay struct {
	GameID      string `json:"game_id"`
	PlayerID    string `json:"player_id"`
	ContextType string `json:"context_type"`
	ContextID   string `json:"context_id,omitempty"`
	Score       int    `json:"score,omitempty"`
	Payload     string `json:"payload,omitempty"`
}

type Handover struct {
	Metadata            string `json:"metadata,omitempty"`
	NewOwnerAppID       string `json:"new_owner_app_id,omitempty"`
	PreviousOwnerAppID  string `json:"previous_owner_app_id,omitempty"`
	RequestedOwnerAppID string `json:"requested_owner_app_id,omitempty"`
	//pageID
}

type Optin struct {
	Ref     string `json:"ref"`
	UserRef string `json:"user_ref"`
}

type PolicyEnforcement struct {
	Action PolicyEnforcementAction `json:"action"`
	Reason string                  `json:"reason,omitempty"` //This field is absent if action is unblock
}

type Postback struct {
	Title    string   `json:"title,omitempty"`
	Payload  string   `json:"payload,omitempty"`
	Referral Referral `json:"referral,omitempty"`
}

type Reaction struct {
	Reaction string         `json:"reaction"`
	Emoji    string         `json:"emoji"`
	Action   ReactionAction `json:"action"`
	Mid      string         `json:"mid"`
}

type MessageRead struct {
	Watermark int `json:"watermark"`
}

type Referral struct {
	Source     ReferralSource `json:"source"`
	Type       string         `json:"type"`
	Ref        string         `json:"ref,omitempty"`
	RefererUri string         `json:"referer_uri,omitempty"`
}
