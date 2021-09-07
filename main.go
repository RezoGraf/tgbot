package main

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/pkg/errors"
	tb "gopkg.in/tucnak/telebot.v2"
)

type User struct {
	ID int `json:"id"`

	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
	IsBot        bool   `json:"is_bot"`

	// Returns only in getMe
	CanJoinGroups   bool `json:"can_join_groups"`
	CanReadMessages bool `json:"can_read_all_group_messages"`
	SupportsInline  bool `json:"supports_inline_queries"`
}

type MessageEntity struct {
	// Specifies entity type.
	Type EntityType `json:"type"`

	// Offset in UTF-16 code units to the start of the entity.
	Offset int `json:"offset"`

	// Length of the entity in UTF-16 code units.
	Length int `json:"length"`

	// (Optional) For EntityTextLink entity type only.
	//
	// URL will be opened after user taps on the text.
	URL string `json:"url,omitempty"`

	// (Optional) For EntityTMention entity type only.
	User *User `json:"user,omitempty"`

	// (Optional) For EntityCodeBlock entity type only.
	Language string `json:"language,omitempty"`
}

type Message struct {
	ID int `json:"message_id"`

	InlineID string `json:"-"`

	// For message sent to channels, Sender will be nil
	Sender *User `json:"from"`

	// Unixtime, use Message.Time() to get time.Time
	Unixtime int64 `json:"date"`

	// Conversation the message belongs to.
	Chat *Chat `json:"chat"`

	// Sender of the message, sent on behalf of a chat.
	SenderChat *Chat `json:"sender_chat"`

	// For forwarded messages, sender of the original message.
	OriginalSender *User `json:"forward_from"`

	// For forwarded messages, chat of the original message when
	// forwarded from a channel.
	OriginalChat *Chat `json:"forward_from_chat"`

	// For forwarded messages, identifier of the original message
	// when forwarded from a channel.
	OriginalMessageID int `json:"forward_from_message_id"`

	// For forwarded messages, signature of the post author.
	OriginalSignature string `json:"forward_signature"`

	// For forwarded messages, sender's name from users who
	// disallow adding a link to their account.
	OriginalSenderName string `json:"forward_sender_name"`

	// For forwarded messages, unixtime of the original message.
	OriginalUnixtime int `json:"forward_date"`

	// For replies, ReplyTo represents the original message.
	//
	// Note that the Message object in this field will not
	// contain further ReplyTo fields even if it
	// itself is a reply.
	ReplyTo *Message `json:"reply_to_message"`

	// Shows through which bot the message was sent.
	Via *User `json:"via_bot"`

	// (Optional) Time of last edit in Unix
	LastEdit int64 `json:"edit_date"`

	// AlbumID is the unique identifier of a media message group
	// this message belongs to.
	AlbumID string `json:"media_group_id"`

	// Author signature (in channels).
	Signature string `json:"author_signature"`

	// For a text message, the actual UTF-8 text of the message.
	Text string `json:"text"`

	// For registered commands, will contain the string payload:
	//
	// Ex: `/command <payload>` or `/command@botname <payload>`
	Payload string `json:"-"`

	// For text messages, special entities like usernames, URLs, bot commands,
	// etc. that appear in the text.
	Entities []MessageEntity `json:"entities,omitempty"`

	// Some messages containing media, may as well have a caption.
	Caption string `json:"caption,omitempty"`

	// For messages with a caption, special entities like usernames, URLs,
	// bot commands, etc. that appear in the caption.
	CaptionEntities []MessageEntity `json:"caption_entities,omitempty"`

	// For an audio recording, information about it.
	Audio *Audio `json:"audio"`

	// For a general file, information about it.
	Document *Document `json:"document"`

	// For a photo, all available sizes (thumbnails).
	Photo *Photo `json:"photo"`

	// For a sticker, information about it.
	Sticker *Sticker `json:"sticker"`

	// For a voice message, information about it.
	Voice *Voice `json:"voice"`

	// For a video note, information about it.
	VideoNote *VideoNote `json:"video_note"`

	// For a video, information about it.
	Video *Video `json:"video"`

	// For a animation, information about it.
	Animation *Animation `json:"animation"`

	// For a contact, contact information itself.
	Contact *Contact `json:"contact"`

	// For a location, its longitude and latitude.
	Location *Location `json:"location"`

	// For a venue, information about it.
	Venue *Venue `json:"venue"`

	// For a poll, information the native poll.
	Poll *Poll `json:"poll"`

	// For a dice, information about it.
	Dice *Dice `json:"dice"`

	// For a service message, represents a user,
	// that just got added to chat, this message came from.
	//
	// Sender leads to User, capable of invite.
	//
	// UserJoined might be the Bot itself.
	UserJoined *User `json:"new_chat_member"`

	// For a service message, represents a user,
	// that just left chat, this message came from.
	//
	// If user was kicked, Sender leads to a User,
	// capable of this kick.
	//
	// UserLeft might be the Bot itself.
	UserLeft *User `json:"left_chat_member"`

	// For a service message, represents a new title
	// for chat this message came from.
	//
	// Sender would lead to a User, capable of change.
	NewGroupTitle string `json:"new_chat_title"`

	// For a service message, represents all available
	// thumbnails of the new chat photo.
	//
	// Sender would lead to a User, capable of change.
	NewGroupPhoto *Photo `json:"new_chat_photo"`

	// For a service message, new members that were added to
	// the group or supergroup and information about them
	// (the bot itself may be one of these members).
	UsersJoined []User `json:"new_chat_members"`

	// For a service message, true if chat photo just
	// got removed.
	//
	// Sender would lead to a User, capable of change.
	GroupPhotoDeleted bool `json:"delete_chat_photo"`

	// For a service message, true if group has been created.
	//
	// You would receive such a message if you are one of
	// initial group chat members.
	//
	// Sender would lead to creator of the chat.
	GroupCreated bool `json:"group_chat_created"`

	// For a service message, true if supergroup has been created.
	//
	// You would receive such a message if you are one of
	// initial group chat members.
	//
	// Sender would lead to creator of the chat.
	SuperGroupCreated bool `json:"supergroup_chat_created"`

	// For a service message, true if channel has been created.
	//
	// You would receive such a message if you are one of
	// initial channel administrators.
	//
	// Sender would lead to creator of the chat.
	ChannelCreated bool `json:"channel_chat_created"`

	// For a service message, the destination (supergroup) you
	// migrated to.
	//
	// You would receive such a message when your chat has migrated
	// to a supergroup.
	//
	// Sender would lead to creator of the migration.
	MigrateTo int64 `json:"migrate_to_chat_id"`

	// For a service message, the Origin (normal group) you migrated
	// from.
	//
	// You would receive such a message when your chat has migrated
	// to a supergroup.
	//
	// Sender would lead to creator of the migration.
	MigrateFrom int64 `json:"migrate_from_chat_id"`

	// Specified message was pinned. Note that the Message object
	// in this field will not contain further ReplyTo fields even
	// if it is itself a reply.
	PinnedMessage *Message `json:"pinned_message"`

	// Message is an invoice for a payment.
	Invoice *Invoice `json:"invoice"`

	// Message is a service message about a successful payment.
	Payment *Payment `json:"successful_payment"`

	// The domain name of the website on which the user has logged in.
	ConnectedWebsite string `json:"connected_website,omitempty"`

	// Inline keyboard attached to the message.
	ReplyMarkup InlineKeyboardMarkup `json:"reply_markup"`

	VoiceChatSchedule *VoiceChatScheduled `json:"voice_chat_scheduled,omitempty"`

	// For a service message, a voice chat started in the chat.
	VoiceChatStarted *VoiceChatStarted `json:"voice_chat_started,omitempty"`

	// For a service message, a voice chat ended in the chat.
	VoiceChatEnded *VoiceChatEnded `json:"voice_chat_ended,omitempty"`

	// For a service message, some users were invited in the voice chat.
	VoiceChatParticipantsInvited *VoiceChatParticipantsInvited `json:"voice_chat_participants_invited,omitempty"`

	// For a service message, represents the content of a service message,
	// sent whenever a user in the chat triggers a proximity alert set by another user.
	ProximityAlert *ProximityAlertTriggered `json:"proximity_alert_triggered,omitempty"`

	// For a service message, represents about a change in auto-delete timer settings.
	AutoDeleteTimer *MessageAutoDeleteTimerChanged `json:"message_auto_delete_timer_changed,omitempty"`
}

type Callback struct {
	ID string `json:"id"`

	// For message sent to channels, Sender may be empty
	Sender *User `json:"from"`

	// Message will be set if the button that originated the query
	// was attached to a message sent by a bot.
	Message *Message `json:"message"`

	// MessageID will be set if the button was attached to a message
	// sent via the bot in inline mode.
	MessageID string `json:"inline_message_id"`

	// Data associated with the callback button. Be aware that
	// a bad client can send arbitrary data in this field.
	Data string `json:"data"`
}

type Location struct {
	// Latitude
	Lat float32 `json:"latitude"`
	// Longitude
	Lng float32 `json:"longitude"`

	// Horizontal Accuracy
	HorizontalAccuracy *float32 `json:"horizontal_accuracy,omitempty"`

	// Period in seconds for which the location will be updated
	// (see Live Locations, should be between 60 and 86400.)
	LivePeriod int `json:"live_period,omitempty"`

	Heading int `json:"heading,omitempty"`

	ProximityAlertRadius int `json:"proximity_alert_radius,omitempty"`
}

type Query struct {
	// Unique identifier for this query. 1-64 bytes.
	ID string `json:"id"`

	// Sender.
	From User `json:"from"`

	// Sender location, only for bots that request user location.
	Location *Location `json:"location"`

	// Text of the query (up to 512 characters).
	Text string `json:"query"`

	// Offset of the results to be returned, can be controlled by the bot.
	Offset string `json:"offset"`

	// ChatType of the type of the chat, from which the inline query was sent.
	ChatType string `json:"chat_type"`
}

type ChosenInlineResult struct {
	From      User      `json:"from"`
	Location  *Location `json:"location,omitempty"`
	ResultID  string    `json:"result_id"`
	Query     string    `json:"query"`
	MessageID string    `json:"inline_message_id"` // inline messages only!
}

type ShippingQuery struct {
	Sender  *User           `json:"from"`
	ID      string          `json:"id"`
	Payload string          `json:"invoice_payload"`
	Address ShippingAddress `json:"shipping_address"`
}

type Order struct {
	Name        string          `json:"name"`
	PhoneNumber string          `json:"phone_number"`
	Email       string          `json:"email"`
	Address     ShippingAddress `json:"shipping_address"`
}

type PreCheckoutQuery struct {
	Sender   *User  `json:"from"`
	ID       string `json:"id"`
	Currency string `json:"currency"`
	Payload  string `json:"invoice_payload"`
	Total    int    `json:"total_amount"`
	OptionID string `json:"shipping_option_id"`
	Order    Order  `json:"order_info"`
}

type ShippingAddress struct {
	CountryCode string `json:"country_code"`
	State       string `json:"state"`
	City        string `json:"city"`
	StreetLine1 string `json:"street_line1"`
	StreetLine2 string `json:"street_line2"`
	PostCode    string `json:"post_code"`
}

type ChatType string

type ChatPhoto struct {
	// File identifiers of small (160x160) chat photo
	SmallFileID       string `json:"small_file_id"`
	SmallFileUniqueID string `json:"small_file_unique_id"`

	// File identifiers of big (640x640) chat photo
	BigFileID       string `json:"big_file_id"`
	BigFileUniqueID string `json:"big_file_unique_id"`
}

type Rights struct {
	CanBeEdited         bool `json:"can_be_edited"`
	CanChangeInfo       bool `json:"can_change_info"`
	CanPostMessages     bool `json:"can_post_messages"`
	CanEditMessages     bool `json:"can_edit_messages"`
	CanDeleteMessages   bool `json:"can_delete_messages"`
	CanInviteUsers      bool `json:"can_invite_users"`
	CanRestrictMembers  bool `json:"can_restrict_members"`
	CanPinMessages      bool `json:"can_pin_messages"`
	CanPromoteMembers   bool `json:"can_promote_members"`
	CanSendMessages     bool `json:"can_send_messages"`
	CanSendMedia        bool `json:"can_send_media_messages"`
	CanSendPolls        bool `json:"can_send_polls"`
	CanSendOther        bool `json:"can_send_other_messages"`
	CanAddPreviews      bool `json:"can_add_web_page_previews"`
	CanManageVoiceChats bool `json:"can_manage_voice_chats"`
	CanManageChat       bool `json:"can_manage_chat"`
}

type ChatLocation struct {
	Location Location `json:"location,omitempty"`
	Address  string   `json:"address,omitempty"`
}

type Chat struct {
	ID int64 `json:"id"`

	// See ChatType and consts.
	Type ChatType `json:"type"`

	// Won't be there for ChatPrivate.
	Title string `json:"title"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`

	// Still shows whether the user is a member
	// of the chat at the moment of the request.
	Still bool `json:"is_member,omitempty"`

	// Returns only in getChat
	Bio              string        `json:"bio,omitempty"`
	Photo            *ChatPhoto    `json:"photo,omitempty"`
	Description      string        `json:"description,omitempty"`
	InviteLink       string        `json:"invite_link,omitempty"`
	PinnedMessage    *Message      `json:"pinned_message,omitempty"`
	Permissions      *Rights       `json:"permissions,omitempty"`
	SlowMode         int           `json:"slow_mode_delay,omitempty"`
	StickerSet       string        `json:"sticker_set_name,omitempty"`
	CanSetStickerSet bool          `json:"can_set_sticker_set,omitempty"`
	LinkedChatID     int64         `json:"linked_chat_id,omitempty"`
	ChatLocation     *ChatLocation `json:"location,omitempty"`
}

type ChatMemberUpdated struct {
	// Chat where the user belongs to.
	Chat Chat `json:"chat"`

	// From which user the action was triggered.
	From User `json:"from"`

	// Unixtime, use ChatMemberUpdated.Time() to get time.Time
	Unixtime int64 `json:"date"`

	// Previous information about the chat member.
	OldChatMember *ChatMember `json:"old_chat_member"`

	// New information about the chat member.
	NewChatMember *ChatMember `json:"new_chat_member"`

	// (Optional) InviteLink which was used by the user to
	// join the chat; for joining by invite link events only.
	InviteLink *ChatInviteLink `json:"invite_link"`
}

type Update struct {
	ID int `json:"update_id"`

	Message            *Message            `json:"message,omitempty"`
	EditedMessage      *Message            `json:"edited_message,omitempty"`
	ChannelPost        *Message            `json:"channel_post,omitempty"`
	EditedChannelPost  *Message            `json:"edited_channel_post,omitempty"`
	Callback           *Callback           `json:"callback_query,omitempty"`
	Query              *Query              `json:"inline_query,omitempty"`
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result,omitempty"`
	ShippingQuery      *ShippingQuery      `json:"shipping_query,omitempty"`
	PreCheckoutQuery   *PreCheckoutQuery   `json:"pre_checkout_query,omitempty"`
	Poll               *Poll               `json:"poll,omitempty"`
	PollAnswer         *PollAnswer         `json:"poll_answer,omitempty"`
	MyChatMember       *ChatMemberUpdated  `json:"my_chat_member,omitempty"`
	ChatMember         *ChatMemberUpdated  `json:"chat_member,omitempty"`
}

type PollType string

type PollOption struct {
	Text       string `json:"text"`
	VoterCount int    `json:"voter_count"`
}

type EntityType string

type ParseMode = string

const (
	EntityMention       EntityType = "mention"
	EntityTMention      EntityType = "text_mention"
	EntityHashtag       EntityType = "hashtag"
	EntityCashtag       EntityType = "cashtag"
	EntityCommand       EntityType = "bot_command"
	EntityURL           EntityType = "url"
	EntityEmail         EntityType = "email"
	EntityPhone         EntityType = "phone_number"
	EntityBold          EntityType = "bold"
	EntityItalic        EntityType = "italic"
	EntityUnderline     EntityType = "underline"
	EntityStrikethrough EntityType = "strikethrough"
	EntityCode          EntityType = "code"
	EntityCodeBlock     EntityType = "pre"
	EntityTextLink      EntityType = "text_link"
)

type Poll struct {
	ID         string       `json:"id"`
	Type       PollType     `json:"type"`
	Question   string       `json:"question"`
	Options    []PollOption `json:"options"`
	VoterCount int          `json:"total_voter_count"`

	// (Optional)
	Closed          bool            `json:"is_closed,omitempty"`
	CorrectOption   int             `json:"correct_option_id,omitempty"`
	MultipleAnswers bool            `json:"allows_multiple_answers,omitempty"`
	Explanation     string          `json:"explanation,omitempty"`
	ParseMode       ParseMode       `json:"explanation_parse_mode,omitempty"`
	Entities        []MessageEntity `json:"explanation_entities"`

	// True by default, shouldn't be omitted.
	Anonymous bool `json:"is_anonymous"`

	// (Mutually exclusive)
	OpenPeriod    int   `json:"open_period,omitempty"`
	CloseUnixdate int64 `json:"close_date,omitempty"`
}

type PollAnswer struct {
	PollID  string `json:"poll_id"`
	User    User   `json:"user"`
	Options []int  `json:"option_ids"`
}

type Poller interface {
	// Poll is supposed to take the bot object
	// subscription channel and start polling
	// for Updates immediately.
	//
	// Poller must listen for stop constantly and close
	// it as soon as it's done polling.
	Poll(b *Bot, updates chan Update, stop chan struct{})
}

type Bot struct {
	Me      *User
	Token   string
	URL     string
	Updates chan Update
	Poller  Poller
	// contains filtered or unexported fields
}

type File struct {
	FileID   string `json:"file_id"`
	UniqueID string `json:"file_unique_id"`
	FileSize int    `json:"file_size"`

	// file on telegram server https://core.telegram.org/bots/api#file
	FilePath string `json:"file_path"`

	// file on local file system.
	FileLocal string `json:"file_local"`

	// file on the internet
	FileURL string `json:"file_url"`

	// file backed with io.Reader
	FileReader io.Reader `json:"-"`
	// contains filtered or unexported fields
}

func (b *Bot) GetFile(file *File) (io.ReadCloser, error) {
	f, err := b.FileByID(file.FileID)
	if err != nil {
		return nil, err
	}

	url := b.URL + "/file/bot" + b.Token + "/" + f.FilePath
	file.FilePath = f.FilePath // saving file path

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, wrapError(err)
	}

	resp, err := b.client.Do(req)
	if err != nil {
		return nil, wrapError(err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, errors.Errorf("telebot: expected status 200 but got %s", resp.Status)
	}

	return resp.Body, nil
}

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:  "TOKEN_HERE",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "hello world")
	})

	b.Start()
}
