package CapVideo
//废弃，之前在windows下直接调用摄像头，后面改成opencv
import (
	"syscall"
	"unsafe"
)

const  cWM_USER  = 0x400
const  cWS_CHILD = 0x40000000
const  cWS_VISIBLE = 0x10000000
const  cSWP_NOMOVE = 0x2
const  cSWP_NOZORDER = 0x4
const  cWM_CAP_DRIVER_CONNECT = cWM_USER + 10
const  cWM_CAP_DRIVER_DISCONNECT = cWM_USER + 11
const  cWM_CAP_SET_CALLBACK_FRAME = cWM_USER + 5
const  cWM_CAP_SET_PREVIEW = cWM_USER + 50
const  cWM_CAP_SET_PREVIEWRATE = cWM_USER + 52
const  cWM_CAP_SET_VIDEOFORMAT = cWM_USER + 45
const  cWM_CAP_GRAB_FRAME_NOSTOP = cWM_USER + 61
const  cWM_CAP_GRAB_FRAME = cWM_USER + 60
const  cWM_CAP_FILE_SAVEDIBA = cWM_USER + 25


var (
	avicap32 = syscall.NewLazyDLL("Avicap32.dll")
	pcapGetDriverDescriptionA = avicap32.NewProc("capGetDriverDescriptionA")
	pcapCreateCaptureWindowA = avicap32.NewProc("capCreateCaptureWindowA")
	pSendMessage = avicap32.NewProc("SendMessage")
)

  func CapGetDriverDescriptionA(lpszName *[100]byte,lpszVer *[100]byte)(bool,error) {
	  ret, _, callErr :=pcapGetDriverDescriptionA.Call(0,uintptr(unsafe.Pointer(lpszName)),100,uintptr(unsafe.Pointer(lpszVer)),100)
	  if(ret==0){
		 return  false,callErr;
	  }
	  return true,nil;
  }

  func CapCreateCaptureWindowA(lpszName *[100]byte,wnd int)(int,error) {
	ret, _, callErr := pcapCreateCaptureWindowA.Call(uintptr(unsafe.Pointer(lpszName)),cWS_VISIBLE + cWS_CHILD,0,0,1,1,uintptr(wnd),0)
	if(ret==0){
		return  0,callErr;
	}
	return int(ret),nil;
}

func SendMessage( hWnd int,  wMsg int, wParam int ,  lParam int) bool{
	ret, _, _ :=pSendMessage.Call(uintptr(hWnd),uintptr(wMsg),uintptr(wParam),uintptr(lParam));
	if(ret==0){
		return  false;
	}
	return true;
}


//连接设备
func  CapDriverConnect( lwnd int, i int) bool {
	return SendMessage(lwnd, cWM_CAP_DRIVER_CONNECT, i, 0);
}

//断开连接
func  CapDriverDisconnect(lwnd int) bool {
	return SendMessage(lwnd, cWM_CAP_DRIVER_DISCONNECT, 0, 0);
}

//设置为预览模式
func  CapPreview( lwnd int,  f int) bool{
	return SendMessage(lwnd, cWM_CAP_SET_PREVIEW, f, 0);
}

//设置预览帧率
func  CapPreviewRate( lwnd int, wMS int) bool{
  	return SendMessage(lwnd, cWM_CAP_SET_PREVIEWRATE, wMS, 0);
}

//设置视频格式
func  CapSetVideoFormat( hCapWnd int, BmpFormat int ,  CapFormatSize int) bool {
	return SendMessage(hCapWnd, cWM_CAP_SET_VIDEOFORMAT, CapFormatSize, BmpFormat);
}

/*

public class CapVideo
{
#region DLL Import Method
[DllImport("avicap32.dll")]
public static extern IntPtr capCreateCaptureWindowA(byte[] lpszWindowName, int dwStyle, int x, int y, int nWidth, int nHeight, IntPtr hWndParent, int nID);
[DllImport("avicap32.dll")]
public static extern bool capGetDriverDescriptionA(short wDriver, byte[] lpszName, int cbName, byte[] lpszVer, int cbVer);
[DllImport("avicap32.dll")]
public static extern int capGetVideoFormat(IntPtr hWnd, IntPtr psVideoFormat, int wSize);
[DllImport("User32.dll")]
public static extern bool SendMessage(IntPtr hWnd, int wMsg, bool wParam, int lParam);
[DllImport("User32.dll")]
public static extern bool SendMessage(IntPtr hWnd, int wMsg, short wParam, int lParam);
[DllImport("User32.dll")]
public static extern bool SendMessage(IntPtr hWnd, int wMsg, short wParam, FrameEventHandler lParam);
[DllImport("User32.dll")]
public static extern bool SendMessage(IntPtr hWnd, int wMsg, int wParam, ref BITMAPINFO lParam);
[DllImport("User32.dll")]
public static extern int SetWindowPos(IntPtr hWnd, int hWndInsertAfter, int x, int y, int cx, int cy, int wFlags);
#endregion

}*/