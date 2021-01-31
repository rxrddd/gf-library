package redisx

import (
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
