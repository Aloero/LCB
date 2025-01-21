package LCB

import (
	"sync"
)

type Bot struct {
	Token        string
	updatesChan  chan Update
	handlers     []Handler
	lastUpdateId int64
	state        map[int64]interface{}
	logs         bool
	Mu           sync.Mutex
}

func NewBot(token string, logs bool) *Bot {
	return &Bot{
		Token:        token,
		updatesChan:  make(chan Update),
		handlers:     []Handler{},
		lastUpdateId: 0,
		state:        make(map[int64]interface{}),
		logs:         logs,
		Mu:           sync.Mutex{},
	}
}

type Utils struct {
	Inline       *InlineKeyboardMarkup
	Reply        *ReplyKeyboardMarkup
	Delete       *DeleteKeyboard
	ReplyMessage *int64
	MessageThreadID *int64
}

type Handler struct {
	Filter   func(update Update) bool
	Callback func(update Update)
}

// ChatJoinRequest представляет собой запрос на вступление в чат.
type ChatJoinRequest struct {
	Chat *Chat `json:"chat"`
	From *User `json:"from"`
	Date int64 `json:"date"`
}

// ChatMember представляет собой информацию о члене чата.
type ChatMember struct {
	User            *User  `json:"user"`
	Status          string `json:"status"` // Например, "member", "administrator", "restricted", "left", "kicked"
	CanBeEdited     bool   `json:"can_be_edited"`
	CanPostMessages bool   `json:"can_post_messages"`
	// Добавьте другие поля при необходимости
}

// ShippingOption представляет собой вариант доставки.
type ShippingOption struct {
	ID     string         `json:"id"`
	Title  string         `json:"title"`
	Prices []LabeledPrice `json:"prices"`
}

// LabeledPrice представляет собой ценовую метку с описанием.
type LabeledPrice struct {
	Label  string `json:"label"`  // Описание товара или услуги
	Amount int64  `json:"amount"` // Сумма в копейках (например, 1000 означает 10.00)
}

// ShippingAddress представляет собой адрес доставки.
type ShippingAddress struct {
	CountryCode string `json:"country_code"`
	State       string `json:"state"`
	City        string `json:"city"`
	StreetLine  string `json:"street_line"`
	PostCode    string `json:"post_code"`
}

// OrderInfo представляет собой информацию о заказе.
type OrderInfo struct {
	Name            string           `json:"name,omitempty"`
	PhoneNumber     string           `json:"phone_number,omitempty"`
	Email           string           `json:"email,omitempty"`
	ShippingAddress *ShippingAddress `json:"shipping_address,omitempty"`
}

// MessageEntity представляет собой сущность сообщения (например, ссылки, хэштеги и т.д.).
type MessageEntity struct {
	Type   string `json:"type"` // Например, "mention", "hashtag", "bot_command", "url", "email", "phone_number", "bold", "italic", "underline", "strikethrough"
	Offset int64  `json:"offset"`
	Length int64  `json:"length"`
	URL    string `json:"url,omitempty"`
	User   *User
}

// Audio представляет собой аудиофайл.
type Audio struct {
	FileID    string     `json:"file_id"`
	Duration  int64      `json:"duration"`
	Performer string     `json:"performer,omitempty"`
	Title     string     `json:"title,omitempty"`
	Thumb     *PhotoSize `json:"thumb,omitempty"`
	MimeType  string     `json:"mime_type,omitempty"`
	FileSize  int64      `json:"file_size,omitempty"`
}

// Document представляет собой документ.
type Document struct {
	FileID   string     `json:"file_id"`
	Thumb    *PhotoSize `json:"thumb,omitempty"`
	FileName string     `json:"file_name,omitempty"`
	MimeType string     `json:"mime_type,omitempty"`
	FileSize int64      `json:"file_size,omitempty"`
}

// Animation представляет собой анимацию.
type Animation struct {
	FileID   string     `json:"file_id"`
	Width    int64      `json:"width"`
	Height   int64      `json:"height"`
	Duration int64      `json:"duration"`
	Thumb    *PhotoSize `json:"thumb,omitempty"`
	FileName string     `json:"file_name,omitempty"`
	MimeType string     `json:"mime_type,omitempty"`
	FileSize int64      `json:"file_size,omitempty"`
}

// Game представляет собой игровое сообщение.
type Game struct {
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Photo       []PhotoSize `json:"photo"`
	Text        string      `json:"text,omitempty"`
	Animation   *Animation  `json:"animation,omitempty"`
}

// PhotoSize представляет собой размер фото.
type PhotoSize struct {
	FileID   string `json:"file_id"`
	Width    int64  `json:"width"`
	Height   int64  `json:"height"`
	FileSize int64  `json:"file_size,omitempty"`
}

// Sticker представляет собой стикер.
type Sticker struct {
	FileID   string     `json:"file_id"`
	Width    int64      `json:"width"`
	Height   int64      `json:"height"`
	Thumb    *PhotoSize `json:"thumb,omitempty"`
	Emoji    string     `json:"emoji,omitempty"`
	FileSize int64      `json:"file_size,omitempty"`
}

// Video представляет собой видеозапись.
type Video struct {
	FileID   string     `json:"file_id"`
	Width    int64      `json:"width"`
	Height   int64      `json:"height"`
	Duration int64      `json:"duration"`
	Thumb    *PhotoSize `json:"thumb,omitempty"`
	MimeType string     `json:"mime_type,omitempty"`
	FileSize int64      `json:"file_size,omitempty"`
}

// Voice представляет собой голосовое сообщение.
type Voice struct {
	FileID   string `json:"file_id"`
	Duration int64  `json:"duration"`
	MimeType string `json:"mime_type,omitempty"`
	FileSize int64  `json:"file_size,omitempty"`
}

// VideoNote представляет собой видеозапись с заметкой.
type VideoNote struct {
	FileID   string     `json:"file_id"`
	Length   int64      `json:"length"`
	Duration int64      `json:"duration"`
	Thumb    *PhotoSize `json:"thumb,omitempty"`
	FileSize int64      `json:"file_size,omitempty"`
}

// Contact представляет собой контакт.
type Contact struct {
	PhoneNumber    string    `json:"phone_number"`              // Номер телефона
	FirstName      string    `json:"first_name"`                // Имя
	LastName       string    `json:"last_name,omitempty"`       // Фамилия (необязательно)
	UserID         int64     `json:"user_id,omitempty"`         // Идентификатор пользователя (необязательно)
	Username       string    `json:"username,omitempty"`        // Имя пользователя (необязательно)
	VCard          string    `json:"vcard,omitempty"`           // Визитная карточка (необязательно)
	Email          string    `json:"email,omitempty"`           // Электронная почта (необязательно)
	Address        string    `json:"address,omitempty"`         // Адрес (необязательно)
	CompanyName    string    `json:"company_name,omitempty"`    // Компания (необязательно)
	JobTitle       string    `json:"job_title,omitempty"`       // Должность (необязательно)
	BirthDate      string    `json:"birth_date,omitempty"`      // Дата рождения (необязательно)
	Website        string    `json:"website,omitempty"`         // Личный сайт (необязательно)
	Notes          string    `json:"notes,omitempty"`           // Дополнительные заметки (необязательно)
	ProfilePicture string    `json:"profile_picture,omitempty"` // Ссылка на изображение профиля (необязательно)
	IsFavorite     bool      `json:"is_favorite,omitempty"`     // Флаг избранного контакта (необязательно)
	IsBlocked      bool      `json:"is_blocked,omitempty"`      // Флаг заблокированного контакта (необязательно)
	CreatedAt      string    `json:"created_at,omitempty"`      // Дата создания контакта (необязательно)
	UpdatedAt      string    `json:"updated_at,omitempty"`      // Дата последнего обновления (необязательно)
	Location       *Location `json:"location,omitempty"`        // Координаты местоположения контакта (необязательно)
}

// Venue представляет собой место.
type Venue struct {
	Location       *Location `json:"location"`
	Title          string    `json:"title"`
	Address        string    `json:"address"`
	FoursquareID   string    `json:"foursquare_id,omitempty"`
	FoursquareType string    `json:"foursquare_type,omitempty"`
}

// Location представляет собой геолокацию.
type Location struct {
	Longitude float64 `json:"longitude"`
	Latitude  float64 `json:"latitude"`
}

type ResponsePostMessage struct {
	Ok     bool `json:"ok"`
	Result struct {
		MessageID int `json:"message_id"`
	} `json:"result"`
}

type Dice struct {
	Emoji string `json:"emoji"`
	Value int    `json:"value"`
}

type From struct {
	ID                      int64  `json:"id"`                                    // Уникальный идентификатор пользователя
	IsBot                   bool   `json:"is_bot"`                                // Является ли пользователь ботом
	FirstName               string `json:"first_name"`                            // Имя пользователя
	LastName                string `json:"last_name,omitempty"`                   // Фамилия пользователя
	Username                string `json:"username,omitempty"`                    // Имя пользователя (никнейм)
	LanguageCode            string `json:"language_code,omitempty"`               // Код языка пользователя
	CanJoinGroups           bool   `json:"can_join_groups,omitempty"`             // Может ли пользователь присоединяться к группам (для ботов)
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages,omitempty"` // Может ли бот читать все сообщения в группе
	SupportsInlineQueries   bool   `json:"supports_inline_queries,omitempty"`     // Поддерживает ли пользователь встроенные запросы (для ботов)
	IsPremium               bool   `json:"is_premium,omitempty"`                  // Является ли пользователь подписчиком премиум-сервиса
	AddedToAttachmentMenu   bool   `json:"added_to_attachment_menu,omitempty"`    // Был ли пользователь добавлен в меню вложений
}

type User struct {
	ID                      int64  `json:"id"`
	IsBot                   bool   `json:"is_bot"`
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name,omitempty"`
	Username                string `json:"username,omitempty"`
	LanguageCode            string `json:"language_code,omitempty"`
	CanJoinGroups           bool   `json:"can_join_groups,omitempty"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages,omitempty"`
	SupportsInlineQueries   bool   `json:"supports_inline_queries,omitempty"`
	IsPremium               bool   `json:"is_premium,omitempty"`
	AddedToAttachmentMenu   bool   `json:"added_to_attachment_menu,omitempty"`
	Bio                     string `json:"bio,omitempty"`
	PhoneNumber             string `json:"phone_number,omitempty"`
	Email                   string `json:"email,omitempty"`
	ProfilePhotoURL         string `json:"profile_photo_url,omitempty"`
}

type TelegramResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type FileResponse struct {
	Ok     bool  `json:"ok"`
	Result *File `json:"result"`
}

type File struct {
	FileID   string `json:"file_id"`
	FilePath string `json:"file_path"`
	FileSize int    `json:"file_size"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text         string      `json:"text"`
	URL          string      `json:"url"`
	CallbackData string      `json:"callback_data"`
	WebApp       *WebAppInfo `json:"web_app,omitempty"`
}

type ReplyKeyboardMarkup struct {
	ReplyKeyboard   [][]ReplyKeyboardButton `json:"keyboard"`
	ResizeKeyboard  bool                    `json:"resize_keyboard"`
	OneTimeKeyboard bool                    `json:"one_time_keyboard"`
}

type ReplyKeyboardButton struct {
	Text string `json:"text"`
}

type WebAppInfo struct {
	URL string `json:"url"`
}

type DeleteKeyboard struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
}

// Update представляет собой одно обновление, получаемое от Telegram.
type Update struct {
	UpdateID           int64               `json:"update_id"`
	Message            *Message            `json:"message,omitempty"`
	EditedMessage      *Message            `json:"edited_message,omitempty"`
	ChannelPost        *Message            `json:"channel_post,omitempty"`
	EditedChannelPost  *Message            `json:"edited_channel_post,omitempty"`
	InlineQuery        *InlineQuery        `json:"inline_query,omitempty"`
	ChosenInlineResult *ChosenInlineResult `json:"chosen_inline_result,omitempty"`
	CallbackQuery      *CallbackQuery      `json:"callback_query,omitempty"`
	ShippingQuery      *ShippingQuery      `json:"shipping_query,omitempty"`
	PreCheckoutQuery   *PreCheckoutQuery   `json:"pre_checkout_query,omitempty"`
	Poll               *Poll               `json:"poll,omitempty"`
	PollAnswer         *PollAnswer         `json:"poll_answer,omitempty"`
	MyChatMember       *ChatMemberUpdated  `json:"my_chat_member,omitempty"`
	ChatMember         *ChatMemberUpdated  `json:"chat_member,omitempty"`
	ChatJoinRequest    *ChatJoinRequest    `json:"chat_join_request,omitempty"`
}

// Message представляет собой сообщение, полученное от пользователя.
type Message struct {
	MessageID              int64               `json:"message_id"`
	From                   *User               `json:"from"`
	SenderChat             *Chat               `json:"sender_chat,omitempty"`
	Chat                   *Chat               `json:"chat"`
	Date                   int64               `json:"date"`
	ForwardFrom            *User               `json:"forward_from,omitempty"`
	ForwardFromChat        *Chat               `json:"forward_from_chat,omitempty"`
	ForwardDate            int64               `json:"forward_date,omitempty"`
	ReplyTo                *Message            `json:"reply_to_message,omitempty"`
	EditDate               int64               `json:"edit_date,omitempty"`
	MediaGroupID           string              `json:"media_group_id,omitempty"`
	AuthorSignature        string              `json:"author_signature,omitempty"`
	Text                   string              `json:"text,omitempty"`
	Entities               []MessageEntity     `json:"entities,omitempty"`
	CaptionEntities        []MessageEntity     `json:"caption_entities,omitempty"`
	Audio                  *Audio              `json:"audio,omitempty"`
	Document               *Document           `json:"document,omitempty"`
	Animation              *Animation          `json:"animation,omitempty"`
	Game                   *Game               `json:"game,omitempty"`
	Photo                  []PhotoSize         `json:"photo,omitempty"`
	Sticker                *Sticker            `json:"sticker,omitempty"`
	Video                  *Video              `json:"video,omitempty"`
	Voice                  *Voice              `json:"voice,omitempty"`
	VideoNote              *VideoNote          `json:"video_note,omitempty"`
	Contact                *Contact            `json:"contact,omitempty"`
	Venue                  *Venue              `json:"venue,omitempty"`
	Location               *Location           `json:"location,omitempty"`
	NewChatMembers         []*User             `json:"new_chat_members,omitempty"`
	LeftChatMember         *User               `json:"left_chat_member,omitempty"`
	NewChatTitle           string              `json:"new_chat_title,omitempty"`
	NewChatPhoto           []PhotoSize         `json:"new_chat_photo,omitempty"`
	DeleteChatPhoto        bool                `json:"delete_chat_photo,omitempty"`
	GroupChatCreated       bool                `json:"group_chat_created,omitempty"`
	SupergroupChatCreated  bool                `json:"supergroup_chat_created,omitempty"`
	ChannelChatCreated     bool                `json:"channel_chat_created,omitempty"`
	MessageAutoDeleteTimer int64               `json:"message_auto_delete_timer_changed,omitempty"`
	MigrateToChatID        int64               `json:"migrate_to_chat_id,omitempty"`
	MigrateFromChatID      int64               `json:"migrate_from_chat_id,omitempty"`
	PinnedMessage          *Message            `json:"pinned_message,omitempty"`
	LinkPreviewOptions     *LinkPreviewOptions `json:"link_preview_options,omitempty"`
	Dice                   *Dice               `json:"dice"`
}

type LinkPreviewOptions struct {
	IsDisabled bool `json:"is_disabled"`
}

// InlineQuery представляет собой запрос на inline-режим.
type InlineQuery struct {
	ID       string    `json:"id"`
	From     *User     `json:"from"`
	Location *Location `json:"location,omitempty"`
	Query    string    `json:"query"`
	Offset   string    `json:"offset"`
}

// ChosenInlineResult представляет собой результат, выбранный пользователем в inline-режиме.
type ChosenInlineResult struct {
	ResultID string    `json:"result_id"`
	From     *User     `json:"from"`
	Location *Location `json:"location,omitempty"`
	Query    string    `json:"query"`
}

// CallbackQuery представляет собой обратный вызов, полученный от нажатия на кнопку.
type CallbackQuery struct {
	ID      string   `json:"id"`
	From    *User    `json:"from"`
	Message *Message `json:"message,omitempty"`
	Chat    *Chat    `json:"chat,omitempty"`
	Data    string   `json:"data,omitempty"`
}

// ShippingQuery представляет собой запрос на доставку.
type ShippingQuery struct {
	ID              string            `json:"id"`
	From            *User             `json:"from"`
	Currency        string            `json:"currency"`
	TotalAmount     int64             `json:"total_amount"`
	InvoicePayload  string            `json:"invoice_payload"`
	ShippingOptions []*ShippingOption `json:"shipping_options,omitempty"`
	Address         *ShippingAddress  `json:"order_info,omitempty"`
}

// PreCheckoutQuery представляет собой запрос перед оформлением заказа.
type PreCheckoutQuery struct {
	ID               string     `json:"id"`
	From             *User      `json:"from"`
	Currency         string     `json:"currency"`
	TotalAmount      int64      `json:"total_amount"`
	InvoicePayload   string     `json:"invoice_payload"`
	ShippingOptionID string     `json:"shipping_option_id,omitempty"`
	OrderInfo        *OrderInfo `json:"order_info,omitempty"`
}

// Poll представляет собой опрос.
type Poll struct {
	ID              string       `json:"id"`
	Question        string       `json:"question"`
	Options         []PollOption `json:"options"`
	IsClosed        bool         `json:"is_closed"`
	TotalVoterCount int64        `json:"total_voter_count"`
}

// PollOption представляет собой вариант ответа в опросе.
type PollOption struct {
	OptionID   int64  `json:"option_id"`
	Text       string `json:"text"`
	VoterCount int64  `json:"voter_count"`
}

// PollAnswer представляет собой ответ на опрос.
type PollAnswer struct {
	PollID   string `json:"poll_id"`
	From     *User  `json:"from"`
	OptionID int64  `json:"option_id"`
}

// Chat представляет собой информацию о чате.
type Chat struct {
	ID                    int64            `json:"id"`                                 // Уникальный идентификатор чата
	Type                  string           `json:"type"`                               // Тип чата ("private", "group", "supergroup", "channel")
	Title                 string           `json:"title,omitempty"`                    // Название чата (для групп, каналов и супергрупп)
	Username              string           `json:"username,omitempty"`                 // Имя пользователя (если чат является каналом или супергруппой)
	FirstName             string           `json:"first_name,omitempty"`               // Имя (для приватных чатов)
	LastName              string           `json:"last_name,omitempty"`                // Фамилия (для приватных чатов)
	Photo                 *ChatPhoto       `json:"photo,omitempty"`                    // Фотография чата (ссылка на изображение)
	Bio                   string           `json:"bio,omitempty"`                      // Биография (например, для личных чатов)
	InviteLink            string           `json:"invite_link,omitempty"`              // Ссылка на приглашение в чат
	Status                string           `json:"status,omitempty"`                   // Статус (например, администратора или участника в группе)
	Description           string           `json:"description,omitempty"`              // Описание чата (для каналов и супергрупп)
	PinnedMessage         *Message         `json:"pinned_message,omitempty"`           // Закрепленное сообщение в чате
	Permissions           *ChatPermissions `json:"permissions,omitempty"`              // Разрешения для участников чата
	SlowModeDelay         int              `json:"slow_mode_delay,omitempty"`          // Интервал замедленного режима (в секундах)
	MessageAutoDeleteTime int              `json:"message_auto_delete_time,omitempty"` // Время автоудаления сообщений
	StickerSetName        string           `json:"sticker_set_name,omitempty"`         // Название набора стикеров
	CanSetStickerSet      bool             `json:"can_set_sticker_set,omitempty"`      // Может ли администратор установить стикеры
	LinkedChatID          int64            `json:"linked_chat_id,omitempty"`           // Идентификатор связанного чата (например, обсуждение канала)
	Location              *ChatLocation    `json:"location,omitempty"`                 // Местоположение для чатов, связанных с местами
	CreatedAt             string           `json:"created_at,omitempty"`               // Дата создания чата
	UpdatedAt             string           `json:"updated_at,omitempty"`               // Дата последнего обновления информации о чате
}

type ChatPermissions struct {
	CanSendMessages       bool `json:"can_send_messages,omitempty"`         // Разрешение на отправку сообщений
	CanSendMediaMessages  bool `json:"can_send_media_messages,omitempty"`   // Разрешение на отправку медиафайлов
	CanSendPolls          bool `json:"can_send_polls,omitempty"`            // Разрешение на создание опросов
	CanSendOtherMessages  bool `json:"can_send_other_messages,omitempty"`   // Разрешение на отправку других типов сообщений
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews,omitempty"` // Разрешение на добавление превью веб-страниц
	CanChangeInfo         bool `json:"can_change_info,omitempty"`           // Разрешение на изменение информации о группе
	CanInviteUsers        bool `json:"can_invite_users,omitempty"`          // Разрешение на приглашение участников
	CanPinMessages        bool `json:"can_pin_messages,omitempty"`          // Разрешение на закрепление сообщений
}

type ChatLocation struct {
	Location Location `json:"location"`          // Координаты местоположения
	Address  string   `json:"address,omitempty"` // Адрес места
}

// ChatPhoto представляет собой фото чата.
type ChatPhoto struct {
	SmallFileID string `json:"small_file_id"`
	BigFileID   string `json:"big_file_id"`
}

// ChatMemberUpdated представляет собой информацию о члене чата.
type ChatMemberUpdated struct {
	Chat          *Chat `json:"chat"`
	From          *User `json:"from"`
	Date          int64 `json:"date"`
	OldChatMember *Chat `json:"old_chat_member"` // Older status of the chat member
	NewChatMember *Chat `json:"new_chat_member"` // New status of the chat member
}
