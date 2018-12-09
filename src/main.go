package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"

	"github.com/robertkrimen/otto"
	"github.com/urfave/cli"
)

// ColorPoint 色点
type ColorPoint struct {
	X int
	Y int
	R uint32
	G uint32
	B uint32
	A uint32
}

// ImageRGBAMap 图片map
type ImageRGBAMap struct {
	Width  int
	Height int
	Data   []ColorPoint
}

// Middleware 定义中间件
type Middleware func(imageMap ImageRGBAMap) (outImageMap ImageRGBAMap, err error)

func main() {
	var err error
	var DeCode, EnCode bool
	var Out, Ipt, JavaScript string
	//实例化一个命令行程序
	app := cli.NewApp()
	//程序名称
	app.Name = "PicTool"
	//程序的用途描述
	app.Usage = "解析JPEG或PNG图片为JSON/还原JSON到PNG"
	//程序的版本号
	app.Version = "0.0.1"
	//设置可接受的参数
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "decode",
			Usage:       "从图片解码为JSON",
			Destination: &DeCode,
		},
		cli.BoolFlag{
			Name:        "encode",
			Usage:       "将JSON编码为图片",
			Destination: &EnCode,
		},
		cli.StringFlag{
			Name:        "output",
			Value:       "outfile",
			Usage:       "输出文件名称",
			Destination: &Out,
		},
		cli.StringFlag{
			Name:        "input",
			Usage:       "输入文件路径",
			Destination: &Ipt,
		},

		cli.StringFlag{
			Name:        "javascript",
			Usage:       "处理中间脚本路径",
			Destination: &JavaScript,
		},
	}
	//该程序执行的代码
	app.Action = func(c *cli.Context) error {

		if DeCode && EnCode {
			return errors.New("-decode -encode 参数同时存在")
		}

		if !DeCode && !EnCode {
			return errors.New("-decode -encode 至少输入一项")
		}

		if Ipt == "" {
			return errors.New("-input 参数错误")
		}

		if DeCode {
			img, err := ReadImageFile(Ipt)
			if err != nil {
				return err
			}
			imgMap := DeCodeRGBA(img)
			if JavaScript != "" {
				imgMap, err = JavaScriptMiddleware(imgMap, JavaScript)
				if err != nil {
					return err
				}
			}
			err = WriteJSONFile(imgMap, Out)
			if err != nil {
				return err
			}
			return nil
		}

		if EnCode {
			imgMap, err := ReadJSONFile(Ipt)
			if err != nil {
				return err
			}
			if JavaScript != "" {
				imgMap, err = JavaScriptMiddleware(imgMap, JavaScript)
				if err != nil {
					return err
				}
			}
			img, err := EnCodeRGBA(imgMap)
			if err != nil {
				return err
			}
			err = WriteImageFile(img, Out)
			if err != nil {
				return err
			}
		}

		return nil
	}
	//启动应用
	if err = app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}

// DeCodeRGBA 读取RGBA
func DeCodeRGBA(img image.Image) (imageMap ImageRGBAMap) {
	fmt.Println("读取图片RGBA中....")
	imageMap = ImageRGBAMap{}
	rectangle := img.Bounds()
	imageMap.Width = rectangle.Max.X
	imageMap.Height = rectangle.Max.Y
	for yindex := rectangle.Min.Y; yindex < rectangle.Max.Y; yindex++ {
		for xindex := rectangle.Min.X; xindex < rectangle.Max.X; xindex++ {
			r, g, b, a := img.At(xindex, yindex).RGBA()
			imageMap.Data = append(imageMap.Data, ColorPoint{
				X: xindex,
				Y: yindex,
				R: r >> 8,
				G: g >> 8,
				B: b >> 8,
				A: a >> 8,
			})
		}
	}
	fmt.Println("读取图片RGBA完成")
	return imageMap
}

// EnCodeRGBA 编码RGBA
func EnCodeRGBA(imgMap ImageRGBAMap) (img image.Image, err error) {
	fmt.Println("将JSON编码为图片中....")
	imgRectangle := image.Rectangle{
		Min: image.Point{
			X: 0,
			Y: 0,
		},
		Max: image.Point{
			X: imgMap.Width,
			Y: imgMap.Height,
		},
	}
	imgRGBA := image.NewRGBA(imgRectangle)
	for _, colorPoint := range imgMap.Data {
		imgRGBA.Set(colorPoint.X, colorPoint.Y, color.RGBA{uint8(colorPoint.R), uint8(colorPoint.G), uint8(colorPoint.B), uint8(colorPoint.A)})
	}
	fmt.Println("将JSON编码为图片完成")
	return imgRGBA.SubImage(imgRectangle), nil
}

// ReadJSONFile 读取
func ReadJSONFile(path string) (imgMap ImageRGBAMap, err error) {
	fmt.Println("读取JSON中....")
	//打开文件
	filePtr, err := os.Open(path)
	if err != nil {
		fmt.Println("打开文件失败....")
		return imgMap, err
	}
	defer filePtr.Close()

	//创建该文件的json解码器
	decoder := json.NewDecoder(filePtr)

	//把解码的结果存在xhz的地址中
	err = decoder.Decode(&imgMap)
	if err != nil {
		fmt.Println("JSON解码错误....")
		return imgMap, err
	}
	fmt.Println("读取JSON完成")
	return imgMap, nil
}

// WriteJSONFile 写入读取
func WriteJSONFile(imgMap ImageRGBAMap, name string) (err error) {
	fmt.Println("写出JSON中....")
	//创建文件（并打开）
	filePtr, err := os.Create(name + ".json")
	if err != nil {
		fmt.Println("创建文件失败....")
		return err
	}
	defer filePtr.Close()

	//创建基于文件的JSON编码器
	encoder := json.NewEncoder(filePtr)

	//将小黑子实例编码到文件中
	err = encoder.Encode(imgMap)
	if err != nil {
		fmt.Println("JSON编码失败....")
		return err
	}
	fmt.Println("写出JSON完成")
	return err
}

// ReadImageFile 读取一个图片文件
func ReadImageFile(path string) (img image.Image, err error) {
	fmt.Println("读取图片中....")
	fileByte, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("读取文件失败....")
		return img, err
	}
	img, _, err = image.Decode(bytes.NewBuffer(fileByte))
	if err != nil {
		fmt.Println("图片解码失败....")
		return img, err
	}
	fmt.Println("读取图片完成....")
	return img, err
}

// WriteImageFile 写入一个图片文件
func WriteImageFile(img image.Image, name string) (err error) {
	fmt.Println("写出图片中....")
	fullPath := name + ".png"
	filePtr, err := os.Create(fullPath)
	if err != nil {
		fmt.Println("文件创建失败....")
		return err
	}
	defer filePtr.Close()
	err = png.Encode(filePtr, img)
	if err != nil {
		fmt.Println("图片编码失败....")
		return err
	}
	fmt.Println("写出图片完成....")
	return err
}

// JavaScriptMiddleware 中间件
func JavaScriptMiddleware(imageMap ImageRGBAMap, path string) (outImageMap ImageRGBAMap, err error) {
	fmt.Println("读取JavaScript脚本....")
	scriptByte, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("文件读取失败....")
		return outImageMap, err
	}

	fmt.Println("执行JavaScript脚本....")
	jsvm := otto.New()

	imageMapJSON, err := json.Marshal(imageMap)
	if err != nil {
		fmt.Println("JSON解码失败....")
		return outImageMap, err
	}
	jsvm.Set("mainCallBack", func(call otto.FunctionCall) otto.Value {
		fmt.Printf("mainCallBack, %s.\n", call.Argument(0).String())
		return otto.Value{}
	})
	useScript := fmt.Sprintf("JSON.stringify(main(%s))", string(imageMapJSON))
	fullScript := fmt.Sprintf("%s; \n %s ", string(scriptByte), useScript)

	value, err := jsvm.Run(fullScript)
	if err != nil {
		fmt.Println("JavaScript执行出错....")
		return outImageMap, err
	}

	valueString, err := value.ToString()
	if err != nil {
		fmt.Println("获取JavaScript执行结果出错....")
		return outImageMap, err
	}

	err = json.Unmarshal([]byte(valueString), &outImageMap)
	if err != nil {
		fmt.Println("解码JavaScript执行结果出错....")
		return outImageMap, err
	}
	fmt.Println("执行JavaScript完毕")
	return outImageMap, err
}
