package protos


type EmptyArgReq struct { 
}

type NormalAck struct {
    Code int `json:"code"`
    Result string `json:"result"`
}

type QueryGateReq struct { 
}

type QueryGateAck struct {
    Code int `json:"code"`
    IP string `json:"ip"`
    Port string `json:"port"`
}

type LoginReq struct { 
    Name string `json:"name"`
}


type RoomEntryAck struct {
	Code int `json:"code"`
    Result string `json:"result"`
}

type OnNewUser struct {    
    Name string `json:"name"`
}

type OnMembers struct {    
    Members []string `json:"members"`
}

type OnUserLeave struct {    
    Name string `json:"name"`
}

type ChatMsg struct {    
    Name string `json:"name"`
    Content string `json:"content"`
}


