# ImageSwap
ImageSwap是一个便于简单图片研究的工具，它可以快速的把一张图片解析成为JSON文件或者将JSON还原成为一张图片。目前解析仅支持png、jpg、jpeg格式

### 示例
```
//解码图片
imageSwap -decode  -input ./image.png -output imagejson  -javascript ./script.js

//编码图片
imageSwap -encode  -input ./image.json -output imagejson 
```

### 基本命令使用
```
GLOBAL OPTIONS:
   --decode            从图片解码为JSON
   --encode            将JSON编码为图片
   --output value      输出文件名称 (default: "outfile")
   --input value       输入文件路径
   --javascript value  处理中间脚本路径
   --help, -h          查看帮助
   --version, -v       查看版本信息
```

### JSON格式示例
```
{
    "Width": 50, //图片宽度
    "Height": 50, //图片高度
    "Data": [
        {
            "X": 0, //像素X位置 
            "Y": 0, //像素Y位置
            "R": 255, //像素R(red)色值
            "G": 255, //像素G(green)色值
            "B": 255, //像素B(blue)色值
            "A": 255  //像素A(Alpha)通道值
        },
        {...},
        ...
    ]
```

### JavaScript 脚本示例
脚本必须有一个main函数作为入口,main函数会接受一个img值，该值为图片的JSON数据，此外，main函数可以将修改的JSON数据返回,最终会影响输出。
```
function main(img) { 
    img.Data.map(function (point) {
    })
    return img
}
```

