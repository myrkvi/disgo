package core

import "github.com/DisgoOrg/disgo/discord"

type VoiceState struct {
	discord.VoiceState
	Disgo  Disgo
	Member *Member
}

// Mute returns if the Member is muted
func (s *VoiceState) Mute() bool {
	return s.GuildMute || s.SelfMute
}

// Deaf returns if the Member is deafened
func (s *VoiceState) Deaf() bool {
	return s.GuildDeaf || s.SelfDeaf
}

// Guild returns the Guild of this VoiceState from the Caches
func (s *VoiceState) Guild() *Guild {
	return s.Disgo.Caches().GuildCache().Get(s.GuildID)
}

// VoiceChannel returns the VoiceChannel of this VoiceState from the Caches
func (s *VoiceState) VoiceChannel() VoiceChannel {
	if s.ChannelID == nil {
		return nil
	}
	return s.Disgo.Caches().VoiceChannelCache().Get(*s.ChannelID)
}
