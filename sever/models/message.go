package models

import "time"

type EmpxMessage struct {
	NodeID int       `json:"nodeID" db:"node_id"`
	Type   int       `json:"type" db:"type"`
	Value  string    `json:"value" db:"message"`
	TS     time.Time `json:"ts" db:"received_at"`
	UserId int       `json:"userId" db:"user_id"`
}

type RawEmpxMessage struct {
	NodeID int    `json:"nodeID" db:"node_id"`
	Type   string `json:"type" db:"type"`
	Value  string `json:"value" db:"message"`
	TS     int64  `json:"ts" db:"received_at"`
	UserId int    `json:"userId" db:"user_id"`
}

type GetMessage struct {
	UserId int `json:"userId" db:"user_id"`
	Type   int `json:"type" db:"type"`
}

type QueryMessages struct {
	UserId    int       `json:"userId" db:"user_id"`
	Type      int       `json:"type" db:"type"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// EMQXMessagePublish 对应 EMQX 事件 "message.publish" 的完整 payload
type EMQXMessagePublish struct {
	UserName          string   `json:"username"`
	Topic             string   `json:"topic"`
	Timestamp         int64    `json:"timestamp"`
	QoS               int      `json:"qos"`
	PublishReceivedAt int64    `json:"publish_received_at"`
	PubProps          PubProps `json:"pub_props"`
	PeerHost          string   `json:"peerhost"`
	Payload           string   `json:"payload"` // 原始 JSON 字符串，可二次反序列化
	Node              string   `json:"node"`
	Metadata          Metadata `json:"metadata"`
	ID                string   `json:"id"`
	Flags             Flags    `json:"flags"`
	Event             string   `json:"event"` // 固定为 "message.publish"
	ClientID          string   `json:"clientid"`
}

type PubProps struct {
	UserProperty map[string]string `json:"User-Property"`
}

type Metadata struct {
	RuleID string `json:"rule_id"`
}

type Flags struct {
	Retain bool `json:"retain"`
	Dup    bool `json:"dup"`
}
