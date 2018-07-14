package main

import (
	"github.com/bradfitz/gomemcache/memcache"
	"fmt"
	"time"
)

func getRecord(mKey string) *memcache.Item {
	fmt.Println("get", mKey)
	item, err := memcacheClient.Get(mKey)

	if err == memcache.ErrCacheMiss {
		fmt.Println("Record not found in memcache")
		return nil
	} else if err != nil {
		PanicOnErr(err)
	}
	return item
}

var (
	memcacheClient *memcache.Client
)

func main() {
	MemcachedAddresses := []string{"127.0.0.1:11211"}
	memcacheClient = memcache.New(MemcachedAddresses...)

	mKey := "record_21"

	item := getRecord(mKey)
	fmt.Println("first get %+v\n", item)

	ttl := 5
	err := memcacheClient.Set(&memcache.Item{
		Key:   mKey,
		Value: []byte("1"),
		// указываем через сколько секунд запись пропадет из кеша
		Expiration: int32(ttl),
	})
	PanicOnErr(err)

	time.Sleep(time.Microsecond)
	//time.Sleep(time.Duration(ttl+1) * time.Second)
	item = getRecord(mKey)

	err = memcacheClient.Add(&memcache.Item{
		Key: mKey,
		Value: []byte("2"),
		Expiration: int32(ttl),
	})

	// если запись не была добавлена, вернётся соответствующая ошибка
	if err == memcache.ErrNotStored {
		fmt.Println("Record not stored")
	} else if err != nil {
		PanicOnErr(err)
	}
	item = getRecord(mKey)
	fmt.Println("third get %+v\n", item)

	afterIncrement, err := memcacheClient.Increment(mKey, uint64(2))
	PanicOnErr(err)
	fmt.Println("afterIncrement by 2 ", mKey, "is", afterIncrement)
	afterDecrement, err := memcacheClient.Decrement(mKey, uint64(1))
	PanicOnErr(err)
	fmt.Println("afterDecrement by 1 ", mKey, "is", afterDecrement)

	//для несуществующей записи инкремент невозможен
	afterIncrement, err = memcacheClient.Increment(mKey + "_not_exist", uint64(1))
	fmt.Println("afterIncrement not existing record ", afterIncrement)
	if err == memcache.ErrCacheMiss {
		fmt.Println("Record not exist")
	} else if err != nil {
		PanicOnErr(err)
	} else {
		fmt.Println("afterDecrement by 1 ", mKey, "is", afterDecrement)
	}

	mkeys := []string{mKey, "record_22"}
	fmt.Println("get multiple", mkeys)
	multipleItems, err := memcacheClient.GetMulti(mkeys)
	PanicOnErr(err)
	fmt.Println("multipleItems", multipleItems)

}

func PanicOnErr(err error) {
	if err != nil {
		panic(err)
	}
}
