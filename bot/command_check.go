package bot

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/telebot.v3"
)

func commandCheck(c telebot.Context) error {
	var err error
	var btn interface{}

	message := ""

	if c.Message().IsReply() {
		miner := getUser(c.Message().ReplyTo.Sender.ID)
		if miner.ID > 0 {
			diff := time.Since(miner.MiningTime)
			if miner.isActive() {
				message = "This user is currently mining. 🚀"
			} else {
				message = fmt.Sprintf("This user is not mining currently, but has mined %d hours and %d minutes ago.\n\nTo continue mining, click the button bellow.", int64(diff.Hours()), int64(diff.Minutes())%int64(diff.Hours()))
				btn = getRestartButton()
			}
		} else {
			message = "This user has never mined.\n\nTo start mining, subscribe to our channel and click the button bellow."
			btn = getAppButton()
		}
		log.Println(prettyPrint(miner))
	} else {
		miner := getUser(c.Message().Sender.ID)
		if miner.ID > 0 {
			diff := time.Since(miner.MiningTime)
			if miner.isActive() {
				message = "You are currently mining. 🚀"
			} else {
				message = fmt.Sprintf("You are not mining currently, but you have mined %d hours and %d minutes ago.\n\nTo continue mining, click the button bellow.", int64(diff.Hours()), int64(diff.Minutes())%int64(diff.Hours()))
				btn = getRestartButton()
			}
		} else {
			message = "You haven't mined so far.\n\nTo start mining, subscribe to our channel and click the button bellow."
			btn = getAppButton()
		}
	}

	_, err = b.Send(c.Chat(), message, btn, telebot.NoPreview)

	return err
}
