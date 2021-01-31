### 购物车业务组件实现

- 使用set和hash实现

- set用于储存每个用户的购物车商品id信息
- hash储存每一个用户-商品的信息

### 提供以下接口
```
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
```
