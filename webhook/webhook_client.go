package webhook

import (
	"context"
	"errors"
	"net/url"
	"strings"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/log"
	"github.com/disgoorg/snowflake/v2"
)

// ErrInvalidWebhookURL is returned when the given webhookURL is invalid
var ErrInvalidWebhookURL = errors.New("invalid webhook URL")

// NewWithURL creates a new Client by parsing the given webhookURL for the ID and Token.
func NewWithURL(webhookURL string, opts ...ConfigOpt) (*Client, error) {
	u, err := url.Parse(webhookURL)
	if err != nil {
		return nil, err
	}

	parts := strings.FieldsFunc(u.Path, func(r rune) bool { return r == '/' })
	if len(parts) != 4 {
		return nil, ErrInvalidWebhookURL
	}

	token := parts[3]
	id, err := snowflake.Parse(parts[2])
	if err != nil {
		return nil, err
	}

	return New(id, token, opts...), nil
}

// New creates a new Client with the given ID, Token and ConfigOpt(s).
func New(id snowflake.ID, token string, opts ...ConfigOpt) *Client {
	config := DefaultConfig()
	config.Apply(opts)

	return &Client{
		ID:                     id,
		Token:                  token,
		Logger:                 config.Logger,
		Rest:                   config.Webhooks,
		restClient:             config.RestClient,
		defaultAllowedMentions: config.DefaultAllowedMentions,
	}
}

// Client is a high level interface for interacting with Discord's Webhooks API.
type Client struct {
	ID                     snowflake.ID
	Token                  string
	Logger                 log.Logger
	Rest                   rest.Webhooks
	restClient             rest.Client
	defaultAllowedMentions *discord.AllowedMentions
}

// URL returns the full Webhook URL
func (c *Client) URL() string {
	return discord.WebhookURL(c.ID, c.Token)
}

// Close closes all connections the Webhook Client has open
func (c *Client) Close(ctx context.Context) {
	c.restClient.Close(ctx)
}

// GetWebhook fetches the current Webhook from discord

func (c *Client) GetWebhook(opts ...rest.RequestOpt) (*discord.IncomingWebhook, error) {
	webhook, err := c.Rest.GetWebhookWithToken(c.ID, c.Token, opts...)
	if incomingWebhook, ok := webhook.(discord.IncomingWebhook); ok && err == nil {
		return &incomingWebhook, nil
	}
	return nil, err
}

// UpdateWebhook updates the current Webhook
func (c *Client) UpdateWebhook(webhookUpdate discord.WebhookUpdateWithToken, opts ...rest.RequestOpt) (*discord.IncomingWebhook, error) {
	webhook, err := c.Rest.UpdateWebhookWithToken(c.ID, c.Token, webhookUpdate, opts...)
	if incomingWebhook, ok := webhook.(discord.IncomingWebhook); ok && err == nil {
		return &incomingWebhook, nil
	}
	return nil, err
}

// DeleteWebhook deletes the current Webhook
func (c *Client) DeleteWebhook(opts ...rest.RequestOpt) error {
	return c.Rest.DeleteWebhookWithToken(c.ID, c.Token, opts...)
}

// CreateMessageInThread creates a new Message from the discord.WebhookMessageCreate in the provided thread
func (c *Client) CreateMessageInThread(messageCreate discord.WebhookMessageCreate, threadID snowflake.ID, opts ...rest.RequestOpt) (*discord.Message, error) {
	if messageCreate.AllowedMentions == nil {
		messageCreate.AllowedMentions = c.defaultAllowedMentions
	}
	return c.Rest.CreateWebhookMessage(c.ID, c.Token, messageCreate, true, threadID, opts...)
}

// CreateMessage creates a new Message from the discord.WebhookMessageCreate
func (c *Client) CreateMessage(messageCreate discord.WebhookMessageCreate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.CreateMessageInThread(messageCreate, 0, opts...)
}

// CreateContent creates a new Message from the provided content
func (c *Client) CreateContent(content string, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.CreateMessage(discord.WebhookMessageCreate{
		Content:         content,
		AllowedMentions: c.defaultAllowedMentions,
	}, opts...)
}

// CreateEmbeds creates a new Message from the provided discord.Embed(s)
func (c *Client) CreateEmbeds(embeds []discord.Embed, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.CreateMessage(discord.WebhookMessageCreate{
		Embeds:          embeds,
		AllowedMentions: c.defaultAllowedMentions,
	}, opts...)
}

// UpdateMessage updates an already sent Webhook Message with the discord.WebhookMessageUpdate
func (c *Client) UpdateMessage(messageID snowflake.ID, messageUpdate discord.WebhookMessageUpdate, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.UpdateMessageInThread(messageID, messageUpdate, 0, opts...)
}

// UpdateMessageInThread updates an already sent Webhook Message with the discord.WebhookMessageUpdate in the provided thread
func (c *Client) UpdateMessageInThread(messageID snowflake.ID, messageUpdate discord.WebhookMessageUpdate, threadID snowflake.ID, opts ...rest.RequestOpt) (*discord.Message, error) {
	if messageUpdate.AllowedMentions == nil {
		messageUpdate.AllowedMentions = c.defaultAllowedMentions
	}
	return c.Rest.UpdateWebhookMessage(c.ID, c.Token, messageID, messageUpdate, threadID, opts...)
}

// UpdateContent updates an already sent Webhook Message with the content
func (c *Client) UpdateContent(messageID snowflake.ID, content string, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.UpdateMessage(messageID, discord.WebhookMessageUpdate{
		Content:         &content,
		AllowedMentions: c.defaultAllowedMentions,
	}, opts...)
}

// UpdateEmbeds updates an already sent Webhook Message with the discord.Embed(s)
func (c *Client) UpdateEmbeds(messageID snowflake.ID, embeds []discord.Embed, opts ...rest.RequestOpt) (*discord.Message, error) {
	return c.UpdateMessage(messageID, discord.WebhookMessageUpdate{
		Embeds:          &embeds,
		AllowedMentions: c.defaultAllowedMentions,
	}, opts...)
}

// DeleteMessage deletes an already sent Webhook Message
func (c *Client) DeleteMessage(messageID snowflake.ID, opts ...rest.RequestOpt) error {
	return c.DeleteMessageInThread(messageID, 0, opts...)
}

// DeleteMessageInThread deletes an already sent Webhook Message in the provided thread
func (c *Client) DeleteMessageInThread(messageID snowflake.ID, threadID snowflake.ID, opts ...rest.RequestOpt) error {
	return c.Rest.DeleteWebhookMessage(c.ID, c.Token, messageID, threadID, opts...)
}
