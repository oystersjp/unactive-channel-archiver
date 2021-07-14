package archiver

import (
	"fmt"
	"github.com/slack-go/slack"
	"strconv"
	"time"
)

type Archiver struct {
	c Config
	api *slack.Client
}

type Config struct {
	ExcludeChannels  []string
	ArchiveLimitDays int
}

type Channel struct {
	ID       string
	Name     string
	IsMember bool
	NumMembers int
}

func NewArchiver(t string, chs []string, d int) *Archiver {
	api := slack.New(t)

	return &Archiver{
		c: Config{
			ExcludeChannels:  chs,
			ArchiveLimitDays: d,
		},
		api: api,
	}
}

func (a Archiver) Exec(summaryCh string) {
	// チャンネル一覧取得
	chs := a.GetChannels()
	// botID取得
	botID := a.GetBotID()
	var archiveChs []string
	var summaryChID string
	// チャンネルごとに処理
	for _, ch := range chs {
		// チャンネル入っていないと履歴取れないのでjoin
		if !ch.IsMember {
			ch.Join(a.api)
		}
		// 投稿用のチャンネルIDを取得しておく
		if ch.Name == summaryCh {
			summaryChID = ch.ID
		}
		// アーカイブ対象はアーカイブする
		limit := time.Now().Add(-1 * time.Duration(a.c.ArchiveLimitDays) * 24 * time.Hour).Unix()
		if ch.DecideArchive(a.api, botID, limit) {
			archiveChs = append(archiveChs, ch.Name)
			ch.Archive(a.api)
		}
	}
	// アーカイブしたチャンネルをポストする
	a.PostSummary(summaryChID, archiveChs)
}

func (a Archiver) GetChannels() []Channel {
	param := &slack.GetConversationsParameters{
		ExcludeArchived: true,
		Types:           []string{"public_channel"},
		Limit:           1000,
	}

	var chs []Channel
	for {
		resChs, cursor, _ := a.api.GetConversations(param)
		for _, resCh := range resChs {
			// 外部共有チャンネルは除外
			if resCh.IsExtShared {
				continue
			}
			chs = append(
				chs,
				Channel{
					ID:         resCh.ID,
					Name:       resCh.Name,
					IsMember:   resCh.IsMember,
					NumMembers: resCh.NumMembers,
				},
			)
		}
		// ページネーション
		if cursor == "" {
			break
		}
		param.Cursor = cursor
	}
	return chs
}

func (a Archiver) GetBotID() string {
	res, err := a.api.AuthTest()
	if err != nil {
		panic("[Error] failed get auth information " + err.Error())
	}
	return res.BotID
}

// apiのrate limit用
const waitTime = time.Second * 2

func (c *Channel) Join(api *slack.Client) {
	_, _, _, err := api.JoinConversation(c.ID)
	if err != nil {
		fmt.Println("[Error] failed channel join " + err.Error())
	}
	time.Sleep(waitTime)
}

func (c *Channel) DecideArchive(api *slack.Client, botID string, oldestUnix int64) bool {
	// 履歴を取得
	h, err := api.GetConversationHistory(
		&slack.GetConversationHistoryParameters{
			ChannelID: c.ID,
			Oldest:    strconv.FormatInt(oldestUnix, 10),
		},
	)
	if err != nil {
		fmt.Println("[Error] failed channel histories " + c.Name + err.Error())
		return false
	}
	time.Sleep(waitTime)
	for _, m := range h.Messages {
		// 自分の投稿は無視
		if m.BotID == botID {
			continue
		}
		// 通常の message, bot_message もしくは unarchive があれば閉じない
		if m.SubType == "" || m.SubType == "bot_message" || m.SubType == "channel_unarchive" {
			return false
		}
	}
	return true
}

func (c *Channel) Archive(api *slack.Client) {
	_, _, err := api.PostMessage(
		c.ID,
		slack.MsgOptionText(
			"このチャンネルは長期間投稿がなかったため自動archiveします。\nもし利用されている場合はお手数ですが手動でunarchiveしてください。",
			true,
		),
	)
	if err != nil {
		fmt.Println("[Error] failed archive post channel: " + c.Name + err.Error())
	}
	err = api.ArchiveConversation(c.ID)
	if err != nil {
		fmt.Println("[Error] failed archive channel: " + c.Name + err.Error())
	}
}

func (a Archiver) PostSummary(postChId string, chs []string) {
	summaryText := "下記チャンネルは長期間activeでなかったためarchiveしました。\n必要なチャンネルの場合は手動でunarchiveしてください。\n"
	for _, ch := range chs {
		summaryText += "#" + ch + "\n"
	}
	_, _, err := a.api.PostMessage(postChId, slack.MsgOptionText(summaryText, true))
	if err != nil {
		panic("[Error] failed post summary" + err.Error())
	}
}
