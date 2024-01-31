package modules

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/vote_bot/log"
	"github.com/vote_bot/models"
)

var bot *tgbotapi.BotAPI
var err error
var chatId int64

func Telegram() {

	chatId = int64(models.Config.Telegram.ChatID)
	bot, err = tgbotapi.NewBotAPI(models.Config.Telegram.BotToken)
	if err != nil {
		log.Logger.Error.Fatal(err)
	}
	log.Logger.Trace.Println("Authorized on account", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
}

// Send a text message
func SendTelegramTextMessage(chatID int, str string) {
	msg := tgbotapi.NewMessage(int64(chatID), str)
	bot.Send(msg)
}

// SendTelegramProposalMessage
func SendTelegramProposalMessage(proposal models.Proposals, explore string) {
	msg := tgbotapi.NewMessage(chatId, "#"+proposal.ChainName+"\nproposalID : "+proposal.ProposalID+"\nProposalType : "+proposal.Type+"\nProposalTitle : "+proposal.Title+
		"\nUpdateTime : "+fmt.Sprint(proposal.UpdateDate.Format(time.RFC3339))+"\nEndTime : "+fmt.Sprint(proposal.VotingEndTime.Format(time.RFC3339))+"\nExplore : "+explore)
	sentMsg, _ := bot.Send(msg)

	// update MsgID to DB
	err := updateMsgID(proposal.ChainName, proposal.ProposalID, sentMsg.MessageID)
	models.ProposalsDB[models.DBFind(proposal.ChainName, proposal.ProposalID)].MsgID = sentMsg.MessageID
	if err != nil {
		SendTelegramTextMessage(models.Config.Telegram.ChatIDAdmin, fmt.Sprintf("update updateMsgID failed: %v\n", err))
		log.Logger.Error.Println("update updateMsgID failed: ", err)
	}

}

// You can implement it according to how long the proposal voting period is left.
// ex : ReplyTelegramButtonMessage("We have 24 hours left.")
func ReplyTelegramButtonMessage(msgID int, content string) {
	replyMsg := tgbotapi.NewMessage(chatId, content)
	replyMsg.ReplyToMessageID = msgID
	bot.Send(replyMsg)
}

// Functions which is update MsgID to DB
func updateMsgID(chainName, proposalID string, msgID int) error {
	proposal := models.ProposalsDB[models.DBFind(chainName, proposalID)]
	proposal.MsgID = msgID
	err := models.UpdateDBProposals(proposal)
	if err != nil {
		log.Logger.Error.Println("update updateMsgID failed: ", err)
		return err
	}
	return nil
}
