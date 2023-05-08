package bot

import (
	"fmt"
	"teneta-tg/internal/entities"
)

const (
	startCommand = "start"

	actAsProviderCommand = "act_as_provider"
	actAsCustomerCommand = "act_as_customer"

	addVCPULimitCommand    = "vcpu"
	addRAMLimitCommand     = "ram"
	addStorageLimitCommand = "storage"
	addNetworkLimitCommand = "network"
	addPortCommand         = "ports"

	aboutCommand = "about"
)

var (
	resourceCommandState = map[string]int{
		addVCPULimitCommand:    addVCPULimitState,
		addRAMLimitCommand:     addRamLimitState,
		addStorageLimitCommand: addStorageLimitState,
		addNetworkLimitCommand: addNetworkLimitState,
		addPortCommand:         addPortsState,
	}
)

func (b *Bot) proceedCommand(user *entities.User, command string) {
	switch command {

	// SYSTEM
	case startCommand:
		b.proceedStartCommand(user)
	case actAsProviderCommand:
		b.execActAsProviderCommand(user)
		b.userService.Save(user)

	// PROVIDER
	case addVCPULimitCommand, addRAMLimitCommand, addStorageLimitCommand, addNetworkLimitCommand, addPortCommand:
		b.proceedAddResourceCommand(user, command, resourceCommandState[command])
		b.userService.Save(user)

		// CUSTOMER
	case actAsCustomerCommand:
		b.proceedActAsCustomerCommand(user)
		b.userService.Save(user)
	default:
		b.messageCh <- &MessageResponse{
			ChatId: user.ChatID,
			Text:   b.translator.Translate("undefined_command", "en", nil),
		}
	}
}

func (b *Bot) proceedStartCommand(user *entities.User) {
	b.response(user, "start_command_response", nil, nil, nil)
}

func (b *Bot) execActAsProviderCommand(user *entities.User) {
	if user.ProviderConfig == nil {
		user.ProviderConfig = &entities.Provider{}
		user.ProviderConfig.ChatID = user.ChatID
	}

	user.State = actAsProviderState

	b.response(user, "act_as_provider_response", nil, nil, nil)
}

func (b *Bot) proceedActAsCustomerCommand(user *entities.User) {
	user.State = actAsCustomerState

	b.response(user, "act_as_customer_response", map[string]interface{}{"count": 0, "costs": 0}, b.actAsCustomerInlineKeyboard(user.Language), nil)
}

func (b *Bot) proceedAddResourceCommand(user *entities.User, t string, state int) {
	if !user.IsProvider() {
		b.response(user, "user_not_registered_as_provider", nil, nil, nil)

		return
	}

	if user.State != actAsProviderState {
		b.response(user, "user_current_context_is_not_provider", nil, nil, nil)

		return
	}

	user.State = state

	b.response(user, fmt.Sprintf("%s_add_start", t), nil, nil, nil)

	return
}