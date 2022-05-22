package share_clipboard

import (
	"fmt"
	"github.com/juju/errors"
	"golang.design/x/hotkey"
	"strings"
)

// Translate eg: ctrl+shift+C
func translate(hotkey string) (mods []hotkey.Modifier, key hotkey.Key, err error) {
	for _, ele := range strings.Split(hotkey, "+") {
		e := strings.ToUpper(ele)
		mod, ok1 := ModifierMap[e]
		if ok1 {
			mods = append(mods, mod)
		}
		k, ok2 := KeyMap[e]
		if ok2 {
			if key != 0 {
				return nil, 0, fmt.Errorf("support one key only")
			}
			key = k
		}
		if !ok1 && !ok2 {
			return nil, 0, fmt.Errorf("wrong key")
		}
	}
	return mods, key, nil
}

func ListenHotKey(hk string, hook func()) error {
	mods, key, err := translate(hk)
	if err != nil {
		return errors.Trace(err)
	}
	HK := hotkey.New(mods, key)
	if err := HK.Register(); err != nil {
		return errors.Trace(err)
	}
	go func() {
		for {
			<-HK.Keydown()
			<-HK.Keyup()
			hook()
		}
	}()
	return nil
}
