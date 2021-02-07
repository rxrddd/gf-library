package like

import (
	"context"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/test/gtest"
	"github.com/rxrddd/gf-library/utils/redisx"
	"testing"
)

func TestNewPostLike(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		testSetup(New(context.Background(), &Option{
			RedisKey: "post_user_like",
		}), t)
		redisx.DelKey(g.Redis(), "post_user_like*")
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
