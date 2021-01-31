### 关注组件实现

- 使用`redis`的`sortset`实现
### 提供以下接口
```
type IFollow interface {
	Follow(userId string, followId string) (err error)                                   //关注
	UnFollow(userId string, followId string) (err error)                                 //取消关注
	FansList(userId string) (list []string, err error)                                   //粉丝列表
	FollowList(userId string) (list []string, err error)                                 //关注列表
	IsMutualFollow(userId string, followId string) (flag bool, err error)                //是否互相关注
	IsFollow(userId string, followId string) (flag bool, err error)                      //是否关注对方
	IsFans(userId string, followId string) (flag bool, err error)                        //是否是对方的粉丝
	FollowCount(userId string) (count int64, err error)                                  //我关注的总数
	FansCount(userId string) (count int64, err error)                                    //我的粉丝总数
	CommonFollow(userId string, followId string) (list []string, err error)              //共同关注列表
	IsCommonFollow(userId string, followId string, findId string) (flag bool, err error) //是否共同关注了某人
}
```

