package xydim

import (
	"fmt"
	"strconv"
	"strings"
)

type XyDimInt struct {
	X int
	Y int
}

func (d *XyDimInt) String() string {
	return fmt.Sprintf("%dx%d", d.X, d.Y)
}

func (d *XyDimInt) Set(s string) error {
	parts := strings.Split(s, "x")
	if len(parts) != 2 {
		return fmt.Errorf("must be WxH")
	}

	var err1, err2 error
	d.X, err1 = strconv.Atoi(parts[0])
	d.Y, err2 = strconv.Atoi(parts[1])

	if err1 != nil || err2 != nil {
		return fmt.Errorf("int conversion error")
	}
	return nil
}

type XyDimFloat struct {
	X float32
	Y float32
}

func (d *XyDimFloat) String() string {
	return fmt.Sprintf("%gx%g", d.X, d.Y)
}

func (d *XyDimFloat) Set(s string) error {
	parts := strings.Split(s, "x")
	if len(parts) != 2 {
		return fmt.Errorf("must be WxH")
	}

	x, err1 := strconv.ParseFloat(parts[0], 32)
	y, err2 := strconv.ParseFloat(parts[1], 32)

	if err1 != nil || err2 != nil {
		return fmt.Errorf("float conversion error")
	}

	d.X = float32(x)
	d.Y = float32(y)
	return nil
}
