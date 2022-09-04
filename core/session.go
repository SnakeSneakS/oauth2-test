package core

import (
	"encoding/gob"

	"github.com/gin-contrib/sessions/cookie"
)

var Store cookie.Store

func InitSessionStore() {
	gob.Register(map[string]interface{}{})
	Store = cookie.NewStore([]byte("secret"))
}
