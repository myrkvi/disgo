package gateway

import (
	"io"
	"time"

	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"

	"github.com/disgoorg/disgo/discord"
)

type Event interface {
	EventType() EventType
}

// EventUnknown is an event that is not known to disgo
type EventUnknown struct {
	Data json.RawMessage
	T    EventType
}

func (e EventUnknown) MarshalJSON() ([]byte, error) {
	return e.Data.MarshalJSON()
}

func (e *EventUnknown) UnmarshalJSON(data []byte) error {
	return e.Data.UnmarshalJSON(data)
}

func (e EventUnknown) EventType() EventType {
	return e.T
}

// EventReady is the event sent by discord when you successfully Identify
type EventReady struct {
	Version          int                        `json:"v"`
	User             discord.OAuth2User         `json:"user"`
	Guilds           []discord.UnavailableGuild `json:"guilds"`
	SessionID        string                     `json:"session_id"`
	ResumeGatewayURL string                     `json:"resume_gateway_url"`
	Shard            [2]int                     `json:"shard,omitempty"`
	Application      discord.PartialApplication `json:"application"`
}

func (EventReady) EventType() EventType { return EventTypeReady }

type EventApplicationCommandPermissionsUpdate struct {
	discord.ApplicationCommandPermissions
}

func (EventApplicationCommandPermissionsUpdate) EventType() EventType {
	return EventTypeApplicationCommandPermissionsUpdate
}

type EventChannelCreate struct {
	discord.GuildChannel
}

func (e *EventChannelCreate) UnmarshalJSON(data []byte) error {
	var v discord.UnmarshalChannel
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	e.GuildChannel = v.Channel.(discord.GuildChannel)
	return nil
}

func (EventChannelCreate) EventType() EventType { return EventTypeChannelCreate }

type EventChannelUpdate struct {
	discord.GuildChannel
	OldGuildChannel discord.GuildChannel `json:"-"`
}

func (e *EventChannelUpdate) UnmarshalJSON(data []byte) error {
	var v discord.UnmarshalChannel
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	e.GuildChannel = v.Channel.(discord.GuildChannel)
	return nil
}

func (EventChannelUpdate) EventType() EventType { return EventTypeChannelUpdate }

type EventChannelDelete struct {
	discord.GuildChannel
}

func (e *EventChannelDelete) UnmarshalJSON(data []byte) error {
	var v discord.UnmarshalChannel
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	e.GuildChannel = v.Channel.(discord.GuildChannel)
	return nil
}

func (EventChannelDelete) EventType() EventType { return EventTypeChannelDelete }

type EventThreadCreate struct {
	discord.GuildThread
	ThreadMember discord.ThreadMember `json:"thread_member"`
}

func (EventThreadCreate) EventType() EventType { return EventTypeThreadCreate }

type EventThreadUpdate struct {
	discord.GuildThread
	OldGuildThread discord.GuildThread `json:"-"`
}

func (EventThreadUpdate) EventType() EventType { return EventTypeThreadUpdate }

type EventThreadDelete struct {
	ID             snowflake.ID        `json:"id"`
	GuildID        snowflake.ID        `json:"guild_id"`
	ParentID       snowflake.ID        `json:"parent_id"`
	Type           discord.ChannelType `json:"type"`
	OldGuildThread discord.GuildThread `json:"-"`
}

func (EventThreadDelete) EventType() EventType { return EventTypeThreadDelete }

type EventThreadListSync struct {
	GuildID    snowflake.ID           `json:"guild_id"`
	ChannelIDs []snowflake.ID         `json:"channel_ids"`
	Threads    []discord.GuildThread  `json:"threads"`
	Members    []discord.ThreadMember `json:"members"`
}

func (EventThreadListSync) EventType() EventType { return EventTypeThreadListSync }

type EventThreadMemberUpdate struct {
	discord.ThreadMember
}

func (EventThreadMemberUpdate) EventType() EventType { return EventTypeThreadMemberUpdate }

type AddedThreadMember struct {
	discord.ThreadMember
	Member   discord.Member    `json:"member"`
	Presence *discord.Presence `json:"presence"`
}

type EventThreadMembersUpdate struct {
	ID               snowflake.ID           `json:"id"`
	GuildID          snowflake.ID           `json:"guild_id"`
	MemberCount      int                    `json:"member_count"`
	AddedMembers     []AddedThreadMember    `json:"added_members"`
	RemovedMemberIDs []snowflake.ID         `json:"removed_member_ids"`
	RemovedMembers   []discord.ThreadMember `json:"-"`
}

func (EventThreadMembersUpdate) EventType() EventType { return EventTypeThreadMembersUpdate }

type EventGuildCreate struct {
	discord.GatewayGuild
}

func (EventGuildCreate) EventType() EventType { return EventTypeGuildCreate }

type EventGuildUpdate struct {
	discord.GatewayGuild
	OldGuild discord.Guild `json:"-"`
}

func (EventGuildUpdate) EventType() EventType { return EventTypeGuildUpdate }

type EventGuildDelete struct {
	discord.GatewayGuild
}

func (EventGuildDelete) EventType() EventType { return EventTypeGuildDelete }

type EventGuildAuditLogEntryCreate struct {
	discord.AuditLogEntry
	GuildID snowflake.ID `json:"guild_id"`
}

func (EventGuildAuditLogEntryCreate) EventType() EventType { return EventTypeGuildAuditLogEntryCreate }

type EventMessageReactionAdd struct {
	UserID          snowflake.ID         `json:"user_id"`
	ChannelID       snowflake.ID         `json:"channel_id"`
	MessageID       snowflake.ID         `json:"message_id"`
	GuildID         *snowflake.ID        `json:"guild_id"`
	Member          *discord.Member      `json:"member"`
	Emoji           discord.PartialEmoji `json:"emoji"`
	MessageAuthorID *snowflake.ID        `json:"message_author_id"`
}

func (e *EventMessageReactionAdd) UnmarshalJSON(data []byte) error {
	type eventMessageReactionAdd EventMessageReactionAdd
	var v eventMessageReactionAdd
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*e = EventMessageReactionAdd(v)
	if e.Member != nil && e.GuildID != nil {
		e.Member.GuildID = *e.GuildID
	}
	return nil
}

func (EventMessageReactionAdd) EventType() EventType { return EventTypeMessageReactionAdd }

type EventMessageReactionRemove struct {
	UserID    snowflake.ID         `json:"user_id"`
	ChannelID snowflake.ID         `json:"channel_id"`
	MessageID snowflake.ID         `json:"message_id"`
	GuildID   *snowflake.ID        `json:"guild_id"`
	Emoji     discord.PartialEmoji `json:"emoji"`
}

func (EventMessageReactionRemove) EventType() EventType { return EventTypeMessageReactionRemove }

type EventMessageReactionRemoveEmoji struct {
	ChannelID snowflake.ID         `json:"channel_id"`
	MessageID snowflake.ID         `json:"message_id"`
	GuildID   *snowflake.ID        `json:"guild_id"`
	Emoji     discord.PartialEmoji `json:"emoji"`
}

func (EventMessageReactionRemoveEmoji) EventType() EventType {
	return EventTypeMessageReactionRemoveEmoji
}

type EventMessageReactionRemoveAll struct {
	ChannelID snowflake.ID  `json:"channel_id"`
	MessageID snowflake.ID  `json:"message_id"`
	GuildID   *snowflake.ID `json:"guild_id"`
}

func (EventMessageReactionRemoveAll) EventType() EventType { return EventTypeMessageReactionRemoveAll }

type EventChannelPinsUpdate struct {
	GuildID             *snowflake.ID `json:"guild_id"`
	ChannelID           snowflake.ID  `json:"channel_id"`
	LastPinTimestamp    *time.Time    `json:"last_pin_timestamp"`
	OldLastPinTimestamp *time.Time    `json:"-"`
}

func (EventChannelPinsUpdate) EventType() EventType { return EventTypeChannelPinsUpdate }

type EventGuildMembersChunk struct {
	GuildID    snowflake.ID       `json:"guild_id"`
	Members    []discord.Member   `json:"members"`
	ChunkIndex int                `json:"chunk_index"`
	ChunkCount int                `json:"chunk_count"`
	NotFound   []snowflake.ID     `json:"not_found"`
	Presences  []discord.Presence `json:"presences"`
	Nonce      string             `json:"nonce"`
}

func (EventGuildMembersChunk) EventType() EventType { return EventTypeGuildMembersChunk }

type EventGuildBanAdd struct {
	GuildID snowflake.ID `json:"guild_id"`
	User    discord.User `json:"user"`
}

func (EventGuildBanAdd) EventType() EventType { return EventTypeGuildBanAdd }

type EventGuildBanRemove struct {
	GuildID snowflake.ID `json:"guild_id"`
	User    discord.User `json:"user"`
}

func (EventGuildBanRemove) EventType() EventType { return EventTypeGuildBanRemove }

type EventGuildEmojisUpdate struct {
	GuildID snowflake.ID    `json:"guild_id"`
	Emojis  []discord.Emoji `json:"emojis"`
}

func (e *EventGuildEmojisUpdate) UnmarshalJSON(data []byte) error {
	type eventGuildEmojisUpdate EventGuildEmojisUpdate
	var v eventGuildEmojisUpdate
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*e = EventGuildEmojisUpdate(v)
	for i := range e.Emojis {
		e.Emojis[i].GuildID = e.GuildID
	}
	return nil
}

func (EventGuildEmojisUpdate) EventType() EventType { return EventTypeGuildEmojisUpdate }

type EventGuildStickersUpdate struct {
	GuildID  snowflake.ID      `json:"guild_id"`
	Stickers []discord.Sticker `json:"stickers"`
}

func (EventGuildStickersUpdate) EventType() EventType { return EventTypeGuildStickersUpdate }

type EventGuildIntegrationsUpdate struct {
	GuildID snowflake.ID `json:"guild_id"`
}

func (EventGuildIntegrationsUpdate) EventType() EventType { return EventTypeGuildIntegrationsUpdate }

type EventGuildMemberAdd struct {
	discord.Member
}

func (EventGuildMemberAdd) EventType() EventType { return EventTypeGuildMemberAdd }

type EventGuildMemberUpdate struct {
	discord.Member
	OldMember discord.Member `json:"-"`
}

func (EventGuildMemberUpdate) EventType() EventType { return EventTypeGuildMemberUpdate }

type EventGuildMemberRemove struct {
	GuildID snowflake.ID `json:"guild_id"`
	User    discord.User `json:"user"`
}

func (EventGuildMemberRemove) EventType() EventType { return EventTypeGuildMemberRemove }

type EventGuildRoleCreate struct {
	GuildID snowflake.ID `json:"guild_id"`
	Role    discord.Role `json:"role"`
}

func (e *EventGuildRoleCreate) UnmarshalJSON(data []byte) error {
	type eventGuildRoleCreate EventGuildRoleCreate
	var v eventGuildRoleCreate
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*e = EventGuildRoleCreate(v)
	e.Role.GuildID = e.GuildID
	return nil
}

func (e *EventGuildRoleCreate) MarshalJSON() ([]byte, error) {
	type eventGuildRoleCreate EventGuildRoleCreate
	e.GuildID = e.Role.GuildID
	return json.Marshal(eventGuildRoleCreate(*e))
}

func (EventGuildRoleCreate) EventType() EventType { return EventTypeGuildRoleCreate }

type EventGuildRoleDelete struct {
	GuildID snowflake.ID `json:"guild_id"`
	RoleID  snowflake.ID `json:"role_id"`
	Role    discord.Role `json:"-"`
}

func (EventGuildRoleDelete) EventType() EventType { return EventTypeGuildRoleDelete }

type EventGuildRoleUpdate struct {
	GuildID snowflake.ID `json:"guild_id"`
	Role    discord.Role `json:"role"`
	OldRole discord.Role `json:"-"`
}

func (e *EventGuildRoleUpdate) UnmarshalJSON(data []byte) error {
	type eventGuildRoleUpdate EventGuildRoleUpdate
	var v eventGuildRoleUpdate
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*e = EventGuildRoleUpdate(v)
	e.Role.GuildID = e.GuildID
	return nil
}

func (e *EventGuildRoleUpdate) MarshalJSON() ([]byte, error) {
	type eventGuildRoleUpdate EventGuildRoleUpdate
	e.GuildID = e.Role.GuildID
	return json.Marshal(eventGuildRoleUpdate(*e))
}

func (EventGuildRoleUpdate) EventType() EventType { return EventTypeGuildRoleUpdate }

type EventGuildScheduledEventCreate struct {
	discord.GuildScheduledEvent
}

func (EventGuildScheduledEventCreate) EventType() EventType {
	return EventTypeGuildScheduledEventCreate
}

type EventGuildScheduledEventUpdate struct {
	discord.GuildScheduledEvent
	OldGuildScheduledEvent discord.GuildScheduledEvent `json:"-"`
}

func (EventGuildScheduledEventUpdate) EventType() EventType {
	return EventTypeGuildScheduledEventUpdate
}

type EventGuildScheduledEventDelete struct {
	discord.GuildScheduledEvent
}

func (EventGuildScheduledEventDelete) EventType() EventType {
	return EventTypeGuildScheduledEventDelete
}

type EventGuildScheduledEventUserAdd struct {
	GuildScheduledEventID snowflake.ID `json:"guild_scheduled_event_id"`
	UserID                snowflake.ID `json:"user_id"`
	GuildID               snowflake.ID `json:"guild_id"`
}

func (EventGuildScheduledEventUserAdd) EventType() EventType {
	return EventTypeGuildScheduledEventUserAdd
}

type EventGuildScheduledEventUserRemove struct {
	GuildScheduledEventID snowflake.ID `json:"guild_scheduled_event_id"`
	UserID                snowflake.ID `json:"user_id"`
	GuildID               snowflake.ID `json:"guild_id"`
}

func (EventGuildScheduledEventUserRemove) EventType() EventType {
	return EventTypeGuildScheduledEventUserRemove
}

// RespondFunc is used to respond to Discord's Outgoing Webhooks
type RespondFunc func(responseType discord.InteractionResponseType, data discord.InteractionResponseData, opt ...rest.RequestOpt) error

type EventInteractionCreate struct {
	discord.Interaction
	Respond RespondFunc `json:"-"`
}

func (e *EventInteractionCreate) UnmarshalJSON(data []byte) error {
	interaction, err := discord.UnmarshalInteraction(data)
	if err != nil {
		return err
	}
	e.Interaction = interaction
	return nil
}

func (e EventInteractionCreate) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.Interaction)
}

func (EventInteractionCreate) EventType() EventType { return EventTypeInteractionCreate }

type EventInviteCreate struct {
	discord.Invite
}

func (EventInviteCreate) EventType() EventType { return EventTypeInviteCreate }

type EventInviteDelete struct {
	ChannelID snowflake.ID  `json:"channel_id"`
	GuildID   *snowflake.ID `json:"guild_id"`
	Code      string        `json:"code"`
}

func (EventInviteDelete) EventType() EventType { return EventTypeInviteDelete }

type EventMessageCreate struct {
	discord.Message
}

func (EventMessageCreate) EventType() EventType { return EventTypeMessageCreate }

type EventMessageUpdate struct {
	discord.Message
	OldMessage discord.Message `json:"-"`
}

func (EventMessageUpdate) EventType() EventType { return EventTypeMessageUpdate }

type EventMessageDelete struct {
	ID         snowflake.ID    `json:"id"`
	ChannelID  snowflake.ID    `json:"channel_id"`
	GuildID    *snowflake.ID   `json:"guild_id,omitempty"`
	OldMessage discord.Message `json:"-"`
}

func (EventMessageDelete) EventType() EventType { return EventTypeMessageDelete }

type EventMessageDeleteBulk struct {
	IDs         []snowflake.ID    `json:"id"`
	ChannelID   snowflake.ID      `json:"channel_id"`
	GuildID     *snowflake.ID     `json:"guild_id,omitempty"`
	OldMessages []discord.Message `json:"-"`
}

func (EventMessageDeleteBulk) EventType() EventType { return EventTypeMessageDeleteBulk }

type EventPresenceUpdate struct {
	discord.Presence
	OldPresence discord.Presence `json:"-"`
}

func (EventPresenceUpdate) EventType() EventType { return EventTypePresenceUpdate }

type EventStageInstanceCreate struct {
	discord.StageInstance
}

func (EventStageInstanceCreate) EventType() EventType { return EventTypeStageInstanceCreate }

type EventStageInstanceUpdate struct {
	discord.StageInstance
	OldStageInstance discord.StageInstance `json:"-"`
}

func (EventStageInstanceUpdate) EventType() EventType { return EventTypeStageInstanceUpdate }

type EventStageInstanceDelete struct {
	discord.StageInstance
}

func (EventStageInstanceDelete) EventType() EventType { return EventTypeStageInstanceDelete }

type EventTypingStart struct {
	ChannelID snowflake.ID    `json:"channel_id"`
	GuildID   *snowflake.ID   `json:"guild_id,omitempty"`
	UserID    snowflake.ID    `json:"user_id"`
	Timestamp time.Time       `json:"timestamp"`
	Member    *discord.Member `json:"member,omitempty"`
	User      discord.User    `json:"user"`
}

func (e *EventTypingStart) UnmarshalJSON(data []byte) error {
	type typingStartEvent EventTypingStart
	var v struct {
		Timestamp int64 `json:"timestamp"`
		typingStartEvent
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*e = EventTypingStart(v.typingStartEvent)
	e.Timestamp = time.Unix(v.Timestamp, 0)
	return nil
}

func (EventTypingStart) EventType() EventType { return EventTypeTypingStart }

type EventUserUpdate struct {
	discord.OAuth2User
	OldUser discord.OAuth2User `json:"-"`
}

func (EventUserUpdate) EventType() EventType { return EventTypeUserUpdate }

type EventVoiceStateUpdate struct {
	discord.VoiceState
	Member        discord.Member     `json:"member"`
	OldVoiceState discord.VoiceState `json:"-"`
}

func (EventVoiceStateUpdate) EventType() EventType { return EventTypeVoiceStateUpdate }

type EventVoiceServerUpdate struct {
	Token    string       `json:"token"`
	GuildID  snowflake.ID `json:"guild_id"`
	Endpoint *string      `json:"endpoint"`
}

func (EventVoiceServerUpdate) EventType() EventType { return EventTypeVoiceServerUpdate }

type EventWebhooksUpdate struct {
	GuildID   snowflake.ID `json:"guild_id"`
	ChannelID snowflake.ID `json:"channel_id"`
}

func (EventWebhooksUpdate) EventType() EventType { return EventTypeWebhooksUpdate }

type EventIntegrationCreate struct {
	discord.Integration
	GuildID snowflake.ID `json:"guild_id"`
}

func (e *EventIntegrationCreate) UnmarshalJSON(data []byte) error {
	type integrationCreateEvent EventIntegrationCreate
	var v struct {
		discord.UnmarshalIntegration
		integrationCreateEvent
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*e = EventIntegrationCreate(v.integrationCreateEvent)

	e.Integration = v.UnmarshalIntegration.Integration
	return nil
}

func (EventIntegrationCreate) EventType() EventType { return EventTypeIntegrationCreate }

type EventIntegrationUpdate struct {
	discord.Integration
	GuildID snowflake.ID `json:"guild_id"`
}

func (e *EventIntegrationUpdate) UnmarshalJSON(data []byte) error {
	type integrationUpdateEvent EventIntegrationUpdate
	var v struct {
		discord.UnmarshalIntegration
		integrationUpdateEvent
	}

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	*e = EventIntegrationUpdate(v.integrationUpdateEvent)

	e.Integration = v.UnmarshalIntegration.Integration
	return nil
}

func (EventIntegrationUpdate) EventType() EventType { return EventTypeIntegrationUpdate }

type EventIntegrationDelete struct {
	ID            snowflake.ID  `json:"id"`
	GuildID       snowflake.ID  `json:"guild_id"`
	ApplicationID *snowflake.ID `json:"application_id"`
}

func (EventIntegrationDelete) EventType() EventType { return EventTypeIntegrationDelete }

type EventAutoModerationRuleCreate struct {
	discord.AutoModerationRule
}

func (EventAutoModerationRuleCreate) EventType() EventType { return EventTypeAutoModerationRuleCreate }

type EventAutoModerationRuleUpdate struct {
	discord.AutoModerationRule
}

func (EventAutoModerationRuleUpdate) EventType() EventType { return EventTypeAutoModerationRuleUpdate }

type EventAutoModerationRuleDelete struct {
	discord.AutoModerationRule
}

func (EventAutoModerationRuleDelete) EventType() EventType { return EventTypeAutoModerationRuleDelete }

type EventAutoModerationActionExecution struct {
	GuildID              snowflake.ID                      `json:"guild_id"`
	Action               discord.AutoModerationAction      `json:"action"`
	RuleID               snowflake.ID                      `json:"rule_id"`
	RuleTriggerType      discord.AutoModerationTriggerType `json:"rule_trigger_type"`
	UserID               snowflake.ID                      `json:"user_id"`
	ChannelID            *snowflake.ID                     `json:"channel_id,omitempty"`
	MessageID            *snowflake.ID                     `json:"message_id,omitempty"`
	AlertSystemMessageID snowflake.ID                      `json:"alert_system_message_id"`
	Content              string                            `json:"content"`
	MatchedKeywords      *string                           `json:"matched_keywords"`
	MatchedContent       *string                           `json:"matched_content"`
}

func (EventAutoModerationActionExecution) EventType() EventType {
	return EventTypeAutoModerationActionExecution
}

type EventRaw struct {
	T       EventType
	Payload io.Reader
}

func (EventRaw) EventType() EventType { return EventTypeRaw }

type EventHeartbeatAck struct {
	LastHeartbeat time.Time
	NewHeartbeat  time.Time
}

func (EventHeartbeatAck) EventType() EventType { return EventTypeHeartbeatAck }
