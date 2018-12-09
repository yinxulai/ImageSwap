function main(img) {
    img.Data.map(function (point) {
        point.R = 0
    })
    return img
}