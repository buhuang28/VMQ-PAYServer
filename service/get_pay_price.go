package service

import (
	"Pay-Server/log"
	"fmt"
	"github.com/shopspring/decimal"
	"strconv"
	"sync"
)

var (
	PriceMap  = make(map[float64]bool)
	PriceLock sync.Mutex
)

func init() {
	quantity := decimal.NewFromFloat(299)
	float := decimal.NewFromFloat(0.01)
	for i := 0; i < 21; i++ {
		sub := quantity.Sub(decimal.NewFromInt(int64(i)).Mul(float))
		f, _ := sub.Float64()
		value, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", f), 64)
		PriceMap[value] = false
	}
}

func GetPayPrice() float64 {
	PriceLock.Lock()
	defer func() {
		e := recover()
		if e != nil {
			log.PrintStackTrace(e)
		}
		PriceLock.Unlock()
	}()
	//false是未被使用
	for k, v := range PriceMap {
		if !v {
			PriceMap[k] = true
			return k
		}
	}
	return 0
}

func RecoverPrice(f float64) bool {
	PriceLock.Lock()
	defer func() {
		e := recover()
		if e != nil {
			log.PrintStackTrace(e)
		}
		PriceLock.Unlock()
	}()
	_, ok := PriceMap[f]
	if ok {
		//价格恢复到可使用
		PriceMap[f] = false
		return true
	}
	return false
}
