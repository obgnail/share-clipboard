// +build windows

package share_clipboard

import (
	"golang.design/x/hotkey"
	"runtime"
	"syscall"
	"time"
	"unsafe"
)

const (
	cfUnicodetext = 13
	gmemMoveable  = 0x0002
)

var (
	user32                     = syscall.MustLoadDLL("user32")
	isClipboardFormatAvailable = user32.MustFindProc("IsClipboardFormatAvailable")
	openClipboard              = user32.MustFindProc("OpenClipboard")
	closeClipboard             = user32.MustFindProc("CloseClipboard")
	emptyClipboard             = user32.MustFindProc("EmptyClipboard")
	getClipboardData           = user32.MustFindProc("GetClipboardData")
	setClipboardData           = user32.MustFindProc("SetClipboardData")

	kernel32     = syscall.NewLazyDLL("kernel32")
	globalAlloc  = kernel32.NewProc("GlobalAlloc")
	globalFree   = kernel32.NewProc("GlobalFree")
	globalLock   = kernel32.NewProc("GlobalLock")
	globalUnlock = kernel32.NewProc("GlobalUnlock")
	lstrcpy      = kernel32.NewProc("lstrcpyW")
)

// waitOpenClipboard opens the clipboard, waiting for up to a second to do so.
func waitOpenClipboard() error {
	started := time.Now()
	limit := started.Add(time.Second)
	var r uintptr
	var err error
	for time.Now().Before(limit) {
		r, _, err = openClipboard.Call(0)
		if r != 0 {
			return nil
		}
		time.Sleep(time.Millisecond)
	}
	return err
}

func GetTextFromClip() (string, error) {
	// LockOSThread ensure that the whole method will keep executing on the same thread from begin to end (it actually locks the goroutine thread attribution).
	// Otherwise if the goroutine switch thread during execution (which is a common practice), the OpenClipboard and CloseClipboard will happen on two different threads, and it will result in a clipboard deadlock.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	if formatAvailable, _, err := isClipboardFormatAvailable.Call(cfUnicodetext); formatAvailable == 0 {
		return "", err
	}
	err := waitOpenClipboard()
	if err != nil {
		return "", err
	}

	h, _, err := getClipboardData.Call(cfUnicodetext)
	if h == 0 {
		_, _, _ = closeClipboard.Call()
		return "", err
	}

	l, _, err := globalLock.Call(h)
	if l == 0 {
		_, _, _ = closeClipboard.Call()
		return "", err
	}

	text := syscall.UTF16ToString((*[1 << 20]uint16)(unsafe.Pointer(l))[:])

	r, _, err := globalUnlock.Call(h)
	if r == 0 {
		_, _, _ = closeClipboard.Call()
		return "", err
	}

	closed, _, err := closeClipboard.Call()
	if closed == 0 {
		return "", err
	}
	return text, nil
}

func GetBytesFromClip() ([]byte, error) {
	res, err := GetTextFromClip()
	if err != nil {
		return nil, err
	}
	return ToBytes(res), nil
}

func SetClipText(text string) error {
	// LockOSThread ensure that the whole method will keep executing on the same thread from begin to end (it actually locks the goroutine thread attribution).
	// Otherwise if the goroutine switch thread during execution (which is a common practice), the OpenClipboard and CloseClipboard will happen on two different threads, and it will result in a clipboard deadlock.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	err := waitOpenClipboard()
	if err != nil {
		return err
	}

	r, _, err := emptyClipboard.Call(0)
	if r == 0 {
		_, _, _ = closeClipboard.Call()
		return err
	}

	data := syscall.StringToUTF16(text)

	// "If the hMem parameter identifies a memory object, the object must have
	// been allocated using the function with the GMEM_MOVEABLE flag."
	h, _, err := globalAlloc.Call(gmemMoveable, uintptr(len(data)*int(unsafe.Sizeof(data[0]))))
	if h == 0 {
		_, _, _ = closeClipboard.Call()
		return err
	}
	defer func() {
		if h != 0 {
			globalFree.Call(h)
		}
	}()

	l, _, err := globalLock.Call(h)
	if l == 0 {
		_, _, _ = closeClipboard.Call()
		return err
	}

	r, _, err = lstrcpy.Call(l, uintptr(unsafe.Pointer(&data[0])))
	if r == 0 {
		_, _, _ = closeClipboard.Call()
		return err
	}

	r, _, err = globalUnlock.Call(h)
	if r == 0 {
		if err.(syscall.Errno) != 0 {
			_, _, _ = closeClipboard.Call()
			return err
		}
	}

	r, _, err = setClipboardData.Call(cfUnicodetext, h)
	if r == 0 {
		_, _, _ = closeClipboard.Call()
		return err
	}
	h = 0 // suppress deferred cleanup
	closed, _, err := closeClipboard.Call()
	if closed == 0 {
		return err
	}
	return nil
}

func SetClipBytes(b []byte) error {
	return SetClipText(ToString(b))
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
