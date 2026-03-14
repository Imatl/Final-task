package integrations

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"

	"supportflow/db/postgre"
)

type Integration struct {
	ID     string            `json:"id"`
	Name   string            `json:"name"`
	Type   string            `json:"type"`
	Status string            `json:"status"`
	Config map[string]string `json:"config,omitempty"`
}

type Channel interface {
	Start() error
	Stop() error
	Status() string
}

var (
	mu       sync.RWMutex
	channels = map[string]Channel{}
	configs  = map[string]*Integration{}
)

func Init(ctx context.Context) {
	rows, err := postgre.ListIntegrations(ctx)
	if err != nil {
		log.Printf("[integrations] load from DB error: %v", err)
		return
	}

	for _, r := range rows {
		if r.Status != "connected" {
			continue
		}
		if err := Connect(ctx, r.ID, r.Type, r.Config); err != nil {
			log.Printf("[integrations] auto-connect %s failed: %v", r.ID, err)
		}
	}
}

func Connect(ctx context.Context, id, channelType string, config map[string]string) error {
	mu.Lock()
	defer mu.Unlock()

	if ch, ok := channels[id]; ok {
		ch.Stop()
	}

	var ch Channel
	switch channelType {
	case "telegram":
		token := config["bot_token"]
		if token == "" {
			return fmt.Errorf("bot_token is required")
		}
		ch = NewTelegramChannel(token)
	default:
		return fmt.Errorf("unsupported channel type: %s", channelType)
	}

	if err := ch.Start(); err != nil {
		return fmt.Errorf("failed to start %s: %w", channelType, err)
	}

	channels[id] = ch
	configs[id] = &Integration{
		ID:     id,
		Name:   config["name"],
		Type:   channelType,
		Status: "connected",
		Config: config,
	}

	if err := postgre.UpsertIntegration(ctx, &postgre.IntegrationRow{
		ID:     id,
		Type:   channelType,
		Name:   config["name"],
		Config: config,
		Status: "connected",
	}); err != nil {
		log.Printf("[integrations] save to DB error: %v", err)
	}

	log.Printf("[integrations] %s (%s) connected", channelType, id)
	return nil
}

func Disconnect(ctx context.Context, id string) error {
	mu.Lock()
	defer mu.Unlock()

	ch, ok := channels[id]
	if !ok {
		return fmt.Errorf("integration %s not found", id)
	}

	ch.Stop()
	delete(channels, id)
	if cfg, ok := configs[id]; ok {
		cfg.Status = "disconnected"
	}

	if err := postgre.UpdateIntegrationStatus(ctx, id, "disconnected"); err != nil {
		log.Printf("[integrations] update DB status error: %v", err)
	}

	log.Printf("[integrations] %s disconnected", id)
	return nil
}

func List() []Integration {
	mu.RLock()
	defer mu.RUnlock()

	var result []Integration
	for _, cfg := range configs {
		c := *cfg
		if ch, ok := channels[cfg.ID]; ok {
			c.Status = ch.Status()
		}
		c.Config = nil
		result = append(result, c)
	}
	return result
}

func SendToCustomer(customerID, text string) error {
	ctx := context.Background()
	extID, err := postgre.GetChannelMapping(ctx, customerID, "telegram")
	if err != nil {
		return fmt.Errorf("no channel can reach customer %s", customerID)
	}

	chatID, err := strconv.ParseInt(extID, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid chat ID for customer %s", customerID)
	}

	mu.RLock()
	defer mu.RUnlock()

	for _, ch := range channels {
		if tg, ok := ch.(*TelegramChannel); ok {
			tg.sendReply(chatID, text)
			return nil
		}
	}

	return fmt.Errorf("no telegram channel active")
}

func GetStatus(id string) string {
	mu.RLock()
	defer mu.RUnlock()

	if ch, ok := channels[id]; ok {
		return ch.Status()
	}
	if cfg, ok := configs[id]; ok {
		return cfg.Status
	}
	return "not_configured"
}
