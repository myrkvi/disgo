package handlers

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/httpserver"
)


var allEventHandlers = []bot.GatewayEventHandler{


	bot.NewGatewayEventHandler(gateway.EventTypeApplicationCommandPermissionsUpdate, gatewayHandlerApplicationCommandPermissionsUpdate),

	bot.NewGatewayEventHandler(gateway.EventTypeAutoModerationRuleCreate, gatewayHandlerAutoModerationRuleCreate),
	bot.NewGatewayEventHandler(gateway.EventTypeAutoModerationRuleUpdate, gatewayHandlerAutoModerationRuleUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeAutoModerationRuleDelete, gatewayHandlerAutoModerationRuleDelete),
	bot.NewGatewayEventHandler(gateway.EventTypeAutoModerationActionExecution, gatewayHandlerAutoModerationActionExecution),

	bot.NewGatewayEventHandler(gateway.EventTypeChannelCreate, gatewayHandlerChannelCreate),
	bot.NewGatewayEventHandler(gateway.EventTypeChannelUpdate, gatewayHandlerChannelUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeChannelDelete, gatewayHandlerChannelDelete),
	bot.NewGatewayEventHandler(gateway.EventTypeChannelPinsUpdate, gatewayHandlerChannelPinsUpdate),

	bot.NewGatewayEventHandler(gateway.EventTypeGuildCreate, gatewayHandlerGuildCreate),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildUpdate, gatewayHandlerGuildUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildDelete, gatewayHandlerGuildDelete),

	bot.NewGatewayEventHandler(gateway.EventTypeGuildAuditLogEntryCreate, gatewayHandlerGuildAuditLogEntryCreate),

	bot.NewGatewayEventHandler(gateway.EventTypeGuildBanAdd, gatewayHandlerGuildBanAdd),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildBanRemove, gatewayHandlerGuildBanRemove),

	bot.NewGatewayEventHandler(gateway.EventTypeGuildEmojisUpdate, gatewayHandlerGuildEmojisUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildStickersUpdate, gatewayHandlerGuildStickersUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildIntegrationsUpdate, gatewayHandlerGuildIntegrationsUpdate),

	bot.NewGatewayEventHandler(gateway.EventTypeGuildMemberAdd, gatewayHandlerGuildMemberAdd),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildMemberRemove, gatewayHandlerGuildMemberRemove),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildMemberUpdate, gatewayHandlerGuildMemberUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildMembersChunk, gatewayHandlerGuildMembersChunk),

	bot.NewGatewayEventHandler(gateway.EventTypeGuildRoleCreate, gatewayHandlerGuildRoleCreate),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildRoleUpdate, gatewayHandlerGuildRoleUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildRoleDelete, gatewayHandlerGuildRoleDelete),

	bot.NewGatewayEventHandler(gateway.EventTypeGuildScheduledEventCreate, gatewayHandlerGuildScheduledEventCreate),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildScheduledEventUpdate, gatewayHandlerGuildScheduledEventUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildScheduledEventDelete, gatewayHandlerGuildScheduledEventDelete),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildScheduledEventUserAdd, gatewayHandlerGuildScheduledEventUserAdd),
	bot.NewGatewayEventHandler(gateway.EventTypeGuildScheduledEventUserRemove, gatewayHandlerGuildScheduledEventUserRemove),

	bot.NewGatewayEventHandler(gateway.EventTypeIntegrationCreate, gatewayHandlerIntegrationCreate),
	bot.NewGatewayEventHandler(gateway.EventTypeIntegrationUpdate, gatewayHandlerIntegrationUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeIntegrationDelete, gatewayHandlerIntegrationDelete),

	bot.NewGatewayEventHandler(gateway.EventTypeInteractionCreate, gatewayHandlerInteractionCreate),

	bot.NewGatewayEventHandler(gateway.EventTypeInviteCreate, gatewayHandlerInviteCreate),
	bot.NewGatewayEventHandler(gateway.EventTypeInviteDelete, gatewayHandlerInviteDelete),
	

	bot.NewGatewayEventHandler(gateway.EventTypePresenceUpdate, gatewayHandlerPresenceUpdate),

	bot.NewGatewayEventHandler(gateway.EventTypeStageInstanceCreate, gatewayHandlerStageInstanceCreate),
	bot.NewGatewayEventHandler(gateway.EventTypeStageInstanceUpdate, gatewayHandlerStageInstanceUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeStageInstanceDelete, gatewayHandlerStageInstanceDelete),

	bot.NewGatewayEventHandler(gateway.EventTypeTypingStart, gatewayHandlerTypingStart),
	bot.NewGatewayEventHandler(gateway.EventTypeUserUpdate, gatewayHandlerUserUpdate),

	bot.NewGatewayEventHandler(gateway.EventTypeVoiceStateUpdate, gatewayHandlerVoiceStateUpdate),
	bot.NewGatewayEventHandler(gateway.EventTypeVoiceServerUpdate, gatewayHandlerVoiceServerUpdate),

	bot.NewGatewayEventHandler(gateway.EventTypeWebhooksUpdate, gatewayHandlerWebhooksUpdate),
}
