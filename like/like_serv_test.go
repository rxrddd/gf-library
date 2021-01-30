package like

import (
	"context"
	"github.com/gogf/gf/test/gtest"
	"github.com/gogf/gf/util/gconv"
	"testing"
)

func TestNewPostLike(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		testSetup(NewLike(context.Background(), &Option{
			RedisKey: "post_user_like",
		}), t)
		testSetup(NewLike(context.Background(), &Option{
			RedisKey: "vip_user_like",
		}), t)
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

func TestNewPostLike_2(t *testing.T) {
	likeId := "1"
	var err error
	opt := &Option{
		RedisKey: "post_user_like",
	}
	like := NewLike(context.Background(), opt)
	gtest.C(t, func(t *gtest.T) {
		for i := 0; i <= 10; i++ {
			err = like.Like(likeId, gconv.String(i))
			t.Assert(err, nil)
		}
		err = like.UnLike(likeId, "2")
		err = like.UnLike(likeId, "12")

	})
	gtest.C(t, func(t *gtest.T) {
		likeId2 := "2"
		for i := 0; i <= 10; i++ {
			err = like.Like(likeId2, gconv.String(i))
			t.Assert(err, nil)
		}
		err = like.UnLike(likeId2, "2")
		err = like.UnLike(likeId2, "7")
		count, err := like.Count("2")
		t.Assert(err, nil)
		t.Assert(count, 9)
	})
}
