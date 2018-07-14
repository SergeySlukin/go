package main //1
import (
	"reflect"
	"fmt"
)

type rect struct {
	width int
	height int
}

type circle struct {
	radius int
}

func Set(v interface{})  {
	rvp := reflect.ValueOf(v)
	if rvp.Kind() != reflect.Ptr {
		panic("Ожидается указатель")
	}
	rv := reflect.Indirect(rvp)
	i := rv.Interface()
	switch i.(type) {
	case rect:
		val := rect{100, 50}
		rv.Set(reflect.ValueOf(val))
	case circle:
		val := circle{200}
		rv.Set(reflect.ValueOf(val))
	default:
		panic("Undefined type")
	}
}

func main()  {
	var rect rect
	Set(&rect)
	fmt.Println(rect)

	var circle circle
	Set(&circle)
	fmt.Println(circle)
}