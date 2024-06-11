package internal

import (
	"github.com/jimyag/pastenotifier/handle"
)

func init() {
	PN.Register(&handle.Timestamp{})
	PN.Register(&handle.IpIsp{})
}
