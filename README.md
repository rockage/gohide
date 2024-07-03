#### 用途
- 用热键隐藏你不想给别人看到的窗口，让它置于后台并在Windows任务栏消失；
- 再次按下热键，窗口恢复；
- 是上班摸鱼、看小说、看视频的好伴侣。
#### 结构
- 功能：GO（with CGO） + Win32 API
- GUI框架：Fyne V2.0
#### 编译
- 首先安装Fyne，而安装它又必须打开CGO和安装一个C编译器（本程序使用MINGW64），具体请看Fyne官网说明：https://docs.fyne.io/started/
- 编译指令：go build -ldflags="-H windowsgui" -o gohide.exe (程序可用但没有ICON)
- 将ICON嵌入EXE: fyne package -icon icon.png
#### 使用方法举例
- 以WIN自带的记事本为例，首先打开Windows 记事本
- 将Windows记事本置于当前窗口
- 按下ALT + 2 键 （记事本已被设置为当前需要隐藏的窗口）
- 按下ALT + 1 键 （记事本被隐藏）
- 再次按下 ALT + 1键 （记事本还原）
- 按下 ALT+3 键 （本程序自身被隐藏，再次按下则恢复）
- ALT+1 ALT+2 ALT+3 为本程序默认按键，此3个按键可自行设置
#### 后记
- 本程序大量使用WIN32 API，因此只能运行在Windows平台下，无法跨平台。
- 本程序只对传统Win传统窗口程序生效，对微软UWP程序无效，例：Win10自带的“记事本”是一个传统Win窗口程序是有效的，而Win10自带的“计算器”则是一个UWP程序，本软件对其无效！
- 本程序对部分屏蔽了全局钩子的游戏软件无效(例如：三国志-战略版 for PC)，下一个版本我会考虑使用"驱动级"全局钩子应对，敬请期待。

#### What is it?
Use hotkeys to hide windows you don't want others to see, placing them in the background and making them disappear from the Windows taskbar.
Press the hotkey again to restore the window.
It's a great companion for slacking off at work, reading novels, or watching videos.
#### Languages
Main: GO (with CGO) + Win32 API
GUI Framework: Fyne V2.0
#### Build
First, install Fyne. To install Fyne, you need to enable CGO and install a C compiler (such as MINGW64). For details, see the Fyne documentation: https://docs.fyne.io/started/
Compile command: go build -ldflags="-H windowsgui" -o gohide.exe (program is usable but without an icon).
Embed the ICON into the EXE: fyne package -icon icon.png
#### Usage Example
Take the built-in Windows Notepad as an example. First, open Windows Notepad.
Bring the Windows Notepad window to the forewindow.
Press ALT + 2 (Notepad is set as the current window to be hidden).
Press ALT + 1 (Notepad is hidden).
Press ALT + 1 again (Notepad is restored).
Press ALT + 3 (This program itself is hidden; press again to restore).
ALT + 1, ALT + 2, and ALT + 3 are the default keys for this program, and these three keys can be customized.
#### Note
This program uses Win32 API, so it can only run on Windows platforms and is not cross-platform.
This program only works on traditional Windows window programs and does not work on Microsoft UWP programs. For example, the built-in "Notepad" in Windows 10 is a traditional Windows window program, so it works. However, the built-in "Calculator" is a UWP program, so this software does not work on it!
This program does not work on some games that block global hooks (e.g., "三国志 - 战略版" for PC). In the next version, I will consider using "driver-level" global hooks to address this issue. Stay tuned.
Feel free to customize the repository URL and other details as needed.
