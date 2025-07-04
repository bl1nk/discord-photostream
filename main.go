package main

import (
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		logger.Error("DISCORD_BOT_TOKEN environment variable not set")
		return
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		logger.Error("Error creating Discord session", "error", err)
		return
	}

	channelID := os.Getenv("DISCORD_CHANNEL_ID")
	if channelID == "" {
		logger.Error("DISCORD_CHANNEL_ID environment variable not set")
		logger.Info("Available guilds and channels:")

		done := make(chan bool)

		// Add a one-time handler for the ready event to list guilds and channels.
		dg.AddHandlerOnce(func(s *discordgo.Session, r *discordgo.Ready) {
			for _, guild := range s.State.Guilds {
				logger.Info("Guild", "name", guild.Name, "id", guild.ID)
				channels, err := s.GuildChannels(guild.ID)
				if err != nil {
					logger.Error("Error getting channels", "guild_id", guild.ID, "error", err)
					continue
				}
				for _, c := range channels {
					if c.Type == discordgo.ChannelTypeGuildText {
						logger.Info("Channel", "name", c.Name, "id", c.ID)
					}
				}
			}
			done <- true
		})

		err = dg.Open()
		if err != nil {
			logger.Error("Error opening connection", "error", err)
			return
		}
		defer dg.Close()

		<-done
		return
	}

	dg.AddHandler(onMessage)

	err = dg.Open()
	if err != nil {
		logger.Error("Error opening connection", "error", err)
		return
	}
	defer dg.Close()

	logger.Info("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

// onMessage is called whenever a new message is postes to a channel that the bot has access to.
func onMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	channelID := os.Getenv("DISCORD_CHANNEL_ID")

	// Only process messages in the configured channel.
	if m.ChannelID != channelID {
		return
	}

	// If the message does not have any image attachments, delete it.
	hasImage := false
	for _, attachment := range m.Attachments {
		// In Discord, image attachments have a ContentType that starts with "image/".
		if strings.HasPrefix(attachment.ContentType, "image/") {
			hasImage = true
			break
		}
	}

	if !hasImage {
		logger.Info("Deleting message without image attachment", "message_id", m.ID)
		err := s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			logger.Error("Error deleting message", "message_id", m.ID, "error", err)
		}
	}
}
