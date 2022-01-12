package entities

import bflag "github.com/vaiktorg/grimoire/bitflag"

type (
	Groups struct {
		Model
		GroupFlags bflag.BitFlag
	}

	Permissions struct {
		Model
		Flags bflag.BitFlag
	}
)
