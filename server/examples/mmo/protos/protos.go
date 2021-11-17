package protos

import (
    "github.com/dfklegend/cell/server/common"
)

// 主要定义协议
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


type RoomEntry struct {
    UId string `json:"uid"`	
    Name string `json:"name"` 
    ServerId string `json:"serverid"`
    NetId common.NetIdType `json:"netid"`
}

type RoomEntryAck struct {
	Code int `json:"code"`
    Result string `json:"result"`
}

type RoomLeave struct { 
    UId string `json:"uid"`
    ServerId string `json:"serverid"`
    NetId common.NetIdType `json:"netid"`
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


