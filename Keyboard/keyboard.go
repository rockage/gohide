package Keyboard

/*
#include <stdio.h>
#include <windows.h>
#include <locale.h>

HWND getWindowName(WCHAR* buffer, int bufferSize) {
	HWND hwndCurrent;
    setlocale(LC_ALL, "chs"); // 设置字符集为简体中文
    WCHAR windowName[512]; // 用于存储窗口标题
    // 获取当前活动窗口的句柄
    hwndCurrent = GetForegroundWindow();
    // 获取窗口标题
    GetWindowTextW(hwndCurrent, windowName, sizeof(windowName) / sizeof(windowName[0]));
	// 将窗口标题拷贝到输出缓冲区 （通过指针 buffer[0] 回传给GO，非直接return)
    wcsncpy(buffer, windowName, bufferSize);
    buffer[bufferSize - 1] = L'\0';  // 确保字符串以 null 结尾
	return hwndCurrent;
}

void procShowWindow(HWND hwnd){
   BOOL isVisible = IsWindowVisible(hwnd); // 窗口是否隐藏？
    if (isVisible) {
    	ShowWindow(hwnd, SW_HIDE);
    } else {

	    ShowWindow(hwnd, SW_SHOW);
    }
}

void setTopMost(HWND hwnd){
	SetWindowPos(hwnd, HWND_TOPMOST, 0, 0, 0, 0, SWP_NOSIZE | SWP_NOMOVE);
}

void setNoTopMost(HWND hwnd){
	SetWindowPos(hwnd, HWND_NOTOPMOST, 0, 0, 0, 0, SWP_NOSIZE | SWP_NOMOVE);
}
*/
import "C" // C 代码结束

import (
	"gohide/Global"
	"reflect"
	"strconv"
	"syscall"
	"time"
	"unicode"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32                  = windows.NewLazySystemDLL("user32.dll")
	procSetWindowsHookEx    = user32.NewProc("SetWindowsHookExW")
	procCallNextHookEx      = user32.NewProc("CallNextHookEx")
	procUnhookWindowsHookEx = user32.NewProc("UnhookWindowsHookEx")
	procGetMessage          = user32.NewProc("GetMessageW")
	procTranslateMessage    = user32.NewProc("TranslateMessage")
	procDispatchMessage     = user32.NewProc("DispatchMessageW")
	procMapVirtualKey       = user32.NewProc("MapVirtualKeyW")
	procGetKeyNameText      = user32.NewProc("GetKeyNameTextW")
	procGetKeyState         = user32.NewProc("GetKeyState")
	procGetAsyncKeyState    = user32.NewProc("GetAsyncKeyState")
	procShowWindow          = user32.NewProc("ShowWindow")
	procFindWindowW         = user32.NewProc("FindWindowW")
	procIsWindowVisible     = user32.NewProc("IsWindowVisible")
	procGetForegroundWindow = user32.NewProc("GetForegroundWindow")
	procGetWindowTextW      = user32.NewProc("GetWindowTextW")
	procSetWindowLong       = user32.NewProc("SetWindowLongW")
	procGetWindowLong       = user32.NewProc("GetWindowLong")
	hHook                   windows.Handle
	HwndSelf                unsafe.Pointer
	HwndCurrent             unsafe.Pointer
	AllKeyData              Global.Data
)

const (
	GWL_STYLE        = -16
	GWL_EXSTYLE      = -20
	WS_EX_TOOLWINDOW = 0x00000080
	WH_KEYBOARD_LL   = 13
	HC_ACTION        = 0
	WM_KEYDOWN       = 0x0100
	WM_KEYUP         = 0x0101
	WM_SYSKEYDOWN    = 0x0104
	WM_SYSKEYUP      = 0x0105
	VK_LCONTROL      = 0xA2
	VK_LMENU         = 0xA4
	VK_LSHIFT        = 0xA0
	VK_RCONTROL      = 0xA3
	VK_RMENU         = 0xA5
	VK_RSHIFT        = 0xA1
	MAPVK_VK_TO_VSC  = 0
)

type KBDLLHOOKSTRUCT struct {
	VkCode      uint32
	ScanCode    uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

type (
	DWORD     uint32
	WPARAM    uintptr
	LPARAM    uintptr
	LRESULT   uintptr
	HANDLE    uintptr
	HINSTANCE HANDLE
	HHOOK     HANDLE
	HWND      HANDLE
)

type POINT struct {
	X, Y int32
}

type MSG struct {
	Hwnd    HWND
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      POINT
}

// 定义回调函数类型
type CallbackFunc func(Global.KeyData)
type CallbackWindowName func(string)

// 向主界面传递消息的两个回调函数
var Callback1 CallbackFunc
var Callback2 CallbackWindowName

func getKeyName(vkCode int) string {
	// 手动处理常用修饰键的键名
	switch vkCode {
	case VK_LCONTROL:
		return "Ctrl"
	case VK_LMENU:
		return "Alt"
	case VK_LSHIFT:
		return "Shift"
	case VK_RCONTROL:
		return "Ctrl"
	case VK_RMENU:
		return "Alt"
	case VK_RSHIFT:
		return "Shift"
	}

	scanCode, _, _ := procMapVirtualKey.Call(uintptr(vkCode), MAPVK_VK_TO_VSC)
	scanCode = scanCode << 16
	if vkCode == VK_LCONTROL || vkCode == VK_RCONTROL || vkCode == VK_LMENU || vkCode == VK_RMENU || vkCode == VK_LSHIFT || vkCode == VK_RSHIFT {
		scanCode |= 0x01000000
	}

	var name [128]uint16
	ret, _, _ := procGetKeyNameText.Call(scanCode, uintptr(unsafe.Pointer(&name[0])), uintptr(len(name)))
	if ret > 0 {
		return syscall.UTF16ToString(name[:])
	}
	return "Unknown Key"
}

var con_keys = map[uint32]bool{
	VK_LCONTROL: false,
	VK_LMENU:    false,
	VK_LSHIFT:   false,
	VK_RCONTROL: false,
	VK_RMENU:    false,
	VK_RSHIFT:   false,
}

func GetKeyState(vk int) bool {
	ret, _, _ := procGetAsyncKeyState.Call(uintptr(vk))
	return ret&0x8000 != 0
}

func KeyboardProc(nCode int32, wParam uintptr, lParam uintptr) uintptr {
	if nCode == HC_ACTION {
		kbdStruct := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam)) //unsafe

		/*
			if wParam == WM_KEYDOWN || wParam == WM_SYSKEYDOWN { // 键盘按下事件：
					if int(kbdStruct.VkCode) == 27 {
						fmt.Println("Press ESC!")
						ReleaseHook()
						os.Exit(0)
					}

			}
		*/

		if wParam == WM_KEYUP || wParam == WM_SYSKEYUP { // 按键抬起事件
			var ret uintptr
			if kbdStruct.VkCode == VK_LCONTROL || // 控制按键
				kbdStruct.VkCode == VK_LMENU ||
				kbdStruct.VkCode == VK_LSHIFT ||
				kbdStruct.VkCode == VK_RCONTROL ||
				kbdStruct.VkCode == VK_RMENU ||
				kbdStruct.VkCode == VK_RSHIFT {
				for key := range con_keys { // 清空控制按键状态位
					con_keys[key] = false
				}

			} else { // 普通按键
				var press_keyData Global.KeyData
				var insert Global.CodeName

				for key := range con_keys {
					ret, _, _ = procGetKeyState.Call(uintptr(key))
					r := ret&0x8000 != 0
					if r {
						insert.Code = int(key)
						insert.Name = getKeyName(int(key))
						press_keyData.Cons = append(press_keyData.Cons, insert)
					}
				}
				press_keyData.Key.Code = int(kbdStruct.VkCode)
				press_keyData.Key.Name = getKeyName(int(kbdStruct.VkCode))
				Callback1(press_keyData) // 主界面按键信息回调
				// K1 显示/隐藏窗口：
				if press_keyData.Key.Code == AllKeyData.K1.Key.Code {
					if reflect.DeepEqual(press_keyData.Cons, AllKeyData.K1.Cons) { // CONS值比较\
						if HwndCurrent != nil {
							C.procShowWindow(C.HWND(HwndCurrent))
						}
					}
				}
				// K2 定义当前窗口为隐藏窗口：
				if press_keyData.Key.Code == AllKeyData.K2.Key.Code {
					if reflect.DeepEqual(press_keyData.Cons, AllKeyData.K2.Cons) {
						const bufferSize = 512
						buffer := make([]C.WCHAR, bufferSize)
						// 调用 C 函数
						hc := C.getWindowName((*C.WCHAR)(unsafe.Pointer(&buffer[0])), C.int(bufferSize))
						HwndCurrent = unsafe.Pointer(hc)
						// 将 C.WCHAR数组:buffer 转换为 Go 字符串 ** C 使用utf16编码，Go使用utf8 **
						utf8String := utf16ToUtf8(buffer)
						hwndStr := strconv.Itoa(int(syscall.Handle(HwndCurrent))) // Hwnd转换为String
						Callback2(utf8String + "[" + hwndStr + "]")               // 主界面窗口标题回调
					}
				}
				// K3: 显示或隐藏程序自身
				if press_keyData.Key.Code == AllKeyData.K3.Key.Code {
					if reflect.DeepEqual(press_keyData.Cons, AllKeyData.K3.Cons) {
						C.procShowWindow(C.HWND(HwndSelf))
					}
				}
			}

		}
	}
	ret, _, _ := procCallNextHookEx.Call(uintptr(hHook), uintptr(nCode), wParam, lParam)
	return ret
}

func utf16ToUtf8(buffer []C.WCHAR) string {

	// 中英文字符宽度不一致将导致UI伸缩
	isChinese := func(char rune) bool {
		return unicode.Is(unicode.Han, char)
	}
	// 将UTF-16编码的字节序列解码为Unicode码点序列
	var utf16Chars []uint16
	for i := 0; i < len(buffer)-1; i++ {
		utf16Chars = append(utf16Chars, uint16(buffer[i]))
		if buffer[i] == 0 {
			break
		}
	}
	runes := utf16.Decode(utf16Chars)
	// 将Unicode码点序列转换为UTF-8编码的字符串
	var count int = 0
	var warpStr string
	for i := 0; i < len(runes)-1; i++ {
		if isChinese(runes[i]) {
			count += 2
		} else {
			count++
		}
		if count >= 13 { // 超过6字符自动换行
			warpStr += "\n"
			count = 0
		}
		warpStr += string(runes[i])

	}
	return warpStr
}

func SetHook() {
	procSetWindowsHookEx.Call(
		WH_KEYBOARD_LL,
		syscall.NewCallback(KeyboardProc),
		0,
		0,
	)
}

func ReleaseHook() {
	procUnhookWindowsHookEx.Call(uintptr(hHook))
}

func MsgLoop() {
	var msg MSG
	for {
		procGetMessage.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
		procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		procDispatchMessage.Call(uintptr(unsafe.Pointer(&msg)))
	}
}

func SetTopMost() {
	C.setTopMost(C.HWND(HwndSelf)) // 将本程序设为顶层
}

func SetNoTopMost() {
	C.setNoTopMost(C.HWND(HwndSelf)) // 将本程序设为顶层
}

func Start() {
	HwndSelf = getGoHideHwnd()

	if hHook != 0 {
		SetHook()
	} else {
		ReleaseHook()
		SetHook()
	}

	defer ReleaseHook()
	MsgLoop()
}

func getGoHideHwnd() unsafe.Pointer {
	// 指定要查找的窗口标题
	windowTitle := "GoHide 摸鱼伴侣@rockage"
	var hwnd uintptr
	// 将窗口标题转换为 LPWSTR 类型
	titlePtr, _ := syscall.UTF16PtrFromString(windowTitle)
	startTime := time.Now()
	for time.Since(startTime) < 10*time.Second {
		hwnd, _, _ = procFindWindowW.Call(uintptr(0), uintptr(unsafe.Pointer(titlePtr)))
		if hwnd != 0 {
			break
		}
		time.Sleep(500 * time.Millisecond) // 每次循环等待500毫秒
	}
	hs := unsafe.Pointer(hwnd)
	C.setTopMost(C.HWND(hs)) // 将本程序设为顶层
	return hs
}
