package shoppingcart

import (
	"fmt"
	"github.com/gogf/gf/database/gredis"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gutil"
	"github.com/gomodule/redigo/redis"
	"github.com/rxrddd/gf-library/utils/redisx"
)

type ICart interface {
	Create(userId string, item Item) (err error)                 //添加一个商品
	Remove(userId string, itemId string) (err error)             //删除一个商品
	Incr(userId string, itemId string) (err error)               //商品加1
	Decr(userId string, itemId string) (err error)               //商品减1
	Clear(userId string) (err error)                             //清除购物车
	List(userId string) (item []*Item, err error)                //购物车列表
	Count(userId string) (count int64, err error)                //购物车合计数量
	HasItem(userId string, itemId string) (flag bool, err error) //是否已经加入了购物车
}

// 单个商品item元素
type Item struct {
	ItemId     string          `json:"item_id"`
	Sku        string          `json:"sku"`
	Spu        string          `json:"spu"`
	Num        int32           `json:"num"`
	SalePrice  float64         `json:"sale_price"`           // 记录加车时候的销售价格
	PostFree   bool            `json:"post_free,omitempty"`  // 是否免邮
	Activities []*ItemActivity `json:"activities,omitempty"` // 参加的活动记录
	CreateTime int64           `json:"create_time"`
}

// 活动
type ItemActivity struct {
	ActID    string `json:"act_id"`    //活动id
	ActType  string `json:"act_type"`  //活动类型
	ActTitle string `json:"act_title"` //活动标题
}

type defaultCart struct {
	redis *gredis.Redis
}

func New(redis ...*gredis.Redis) ICart {
	var r *gredis.Redis
	if len(redis) > 0 {
		r = redis[0]
	} else {
		r = g.Redis()
	}
	return &defaultCart{
		redis: r,
	}
}

func (l *defaultCart) Create(userId string, item Item) (err error) {
	if flag, _ := l.HasItem(userId, item.ItemId); flag {
		return l.Incr(userId, item.ItemId)
	}
	return l.create(userId, item)
}

func (l *defaultCart) Remove(userId string, itemId string) (err error) {
	if flag, _ := l.HasItem(userId, itemId); !flag {
		return nil
	}
	return redisx.Multi(l.redis.Conn(), func(con redis.Conn) {
		con.Send("DEL", l.getHashKey(userId, itemId))
		con.Send("SREM", l.getItemsSetKey(userId), itemId)
	})
}

func (l *defaultCart) Incr(userId string, itemId string) (err error) {
	if flag, _ := l.HasItem(userId, itemId); !flag {
		return nil
	}
	_, err = l.redis.Do("HINCRBY", l.getHashKey(userId, itemId), "num", 1)
	return err
}

func (l *defaultCart) Decr(userId string, itemId string) (err error) {
	if flag, _ := l.HasItem(userId, itemId); !flag {
		return nil
	}
	v, err := l.redis.DoVar("HINCRBY", l.getHashKey(userId, itemId), "num", -1)
	if err != nil {
		return err
	}
	if v.Int64() <= 0 {
		return l.Remove(userId, itemId)
	}
	return nil
}

func (l *defaultCart) Clear(userId string) (err error) {
	keys, _ := l.redis.DoVar("KEYS", l.getAllHashKey(userId))
	return redisx.Multi(l.redis.Conn(), func(con redis.Conn) {
		con.Send("DEL", l.getItemsSetKey(userId))
		for _, key := range keys.Strings() {
			con.Send("DEL", key)
		}
	})
}

func (l *defaultCart) List(userId string) (items []*Item, err error) {
	var itms []*Item
	i, err := l.redis.DoVar("SMEMBERS", l.getItemsSetKey(userId))
	if err != nil {
		return nil, err
	}
	for _, value := range i.Strings() {
		doVar, err := l.redis.DoVar("HGETALL", l.getHashKey(userId, value))
		var itm *Item
		if err != nil {
			g.Log().Error(gerror.Wrap(err, "获取购物车列表错误"))
			break
		}
		err = doVar.Struct(&itm)
		if err == nil {
			itms = append(itms, itm)
		}
	}
	return itms, nil
}

func (l *defaultCart) Count(userId string) (count int64, err error) {
	i, err := l.redis.DoVar("SCARD", l.getItemsSetKey(userId))
	if err != nil {
		return 0, err
	}
	return i.Int64(), nil
}
func (l *defaultCart) HasItem(userId string, itemId string) (flag bool, err error) {
	i, err := l.redis.DoVar("SISMEMBER", l.getItemsSetKey(userId), itemId)
	if err != nil {
		return false, err
	}
	return i.Bool(), err
}

func (l *defaultCart) create(userId string, item Item) error {
	return redisx.Multi(l.redis.Conn(), func(con redis.Conn) {
		args := make([]interface{}, 0)
		args = append(args, l.getHashKey(userId, item.ItemId))
		args = append(args, gutil.MapToSlice(gconv.Map(item))...)
		con.Send("HMSET", args...)
		con.Send("SADD", l.getItemsSetKey(userId), item.ItemId)
	})
}

func (l *defaultCart) getHashKey(userId string, itemId string) string {
	return fmt.Sprintf("cart:user_%s:item_%s", userId, itemId)
}
func (l *defaultCart) getAllHashKey(userId string) string {
	return fmt.Sprintf("cart:user_%s:*", userId)
}
func (l *defaultCart) getItemsSetKey(userId string) string {
	return fmt.Sprintf("cart:user_%s", userId)
}
