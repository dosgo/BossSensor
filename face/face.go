package face

import (
	"github.com/esimov/pigo/core"
	"io/ioutil"
	"os"
	"image"
	"image/jpeg"
	"github.com/Comdex/imgo"
	"log"
	"image/color"
)

/*提取人脸*/
func CollectFace(facefinder string,fname string) string {
	//读取文件
	src, err := pigo.GetImage(fname)//文件
	if err != nil {
		return "";
	}
	//获取宽高
	pixels := pigo.RgbToGrayscale(src)
	cols, rows := src.Bounds().Max.X, src.Bounds().Max.Y

	//config
	cParams := pigo.CascadeParams{
		MinSize:     20,
		MaxSize:     20000,
		ShiftFactor: 0.1,
		ScaleFactor: 1.1,
	}
	imgParams := pigo.ImageParams{
		Pixels: pixels,
		Rows:   rows,
		Cols:   cols,
		Dim:    cols,
	}
	//load face lib.
	cascadeFile, err := ioutil.ReadFile(facefinder)//人脸库文件
	if err != nil {
		return "";
	}
	pigo := pigo.NewPigo()
	// Unpack the binary file. This will return the number of cascade trees,
	// the tree depth, the threshold and the prediction from tree's leaf nodes.
	classifier, err := pigo.Unpack(cascadeFile)
	if err != nil {
		return "";
	}
	// Run the classifier over the obtained leaf nodes and return the detection results.
	// The result contains quadruplets representing the row, column, scale and detection score.
	faces := classifier.RunCascade(imgParams, cParams)
	// Calculate the intersection over union (IoU) of two clusters.
	faces = classifier.ClusterDetections(faces, 0.18)
	if(len(faces)==0){
		return "";
	}
	//裁剪
	fIn, _ := os.Open(fname)
	origin, _, err := image.Decode(fIn)
	if err != nil {
		return "";
	}
	out, _ := os.Create(fname+".face")
	for _, face := range faces {
		img := origin.(*image.YCbCr)
		var statx,staty,endx,endy int;
		statx=face.Col-face.Scale/2;
		staty=face.Row-face.Scale/2;
		endx=statx+face.Scale;
		endy=staty+face.Scale;
		subImg := img.SubImage(image.Rect(statx, staty, endx,  endy)).(*image.YCbCr)
		jpeg.Encode(out, subImg, &jpeg.Options{100})
	}
	//换成灰度图
	ToGrayScale(fname+".face");
	return fname+".face";
}

/*对比人脸*/
func DiffFace(srcface string,fnameface string){
	//人脸对比
	cos,err:=imgo.CosineSimilarity(srcface,fnameface)
	if err!=nil{
		log.Println("diff error")
	}
	if(cos>0.88){

	}
}


/*转灰度图*/
func ToGrayScale(fname string) bool{
	file, err :=  os.Open(fname)
	if err != nil {
		return false;
	}
	m, _, err := image.Decode(file);
	if err != nil {
		return false;
	}
	file.Close();//关闭
	bounds := m.Bounds()
	mgrey := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.At(x, y).RGBA()
			r = r >> 8
			g = g >> 8
			b = b >> 8
			v := uint8((float32(r)*299 + float32(g)*587 + float32(b)*114) / 1000)
			mgrey.Set(x, y, color.Gray{v})
		}
	}
	file1, err :=  os.Create(fname)
	if err != nil {
		return false;
	}
	jpeg.Encode(file1, mgrey, &jpeg.Options{100})
	defer file1.Close()
	return true;
}