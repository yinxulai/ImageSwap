// ./bin/cli --decode --input ./back/outfile.png --javascript ./example/decode.js
// 解密
function main(img) {
  var context = [];
  var flag = 'f1f1f1f1f1f1f1f1f1f1';

  img.Data.forEach(function (point) {
    // 找到标示为 f1f1f1f1f1f1f1f1f1f1 就不找了
    if (context.join('').indexOf(flag) === -1) {
      context.push(Number(point.R).toString(16));
    }
  });

  // 裁掉标示为剩下的转字符串
  var contextStr = context.slice(0, context.length - 10).join(',')
  console.log('口令是：', HexToString(contextStr));
  return img;
}

function HexToString(value) {
  var val = "";
  var arr = value.split(",");

  for (var i = 0; i < arr.length; i++) {
    val += String.fromCharCode(parseInt(arr[i], 16));
  }

  return val;
}
