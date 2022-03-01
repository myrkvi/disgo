package core

import (
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/rest"
	"github.com/DisgoOrg/snowflake"
)

type AuditLog struct {
	discord.AuditLog
	GuildScheduledEvents map[snowflake.Snowflake]*GuildScheduledEvent
	Integrations         map[snowflake.Snowflake]Integration
	Threads              map[snowflake.Snowflake]GuildThread
	Users                map[snowflake.Snowflake]*User
	Webhooks             map[snowflake.Snowflake]Webhook
	GuildID              snowflake.Snowflake
	FilterOptions        AuditLogFilterOptions
	Bot                  *Bot
}

func (l *AuditLog) Guild() *Guild {
	return l.Bot.Caches.Guilds().Get(l.GuildID)
}

// AuditLogFilterOptions fields used to filter audit-log retrieving
type AuditLogFilterOptions struct {
	UserID     snowflake.Snowflake
	ActionType discord.AuditLogEvent
	Before     snowflake.Snowflake
	Limit      int
}

// Before gets new AuditLog(s) from Discord before the last one
func (l *AuditLog) Before(opts ...rest.RequestOpt) (*AuditLog, error) {
	before := snowflake.Snowflake("")
	if len(l.Entries) > 0 {
		before = l.Entries[len(l.Entries)-1].ID
	}
	auditLog, err := l.Bot.RestServices.AuditLogService().GetAuditLog(l.GuildID, l.FilterOptions.UserID, l.FilterOptions.ActionType, before, l.FilterOptions.Limit, opts...)
	if err != nil {
		return nil, err
	}
	return l.Bot.EntityBuilder.CreateAuditLog(l.GuildID, *auditLog, l.FilterOptions, CacheStrategyNoWs), nil
}
