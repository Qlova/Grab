package grab

import "os"
import "image"
import _ "image/png"
import _ "image/jpeg"

type Image struct {
	Data []byte
	Stride int
	Width, Height int
}

func (img *Image) Load(filepath string) error {
	infile, err := os.Open(filepath)
    if err != nil {
        // replace this with real error handling
        return err
    }
    defer infile.Close()

    // Decode will figure out what type of image is in the file on its own.
    // We just have to be sure all the image packages we want are imported.
    src, _, err := image.Decode(infile)
    if err != nil {
        // replace this with real error handling
       return err
    }

    // Create a new grayscale image
    bounds := src.Bounds()
    w, h := bounds.Max.X, bounds.Max.Y
    rgba := image.NewRGBA(image.Rect(0,0, w, h))
    for x := 0; x < w; x++ {
        for y := 0; y < h; y++ {
            oldColor := src.At(x, y)
            grayColor := rgba.ColorModel().Convert(oldColor)
            rgba.Set(x, y, grayColor)
        }
    }
    
    img.Data = rgba.Pix
    img.Stride = int(rgba.Stride)
    img.Width, img.Height = int(w), int(h)
    return nil
}
