package ui

import (
	"github.com/ardanlabs/usdl/chat/api/frontends/client/app"
	"github.com/ethereum/go-ethereum/common"
)

type App interface {
	SendMessageHandler(to common.Address, msg []byte) error
	Contacts() []app.User
	QueryContactByID(id common.Address) (app.User, error)
}
