// 加密
function main(img) {
  var context = StringToHex('66862039').split(',');
  console.log(context);

  img.Data.forEach(function (point, index) {
    if (index <= context.length + 10) {
      // 汉字占 4 个字节
      var charCode = context[index];
      point.R = parseInt(charCode || 'f1', 16);
    }
  });

  return img;
}

function StringToHex(value) {
  var val = "";
  for (var i = 0; i < value.length; i++) {
    if (val == "")
      val = value.charCodeAt(i).toString(16);
    else
      val += "," + value.charCodeAt(i).toString(16);
  }
  return val;
}
