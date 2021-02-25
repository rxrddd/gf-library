package shoppingcart

import (
	"encoding/json"
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
	Create(userId string, item Item) (err error)                  //添加一个商品
	Remove(userId string, itemId string) (err error)              //删除一个商品
	Removes(userId string, itemId []string) (err error)           //删除多个商品
	Incr(userId string, itemId string) (err error)                //商品加1
	Decr(userId string, itemId string) (err error)                //商品减1
	Clear(userId string) (err error)                              //清除购物车
	List(userId string) (item []*Item, err error)                 //购物车列表
	GetItem(userId string, itemId string) (item *Item, err error) //获取一个购物车详情
	Count(userId string) (count int64, err error)                 //购物车合计数量
	HasItem(userId string, itemId string) (flag bool, err error)  //是否已经加入了购物车
}

// 单个商品item元素
type Item struct {
	ItemId     string      `json:"item_id"`
	Sku        string      `json:"sku"`
	Spu        string      `json:"spu"`
	Num        int64       `json:"num"`
	SalePrice  float64     `json:"sale_price"` // 记录加车时候的销售价格
	CreateTime int64       `json:"create_time"`
	CustomAttr g.MapStrAny `json:"custom_attr"` //自定义数据
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
	conn := l.redis.Conn()
	defer conn.Close()
	return redisx.Multi(conn, func(con redis.Conn) {
		con.Send("DEL", l.getHashKey(userId, itemId))
		con.Send("SREM", l.getItemsSetKey(userId), itemId)
	})
}

func (l *defaultCart) Removes(userId string, itemId []string) (err error) {
	if len(itemId) <= 0 {
		return
	}
	conn := l.redis.Conn()
	defer conn.Close()
	return redisx.Multi(conn, func(con redis.Conn) {
		for _, value := range itemId {
			if flag, _ := l.HasItem(userId, value); !flag {
				return
			}
			con.Send("DEL", l.getHashKey(userId, value))
			con.Send("SREM", l.getItemsSetKey(userId), value)
		}
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
	conn := l.redis.Conn()
	defer conn.Close()
	return redisx.Multi(conn, func(con redis.Conn) {
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
		if err != nil {
			g.Log().Error(gerror.Wrap(err, "获取购物车列表错误"))
			break
		}
		item, err := l.convMapToItem(doVar)
		if err != nil {
			g.Log().Error(gerror.Wrap(err, "redis数据转换为item数据错误"))
			break
		}
		itms = append(itms, item)
	}
	return itms, nil
}

func (l *defaultCart) GetItem(userId string, itemId string) (item *Item, err error) {
	if flag, err := l.HasItem(userId, itemId); !flag {
		return item, err
	}
	doVar, err := l.redis.DoVar("HGETALL", l.getHashKey(userId, itemId))
	return l.convMapToItem(doVar)
}

func (l *defaultCart) convMapToItem(doVar *g.Var) (item *Item, err error) {
	itemMap := doVar.Map()
	var customAttr g.MapStrAny
	if err = json.Unmarshal(gconv.Bytes(itemMap["custom_attr"]), &customAttr); err != nil {
		g.Log().Error(gerror.Wrap(err, "解析custom_attr错误"))
		return item, err
	}
	return &Item{
		ItemId:     gconv.String(itemMap["item_id"]),
		Sku:        gconv.String(itemMap["sku"]),
		Spu:        gconv.String(itemMap["spu"]),
		Num:        gconv.Int64(itemMap["num"]),
		SalePrice:  gconv.Float64(itemMap["sale_price"]),
		CreateTime: gconv.Int64(itemMap["create_time"]),
		CustomAttr: customAttr,
	}, nil
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
	conn := l.redis.Conn()
	defer conn.Close()
	return redisx.Multi(conn, func(con redis.Conn) {
		args := make([]interface{}, 0)
		args = append(args, l.getHashKey(userId, item.ItemId))
		jsonStr, _ := json.Marshal(item.CustomAttr)
		itemMap := gconv.Map(item)
		itemMap["custom_attr"] = jsonStr
		args = append(args, gutil.MapToSlice(itemMap)...)
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
