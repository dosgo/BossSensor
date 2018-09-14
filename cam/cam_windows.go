package cam

import "C"
import (
	"syscall"
	"log"
	"unsafe"
	"math/rand"
	"time"
)


var (
	libopencv_videoio341 = syscall.NewLazyDLL("./opencv_world300.dll")
	PcvCreateCameraCapture = libopencv_videoio341.NewProc("cvCreateCameraCapture")
	PcvQueryFrame = libopencv_videoio341.NewProc("cvQueryFrame")
	PcvReleaseCapture = libopencv_videoio341.NewProc("cvReleaseCapture")
	PcvShowImage = libopencv_videoio341.NewProc("cvShowImage")
	PcvDestroyWindow = libopencv_videoio341.NewProc("cvDestroyWindow")
	PcvNamedWindow = libopencv_videoio341.NewProc("cvNamedWindow")
	PcvSaveImgae = libopencv_videoio341.NewProc("cvSaveImage")
	PcvCloneImage = libopencv_videoio341.NewProc("cvCloneImage")
	PcvSetCaptureProperty = libopencv_videoio341.NewProc("cvSetCaptureProperty")
)

var (
	CV_CAP_PROP_DC1394_OFF         = -4 //turn the feature off (not controlled manually nor automatically)
	CV_CAP_PROP_DC1394_MODE_MANUAL = -3 //set automatically when a value of the feature is set by the user
	CV_CAP_PROP_DC1394_MODE_AUTO = -2
	CV_CAP_PROP_DC1394_MODE_ONE_PUSH_AUTO = -1
	CV_CAP_PROP_POS_MSEC       =0
	CV_CAP_PROP_POS_FRAMES     =1
	CV_CAP_PROP_POS_AVI_RATIO  =2
	CV_CAP_PROP_FRAME_WIDTH    =3
	CV_CAP_PROP_FRAME_HEIGHT   =4
	CV_CAP_PROP_FPS            =5
	CV_CAP_PROP_FOURCC         =6
	CV_CAP_PROP_FRAME_COUNT    =7
	CV_CAP_PROP_FORMAT         =8
	CV_CAP_PROP_MODE           =9
	CV_CAP_PROP_BRIGHTNESS    =10
	CV_CAP_PROP_CONTRAST      =11
	CV_CAP_PROP_SATURATION    =12
	CV_CAP_PROP_HUE           =13
	CV_CAP_PROP_GAIN          =14
	CV_CAP_PROP_EXPOSURE      =15
	CV_CAP_PROP_CONVERT_RGB   =16
	CV_CAP_PROP_WHITE_BALANCE_BLUE_U =17
	CV_CAP_PROP_RECTIFICATION =18
	CV_CAP_PROP_MONOCHROME    =19
	CV_CAP_PROP_SHARPNESS     =20
	CV_CAP_PROP_AUTO_EXPOSURE =21
)

/*
*camdev string 在windows下无效为了跟linux通用
*/
func StartCapture(camdev string,fps int,call func(string)) {
	//读取摄像头
	ret,_,err := PcvCreateCameraCapture.Call(uintptr(0));
	if(ret==0){
		log.Println(err)
		return ;
	}
	//释放视频
	defer PcvReleaseCapture.Call(uintptr(unsafe.Pointer(&ret)));
	//设置帧率(测试无效,原因不明)
	//	pcvSetCaptureProperty.Call(uintptr(unsafe.Pointer(ret)),uintptr(CV_CAP_PROP_FPS),uintptr(10))
	for {
		frame,_,err1 := PcvQueryFrame.Call(uintptr(unsafe.Pointer(ret)));
		if(frame==0) {
			log.Println(err1)
		}
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		if(r.Intn(fps)==0) {
			fname := "./img/" + GetRandomString(10) + ".jpg"
			PcvSaveImgae.Call(uintptr(unsafe.Pointer(&[]byte(fname)[0])), uintptr(unsafe.Pointer(frame)), uintptr(0))
			go call(fname);//回调
		}
	}
}


func  GetRandomString(l int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}


