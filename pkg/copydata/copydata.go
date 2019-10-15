package copydata

import (
	"encoding/json"
	"fmt"
	"github.com/lxn/win"
	"reflect"
	"syscall"
	"unicode/utf16"
	"unsafe"
)

func SendMessage(hWndSource, hWndTarget win.HWND, wParam uintptr, b []byte) uintptr {
	header := *(*reflect.SliceHeader)(unsafe.Pointer(&b))
	cd := copyData{
		CbData: uint32(header.Len),
		LpData: header.Data,
		DwData: uintptr(hWndSource),
	}
	return win.SendMessage(hWndTarget, win.WM_COPYDATA, wParam, uintptr(unsafe.Pointer(&cd)))
}

type WndClass struct {
	Src, Dest string
}

func (x WndClass) SendString(msg uintptr, s string) bool {
	return x.sendMsg(msg, utf16FromString(s))
}

func (x WndClass) SendJson(msg uintptr, param interface{}) bool {
	b, err := json.Marshal(param)
	if err != nil {
		panic(err)
	}
	return x.SendString(msg, string(b))
}

func (x WndClass) sendMsg(msg uintptr, b []byte) bool {
	hWndSrc := findWindowClass(x.Src)
	hWndDest := findWindowClass(x.Dest)
	if hWndSrc == 0 || hWndDest == 0 {
		return false
	}
	return SendMessage(hWndSrc, hWndDest, msg, b) != 0
}

func findWindowClass(className string) win.HWND {
	ptrClassName := mustUTF16PtrFromString(className)
	return win.FindWindow(ptrClassName, nil)
}

type copyData struct {
	DwData uintptr
	CbData uint32
	LpData uintptr
}

//func GetData(ptr unsafe.Pointer) (uintptr, []byte) {
//	cd := (*CopyData)(ptr)
//	p := PtrSliceFrom(unsafe.Pointer(cd.LpData), int(cd.CbData))
//	return cd.DwData, *(*[]byte)(p)
//}

//func PtrSliceFrom(p unsafe.Pointer, s int) unsafe.Pointer {
//	return unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(p), Len: s, Cap: s})
//}

func utf16FromString(s string) (b []byte) {
	for i := 0; i < len(s); i++ {
		if s[i] == 0 {
			panic(fmt.Sprintf("%q[%d] is 0", s, i))
		}
	}
	for _, v := range utf16.Encode([]rune(s)) {
		b = append(b, byte(v), byte(v>>8))
	}
	return
}

func mustUTF16PtrFromString(s string) *uint16 {
	p, err := syscall.UTF16PtrFromString(s)
	if err != nil {
		panic(err)
	}
	return p
}
