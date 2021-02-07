package like

import (
	"context"
	"fmt"
	"github.com/gogf/gf/database/gredis"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gutil"
	"github.com/gomodule/redigo/redis"
	"github.com/rxrddd/gf-library/utils/redisx"
)

//点赞接口
type ILike interface {
	Like(likeId string, userId string) (err error)              //点赞
	UnLike(likeId string, userId string) (err error)            //取消赞
	IsLike(likeId string, userId string) (flag bool, err error) //是否点赞
	Count(likeId string) (count int64, err error)               //文章的点赞总数
}

const like = "like"
const unLike = "unlike"

//文章点赞
type defaultLike struct {
	opt   *Option
	redis *gredis.Redis
	commonAction
}

type Option struct {
	RedisKey       string                                         //redis key前缀
	CustomAttrFunc func(likeId string, userId string) g.MapStrAny //自定义增量属性hash的数据 用于异步入库时使用
}

func New(ctx context.Context, opt *Option) ILike {
	return &defaultLike{
		opt:          opt,
		commonAction: commonAction{RedisKey: opt.RedisKey},
		redis:        g.Redis().Ctx(ctx),
	}
}

type commonAction struct {
	RedisKey string
}

//增量变动的文章
func (l *commonAction) getChangeLikeKey() string {
	return fmt.Sprintf("%s_like_change_set", l.RedisKey)
}

//增量变动文章的用户
func (l *commonAction) getChangeLikeUserKey(likeId string) string {
	return fmt.Sprintf("%s_like_change_user_set_%s", l.RedisKey, likeId)
}

//每个文章的set
func (l *commonAction) getSetKey(likeId string) string {
	return fmt.Sprintf("%s_like_set_%s", l.RedisKey, likeId)
}

//每个文章的点赞数
func (l *commonAction) getCounterKey(likeId string) string {
	return fmt.Sprintf("%s_like_%s_counter", l.RedisKey, likeId)
}

//增量点赞的具体内容 比如点赞时间 点赞状态之内的
func (l *commonAction) getAttrKey(likeId string, userId string) string {
	return fmt.Sprintf("%s_like_attr_%s_%s", l.RedisKey, likeId, userId)
}
func (l *defaultLike) Like(likeId string, userId string) (err error) {
	flag, err := l.IsLike(likeId, userId)
	if err != nil {
		return err
	}
	if flag {
		return nil
	}
	return l.like(likeId, userId)
}

func (l *defaultLike) like(likeId string, userId string) (err error) {
	conn := l.redis.Conn()
	defer conn.Close()
	return redisx.Multi(conn, func(c redis.Conn) {
		c.Send("SADD", l.getChangeLikeKey(), likeId)
		c.Send("SADD", l.getChangeLikeUserKey(likeId), userId)
		c.Send("INCR", l.getCounterKey(likeId))
		c.Send("SADD", l.getSetKey(likeId), userId)
		attrs := l.handleAttr(likeId, userId, true)
		c.Send("HMSET", attrs...)
	})
}

func (l *defaultLike) handleAttr(likeId string, userId string, isLike bool) []interface{} {
	args := g.MapStrAny{}
	if isLike {
		args["createTime"] = gtime.Now().String()
		args["modifyTime"] = gtime.Now().String()
		args["status"] = like
	} else {
		args["modifyTime"] = gtime.Now().String()
		args["status"] = unLike
	}
	args["likeId"] = likeId
	args["userId"] = userId
	if l.opt.CustomAttrFunc != nil {
		gutil.MapMerge(args, l.opt.CustomAttrFunc(likeId, userId))
	}
	arr := make([]interface{}, 0)
	arr = append(arr, l.getAttrKey(likeId, userId))
	arr = append(arr, gutil.MapToSlice(args)...)
	return arr
}
func (l *defaultLike) UnLike(likeId string, userId string) (err error) {
	flag, err := l.IsLike(likeId, userId)
	if err != nil {
		return err
	}
	if !flag {
		return nil
	}
	return l.unLike(likeId, userId)
}
func (l *defaultLike) IsLike(likeId string, userId string) (flag bool, err error) {
	i, err := l.redis.DoVar("SISMEMBER", l.getSetKey(likeId), userId)
	if err != nil {
		return false, err
	}
	return i.Bool(), err
}
func (l *defaultLike) Count(likeId string) (count int64, err error) {
	c, err := l.redis.DoVar("GET", l.getCounterKey(likeId))
	if err != nil {
		return 0, err
	}
	return c.Int64(), err
}

func (l *defaultLike) unLike(likeId string, userId string) (err error) {
	conn := l.redis.Conn()
	defer conn.Close()
	return redisx.Multi(conn, func(c redis.Conn) {
		c.Send("SADD", l.getChangeLikeKey(), likeId)
		c.Send("SADD", l.getChangeLikeUserKey(likeId), userId)
		c.Send("DECR", l.getCounterKey(likeId))
		c.Send("SREM", l.getSetKey(likeId), userId)
		attrs := l.handleAttr(likeId, userId, false)
		c.Send("HMSET", attrs...)
	})
}
