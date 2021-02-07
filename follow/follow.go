package follow

import (
	"fmt"
	"github.com/gogf/gf/database/gredis"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gomodule/redigo/redis"
	"github.com/rxrddd/gf-library/utils/redisx"
)

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

//每个人有两个set 一个是自己的关注 一个是自己的粉丝

type defaultFollow struct {
	redis *gredis.Redis
}

func New() IFollow {
	return &defaultFollow{
		redis: g.Redis(),
	}
}
func (l *defaultFollow) getFansSetKey(userId string) string {
	return fmt.Sprintf("fans:user_%s", userId)
}
func (l *defaultFollow) getFollowSetKey(userId string) string {
	return fmt.Sprintf("follow:user_%s", userId)
}

func (l *defaultFollow) Follow(userId string, followId string) (err error) {
	if userId == followId {
		return gerror.New("不能对自己操作")
	}
	flag, err := l.IsFollow(userId, followId)
	if err != nil {
		return err
	}
	if flag {
		return nil
	}
	conn := l.redis.Conn()
	defer conn.Close()
	return redisx.Multi(conn, func(con redis.Conn) {
		con.Send("ZADD", l.getFollowSetKey(userId), gtime.Now().Timestamp(), followId)
		con.Send("ZADD", l.getFansSetKey(followId), gtime.Now().Timestamp(), userId)
	})
}

func (l *defaultFollow) UnFollow(userId string, followId string) (err error) {
	if userId == followId {
		return gerror.New("不能对自己操作")
	}
	flag, err := l.IsFollow(userId, followId)
	if err != nil {
		return err
	}
	if !flag {
		return nil
	}
	conn := l.redis.Conn()
	defer conn.Close()
	return redisx.Multi(conn, func(con redis.Conn) {
		con.Send("ZREM", l.getFollowSetKey(userId), followId)
		con.Send("ZREM", l.getFansSetKey(followId), userId)
	})
}

func (l *defaultFollow) FansList(userId string) (list []string, err error) {
	fans, err := l.redis.DoVar("ZRANGE", l.getFansSetKey(userId), 0, -1)
	var ls []string
	if err != nil {
		return ls, err
	}
	return fans.Strings(), nil
}

func (l *defaultFollow) FollowList(userId string) (list []string, err error) {
	fans, err := l.redis.DoVar("ZRANGE", l.getFollowSetKey(userId), 0, -1)
	var ls []string
	if err != nil {
		return ls, err
	}
	return fans.Strings(), nil
}

func (l *defaultFollow) IsMutualFollow(userId string, followId string) (flag bool, err error) {
	us, _ := l.redis.DoVar("ZRANK", l.getFollowSetKey(userId), followId)
	fs, _ := l.redis.DoVar("ZRANK", l.getFollowSetKey(followId), userId)
	return !us.IsNil() && !fs.IsNil(), nil
}

func (l *defaultFollow) IsFollow(userId string, followId string) (flag bool, err error) {
	us, _ := l.redis.DoVar("ZRANK", l.getFollowSetKey(userId), followId)
	return !us.IsNil(), nil
}
func (l *defaultFollow) IsFans(userId string, followId string) (flag bool, err error) {
	us, _ := l.redis.DoVar("ZRANK", l.getFansSetKey(userId), followId)
	return !us.IsNil(), nil
}
func (l *defaultFollow) FollowCount(userId string) (count int64, err error) {
	c, err := l.redis.DoVar("ZCARD", l.getFollowSetKey(userId))
	if err != nil {
		return 0, err
	}
	return c.Int64(), nil
}

func (l *defaultFollow) FansCount(userId string) (count int64, err error) {
	c, err := l.redis.DoVar("ZCARD", l.getFansSetKey(userId))
	if err != nil {
		return 0, err
	}
	return c.Int64(), nil
}

func (l *defaultFollow) CommonFollow(userId string, followId string) (list []string, err error) {
	key := fmt.Sprintf("common_follow_%s_%s", userId, followId)
	_, err = l.redis.DoVar("ZINTERSTORE", key, 2, l.getFollowSetKey(userId), l.getFollowSetKey(followId))
	if err != nil {
		return list, err
	}
	us, err := l.redis.DoVar("ZRANGE", key, 0, -1)
	if err != nil {
		return list, err
	}
	return us.Strings(), nil
}
func (l *defaultFollow) IsCommonFollow(userId string, followId string, findId string) (flag bool, err error) {
	us, _ := l.redis.DoVar("ZRANK", l.getFollowSetKey(userId), findId)
	fs, _ := l.redis.DoVar("ZRANK", l.getFollowSetKey(followId), findId)
	return !us.IsNil() && !fs.IsNil(), nil
}
