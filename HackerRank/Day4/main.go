package main

//https://www.hackerrank.com/challenges/30-class-vs-instance/problem

import "fmt"

type person struct {
	age int
}

func (p person) NewPerson(initialAge int) person  {

	if initialAge < 0 {
		p.age = 0
		fmt.Println("Age is not valid, setting age to 0.")
	} else {
		p.age = initialAge
	}

	return p
}

func (p person) amIOld() {
	switch {
	case p.age < 13:
		fmt.Println("You are young.")
	case p.age >= 13 && p.age < 18:
		fmt.Println("You are a teenager.")
	default:
		fmt.Println("You are old.")
	}
}

func (p person) yearPasses() person {
	p.age++
	return p
}
