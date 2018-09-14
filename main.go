package main

import "C"
import (
	_ "image/jpeg"
	"github.com/widuu/goini"
	"log"
	"BossSensor/cam"
	"BossSensor/face"
)

type FaceDetctor struct {
	cascadeFile  string
	minSize      int
	maxSize      int
	shiftFactor  float64
	scaleFactor  float64
	iouThreshold float64
}

//保存要识别的头像库
var srcface string;
//头像库文件
var facefinder ="./facefinder";
//摄像头设备
var camdev="/dev/video0";




func captureCall(fname string){
	//提前人脸
	src1face:=face.CollectFace(facefinder,fname);
	//人脸对比
	face.DiffFace(srcface,src1face);
}

func main() {
	if(srcface==""){
		log.Printf("srcface error\r\n")
		return ;
	}
	//启动录像
	cam.StartCapture(camdev,25,captureCall);
}

/*初始化*/
func init(){
	//读取配置文件
	conf := goini.SetConfig("./conf.ini")
	confarr := conf.ReadList()
	var tmpface string;
	for index := 0; index < len(confarr); index++ {
		confmap := confarr[index]
		if _, ok := confmap["config"]; ok {
			//原始的头像
			if _, ok := confmap["config"]["srcface"]; ok {
				tmpface = confmap["config"]["srcface"]
			}
			//人脸库文件
			if _, ok := confmap["config"]["facefinder"]; ok {
				facefinder = confmap["config"]["facefinder"]
			}
			//摄像头
			if _, ok := confmap["config"]["camdev"]; ok {
				camdev = confmap["config"]["camdev"]
			}
		}
	}
	//初始化人脸库
	srcface=face.CollectFace(facefinder,tmpface);
}



