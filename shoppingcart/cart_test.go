package shoppingcart

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
	"testing"
)

const userId = "1"
const itemId1 = "1"
const itemId2 = "2"
const itemId3 = "3"

var item = Item{
	ItemId:     itemId1,
	Sku:        "111",
	Spu:        "222",
	Num:        1,
	SalePrice:  100,
	CreateTime: gtime.Now().Timestamp(),
	CustomAttr: g.MapStrAny{
		"test":  1,
		"test2": 2,
	},
}
var item2 = Item{
	ItemId:     itemId2,
	Sku:        "111",
	Spu:        "222",
	Num:        1,
	SalePrice:  100,
	CreateTime: gtime.Now().Timestamp(),
	CustomAttr: g.MapStrAny{
		"test":  1,
		"test2": 2,
	},
}

func TestDefaultCart(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cart := setUp()
		t.AssertNE(cart, nil)
		t.Log("New")
		var err error
		t.Log("Create")
		err = cart.Create(userId, item)
		t.Assert(err, nil)
		err = cart.Create(userId, item2)
		t.Assert(err, nil)

		err = cart.Create(userId, item)
		t.Assert(err, nil)
		err = cart.Create(userId, item2)
		t.Assert(err, nil)
		t.Log("Count")
		count, err := cart.Count(userId)
		t.Assert(err, nil)
		t.Assert(count, 2)

		t.Log("Incr")
		err = cart.Incr(userId, itemId2)
		t.Assert(err, nil)
		t.Log("HasItem true")
		flag, err := cart.HasItem(userId, itemId1)
		t.Assert(err, nil)
		t.Assert(flag, true)

		t.Log("Decr")
		err = cart.Decr(userId, itemId1)
		t.Assert(err, nil)

		t.Log("HasItem false")
		flag, err = cart.HasItem(userId, itemId3)
		t.Assert(err, nil)
		t.Assert(flag, false)

		t.Log("HasItem List")
		list, err := cart.List(userId)
		t.Assert(err, nil)
		t.Assert(len(list), 2)

		err = cart.Remove(userId, itemId1)
		t.Assert(err, nil)

		flag, err = cart.HasItem(userId, itemId1)
		t.Assert(err, nil)
		t.Assert(flag, false)
		t.Log("Clear")
		err = cart.Clear(userId)
		t.Assert(err, nil)

		t.Log("List 0")
		list, err = cart.List(userId)
		t.Assert(err, nil)
		t.Assert(len(list), 0)

	})
}
func TestDefaultCart_2(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cart := setUp()
		_ = cart.Create(userId, item)
		_ = cart.Incr(userId, itemId1)

		_ = cart.Decr(userId, itemId1)
		_ = cart.Decr(userId, itemId1)
		_ = cart.Decr(userId, itemId1)
		flag, err := cart.HasItem(userId, itemId1)
		t.Assert(err, nil)
		t.Assert(flag, false)
	})
}
func TestDefaultCart_Removes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		cart := setUp()
		_ = cart.Create(userId, item)
		_ = cart.Create(userId, item2)
		err := cart.Removes(userId, []string{itemId1})
		t.Assert(err, nil)
		flag, err := cart.HasItem(userId, itemId1)
		flag2, err := cart.HasItem(userId, itemId2)
		t.Assert(err, nil)
		t.Assert(flag, false)
		t.Assert(flag2, true)
	})
}
func setUp() ICart {
	return New()
}
