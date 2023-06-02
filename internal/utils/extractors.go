package utils

import (
	"errors"
	"net/http"
	"strings"
)

func ExtractIpAddr(r *http.Request) (string, bool) {
	if len(r.Referer()) != 0 {
		return strings.Split(r.Referer(), ":")[0], true
	}

	if len(r.RemoteAddr) != 0 {
		return strings.Split(r.RemoteAddr, ":")[0], true
	}

	return "", false
}

// функция на вход принимает параметр из тела запроса и возвращает (ос, версия ос, ошибка)
func SplitOsAndVersion(str string) (string, string, error) {
	// minimal string - "IOS 1.1.1"
	//					 123456789
	if len(str) < 9 {
		return "", "", errors.New(`can't extract data from string: "` + str + `"`)
	}

	OsVer := strings.Split(str, " ")
	if len(OsVer) != 2 {
		return "", "", errors.New(`invalid format, can't extract os type and os version from string: "` + str + `"`)
	}

	// так делать не нужно, в идеале нужно сделать проверку версий совсем подругому для каждого типа ос
	// так как разные ос имеют свои форматы версий, но для упрощения оставлю так
	if nums := strings.Split(str, "."); len(nums) != 3 {
		return "", "", errors.New(`invalid os version format: "` + str + `"`)
	}
	return OsVer[0], OsVer[1], nil
}
