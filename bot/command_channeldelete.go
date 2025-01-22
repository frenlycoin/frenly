package bot

import (
	"fmt"
	"log"
	"strings"

	"gopkg.in/telebot.v3"
)

func commandChannelDelete(c telebot.Context) error {
	u := getUser(c.Sender().ID)
	var err error
	msg := "Successfully removed faulty boost."
	var btn *telebot.ReplyMarkup

	url := strings.Replace(c.Data(), "\f", "", -1)
	urlSplit := strings.Split(url, "/")
	channelLink := urlSplit[1]
	postId := urlSplit[2]
	ch := getChannelByLink(channelLink)

	log.Printf("channel_id = %d AND telegram_id = %s", ch.ID, postId)

	if !postExists(postId, int(ch.TelegramId)) {
		p := &Post{}
		err := db.Where("channel_id = ? AND telegram_id = ?", ch.ID, postId).First(p).Error
		if err != nil {
			loge(err)
		}

		log.Println(prettyPrint(p))

		err = db.Delete(p).Error
		if err != nil {
			loge(err)
		}
	} else {
		msg = "This boost is not missing."
	}

	unb := u.getUnboosted()

	if len(unb) > 0 {
		btn = getButtonsBoost(unb[0].Name, unb[0].Link)
		msg += fmt.Sprintf("\n\nBoosts Left: %d\n\nClick the button bellow for the next boost:", len(unb))
	} else {
		msg += "\n\nYou have no more boosts available. 👍"
	}

	notifyWithButton(msg, c.Sender().ID, btn)

	return err
}
