package model

type MessageGateway struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	UID      uint   `gorm:"index" json:"uid"`
	Name     string `json:"name"`
	Token    string `gorm:"uniqueIndex;size:67" json:"token"`
	CreateAt int64  `json:"create_at"`
}

type MessageGatewayRequest struct {
	Name string `json:"name"`
}
