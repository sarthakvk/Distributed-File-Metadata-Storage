package httpd

import (
	"fmt"
	"net/http"

	keystore_adapters "github.com/sarthakvk/hex-app/adapters/keystore_adapter"
	"github.com/sarthakvk/hex-app/adapters/logging"
)

var (
	Keystore keystore_adapters.AbstractKeyStore
	logger   = logging.GetHttpdLogger()
)

func RunServer(store keystore_adapters.AbstractKeyStore, port int) {
	logger.Debug("Running HTTP daemon")
	Keystore = store
	for _, url := range urls {
		http.HandleFunc(url.Pattern, url.Handler)
	}
	addr := fmt.Sprintf(":%d", port)
	http.ListenAndServe(addr, nil)
}
