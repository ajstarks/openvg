package main

import (
	"flag"
	"fmt"
	"github.com/ajstarks/openvg"
	"image"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"syscall"
)


func get_image_conf(image_path string) (image_conf image.Config, err error) {
	image_file, err := os.Open(image_path)
	if err != nil {
		return image_conf, err
	}
	defer image_file.Close()
	image_conf, _, err = image.DecodeConfig(image_file)
	if err != nil {
		return image_conf, err
	}
	image_file.Close()
	return
}

func resize_image(w, h int, image_path string) (image_path_resized string, err error) {
	// One option is to use github.com/nfnt/resize,
	//  but it seem to be significantly slower than imagemagick.
	size_str := fmt.Sprintf("%dx%d", w, h)
	image_path_dst := fmt.Sprintf("%s.cache-%s", image_path, size_str)
	_, err = os.Stat(image_path_dst)
	if err != nil {
		cmd := exec.Command("convert", "-resize", size_str, image_path, image_path_dst)
		err := cmd.Run()
		if err != nil {
			return image_path, err
		}
	}
	return image_path_dst, nil
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stdout,
			`usage: %s [ opts ] image-path

Set specified image as an overlay screen via OpenVG lib.
Default is to put this image to the center.

`, os.Args[0])
		flag.PrintDefaults()
	}

	var resize bool
	var bg_color string
	flag.BoolVar(&resize, "resize", false, `Resize image to fit the screen (using "convert" binary.`)
	flag.StringVar(&bg_color, "bg-color", "",
		`Background color to use with centered image, in RRGGBB (hex, e.g. "aabbcc") format.`)

	flag.Parse()
	if flag.NArg() != 1 {
		log.Print("ERROR: Exactly one image-path argument must be specified.")
		flag.Usage()
	}
	var r, g, b uint8
	var err error
	if len(bg_color) == 0 {
		r, g, b = 0, 0, 0
	} else if len(bg_color) == 6 {
		n, err := strconv.ParseUint(bg_color[:2], 16, 8)
		if err == nil {
			r = uint8(n)
			n, err = strconv.ParseUint(bg_color[2:4], 16, 8)
		}
		if err == nil {
			g = uint8(n)
			n, err = strconv.ParseUint(bg_color[4:6], 16, 8)
		}
		if err == nil {
			b = uint8(n)
		}
	}
	if err != nil {
		log.Fatalf("ERROR: Failed to parse bg-color value (%v): %v", bg_color, err)
	}
	image_path := flag.Args()[0]

	exit_code := 0
	defer func() {
		os.Exit(exit_code)
	}()


	openvg.SaveTerm()
	w, h := openvg.Init()
	openvg.RawTerm()
	defer openvg.Finish()
	defer openvg.RestoreTerm()

	image_conf, err := get_image_conf(image_path)
	if err == nil && resize && (image_conf.Width != w || image_conf.Height != h) {
		image_path, err = resize_image(w, h, image_path)
	}
	if err != nil {
		log.Printf("ERROR: Failed to process image (%v): %v", image_path, err)
		exit_code = 1
	} else {
		sig_chan := make(chan os.Signal, 1)
		signal.Notify( sig_chan, os.Interrupt, os.Kill,
			syscall.SIGHUP, syscall.SIGTERM, syscall.SIGALRM )

		openvg.Start(w, h)
		openvg.Background(r, g, b)
		if resize {
			openvg.Image(0, 0, w, h, image_path)
		} else {
			x, y := openvg.VGfloat(w) / 2 - openvg.VGfloat(image_conf.Width) / 2,
				openvg.VGfloat(h) / 2 - openvg.VGfloat(image_conf.Height) / 2
			openvg.Image(x, y, image_conf.Width, image_conf.Height, image_path)
		}
		openvg.End()

		_ = <-sig_chan
	}
}
