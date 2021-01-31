### 通用点赞业务组件实现

**使用set和hash实现**
### 点赞
- 用于增量储存变动的文章id set
```
fmt.Sprintf("%s_like_change_set", l.RedisKey) 
```
- 用于增量储存变动的文章id对应的用户id set
```
fmt.Sprintf("%s_like_change_user_set_%s", l.RedisKey, likeId) 
```

- 用于增量储存文章-用户的点赞属性 如点赞时间,点赞类型(点赞,取消赞)
```
fmt.Sprintf("%s_like_attr_%s_%s", l.RedisKey, likeId, userId) 
```

- 用于文章的所有点赞用户id
```
fmt.Sprintf("%s_like_set_%s", l.RedisKey, likeId) 
```

- 用于文章的点赞总数
```
fmt.Sprintf("%s_like_%s_counter", l.RedisKey, likeId) 
```

### 提供以下接口
```
type ILike interface {
	Like(likeId string, userId string) (err error) //点赞
	UnLike(likeId string, userId string) (err error) //取消赞
	IsLike(likeId string, userId string) (flag bool, err error) //是否点赞
	Count(likeId string) (count int64, err error) //文章的点赞总数
}
```
### 点赞异步入库

> 使用方式:传入一个实现`IStorageAction`的对象就可以进行入库操作
`NewStorage(redisKey, new(defaultStorageAction)).AsyncStorage()`

```
type IStorageAction interface {
	Like(insertList []InsertData)            //入库为赞的操作
	UnLike(delData DelData)                  //入库为取消赞时的操作
	InsOrUpCount(likeId string, count int64) //入库修改数量
}
```
