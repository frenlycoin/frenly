package bot

import (
	"fmt"
	"log"

	"gopkg.in/telebot.v3"
)

func commandUserInfo(c telebot.Context) error {
	m := c.Message()
	var os *telebot.User

	if m.IsForwarded() && m.Private() {
		if m.OriginalSender != nil {
			os = m.OriginalSender
		} else {
			log.Println("Original sender is nil")
			return nil
		}
	} else {
		logs("Message is not forwarded or not private")
		return nil
	}

	// Look up the user in the database
	u := getUser(os.ID)
	if u == nil || u.ID == 0 {
		log.Printf("User not found in database: %d", os.ID)
		return nil
	}

	// Format the user info same as cashout output
	username := u.Name
	if len(u.Code) > 0 && u.Code != u.AddressWithdraw {
		username = fmt.Sprintf("@%s", u.Code)
	} else if u.Code == u.AddressWithdraw {
		username = u.Name
	}

	createdAt := "unknown"
	if !u.CreatedAt.IsZero() {
		createdAt = u.CreatedAt.Format("Jan 2, 2006")
	}

	frenAmount := float64(u.TMU) / float64(Mul9)
	compounds := u.CompoundCount
	depositAddress := u.AddressDeposit
	if len(depositAddress) == 0 {
		depositAddress = u.AddressWithdraw
	}
	if len(depositAddress) == 0 {
		depositAddress = "unknown"
	}

	withdrawAddress := u.AddressWithdraw
	if len(withdrawAddress) == 0 {
		withdrawAddress = "unknown"
	}

	msg := fmt.Sprintf(lUserInfo, username, createdAt, u.CycleCountTotal, compounds, formatNumber(frenAmount), depositAddress, depositAddress, withdrawAddress, withdrawAddress)

	// Send the message to the private chat
	err := c.Send(msg, telebot.NoPreview)
	if err != nil {
		loge(err)
	}

	return err
}
