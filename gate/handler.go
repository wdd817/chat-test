package gate

import (
	"chat-test/app"
	"chat-test/conf"
	"chat-test/log"
	"chat-test/proto"
	"fmt"
	"strings"
)

func (a *agent) checkLogined() bool {
	if a.userID == 0 {
		a.WriteMsg("please input your name to login first")
		return false
	}
	return true
}

func (gate *Gateway) handlePopular(msg string, arg interface{}) {
	a := arg.(*agent)
	if !a.checkLogined() {
		return
	}

	gate.popularReq(msg, a)
}

func (gate *Gateway) handleStats(msg string, arg interface{}) {
	a := arg.(*agent)
	if !a.checkLogined() {
		return
	}

	gate.statsReq(msg, a)
}

func (gate *Gateway) handleChat(msg string, arg interface{}) {
	a := arg.(*agent)

	if a.userID == 0 {
		gate.loginReq(msg, a)
		return
	}

	gate.wordFilteredReq(msg, a, func(resp *proto.WordFilterResp) {
		gate.chatReq(msg, a)
	})
}

func (gate *Gateway) wordFilteredReq(msg string, a *agent, f func(resp *proto.WordFilterResp)) {
	gate.skeleton.Async(app.Instance().GetChanRPC(conf.WordFilterModule), &proto.WordFilterReq{
		Content: msg,
	}, func(resp0 interface{}, err error) {
		if err != nil {
			log.Error("word filtered failed: %v", err)
			return
		}

		f(resp0.(*proto.WordFilterResp))
	})
}

func (gate *Gateway) loginReq(msg string, a *agent) {
	gate.wordFilteredReq(msg, a, func(resp *proto.WordFilterResp) {
		if resp.IsFiltered {
			a.WriteMsg("You cannot use profanity name")
			return
		}

		a.logging = NextID()
		gate.skeleton.Async(app.Instance().GetChanRPC(conf.ChatModule), &proto.LoginReq{
			UserID: a.logging,
			Name: msg,
		}, func(resp0 interface{}, err error) {
			if err != nil {
				log.Error("login failed: %v", err)
				return
			}

			resp := resp0.(*proto.LoginResp)
			if !resp.Success {
				a.WriteMsg("Name already exists, please try again:")
			} else {
				a.WriteMsg([]byte("You have logined"))

				gate.agents[a.logging] = a
				a.userID = a.logging
				a.logging = 0
				for _, msg := range resp.Histories {
					a.WriteMsg(msg)
				}
			}
		})
	})
}

func (gate *Gateway) popularReq(_ string, a *agent) {
	gate.skeleton.Async(app.Instance().GetChanRPC(conf.ChatModule), &proto.PopularReq{
	}, func(resp0 interface{}, err error) {
		if err != nil {
			log.Error("get popular failed: %v", err)
			return
		}

		resp := resp0.(*proto.PopularResp)
		if len(resp.Word) == 0 {
			a.WriteMsg("Current has NO popular world!")
		} else {
			a.WriteMsg(fmt.Sprintf("Current popular word: %v", resp.Word))
		}
	})
}

func (gate *Gateway) statsReq(msg string, a *agent) {
	segs := strings.Split(msg, " ")
	if len(segs) < 2 {
		a.WriteMsg("empty stats parameter, try again")
		return
	}

	gate.skeleton.Async(app.Instance().GetChanRPC(conf.ChatModule), &proto.StatsReq{
		Name: segs[1],
	}, func(resp0 interface{}, err error) {
		if err != nil {
			log.Error("get stats failed: %v", err)
			return
		}

		resp := resp0.(*proto.StatsResp)
		if len(resp.TimeStr) == 0 {
			a.WriteMsg(fmt.Sprintf("Cannot find user: %v", msg))
		} else {
			a.WriteMsg(resp.TimeStr)
		}
	})
}

func (gate *Gateway) chatReq(msg string, a *agent) {
	gate.skeleton.Async(app.Instance().GetChanRPC(conf.ChatModule), &proto.ChatMsgReq{
		UserID: a.userID,
		Content: msg,
	}, func(resp0 interface{}, err error) {
		if err != nil {
			log.Error("send chat msg failed: %v", err)
			return
		}
	})
}

func (gate *Gateway) logoutReq(msg string, a *agent) {
	loginID := a.userID
	if loginID == 0 {
		loginID = a.logging
	} else {
		delete(gate.agents, loginID)
	}
	gate.skeleton.Async(app.Instance().GetChanRPC(conf.ChatModule), &proto.LogoutReq{
		UserID: loginID,
	}, func(resp0 interface{}, err error) {
		if err != nil {
			log.Error("logout failed: %v", err)
			return
		}
	})
}

func (gate *Gateway) OnBroadcastReq(req0 interface{}) interface{} {
	req := req0.(*proto.BroadcastReq)
	for _, a := range gate.agents {
		a.WriteMsg(req.Content)
	}
	return &proto.BroadcastResp{}
}