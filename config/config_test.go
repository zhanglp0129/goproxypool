package config_test

import (
	"fmt"
	"github.com/zhanglp0129/goproxypool/config"
	"testing"
)

var CFG = config.CFG

func TestParseConfig(t *testing.T) {
	fmt.Println(CFG)
}
