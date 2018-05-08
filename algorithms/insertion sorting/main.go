package main

//сортировка вставками, Kormen - 40

import "fmt"

func main() {
	a := []int{5, 2, 4, 6, 1, 3, 4, 5, 62, 52, 4, 1, 621, 54}
	for k, v := range a {
		i := k - 1
		for ; i >= 0 && a[i] > v; i-- {
			a[i+1] = a[i]
		}
		a[i+1] = v
	}
	fmt.Println(a)
}
