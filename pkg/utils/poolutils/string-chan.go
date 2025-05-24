package poolutils

import (
	"sync"
)

var stringChanPool *sync.Pool

func GetStringChan() chan string {
	return stringChanPool.Get().(chan string)
}

func PutStringChan(c chan string) {
	stringChanPool.Put(c)
}
