package integrations

import (
	"context"
	"fmt"
	"log"
	"sync"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"supportflow/core/constants"
	"supportflow/core/structs"
	"supportflow/db/postgre"
	"supportflow/services/ai"
)

type TelegramChannel struct {
	token  string
	bot    *tgbotapi.BotAPI
	status string
	cancel context.CancelFunc
	mu     sync.RWMutex
}

func NewTelegramChannel(token string) *TelegramChannel {
	return &TelegramChannel{
		token:  token,
		status: "connecting",
	}
}

func (t *TelegramChannel) Start() error {
	bot, err := tgbotapi.NewBotAPI(t.token)
	if err != nil {
		t.mu.Lock()
		t.status = "error"
		t.mu.Unlock()
		return fmt.Errorf("invalid bot token: %w", err)
	}

	t.mu.Lock()
	t.bot = bot
	t.status = "connected"
	t.mu.Unlock()

	log.Printf("[telegram] bot @%s connected", bot.Self.UserName)

	ctx, cancel := context.WithCancel(context.Background())
	t.mu.Lock()
	t.cancel = cancel
	t.mu.Unlock()

	go t.poll(ctx)
	return nil
}

func (t *TelegramChannel) Stop() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.cancel != nil {
		t.cancel()
	}
	if t.bot != nil {
		t.bot.StopReceivingUpdates()
	}
	t.status = "disconnected"
	log.Println("[telegram] bot disconnected")
	return nil
}

func (t *TelegramChannel) Status() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.status
}

func (t *TelegramChannel) poll(ctx context.Context) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 30

	updates := t.bot.GetUpdatesChan(u)

	for {
		select {
		case <-ctx.Done():
			return
		case update, ok := <-updates:
			if !ok {
				return
			}
			if update.Message == nil {
				continue
			}
			go t.handleMessage(update.Message)
		}
	}
}

func (t *TelegramChannel) handleMessage(msg *tgbotapi.Message) {
	ctx := context.Background()
	text := msg.Text
	if text == "" {
		return
	}

	log.Printf("[telegram] message from %s (@%s): %s", msg.From.FirstName, msg.From.UserName, truncate(text, 50))

	name := msg.From.FirstName
	if msg.From.LastName != "" {
		name += " " + msg.From.LastName
	}
	email := fmt.Sprintf("tg-%d@telegram.user", msg.From.ID)
	phone := ""
	if msg.Contact != nil && msg.Contact.PhoneNumber != "" {
		phone = msg.Contact.PhoneNumber
	}

	customer, err := postgre.FindOrCreateCustomer(ctx, email, name, "free", phone)
	if err != nil {
		log.Printf("[telegram] find/create customer error: %v", err)
		return
	}

	if err := postgre.UpsertChannelMapping(ctx, customer.ID, "telegram", fmt.Sprintf("%d", msg.Chat.ID)); err != nil {
		log.Printf("[telegram] save channel mapping error: %v", err)
	}

	ticketID := findOpenTicket(ctx, customer.ID)
	if ticketID == "" {
		subject := truncate(text, 100)
		ticket := &structs.Ticket{
			CustomerID: customer.ID,
			Subject:    subject,
			Channel:    "telegram",
			Status:     constants.TicketStatusOpen,
			Priority:   constants.PriorityMedium,
			Category:   "general",
		}
		if err := postgre.CreateTicket(ctx, ticket); err != nil {
			log.Printf("[telegram] create ticket error: %v", err)
			return
		}
		ticketID = ticket.ID
		log.Printf("[telegram] created ticket %s for %s", ticketID, customer.ID)
	}

	custMsg := &structs.Message{
		TicketID: ticketID,
		Role:     constants.RoleCustomer,
		Content:  text,
	}
	if err := postgre.CreateMessage(ctx, custMsg); err != nil {
		log.Printf("[telegram] save message error: %v", err)
	}

	resp, err := ai.ProcessMessage(ctx, ticketID, text, msg.From.LanguageCode)
	if err != nil {
		log.Printf("[telegram] AI error: %v", err)
		t.sendReply(msg.Chat.ID, "Sorry, something went wrong. Please try again.")
		return
	}

	aiMsg := &structs.Message{
		TicketID: ticketID,
		Role:     constants.RoleAI,
		Content:  resp.Message,
	}
	if err := postgre.CreateMessage(ctx, aiMsg); err != nil {
		log.Printf("[telegram] save AI message error: %v", err)
	}

	if _, err := ai.AnalyzeTicket(ctx, ticketID, text); err != nil {
		log.Printf("[telegram] analyze error: %v", err)
	}

	t.sendReply(msg.Chat.ID, resp.Message)
}

func (t *TelegramChannel) sendReply(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"
	if _, err := t.bot.Send(msg); err != nil {
		msg.ParseMode = ""
		if _, err := t.bot.Send(msg); err != nil {
			log.Printf("[telegram] send error: %v", err)
		}
	}
}

func findOpenTicket(ctx context.Context, customerID string) string {
	tickets, _, err := postgre.ListTickets(ctx, structs.TicketFilter{Status: "open"})
	if err != nil {
		return ""
	}
	for _, t := range tickets {
		if t.CustomerID == customerID {
			return t.ID
		}
	}
	tickets2, _, err := postgre.ListTickets(ctx, structs.TicketFilter{Status: "in_progress"})
	if err != nil {
		return ""
	}
	for _, t := range tickets2 {
		if t.CustomerID == customerID {
			return t.ID
		}
	}
	return ""
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
