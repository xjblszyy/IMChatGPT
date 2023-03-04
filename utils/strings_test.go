package utils

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestContainKeyword(t *testing.T) {
	keyword := "小哈"
	content := "你好 小哈"
	res := ContainKeyword(keyword, content)
	assert.True(t, res)
}
