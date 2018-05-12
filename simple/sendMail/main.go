package main

import (
	"net/smtp"
	"log"
)

func main()  {
	auth := smtp.PlainAuth("", "testleads2018@gmail.com", "Nhst27nsuy7", "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, "testleads2018@gmail.com", []string{"sergeyslukin90@gmail.com"}, []byte("test"));
	if err != nil {
		log.Fatal(err)
	}
}
