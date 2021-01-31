package follow

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/test/gtest"
	"github.com/rxrddd/gf-library/utils/redisx"
	"testing"
)

var (
	userId1 = "用户_1"
	userId2 = "用户_2"
	userId3 = "用户_3"
	userId4 = "用户_4"
)

func TestNewDefaultFollow(t *testing.T) {
	delKeys()
	gtest.C(t, func(t *gtest.T) {
		var err error
		fl := NewDefaultFollow()
		t.Log("Follow")
		err = fl.Follow(userId1, userId2)
		t.Assert(err, nil)

		t.Log("IsFollow")
		flag, err := fl.IsFollow(userId1, userId2)
		t.Assert(err, nil)
		t.Assert(flag, true)

		t.Log("FollowCount")
		count, err := fl.FollowCount(userId1)
		t.Assert(err, nil)
		t.Assert(count, 1)

		t.Log("IsFans")
		flag, err = fl.IsFans(userId2, userId1)
		t.Assert(err, nil)
		t.Assert(flag, true)

		t.Log("FollowCount")
		count, err = fl.FollowCount(userId1)
		t.Assert(err, nil)
		t.Assert(count, 1)

		t.Log("FansCount")
		count, err = fl.FansCount(userId2)
		t.Assert(err, nil)
		t.Assert(count, 1)

		t.Log("FollowList")
		list, err := fl.FollowList(userId1)
		t.Assert(err, nil)
		t.Assert(len(list), 1)

		t.Log("FansList")
		list, err = fl.FansList(userId2)
		t.Assert(err, nil)
		t.Assert(len(list), 1)

		t.Log("IsMutualFollow false")
		flag, err = fl.IsMutualFollow(userId1, userId2)
		t.Assert(err, nil)
		t.Assert(flag, false)

		t.Log("Follow userId2 - userId1")
		err = fl.Follow(userId2, userId1)
		t.Assert(err, nil)

		t.Log("IsMutualFollow true")
		flag, err = fl.IsMutualFollow(userId1, userId2)
		t.Assert(err, nil)
		t.Assert(flag, true)

		t.Log("IsMutualFollow true")
		flag, err = fl.IsMutualFollow(userId2, userId1)
		t.Assert(err, nil)
		t.Assert(flag, true)

		t.Log("CommonFollow start")
		err = fl.Follow(userId1, userId3)
		t.Assert(err, nil)
		t.Assert(flag, true)

		err = fl.Follow(userId2, userId3)
		t.Assert(err, nil)
		t.Assert(flag, true)

		lists, err := fl.CommonFollow(userId1, userId2)

		t.Assert(err, nil)
		t.Assert(len(lists), 1)
		t.Log("CommonFollow end")

		t.Log("IsCommonFollow false")
		flag, err = fl.IsCommonFollow(userId1, userId2, userId4)
		t.Assert(err, nil)
		t.Assert(flag, false)

		t.Log("IsCommonFollow true")
		flag, err = fl.IsCommonFollow(userId1, userId2, userId3)
		t.Assert(err, nil)
		t.Assert(flag, true)

		t.Log("UnFollow")
		err = fl.UnFollow(userId1, userId2)
		t.Assert(err, nil)

		t.Log("重复 UnFollow")
		err = fl.UnFollow(userId1, userId2)
		t.Assert(err, nil)

		t.Log("check UnFollow")
		flag, err = fl.IsFollow(userId1, userId2)
		t.Assert(err, nil)
		t.Assert(flag, false)

		t.Log("Follow/UnFollow self")
		err = fl.Follow(userId1, userId1)
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "不能对自己操作")
		err = fl.UnFollow(userId1, userId1)
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "不能对自己操作")

		err = fl.UnFollow(userId1, userId3)
		flag, err = fl.IsCommonFollow(userId1, userId2, userId3)
		t.Assert(err, nil)
		t.Assert(flag, false)

		lists, err = fl.CommonFollow(userId1, userId2)
		t.Assert(err, nil)
		t.Assert(len(lists), 0)
	})
	delKeys()
}

func delKeys() {
	redis := g.Redis()
	redisx.DelKey(redis, "fans:*")
	redisx.DelKey(redis, "follow:*")
}
