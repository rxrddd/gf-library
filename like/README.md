### 通用点赞业务组件实现

**使用set和hash实现**

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

