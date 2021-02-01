package like

import (
	"context"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/test/gtest"
	"github.com/gomodule/redigo/redis"
	"github.com/rxrddd/gf-library/utils/redisx"
	"testing"
)

func TestNewPostLike(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		testSetup(New(context.Background(), &Option{
			RedisKey: "post_user_like",
		}), t)
		delKey("post_user_like*")
	})
}
func testSetup(like ILike, t *gtest.T) {
	var err error
	likeId := "1"
	userId := "1"
	flag, err := like.IsLike(likeId, userId)
	t.Assert(err, nil)
	t.Assert(flag, false)
	t.Log("==============isLike==============")
	err = like.Like(likeId, userId)
	t.Assert(err, nil)
	t.Log("==============Like==============")
	flag, err = like.IsLike(likeId, userId)
	t.Assert(err, nil)
	t.Assert(flag, true)
	t.Log("==============isLike==============")
	i, err := like.Count(likeId)
	t.Assert(err, nil)
	t.Assert(i, 1)
	t.Log("==============Count==============")
	err = like.UnLike(likeId, userId)
	t.Assert(err, nil)
	t.Log("==============UnLike==============")
	flag, err = like.IsLike(likeId, userId)
	t.Assert(err, nil)
	t.Assert(flag, false)
	t.Log("==============isLike==============")
	i, err = like.Count(likeId)
	t.Assert(err, nil)
	t.Assert(i, 0)
	t.Log("==============Count END==============")
}
func delKey(redisKey string) {
	keys, _ := g.Redis().DoVar("KEYS", redisKey)
	redisx.Multi(g.Redis().Conn(), func(con redis.Conn) {
		for _, key := range keys.Strings() {
			con.Send("DEL", key)
		}
	})
}
