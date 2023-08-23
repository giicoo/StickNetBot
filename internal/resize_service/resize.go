package resizeService

import (
	"fmt"
	"os"
	"time"

	"github.com/h2non/bimg"
)

type ResizeService struct {
	path     string
	time_out time.Duration
}

func NewResizeService(path string, time_out time.Duration) *ResizeService {
	return &ResizeService{
		path:     path,
		time_out: time_out,
	}
}

func (r *ResizeService) ResizeImage(file_name string, codec string) error {
	buffer, err := bimg.Read(fmt.Sprintf("%v.%v", file_name, codec))
	if err != nil {
		return fmt.Errorf("read img '%v': %v", file_name, err)
	}

	if codec == "jpg" || codec == "jpeg" {
		buffer, err = bimg.NewImage(buffer).Convert(bimg.PNG)
		if err != nil {
			return fmt.Errorf("convert to png '%v': %v", file_name, err)
		}
	}

	size, err := bimg.NewImage(buffer).Size()
	if err != nil {
		return fmt.Errorf("read size of '%v': %v", file_name, err)
	}
	var size_x float32 = float32(size.Width)
	var size_y float32 = float32(size.Height)

	if size_x > 512 || size_y > 512 {
		x, y := 512, 512
		if size_x > 512 && size_y > 512 {
			k := size_y / 512
			if size_x > size_y {
				k = size_x / 512
			}
			x, y = int(size_x/k), int(size_y/k)
		}
		if size_x > size_y {
			k := size_x / 512
			x, y = int(size_x/k), int(size_y/k)
		}
		if size_y > size_x {
			k := size_y / 512
			x, y = int(size_x/k), int(size_y/k)
		}

		buffer, err = bimg.NewImage(buffer).ForceResize(x, y)
		if err != nil {
			return fmt.Errorf("resize img '%v': %v", file_name, err)
		}
	}
	_ = bimg.NewImage(buffer).Length()
	// TODO: compress img if length > 512000

	bimg.Write(fmt.Sprintf("new_%v.png", file_name), buffer)
	go r.DeleteFilesAfterTimer(file_name)
	fmt.Println("delete timer start")
	return nil
}

// Ignore errs because err will be if file don't exist
func (r *ResizeService) DeleteFiles(file_name string) {
	os.Remove(fmt.Sprintf("%v.png", file_name))
	os.Remove(fmt.Sprintf("%v.jpg", file_name))
	os.Remove(fmt.Sprintf("%v.jpeg", file_name))
	os.Remove(fmt.Sprintf("new_%v.png", file_name))
}

func (r *ResizeService) DeleteFilesAfterTimer(file_name string) {
	select {
	case <-time.After(r.time_out):
		r.DeleteFiles(file_name)
	}
}

// Ignore errs because err will be if file don't exist
func (r *ResizeService) DeleteFiles(file_name string) {
	os.Remove(fmt.Sprintf("%v.png"))
	os.Remove(fmt.Sprintf("%v.jpg"))
	os.Remove(fmt.Sprintf("%v.jpeg"))
	os.Remove(fmt.Sprintf("new_%v.png"))
}

func (r *ResizeService) DeleteFilesAfterTimer(file_name string) {
	select {
	case <-time.After(r.time_out):
		r.DeleteFiles(file_name)
	}
}
