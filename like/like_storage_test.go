package like

import (
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/test/gtest"
	"testing"
)

func TestNewStorage(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		s := &defaultStorageAction{
			db: g.DB(),
		}
		key := "post_user_like"
		NewStorage(key, s).AsyncStorage()
	})
}

//测试使用的接口实现
type defaultStorageAction struct {
	db gdb.DB
}

func (l *defaultStorageAction) Like(insertList []InsertData) {
	batchData := make([]g.MapStrAny, 0)
	for _, value := range insertList {
		batchData = append(batchData, g.MapStrAny{
			"like_id":     value.LikeId,
			"user_id":     value.UserId,
			"create_time": value.createTime,
		})
	}
	_, err := l.db.Table("user_like_post").Batch(100).Data(batchData).Insert()
	if err != nil {
		g.Log().Error(gerror.Wrap(err, "赞：操作数据库错误"))
	}
}
func (l *defaultStorageAction) UnLike(delData DelData) {
	_, err := l.db.Table("user_like_post").Unscoped().Delete(" like_id = ? and user_id in (?)", delData.LikeId, delData.UserIds)
	if err != nil {
		g.Log().Error(gerror.Wrap(err, fmt.Sprintf("取消赞:操作数据库错误 likeId:%s ", delData.LikeId)))
	}
}

func (l *defaultStorageAction) InsOrUpCount(likeId string, count int64) {
	ct, err := l.db.Table("post_like_count").FindCount("like_id = ?", likeId)
	if ct <= 0 {
		_, err = l.db.Table("post_like_count").Data(g.MapStrAny{
			"like_id": likeId,
			"count":   count,
		}).Insert()
	} else {
		_, err = l.db.Table("post_like_count").Where("like_id = ?", likeId).Data(g.MapStrAny{
			"count": count,
		}).Update()
	}
	if err != nil {
		g.Log().Error(gerror.Wrap(err, fmt.Sprintf("更新点赞总数错误 likeId:%s ", likeId)))
	}
}
