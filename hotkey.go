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

var ModifierMap = map[string]hotkey.Modifier{
	"ALT":   hotkey.ModAlt,
	"CTRL":  hotkey.ModCtrl,
	"SHIFT": hotkey.ModShift,
	"WIN":   hotkey.ModWin,
}

var KeyMap = map[string]hotkey.Key{
	"SPACE": hotkey.KeySpace,
	"0":     hotkey.Key0,
	"1":     hotkey.Key1,
	"2":     hotkey.Key2,
	"3":     hotkey.Key3,
	"4":     hotkey.Key4,
	"5":     hotkey.Key5,
	"6":     hotkey.Key6,
	"7":     hotkey.Key7,
	"8":     hotkey.Key8,
	"9":     hotkey.Key9,
	"A":     hotkey.KeyA,
	"B":     hotkey.KeyB,
	"C":     hotkey.KeyC,
	"D":     hotkey.KeyD,
	"E":     hotkey.KeyE,
	"F":     hotkey.KeyF,
	"G":     hotkey.KeyG,
	"H":     hotkey.KeyH,
	"I":     hotkey.KeyI,
	"J":     hotkey.KeyJ,
	"K":     hotkey.KeyK,
	"L":     hotkey.KeyL,
	"M":     hotkey.KeyM,
	"N":     hotkey.KeyN,
	"O":     hotkey.KeyO,
	"P":     hotkey.KeyP,
	"Q":     hotkey.KeyQ,
	"R":     hotkey.KeyR,
	"S":     hotkey.KeyS,
	"T":     hotkey.KeyT,
	"U":     hotkey.KeyU,
	"V":     hotkey.KeyV,
	"W":     hotkey.KeyW,
	"X":     hotkey.KeyX,
	"Y":     hotkey.KeyY,
	"Z":     hotkey.KeyZ,

	"RETURN": hotkey.KeyReturn,
	"ESCAPE": hotkey.KeyEscape,
	"DELETE": hotkey.KeyDelete,
	"TAB":    hotkey.KeyTab,

	"LEFT":  hotkey.KeyLeft,
	"RIGHT": hotkey.KeyRight,
	"UP":    hotkey.KeyUp,
	"DOWN":  hotkey.KeyDown,

	"F1":  hotkey.KeyF1,
	"F2":  hotkey.KeyF2,
	"F3":  hotkey.KeyF3,
	"F4":  hotkey.KeyF4,
	"F5":  hotkey.KeyF5,
	"F6":  hotkey.KeyF6,
	"F7":  hotkey.KeyF7,
	"F8":  hotkey.KeyF8,
	"F9":  hotkey.KeyF9,
	"F10": hotkey.KeyF10,
	"F11": hotkey.KeyF11,
	"F12": hotkey.KeyF12,
	"F13": hotkey.KeyF13,
	"F14": hotkey.KeyF14,
	"F15": hotkey.KeyF15,
	"F16": hotkey.KeyF16,
	"F17": hotkey.KeyF17,
	"F18": hotkey.KeyF18,
	"F19": hotkey.KeyF19,
	"F20": hotkey.KeyF20,
}
