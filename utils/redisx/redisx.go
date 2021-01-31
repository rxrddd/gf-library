package redisx

import (
	"github.com/gogf/gf/database/gredis"
	"github.com/gomodule/redigo/redis"
)

//redisgo 事务
func Multi(conn redis.Conn, fc func(con redis.Conn)) error {
	conn.Send("MULTI")
	fc(conn)
	if _, err := conn.Do("EXEC"); err != nil {
		return err
	}
	return nil
}

func DelKey(conn *gredis.Redis, redisKey string) {
	keys, _ := conn.DoVar("KEYS", redisKey)
	_ = Multi(conn.Conn(), func(con redis.Conn) {
		for _, key := range keys.Strings() {
			con.Send("DEL", key)
		}
	})
}
