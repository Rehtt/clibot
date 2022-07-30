/**
 * @Author: dsreshiram@gmail.com
 * @Date: 2022/7/30 上午 09:14
 */

package cmdbot

import (
	"sync"
)

const (
	RoleAdmin    = "admin"
	RoleOrdinary = "ordinary"
)

var (
	roleMap = make(map[int64]string)
	roleMu  sync.RWMutex
)

func GetRole(uin int64) string {
	roleMu.RLock()
	defer roleMu.RUnlock()
	if role, ok := roleMap[uin]; ok {
		return role
	}
	return RoleOrdinary
}

func GetRoleUins(role string) []int64 {
	var list []int64
	roleMu.RLock()
	defer roleMu.RUnlock()
	for u, r := range roleMap {
		if r == role {
			list = append(list, u)
		}
	}
	return list
}

func SetRole(uin int64, role string) {
	roleMu.Lock()
	defer roleMu.Unlock()
	roleMap[uin] = role
}

func HasPermission(uin int64, p []string) bool {
	if len(p) == 0 {
		return true
	}
	role := GetRole(uin)
	for _, v := range p {
		if v == role {
			return true
		}
	}
	return false
}
