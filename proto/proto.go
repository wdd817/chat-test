package proto

type WordFilterReq struct {
	Content string
}

type WordFilterResp struct {
	IsFiltered bool
	Filtered   string
}

type ChatMsgReq struct {
	UserID  int
	Content string
}

type ChatMsgResp struct {
}

type LoginReq struct {
	UserID int
	Name   string
}

type LoginResp struct {
	Success   bool
	Histories []string
}

type LogoutReq struct {
	UserID int
}

type LogoutResp struct {
}

type PopularReq struct {
}

type PopularResp struct {
	Word string
}

type StatsReq struct {
	Name string
}

type StatsResp struct {
	TimeStr string
}

type BroadcastReq struct {
	Content string
}

type BroadcastResp struct {
}
