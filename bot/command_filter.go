package bot

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gopkg.in/telebot.v3"
	"mvdan.cc/xurls"
)

func commandFilter(c telebot.Context) error {
	var err error
	m := c.Message()

	// Handle the persistent reply keyboard buttons
	if m.Private() {
		if m.Text == "⚪️ Restart Mining" {
			u := getUser(c.Sender().ID)
			rb := getRestartButtons(c)
			if time.Since(u.MiningTime).Minutes() > 1410 {
				u.MiningTime = time.Now()
				u.LastNotification = time.Now()
				u.CycleCount++
				if err := db.Save(u).Error; err != nil {
					loge(err)
				}
				b.Send(c.Sender(), lCycleRestarted, rb, getRestartReplyKeyboard(), telebot.NoPreview)
			} else {
				b.Send(c.Sender(), lCycleRunning, rb, getRestartReplyKeyboard(), telebot.NoPreview)
			}
			return nil
		} else if m.Text == "🔄 Compound" {
			return commandCompound(c)
		}
	}

	// Call commandUserInfo if message is private and forwarded
	if m.Private() && m.IsForwarded() {
		err = commandUserInfo(c)
		if err != nil {
			logs(err.Error())
		}
	}

	text := c.Message().Text

	log.Println(text)

	group, err := b.ChatByID(m.Chat.ID)
	if err != nil {
		logs(err.Error())
	}

	cm, err := b.ChatMemberOf(group, m.Sender)
	if err != nil {
		logs(err.Error())
	}

	// lnk, err := regexp.MatchString("^https?:\\/\\/(?:www\\.)?[-a-zA-Z0-9@:%._\\+~#=]{1,256}\\.[a-zA-Z0-9()]{1,6}\\b(?:[-a-zA-Z0-9()@:%_\\+.~#?&\\/=]*)$", text)
	// if err != nil {
	// 	logs(err.Error())
	// }

	// url, err := xurls.Relaxed.FindString(text)
	url := xurls.Relaxed.FindString(text)
	log.Println(url)
	lnk := len(url) > 0

	if cm.Role != telebot.Administrator &&
		cm.Role != telebot.Creator &&
		m.Chat.ID == Group &&
		!m.Sender.IsBot &&
		(strings.Contains(text, "@") ||
			strings.Contains(text, "http://") ||
			strings.Contains(text, "https://") ||
			lnk) {

		log.Println("ban")

		lg := fmt.Sprintf("Banned: %s\n\n%s", m.Sender.Username, m.Text)
		logTelegramSilent(lg)

		err = b.Delete(m)
		if err != nil {
			logs(err.Error())
		}

		err := b.Restrict(m.Chat, cm)
		if err != nil {
			logs(err.Error())
		}
	}

	return err
}
