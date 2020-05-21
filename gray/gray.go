package gray

import (
	"image"
	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"log"
	"os"
)

// Load A Image from file
func LoadImgFile(fpath string) (image.Image, error) {
	f, err := os.Open(fpath)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return img, nil
}

// Make it Binaryzation
func MakeBinaryzation(img image.Image, shreshold int) image.Image {

	newImg := image.NewRGBA(img.Bounds())

	for y := img.Bounds().Min.Y; y <= img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x <= img.Bounds().Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			//log.Println("R:", r, "G:", g, "B:", b, "A:", a)
			rr, gg, bb := uint8(r>>8), uint8(g>>8), uint8(b>>8)
			//log.Println("Turned R:", rr, "G:", gg, "B:", bb, "A:", aa)
			ggg := float32(rr)*0.299 + float32(gg)*0.587 + float32(bb)*0.114

			if isAboveShreshold(float32(shreshold), ggg) {
				newImg.SetRGBA(x, y, color.RGBA{255, 255, 255, 255})
			} else {
				newImg.SetRGBA(x, y, color.RGBA{0, 0, 0, 255})

			}

		}
	}

	return newImg.SubImage(newImg.Bounds())
}

// isAboveShreshold
func isAboveShreshold(valve float32, gray float32) bool {
	if gray > valve {
		return true
	}
	return false
}

// Save png to file.
func SavePNG(fpath string, img image.Image) bool {
	if img == nil {
		return false
	}
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Println("Failed to create image file:", err.Error())
		return false
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		log.Println("Failed to encode image:", err.Error())
		return false

	}

	//fmt.Scanln()
	return true
}

// Save jpg to file with 100% quality.
func SaveJPEG(fpath string, img image.Image) bool {
	if img == nil {
		return false
	}
	f, err := os.OpenFile(fpath, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Println("Failed to encode image:", err.Error())
		return false
	}
	defer f.Close()

	jpeg.Encode(f, img, &jpeg.Options{
		100,
	})

	return true
	//fmt.Scanln()
}

// make it ascii
func ToASCIIFull(img image.Image) string {
	if img == nil {
		return ""
	}
	var str = ""
	for y := img.Bounds().Min.Y; y <= img.Bounds().Max.Y; y = y + 2 {
		for x := img.Bounds().Min.X; x <= img.Bounds().Max.X; x++ {
			//_, _, _, a := img.At(x, y).RGBA()
			//log.Println(a)
			//aa := uint8(a >> 8)
			str += GetChar(img, x, y)
			if x == img.Bounds().Max.X {
				str += "\n"
			}

		}

	}
	return str
}

func ToImageRatio(ratio int, img image.Image) image.Image {
	if img == nil {
		return nil
	}
	var ratioRect = image.Rect(img.Bounds().Min.X, img.Bounds().Min.Y, (img.Bounds().Min.X+img.Bounds().Dx())/ratio, (img.Bounds().Min.Y+img.Bounds().Dy())/ratio)
	var ratioImg = image.NewRGBA(ratioRect)
	//var stepX = 0
	//var stepY = 0
	for y := ratioImg.Bounds().Min.Y; y < ratioImg.Bounds().Max.Y; y++ {

		for x := ratioImg.Bounds().Min.X; x < ratioImg.Bounds().Max.X; x++ {
			var r, g, b, a uint32
			//count := 0
			//log.Println("x:", x, "y:", y)
			for y_next := y * ratio; y_next < y*ratio+ratio; y_next++ {

				for x_next := x * ratio; x_next < x*ratio+ratio; x_next++ {
					r_next, g_next, b_next, a_next := img.At(x_next, y_next).RGBA()
					//log.Println("R:", r, "G:", g, "B:", b, "A:", a)
					r = r + r_next
					g = g + g_next
					b = b + b_next
					a = a + a_next

					//count++
					//log.Println("xnext:", x_next, "y_next", y_next)
				}

			}
			//log.Println("count:", count)

			devide := uint32(ratio * ratio)
			r, g, b, a = r/devide, g/devide, b/devide, a/devide

			ratioImg.SetRGBA(x, y, color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)})

		}

	}
	return ratioImg
}

// Get Charactor by position
func GetChar(img image.Image, x, y int) string {
	minX := img.Bounds().Min.X
	maxX := img.Bounds().Max.X
	minY := img.Bounds().Min.Y
	maxY := img.Bounds().Max.Y

	/* Border */
	if x == maxX || x == minX || y == maxY || y == minY {
		return "."
	}
	/* Normal */

	topAlpha := !GetPositive(img, x, y-1)
	//topRightAlpha := GetAlpha(img, x+1, y-1)
	rightAlpha := !GetPositive(img, x+1, y)
	//bottomRightAlpha := GetAlpha(img, x+1, y+1)
	bottomAlpha := !GetPositive(img, x, y+1)
	//bottomLeftAlpha := GetAlpha(img, x-1, y+1)
	leftAlpha := !GetPositive(img, x-1, y)
	//topLeftAlpha := GetAlpha(img, x-1, y-1)

	// Around
	if topAlpha && rightAlpha && bottomAlpha && leftAlpha {
		return "#"
	}

	// Right Half
	if topAlpha && rightAlpha && bottomAlpha {
		return "]"
	}

	// Left Half
	if topAlpha && leftAlpha && bottomAlpha {
		return "["
	}

	// Top Half
	if leftAlpha && topAlpha && rightAlpha {
		return "^"
	}

	// Bottom Half
	if leftAlpha && bottomAlpha && rightAlpha {
		return "_"
	}

	// Top Left
	if topAlpha && leftAlpha {
		return "`"
	}

	// Top Right
	if topAlpha && rightAlpha {
		return "\\"
	}

	// Bottom Right
	if bottomAlpha && rightAlpha {
		return "/"
	}

	// Bottom Left
	if bottomAlpha && leftAlpha {
		return "\\"
	}

	// Top Left
	if topAlpha && leftAlpha {
		return "/"
	}

	// Top
	if topAlpha {
		return "^"
	}

	// Right
	if rightAlpha {
		return ">"
	}

	// Bottom
	if bottomAlpha {
		return "_"
	}

	// Left
	if leftAlpha {
		return "<"
	}
	return " "
}

// Get alpha channel value
func GetAlpha(img image.Image, x, y int) uint8 {
	_, _, _, a := img.At(x, y).RGBA()
	return uint8(a >> 8)
}

// Get alpha channel value
func GetPositive(img image.Image, x, y int) bool {
	r, g, b, _ := img.At(x, y).RGBA()
	if r+g+b > 0 {
		return true
	}
	return false
}
