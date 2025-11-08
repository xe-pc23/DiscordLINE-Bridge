package services

import (
	"line-discord-bridge/config"
	"log"

	"github.com/bwmarrin/discordgo"
)

type DiscordService struct {
	session *discordgo.Session
}

var DiscordServiceInstance *DiscordService

// InitDiscordService initializes Discord bot
func InitDiscordService() error {
	session, err := discordgo.New("Bot " + config.AppConfig.DiscordBotToken)
	if err != nil {
		return err
	}

	DiscordServiceInstance = &DiscordService{
		session: session,
	}

	// メッセージ受信ハンドラーを登録
	session.AddHandler(messageCreate)

	// Intent設定（DMを受信するために必要）
	session.Identify.Intents = discordgo.IntentsDirectMessages | discordgo.IntentsGuildMessages

	// Bot起動
	err = session.Open()
	if err != nil {
		return err
	}

	log.Println("Discord Bot started")
	return nil
}

// SendDM sends a DM to specified user
func (s *DiscordService) SendDM(userID string, message string) error {
	// DMチャンネルを作成または取得
	channel, err := s.session.UserChannelCreate(userID)
	if err != nil {
		log.Printf("Failed to create DM channel: %v", err)
		return err
	}

	// メッセージ送信
	_, err = s.session.ChannelMessageSend(channel.ID, message)
	if err != nil {
		log.Printf("Failed to send Discord message: %v", err)
		return err
	}

	log.Printf("Sent Discord DM: %s", message)
	return nil
}

// Close closes Discord session
func (s *DiscordService) Close() {
	if s.session != nil {
		s.session.Close()
	}
}

// messageCreate handles incoming Discord messages
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Bot自身のメッセージは無視
	if m.Author.ID == s.State.User.ID {
		return
	}

	// DM以外は無視（ChannelTypeがDMの場合のみ処理）
	channel, err := s.Channel(m.ChannelID)
	if err != nil || channel.Type != discordgo.ChannelTypeDM {
		return
	}

	// 設定されたユーザーからのメッセージのみ処理
	if m.Author.ID != config.AppConfig.DiscordUserID {
		log.Printf("Ignored message from unknown user: %s", m.Author.ID)
		return
	}

	log.Printf("Received Discord DM: %s", m.Content)

	// LINEに転送
	if LineServiceInstance != nil {
		lastUserID := LineServiceInstance.GetLastUserID()
		if lastUserID != "" {
			err := LineServiceInstance.PushMessage(lastUserID, m.Content)
			if err != nil {
				log.Printf("Failed to forward to LINE: %v", err)
			}
		} else {
			log.Println("No LINE user to forward message to")
		}
	}
}
