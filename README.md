## 共享剪切板工具

### usage

server：

```go
package main

import (
	clipboard "github.com/obgnail/share-clipboard"
)

func main() {
  addr := "192.168.3.3:8899"
  clipboard.ServerRun(addr)
}
```

peer：

```go
package main

import (
	"github.com/juju/errors"
	clipboard "github.com/obgnail/share-clipboard"
	log "github.com/sirupsen/logrus"
	"golang.design/x/hotkey/mainthread"
)

func main() { mainthread.Init(fn) }

func fn() {
  addr   := "192.168.3.3:8899"
  sendHK := "Alt+C"
  loadHK := "Alt+V"

	err = clipboard.ListenHotKey(sendHK, func() {
		clipboard.PeerSendClipboard(addr)
	})
	if err != nil {
		log.Error(errors.ErrorStack(err))
	}
	err = clipboard.ListenHotKey(loadHK, func() {
		clipboard.PeerLoadClipboard(addr)
	})
	if err != nil {
		log.Error(errors.ErrorStack(err))
	}

	forever := make(chan struct{}, 1)
	<-forever
}
```

