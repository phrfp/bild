package extractfeature

import (

  "image"
  _"fmt"
  "github.com/phrfp/bild/clone"
  "math"
	"github.com/phrfp/bild/parallel"
)

func VerticalLinePositions(src* image.Gray16, h1, h2 int) []uint16  {

  //create subimage and pass to peak finder.
  subRect := image.Rect(src.Bounds().Min.X, h1, src.Bounds().Max.X, h2)
  subimage := src.SubImage(subRect)

  pair := maxpeakvert1D( subimage )
  avgpos := avg(pair)
  return avgpos

}

func HorizontalLinePosition(src* image.Gray16, w1, w2, threshold int) []uint16  {

//  fmt.Println("-----------HLPos---------------")
  //create subimage and pass to peak finder.
  subRect := image.Rect(w1, src.Bounds().Min.Y, w2, src.Bounds().Max.Y)
  subimage := src.SubImage(subRect)

  pair := maxpeakhorz1D( subimage, threshold )
//  fmt.Println(pair)
  avgpos := avg(pair)
//  fmt.Println(avgpos)
  return avgpos

}


func avg( array []uint16 ) []uint16 {

  var total1 float64 = 0
  var total2 float64 = 0
  // fmt.Println("---------avg----------")
  // fmt.Println(array)

  for i := 0; i < len(array)/2; i++ {
    ipos := i*2
    total1 += float64(array[ipos])
    total2 += float64(array[ipos+1])
  }

  total1 = total1 / float64(len(array)/2)
  total2 = total2 / float64(len(array)/2)

  tavg := make([]uint16, 2)
  tavg[0] = uint16( math.Ceil(total1) )
  tavg[1] = uint16( math.Ceil(total2) )
  return tavg

}

func maxpeakvert1D( img image.Image ) []uint16 {

//  fmt.Println("---------maxpeak1D----------")


  src := clone.AsGray16(img)
//  fmt.Println(src)
  w := src.Bounds().Dx()
  h := src.Bounds().Dy()

  slicemax:= make([]uint16, 2*h)

  parallel.Line(h, func(start, end int) {
    // fmt.Println("---------Next----------")
    // fmt.Println("Start: ", start)
    // fmt.Println("End: ", end)

    for y := start; y < end; y++ {
    //  fmt.Println("---------Row: ", y)
      var tmax1 uint16 = 0
      var tmax2 uint16 = 0
      var pos1 uint16 = 0
      var pos2 uint16 = 0
      var platCnt uint16 = 0
      for x := 0; x < w; x++ {
        if x > 0 && x < w-1 {
          ipos := y*src.Stride + x*2

          tpix := uint16(src.Pix[ipos+0])<<8 | uint16(src.Pix[ipos+1])

          if (tpix > tmax1 || tpix > tmax2) { // is it woorth evaluating the peak
            tpix_p1 := uint16(src.Pix[ipos+2])<<8 | uint16(src.Pix[ipos+3])
            tpix_m1 := uint16(src.Pix[ipos-2])<<8 | uint16(src.Pix[ipos-1])

            // fmt.Println(tpix_p1)
            // fmt.Println(tpix_m1)
            if (tpix > tpix_p1 && tpix > tpix_m1) { //check that we have a peak
              if (tpix > tmax1 && tpix > tmax2) { //check that its local max
                pos2 = pos1 // downgrade last highest point
                tmax2 = tmax1
                tmax1 = uint16(tpix)
                pos1 = uint16(y)
              } else if ( (tpix > tmax2) && (uint16(y) > pos1+50) ){
                tmax2 = uint16(tpix)
                pos2 = uint16(x)
              }
            } else if (tpix > tpix_m1 && tpix == tpix_p1 ) {  //find rising edge
                platCnt = 1
            //    fmt.Println("Rising")
            } else if (tpix == tpix_m1 && tpix == tpix_p1) { // on plataux
               platCnt += 1
            //   fmt.Println("Plataux")
            } else if (tpix == tpix_m1 && tpix > tpix_p1) { // falling edge
            //  fmt.Println("Falling")
              if (tpix > tmax1 && tpix > tmax2) { //check that its local max
                pos2 = pos1 // downgrade last highest point
                tmax2 = tmax1
                tmax1 = uint16(tpix)
                pos1 = uint16(y)
            //    fmt.Println(int(math.Ceil(float64(platCnt)/2)))
                pos1 = uint16((x - int(math.Ceil(float64(platCnt)/2))))
              } else if ( (tpix > tmax2) && (uint16(y) > pos1+50) ) {
                tmax2 = uint16(tpix)
                pos2 = uint16((x - int(math.Ceil(float64(platCnt)/2))))
              }
              platCnt = 0
            }
          }
        }
      }

      ypos := y*2
      if pos1 > pos2 {
        slicemax[ypos] = pos2
        slicemax[ypos+1] = pos1
      } else {
        slicemax[ypos] = pos1
        slicemax[ypos+1] = pos2
      }
    }
  })

  return slicemax

}


func maxpeakhorz1D( img image.Image, threshold int ) []uint16 {

  // fmt.Println("---------maxpeak1D----------")


  src := clone.AsGray16(img)
  // fmt.Println(src)
  w := src.Bounds().Dx()
  h := src.Bounds().Dy()

  slicemax:= make([]uint16, 2*w)

  parallel.Line(w, func(start, end int) {
    // fmt.Println("---------Next----------")
    // fmt.Println("Start: ", start)
    // fmt.Println("End: ", end)
    for x := start; x < end; x++ {
      var tmax1 uint16 = 0
      var tmax2 uint16 = 0
      var pos1 uint16 = 0
      var pos2 uint16 = 0
      var platCnt uint16 = 0
      for y := 0; y < h; y++ {
    //  fmt.Println("---------Col: ", y)

        if y > 0 && y < h-1 {
          ipos := y*src.Stride + x*2
          ipos_p1 := (y+1)*src.Stride + x*2
          ipos_m1 := (y-1)*src.Stride + x*2

          tpix := uint16(src.Pix[ipos+0])<<8 | uint16(src.Pix[ipos+1])

          if ( (tpix > tmax1 || tpix > tmax2) && tpix > uint16(threshold)) { // is it woorth evaluating the peak - add threshold
            tpix_p1 := uint16(src.Pix[ipos_p1+0])<<8 | uint16(src.Pix[ipos_p1+1])
            tpix_m1 := uint16(src.Pix[ipos_m1+0])<<8 | uint16(src.Pix[ipos_m1+1])
            // fmt.Println(tpix_p1)
            // fmt.Println(tpix_m1)
            if (tpix > tpix_p1 && tpix > tpix_m1) { //check that we have a peak
              if (tpix > tmax1 && tpix > tmax2) { //check that its local max

                pos2 = pos1 // downgrade last highest point
                tmax2 = tmax1
                tmax1 = uint16(tpix)
                pos1 = uint16(y)
          //      fmt.Println(pos1)
              } else if ( (tpix > tmax2) && (uint16(y) > pos1+50) ) {
                tmax2 = uint16(tpix)
                pos2 = uint16(y)
        //        fmt.Println(pos2)
              }
            } else if (tpix > tpix_m1 && tpix == tpix_p1 ) {  //find rising edge
                platCnt = 1
            //    fmt.Println("Rising")
            } else if (tpix == tpix_m1 && tpix == tpix_p1) { // on plataux
               platCnt += 1
            //   fmt.Println("Plataux")
            } else if (tpix == tpix_m1 && tpix > tpix_p1) { // falling edge
            //  fmt.Println("Falling")
              if (tpix > tmax1 && tpix > tmax2) { //check that its local max

            //    fmt.Println(int(math.Ceil(float64(platCnt)/2)))
                pos2 = pos1 // downgrade last highest point
                tmax2 = tmax1
                tmax1 = uint16(tpix)
                pos1 = uint16((y - int(math.Ceil(float64(platCnt)/2))))
      //          fmt.Println(pos1)
              } else if ( (tpix > tmax2) && (uint16(y) > pos1+50) )  {
                tmax2 = uint16(tpix)
                pos2 = uint16((y - int(math.Ceil(float64(platCnt)/2))))
              }
              platCnt = 0
            }
      //      fmt.Println("tpix: ", tpix, " tmax1: ", tmax1, " tmax2: ", tmax2)
          }

        }
      }
      // fmt.Println(pos1)
      // fmt.Println(pos2)
      xpos := x*2
      if pos1 > pos2 {
        slicemax[xpos] = pos2
        slicemax[xpos+1] = pos1
      } else {
        slicemax[xpos] = pos1
        slicemax[xpos+1] = pos2
      }
    }
  })

  return slicemax

}
