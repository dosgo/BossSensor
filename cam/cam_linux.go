package cam


//相机格式参数
const (
	V4L2_PIX_FMT_PJPG = 0x47504A50
	V4L2_PIX_FMT_MJPEG = 0x47504A4D
	V4L2_PIX_FMT_YUYV = 0x56595559
)
var supportedFormats = map[webcam.PixelFormat]bool{
	V4L2_PIX_FMT_PJPG: true,
	V4L2_PIX_FMT_YUYV: true,
}

func toImage( format webcam.PixelFormat, b []byte, w, h int) string {
	var ratio image.YCbCrSubsampleRatio
	switch format {
	case V4L2_PIX_FMT_YUYV:
		ratio = image.YCbCrSubsampleRatio422
	}
	pixs := len(b) / 4 * 2
	img := image.NewYCbCr(image.Rect(0, 0, w, h), ratio)
	img.YStride = w
	img.CStride = w / 2
	img.Y = make([]byte, pixs)
	img.Cb = make([]byte, pixs/2)
	img.Cr = make([]byte, pixs/2)

	for i := 0; i < pixs/2; i++ {
		img.Y[i*2] = b[i*4]
		img.Cb[i] = b[i*4+1]
		img.Y[i*2+1] = b[i*4+2]
		img.Cr[i] = b[i*4+3]
	}
	fname:="./img/"+GetRandomString(10)+".jpg"
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, nil)
	if err != nil {
		fmt.Println(err)
	}
	ioutil.WriteFile(fname, buf.Bytes(), 0644)
	return fname
}

func StartCapture(camdev string,fps int,call func(string)){
	cam, err := webcam.Open(camdev) // Open webcam
	if err != nil {
		log.Printf("Can't find camera %s\r\n",camdev)
		return ;
	}
	defer cam.Close()
	format_desc := cam.GetSupportedFormats()
	//获取相机支持的格式
	var format webcam.PixelFormat

	for f, _:= range format_desc {
		//默认第一种格式
		if(f==V4L2_PIX_FMT_YUYV){
			log.Println("...")
			format=V4L2_PIX_FMT_YUYV
		}
	}
	//没有找到格式
	if(format==0){
		return ;
	}
	//支持的分辨率
	frames := cam.GetSupportedFrameSizes(format)
	var size *webcam.FrameSize
	//默认第一种分辨率
	size=&frames[0];
	log.Println(frames)
	//设置相机格式跟分辨率
	f, w, h, err := cam.SetImageFormat(format, uint32(size.MaxWidth), uint32(size.MaxHeight))
	if err != nil {
		log.Println("SetImageFormat return error", err)
		return
	}

	err = cam.StartStreaming()
	if err != nil { panic(err.Error()) }
	for {
		err = cam.WaitForFrame(2000)
		switch err.(type) {
		case nil:
		case *webcam.Timeout:
			continue
		default:
			log.Println("err")
		}
		frame, err := cam.ReadFrame()
		if len(frame) != 0 && err==nil {
			//不要全采样，丢掉一些帧
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			if(r.Intn(fps)==0) {
				fname:=toImage(f,frame,int(w),int(h));
				go call(fname);//回调
			}

		}
	}
}
