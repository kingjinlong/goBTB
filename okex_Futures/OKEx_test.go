package okcoin

import (
	"fmt"
	"net/http"
	"testing"

	. "github.com/BTB/goBTB/base"
	"github.com/stretchr/testify/assert"
)

var (
	apiKey    = ""
	secretKey = ""
	okex      = NewOKEx(http.DefaultClient, apiKey, secretKey)
)

func TestOKEx_GetFutureDepth(t *testing.T) {
	// https: //www.okex.com/api/v1/depth.do?symbol=ltc_btc
	return
	dep, err := okex.GetFutureDepth(BTC_USD, THIS_WEEK_CONTRACT, 1)
	assert.Nil(t, err)
	t.Log(dep)
}

func TestOKEx_GetFutureUserinfo(t *testing.T) {
	fmt.Println("qry money")
	// https: //www.okex.com/api/v1/depth.do?symbol=ltc_btc
	dep, err := okex.GetFutureUserinfo()
	if dep != nil {
		fmt.Println(dep)
	}
	assert.Nil(t, err)
	t.Log(dep)
}
