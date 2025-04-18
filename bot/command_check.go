package bot

import (
	"fmt"
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
				message = fmt.Sprintf("This user is currently mining. 🚀\n\nStaked Amount: <code>%.9f FREN</code>\nReward Amount: <code>%.9f FREN</code>\nHealth: <code>%d%%</code>", float64(miner.TMU)/float64(Mul9), float64(miner.rewards(true))/float64(Mul9), miner.health())
				unb := miner.getUnboosted()
				if len(unb) > 0 {
					message += "\n\n<b><u>Boost your miner by clicking the button bellow and then boost button under each post that bot leads you to!</u></b>"
					btn = getButtonLink("Boost Miner 🚀", unb[0].Link)
				}
			} else {
				message = fmt.Sprintf("This user is not mining currently, but has mined %d hours and %d minutes ago.\n\nTo continue mining, click the button bellow.", int64(diff.Hours()), int64(diff.Minutes())%60)
				btn = getRestartButtonChannel()
			}
		} else {
			message = "This user has never mined.\n\nTo start mining, subscribe to our channel and click the button bellow."
			btn = getAppButton()
		}
	} else {
		miner := getUser(c.Message().Sender.ID)
		if miner.ID > 0 {
			diff := time.Since(miner.MiningTime)
			if miner.isActive() {
				message = fmt.Sprintf("You are currently mining. 🚀\n\nStaked Amount: <code>%.9f FREN</code>\nReward Amount: <code>%.9f FREN</code>\nHealth: <code>%d%%</code>", float64(miner.TMU)/float64(Mul9), float64(miner.rewards(true))/float64(Mul9), miner.health())
				unb := miner.getUnboosted()
				if len(unb) > 0 {
					message += "\n\n<b><u>Boost your miner by clicking the button bellow and then boost button under each post that bot leads you to!</u></b>"
					btn = getButtonLink("Boost Miner 🚀", unb[0].Link)
				}
			} else {
				message = fmt.Sprintf("You are not mining currently, but you have mined %d hours and %d minutes ago.\n\nTo continue mining, click the button bellow.", int64(diff.Hours()), int64(diff.Minutes())%60)
				btn = getRestartButtonChannel()
			}
		} else {
			message = "You haven't mined so far.\n\nTo start mining, subscribe to our channel and click the button bellow."
			btn = getAppButton()
		}
	}

	if btn != nil {
		_, err = b.Send(c.Chat(), message, btn, telebot.NoPreview)
	} else {
		_, err = b.Send(c.Chat(), message, telebot.NoPreview)
	}

	return err
}
