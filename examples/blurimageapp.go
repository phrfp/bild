package main


import (

  "github.com/phrfp/bild/blur"
  "github.com/phrfp/bild/imgio"


)

func main() {
  img, err := imgio.Open("origImg.png")
  if err != nil {
      panic(err)
  }

  testSmooth := blur.GaussianG16(img,3.0)


  if err := imgio.Save("smoothImg", testSmooth, imgio.PNG); err != nil {
      panic(err)
  }

}
