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
	if len(keys.Strings()) <= 0 {
		return
	}
	c := conn.Conn()
	defer c.Close()
	_ = Multi(c, func(con redis.Conn) {
		for _, key := range keys.Strings() {
			con.Send("DEL", key)
		}
	})
}
