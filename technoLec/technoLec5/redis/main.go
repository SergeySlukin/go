package main

import (
	"github.com/garyburd/redigo/redis"
	"fmt"
	"time"
)

var c redis.Conn

func getRecord(mkey string) string {
	fmt.Println("get", mkey)
	item, err := redis.String(c.Do("GET", mkey))
	if err == redis.ErrNil {
		fmt.Println("Record not found is redis (return value is nil)")
		return ""
	} else if err != nil {
		PanicOnError(err)
	}
	return item
}

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var err error
	c, err = redis.DialURL("redis://user:@localhost:6379/0")
	PanicOnError(err)
	defer c.Close()

	mkey := "record_21"
	item := getRecord(mkey)
	fmt.Printf("first get %+v\n", item)

	ttl := 50
	result, err := redis.String(c.Do("SET", mkey, 1, "EX", ttl))
	PanicOnError(err)
	if result != "OK" {
		panic("result not ok: " + result)
	}

	time.Sleep(time.Millisecond)
	item = getRecord(mkey)
	fmt.Printf("second get %+v\n", item)

	n, _ := redis.Int(c.Do("INCRBY", mkey, 2))
	fmt.Println("INCRBY by 2 ", mkey, "is", n)

	n, _ = redis.Int(c.Do("DECRBY", mkey, 1))
	fmt.Println("DECRBY by 1 ", mkey, "is", n)

	//если записи не было - редис создаст
	n, err = redis.Int(c.Do("INCR", mkey+"_not_exist"))
	fmt.Println("INCR (default by 1) ", mkey+"_not_exist", "is", n)
	PanicOnError(err)

	keys := []interface{}{mkey, mkey + "_not_exist", "sure_not_exist"}
	reply, err := redis.Strings(c.Do("MGET", keys...))
	PanicOnError(err)
	fmt.Println(reply)
}
