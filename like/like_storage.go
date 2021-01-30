package like

import (
	"github.com/gogf/gf/database/gredis"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
)

type ICacheLike interface {
	AsyncStorage() //异步入库
}
type defaultStorage struct {
	redisKey string
	redis    *gredis.Redis
	action   IStorageAction
	commonAction
}
type DelData struct {
	LikeId  string
	UserIds []string
}

type InsertData struct {
	LikeId     string
	UserId     string
	createTime string
	CustomAttr g.MapStrAny
}

type IStorageAction interface {
	Like(insertList []InsertData)            //入库为赞的操作
	UnLike(delData DelData)                  //入库为取消赞时的操作
	InsOrUpCount(likeId string, count int64) //入库修改数量
}

//入库操作方法 传入IStorageAction进行mysql入库操作
func NewStorage(redisKey string, action IStorageAction) ICacheLike {
	return &defaultStorage{
		redisKey:     redisKey,
		action:       action,
		redis:        g.Redis(),
		commonAction: commonAction{RedisKey: redisKey},
	}
}

func (l *defaultStorage) AsyncStorage() {
	likeIds, err := l.redis.DoVar("SMEMBERS", l.getChangeLikeKey())
	if err != nil {
		g.Log().Error(err)
		return
	}
	for _, likeId := range likeIds.Strings() {
		users, err := l.redis.DoVar("SMEMBERS", l.getChangeLikeUserKey(likeId))
		if err != nil {
			g.Log().Error(err)
			break
		}
		var insertData = make([]InsertData, 0)
		var delData DelData
		for _, userId := range users.Strings() {
			info, err := l.redis.DoVar("HGETALL", l.getAttrKey(likeId, userId))
			if err != nil {
				g.Log().Error(err)
				break
			}
			attrInfo := info.Map()
			if gconv.String(attrInfo["status"]) == like {
				insertData = append(insertData, InsertData{
					LikeId:     likeId,
					UserId:     userId,
					createTime: gconv.String(attrInfo["createTime"]),
					CustomAttr: attrInfo,
				})
			} else {
				delData.LikeId = likeId
				delData.UserIds = append(delData.UserIds, userId)
			}
			_, _ = l.redis.DoVar("DEL", l.getAttrKey(likeId, userId))
			_, _ = l.redis.DoVar("SREM", l.getChangeLikeUserKey(likeId), userId)
		}
		if len(insertData) > 0 {
			l.action.Like(insertData)
		}
		if delData.LikeId != "" {
			l.action.UnLike(delData)
		}
		c, _ := l.redis.DoVar("GET", l.getCounterKey(likeId))
		l.action.InsOrUpCount(likeId, c.Int64())
		_, _ = l.redis.DoVar("SREM", l.getChangeLikeKey(), likeId)
	}
}
