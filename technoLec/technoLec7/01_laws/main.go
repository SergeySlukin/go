package main

import (
	"fmt"
	"reflect"
)

func main() {
	x := 3.4
	// TypeOf() принимает на вход interface{}, в этом месте будет аалокация
	fmt.Printf("reflect.Type %+v\n:", reflect.TypeOf(x).String())

	// reflect.Value != значению переданному на вход
	fmt.Println("reflect.Value:", reflect.ValueOf(x).String())

	v := reflect.ValueOf(x)
	fmt.Println("Type value:", v.Type())
	fmt.Println("Type float64:", v.Kind() == reflect.Float64)
	fmt.Println("Value:", v.Float())

	type MyInt int
	var c MyInt = 7
	v = reflect.ValueOf(c)
	fmt.Println("kind is int: ", v.Kind() == reflect.Int)

	y := v.Interface().(MyInt)
	fmt.Println("Value wrapper", v, "Value", y)
	
	access()
}

func access()  {
	var x float64 = 3.4
	v := reflect.ValueOf(x)
	fmt.Println("settability of v:", v.CanSet())

	// чтобы иметь мозможность изменить значение, нам потребуется ссылка
	p := reflect.ValueOf(&x)
	fmt.Println("type of p:", p.Type())
	fmt.Println("settability of p:", p.CanSet())

	// Теперь, используя Elem мы получим value, лежащее по ссылке
	v = p.Elem()
	fmt.Println("settability of v:", v.CanSet())

	v.SetFloat(7.1)
	fmt.Println(v.Interface())
	fmt.Println(x)
}
