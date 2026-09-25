package cache

import "fmt"

const prefixUser = "user"

func UserKey(id uint) string {
	return fmt.Sprintf("%s:%d", prefixUser, id)
}

func UserWithOrdersKey(id uint) string {
	return fmt.Sprintf("%s:%d:with_orders", prefixUser, id)
}
