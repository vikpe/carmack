package embed_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vikpe/carmack/internal/pkg/carmack/embed"
	"github.com/vikpe/serverstat/qserver"
	"github.com/vikpe/serverstat/qserver/convert"
)

func TestFromServer(t *testing.T) {
	t.Run("mvdsv", func(t *testing.T) {
		server := qserver.GenericServer{Version: "mvdsv"}
		expect := embed.FromMvdsvServer(convert.ToMvdsv(server))
		assert.Equal(t, expect, embed.FromServer(server))
	})

	t.Run("qtv", func(t *testing.T) {
		server := qserver.GenericServer{Version: "qtv"}
		expect := embed.FromQtvServer(convert.ToQtv(server))
		assert.Equal(t, expect, embed.FromServer(server))
	})

	t.Run("qwfwd", func(t *testing.T) {
		server := qserver.GenericServer{Version: "qwfwd"}
		expect := embed.FromQwfwdServer(convert.ToQwfwd(server))
		assert.Equal(t, expect, embed.FromServer(server))
	})

	t.Run("generic", func(t *testing.T) {
		server := qserver.GenericServer{Version: "unknown"}
		expect := embed.FromGenericServer(server)
		assert.Equal(t, expect, embed.FromServer(server))
	})
}
