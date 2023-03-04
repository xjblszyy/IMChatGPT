package utils

import (
	"strings"
)

func ContainKeyword(keyword, text string) bool {
	if keyword == "" {
		return false
	}
	return strings.Contains(
		strings.ToLower(text),
		strings.ToLower(keyword),
	)
}

func ContainBlackList(nickname string, blackList []string) bool {
	if len(blackList) == 0 {
		return false
	}
	for _, data := range blackList {
		if data == nickname {
			return true
		}
	}
	return false
}

func GetQuestionFromMsg(keyword, msg string) string {
	question := strings.Replace(msg, strings.ToUpper(keyword), "", -1)
	question = strings.Replace(question, strings.ToLower(keyword), "", -1)
	return question
}
