package okcoin

import (
	"net/http"
	"testing"

	. "github.com/BTB/goBTB/base"
	"github.com/stretchr/testify/assert"
)

var (
	okex = NewOKEx(http.DefaultClient, "", "")
)

func TestOKEx_GetFutureDepth(t *testing.T) {
	// https: //www.okex.com/api/v1/depth.do?symbol=ltc_btc

	dep, err := okex.GetFutureDepth(BTC_USD, THIS_WEEK_CONTRACT, 1)
	assert.Nil(t, err)
	t.Log(dep)
}
