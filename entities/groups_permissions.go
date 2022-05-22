package entities

import "github.com/vaiktorg/grimoire/bitflag"

type (
	Groups struct {
		Model
		GroupFlags bitflag.BitFlag
	}

	Permissions struct {
		Model
		Flags bitflag.BitFlag
	}
)
