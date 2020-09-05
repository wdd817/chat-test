package chat

import (
	"chat-test/app"
	"chat-test/conf"
	"chat-test/log"
	"chat-test/proto"
	"fmt"
	"time"
)

func (chat *Chat) OnChatMsgReq(req0 interface{}) interface{} {
	req := req0.(*proto.ChatMsgReq)

	user := chat.users[req.UserID]
	m := &msg{
		sender:  user.Name,
		content: req.Content,
		timestamp: time.Now().Unix(),
	}

	chat.history.Add(m)
	chat.recent.Add(m)

	chat.broadCast(fmt.Sprintf("[%s] [%s] %s",
		time.Unix(m.timestamp, 0).String(),
		m.sender,
		user.Name,
	))
	return &proto.ChatMsgResp{}
}

func (chat *Chat) OnPopularReq(_ interface{}) interface{} {
	return &proto.PopularResp{
		Word: chat.recent.Popular(),
	}
}

func (chat *Chat) OnLoginReq(req0 interface{}) interface{} {
	req := req0.(*proto.LoginReq)
	_, exist := chat.userIDs[req.Name]
	if exist {
		return &proto.LoginResp{Success: false}
	}

	chat.userIDs[req.Name] = req.UserID
	chat.users[req.UserID] = &User{
		ID:      req.UserID,
		Name:    req.Name,
		LoginTS: time.Now().Unix(),
	}

	var histories []string
	chat.history.Range(func(m *msg) {
		histories = append(histories, fmt.Sprintf("[%s] [%s] %s",
			time.Unix(m.timestamp, 0).String(),
			m.sender,
			m.content,
		))
	})

	chat.broadCast(fmt.Sprintf("[%s] user [%s] login",
		time.Now().String(),
		req.Name))

	return &proto.LoginResp{
		Success:   true,
		Histories: histories,
	}
}

func (chat *Chat) OnStatsReq(req0 interface{}) interface{} {
	req := req0.(*proto.StatsReq)
	targetID, exist := chat.userIDs[req.Name]

	var timeStr string
	if exist {
		target := chat.users[targetID]
		diff := time.Now().Unix() - target.LoginTS
		timeStr = formatTime(diff)
	}

	return &proto.StatsResp{
		TimeStr: timeStr,
	}
}

func (chat *Chat) OnLogoutReq(req0 interface{}) interface{} {
	req := req0.(*proto.LogoutReq)
	user := chat.users[req.UserID]
	if user != nil {
		delete(chat.userIDs, user.Name)
		delete(chat.users, user.ID)
		chat.broadCast(fmt.Sprintf("[%s] user [%s] logout",
			time.Now().String(),
			user.Name))
	}
	return &proto.LogoutResp{}
}

func formatTime(seconds int64) string {
	day := seconds / 86400
	hour := seconds % 86400 / 3600
	min := seconds % 3600 / 60
	sec := seconds % 60
	return fmt.Sprintf("%02dd %02dh %02dm %02ds", day, hour, min, sec)
}

func (chat *Chat) broadCast(msg string) {
	chat.skeleton.Async(app.Instance().GetChanRPC(conf.GateModule), &proto.BroadcastReq{
		Content: msg,
	}, func(resp0 interface{}, err error) {
		if err != nil {
			log.Error("broadcast call failed: %v", err)
			return
		}
	})
}
