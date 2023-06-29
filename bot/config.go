package bot

import (
	"fmt"

	"github.com/disgoorg/log"

	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgo/httpserver"
	"github.com/disgoorg/disgo/internal/tokenhelper"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/sharding"
	"github.com/disgoorg/disgo/voice"
)

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		Logger:               log.Default(),
		MemberChunkingFilter: MemberChunkingFilterNone,
	}
}

// Config lets you configure your Client instance.
type Config struct {
	Logger log.Logger

	RestClient           rest.Client
	RestClientConfigOpts []rest.ConfigOpt
	Rest                 rest.Rest

	EventManager           EventManager
	EventManagerConfigOpts []EventManagerConfigOpt

	VoiceManager           voice.Manager
	VoiceManagerConfigOpts []voice.ManagerConfigOpt

	Gateway           gateway.Gateway
	GatewayConfigOpts []gateway.ConfigOpt

	ShardManager           sharding.Manager
	ShardManagerConfigOpts []sharding.ConfigOpt

	HTTPServer           httpserver.Server
	PublicKey            string
	HTTPServerConfigOpts []httpserver.ConfigOpt

	Caches          cache.Caches
	CacheConfigOpts []cache.ConfigOpt

	MemberChunkingManager MemberChunkingManager
	MemberChunkingFilter  MemberChunkingFilter
}

// ConfigOpt is a type alias for a function that takes a Config and is used to configure your Client.
type ConfigOpt func(config *Config)

// Apply applies the given ConfigOpt(s) to the Config
func (c *Config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// WithLogger lets you inject your own logger implementing log.Logger.
func WithLogger(logger log.Logger) ConfigOpt {
	return func(config *Config) {
		config.Logger = logger
	}
}

// WithRestClient lets you inject your own rest.Client.
func WithRestClient(restClient rest.Client) ConfigOpt {
	return func(config *Config) {
		config.RestClient = restClient
	}
}

// WithRestClientConfigOpts let's you configure the default rest.Client.
func WithRestClientConfigOpts(opts ...rest.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.RestClientConfigOpts = append(config.RestClientConfigOpts, opts...)
	}
}

// WithRest lets you inject your own rest.Rest.
func WithRest(rest rest.Rest) ConfigOpt {
	return func(config *Config) {
		config.Rest = rest
	}
}

// WithEventManager lets you inject your own EventManager.
func WithEventManager(eventManager EventManager) ConfigOpt {
	return func(config *Config) {
		config.EventManager = eventManager
	}
}

// WithEventManagerConfigOpts lets you configure the default EventManager.
func WithEventManagerConfigOpts(opts ...EventManagerConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.EventManagerConfigOpts = append(config.EventManagerConfigOpts, opts...)
	}
}

// WithEventListeners adds the given EventListener(s) to the default EventManager.
func WithEventListeners(eventListeners ...EventListener) ConfigOpt {
	return func(config *Config) {
		config.EventManagerConfigOpts = append(config.EventManagerConfigOpts, WithListeners(eventListeners...))
	}
}

// WithEventListenerFunc adds the given func(c *Client, e E) to the default EventManager.
func WithEventListenerFunc[E gateway.Event](f func(c *Client, e E)) ConfigOpt {
	return WithEventListeners(NewListenerFunc(f))
}

// WithEventListenerChan adds the given chan<- E to the default EventManager.
func WithEventListenerChan[E gateway.Event](c chan<- E) ConfigOpt {
	return WithEventListeners(NewListenerChan(c))
}

// WithGateway lets you inject your own gateway.Gateway.
func WithGateway(gateway gateway.Gateway) ConfigOpt {
	return func(config *Config) {
		config.Gateway = gateway
	}
}

// WithDefaultGateway creates a gateway.Gateway with sensible defaults.
func WithDefaultGateway() ConfigOpt {
	return func(config *Config) {
		config.GatewayConfigOpts = append(config.GatewayConfigOpts, func(_ *gateway.Config) {})
	}
}

// WithGatewayConfigOpts lets you configure the default gateway.Gateway.
func WithGatewayConfigOpts(opts ...gateway.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.GatewayConfigOpts = append(config.GatewayConfigOpts, opts...)
	}
}

// WithShardManager lets you inject your own sharding.Manager.
func WithShardManager(shardManager sharding.Manager) ConfigOpt {
	return func(config *Config) {
		config.ShardManager = shardManager
	}
}

// WithDefaultShardManager creates a sharding.Manager with sensible defaults.
func WithDefaultShardManager() ConfigOpt {
	return func(config *Config) {
		config.ShardManagerConfigOpts = append(config.ShardManagerConfigOpts, func(_ *sharding.Config) {})
	}
}

// WithShardManagerConfigOpts lets you configure the default sharding.Manager.
func WithShardManagerConfigOpts(opts ...sharding.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.ShardManagerConfigOpts = append(config.ShardManagerConfigOpts, opts...)
	}
}

// WithHTTPServer lets you inject your own httpserver.Server.
func WithHTTPServer(httpServer httpserver.Server) ConfigOpt {
	return func(config *Config) {
		config.HTTPServer = httpServer
	}
}

// WithHTTPServerConfigOpts lets you configure the default httpserver.Server.
func WithHTTPServerConfigOpts(publicKey string, opts ...httpserver.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.PublicKey = publicKey
		config.HTTPServerConfigOpts = append(config.HTTPServerConfigOpts, opts...)
	}
}

// WithCaches lets you inject your own cache.Caches.
func WithCaches(caches cache.Caches) ConfigOpt {
	return func(config *Config) {
		config.Caches = caches
	}
}

// WithCacheConfigOpts lets you configure the default cache.Caches.
func WithCacheConfigOpts(opts ...cache.ConfigOpt) ConfigOpt {
	return func(config *Config) {
		config.CacheConfigOpts = append(config.CacheConfigOpts, opts...)
	}
}

// WithMemberChunkingManager lets you inject your own MemberChunkingManager.
func WithMemberChunkingManager(memberChunkingManager MemberChunkingManager) ConfigOpt {
	return func(config *Config) {
		config.MemberChunkingManager = memberChunkingManager
	}
}

// WithMemberChunkingFilter lets you configure the default MemberChunkingFilter.
func WithMemberChunkingFilter(memberChunkingFilter MemberChunkingFilter) ConfigOpt {
	return func(config *Config) {
		config.MemberChunkingFilter = memberChunkingFilter
	}
}

func buildClient(token string, config Config, os string, name string, github string, version string) (*Client, error) {
	if token == "" {
		return nil, discord.ErrNoBotToken
	}
	id, err := tokenhelper.IDFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("error while getting application id from token: %w", err)
	}
	client := &Client{
		Token:         token,
		Logger:        config.Logger,
		ApplicationID: *id,
	}

	if config.RestClient == nil {
		// prepend standard user-agent. this can be overridden as it's appended to the front of the slice
		config.RestClientConfigOpts = append([]rest.ConfigOpt{
			rest.WithUserAgent(fmt.Sprintf("DiscordBot (%s, %s)", github, version)),
			rest.WithLogger(client.Logger),
			func(config *rest.Config) {
				config.RateRateLimiterConfigOpts = append([]rest.RateLimiterConfigOpt{rest.WithRateLimiterLogger(client.Logger)}, config.RateRateLimiterConfigOpts...)
			},
		}, config.RestClientConfigOpts...)

		config.RestClient = rest.NewClient(client.Token, config.RestClientConfigOpts...)
	}

	if config.Rest == nil {
		config.Rest = rest.New(config.RestClient)
	}
	client.Rest = config.Rest

	if config.VoiceManager == nil {
		config.VoiceManager = voice.NewManager(client.UpdateVoiceState, *id, append([]voice.ManagerConfigOpt{voice.WithLogger(client.Logger)}, config.VoiceManagerConfigOpts...)...)
	}
	client.VoiceManager = config.VoiceManager

	if config.EventManager == nil {
		config.EventManager = NewEventManager(client, config.EventManagerConfigOpts...)
	}
	client.EventManager = config.EventManager

	if config.Gateway == nil && len(config.GatewayConfigOpts) > 0 {
		var gatewayRs *discord.Gateway
		gatewayRs, err = client.Rest.GetGateway()
		if err != nil {
			return nil, err
		}

		config.GatewayConfigOpts = append([]gateway.ConfigOpt{
			gateway.WithURL(gatewayRs.URL),
			gateway.WithLogger(client.Logger),
			gateway.WithOS(os),
			gateway.WithBrowser(name),
			gateway.WithDevice(name),
			func(config *gateway.Config) {
				config.RateRateLimiterConfigOpts = append([]gateway.RateLimiterConfigOpt{gateway.WithRateLimiterLogger(client.Logger)}, config.RateRateLimiterConfigOpts...)
			},
		}, config.GatewayConfigOpts...)

		config.Gateway = gateway.New(token, client.HandleEvent, nil, config.GatewayConfigOpts...)
	}
	client.Gateway = config.Gateway

	if config.ShardManager == nil && len(config.ShardManagerConfigOpts) > 0 {
		var gatewayBotRs *discord.GatewayBot
		gatewayBotRs, err = client.Rest.GetGatewayBot()
		if err != nil {
			return nil, err
		}

		shardIDs := make([]int, gatewayBotRs.Shards)
		for i := 0; i < gatewayBotRs.Shards-1; i++ {
			shardIDs[i] = i
		}

		config.ShardManagerConfigOpts = append([]sharding.ConfigOpt{
			sharding.WithShardCount(gatewayBotRs.Shards),
			sharding.WithShardIDs(shardIDs...),
			sharding.WithGatewayConfigOpts(
				gateway.WithURL(gatewayBotRs.URL),
				gateway.WithLogger(client.Logger),
				gateway.WithOS(os),
				gateway.WithBrowser(name),
				gateway.WithDevice(name),
				func(config *gateway.Config) {
					config.RateRateLimiterConfigOpts = append([]gateway.RateLimiterConfigOpt{gateway.WithRateLimiterLogger(client.Logger)}, config.RateRateLimiterConfigOpts...)
				},
			),
			sharding.WithLogger(client.Logger),
			func(config *sharding.Config) {
				config.RateRateLimiterConfigOpts = append([]sharding.RateLimiterConfigOpt{sharding.WithRateLimiterLogger(client.Logger), sharding.WithMaxConcurrency(gatewayBotRs.SessionStartLimit.MaxConcurrency)}, config.RateRateLimiterConfigOpts...)
			},
		}, config.ShardManagerConfigOpts...)

		config.ShardManager = sharding.New(token, client.HandleEvent, config.ShardManagerConfigOpts...)
	}
	client.ShardManager = config.ShardManager

	if config.HTTPServer == nil && config.PublicKey != "" {
		config.HTTPServerConfigOpts = append([]httpserver.ConfigOpt{
			httpserver.WithLogger(client.Logger),
		}, config.HTTPServerConfigOpts...)

		config.HTTPServer = httpserver.New(config.PublicKey, client.HandleEvent, config.HTTPServerConfigOpts...)
	}
	client.HTTPServer = config.HTTPServer

	if config.MemberChunkingManager == nil {
		config.MemberChunkingManager = NewMemberChunkingManager(client, config.Logger, config.MemberChunkingFilter)
	}
	client.MemberChunkingManager = config.MemberChunkingManager

	if config.Caches == nil {
		config.Caches = cache.New(config.CacheConfigOpts...)
	}
	client.Caches = config.Caches

	return client, nil
}
