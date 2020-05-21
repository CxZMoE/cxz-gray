package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"cxz.moe/cxz-gray/gray"
)

var (
	inputFile  string
	outputFile string
	ratio      int
	//imgtype    string
	shreshold int
)

const (
	FileNotExist        = 1 // Code 1 File not exist
	WrongRatioValue     = 2 // Code 2 Wrong ratio value
	WrongShresholdValue = 3 // Code 3 Wrong shreshold value
)

func init() {
	flag.StringVar(&inputFile, "i", "", "<input file>")
	flag.StringVar(&outputFile, "o", "output.png", "<output file>")
	flag.IntVar(&ratio, "r", 3, "<scale ratio>")
	flag.IntVar(&shreshold, "s", 180, "Range:[0-255]\n<alpha channel shreshold filter value.>")
	//flag.StringVar(&imgtype, "t", "png", "\n<the type of output file>")
	flag.Parse()

	log.SetPrefix("[INFO] ")
	_, err := os.Open(inputFile)
	if os.IsNotExist(err) {
		log.Println("Input file does not exist.")
		os.Exit(FileNotExist)
	}
	if ratio < 1 {
		log.Println("Wrong ratio value.")
		os.Exit(WrongRatioValue)
	}

	if shreshold < 0 || shreshold > 255 {
		log.Println("Wrong shreshold value.")
		os.Exit(WrongShresholdValue)
	}

}

func main() {
	img, _ := gray.LoadImgFile(inputFile)
	if img == nil {
		log.Println("Failed to load file.")
		return
	}
	log.Println("Scalling Image...")
	img2 := gray.ToImageRatio(ratio, img)

	log.Println("Making it Binaryzation...")
	img2 = gray.MakeBinaryzation(img2, shreshold)

	log.Println("Generating ASCII text...")
	str := gray.ToASCIIFull(img2)
	if path.Ext(outputFile) == "jpeg" || path.Ext(outputFile) == "jpg" {
		gray.SaveJPEG(outputFile, img2)

	} else {
		gray.SavePNG(outputFile, img2)
		fpath := path.Base(outputFile)
		fpath = strings.Split(fpath, path.Ext(fpath))[0] + ".txt"
		if ioutil.WriteFile(fpath, []byte(str), 0755) != nil {
			log.Println("Failed to write ascii file.")
			return
		}
		log.Println("Output ASCII Image:", outputFile)
		log.Println("Output ASCII Text:", fpath)
	}

}
