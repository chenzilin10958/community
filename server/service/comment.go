package services

import (
	"errors"
	mapset "github.com/deckarep/golang-set/v2"
	"github.com/gin-gonic/gin"
	"xhyovo.cn/community/server/model"
	"xhyovo.cn/community/server/service/event"
)

type CommentsService struct {
	ctx *gin.Context
}

func NewCommentService(ctx *gin.Context) *CommentsService {
	return &CommentsService{ctx: ctx}
}

// 发布评论
func (a *CommentsService) Comment(comment *model.Comments) error {

	articleId := comment.BusinessId
	if f := articleDao.ExistById(articleId); !f {
		return errors.New("文章不存在")
	}

	// 父评论是否存在
	if comment.ParentId != 0 && commentDao.ExistById(comment.ParentId, comment.FromUserId, comment.BusinessId, comment.RootId) {
		return errors.New("回复评论不存在")
	}

	commentDao.AddComment(comment)
	var subscriptionService SubscriptionService

	subscriptionService.Do(&model.Subscriptions{EventId: event.CommentUpdateEvent, BusinessId: comment.BusinessId})
	return nil
}

// 删除评论
func (a *CommentsService) DeleteComment(id, userId int) bool {

	return commentDao.Delete(id, userId) == 1
}

// 查询文章下的评论
func (*CommentsService) GetCommentsByArticleID(page, limit, businessId int) ([]*model.Comments, int64) {

	var parentComments []*model.Comments
	childCommentsMap := make(map[int][]*model.Comments)
	comments, count := commentDao.GetCommentsByArticleID(page, limit, businessId)
	if count == 0 {
		return parentComments, 0
	}
	parentIds := make([]int, len(comments))
	userIds := make([]int, len(comments))
	// 收集根评论
	for i := range comments {
		comment := comments[i]
		if comment.ParentId == 0 {
			parentComments = append(parentComments, comment)
			parentIds = append(parentIds, comment.ID)
		} else {
			childCommentsMap[comment.RootId] = append(childCommentsMap[comment.RootId], comment)
		}
		userIds = append(userIds, comment.FromUserId)
	}

	setCommentUserNameAndArticleTitle(comments)

	ChildCommentNumberMap := commentDao.GetCommentsCountByRootId(parentIds)
	for i := range parentComments {
		parentComments[i].ChildComments = childCommentsMap[parentComments[i].RootId]
		parentComments[i].ChildCommentNumber = ChildCommentNumberMap[parentComments[i].RootId]
	}

	return parentComments, count
}

func setCommentUserNameAndArticleTitle(comments []*model.Comments) {
	userIds := mapset.NewSetWithSize[int](len(comments))
	articleIds := mapset.NewSetWithSize[int](len(comments))
	for i := range comments {
		comment := comments[i]
		articleIds.Add(comment.BusinessId)
		userIds.Add(comment.FromUserId)
		userIds.Add(comment.ToUserId)
	}

	if userIds.IsEmpty() {
		return
	}
	var u UserService
	userNameMap := u.ListByIdsSelectIdNameMap(userIds.ToSlice())

	var a ArticleService
	articleTitleMap := a.ListByIdsSelectIdTitleMap(articleIds.ToSlice())

	for i := range comments {
		comment := comments[i]
		comment.ArticleTitle = articleTitleMap[comment.BusinessId]
		comment.FromUserName = userNameMap[comment.FromUserId]
		if comment.ParentId != 0 {
			comment.ToUserName = userNameMap[comment.ToUserId]
		}
	}
}

// 查询文章下的所有评论(可指定)
func (*CommentsService) GetAllCommentsByArticleID(page, limit, userId, businessId int) ([]*model.Comments, int64) {
	comments, count := commentDao.GetAllCommentsByArticleID(page, userId, limit, businessId)
	if count == 0 {
		return comments, count
	}
	setCommentUserNameAndArticleTitle(comments)
	return comments, count
}

// 查询指定评论下的评论
func (*CommentsService) GetCommentsByRootID(page, limit, rootId int) ([]*model.Comments, int64) {

	var parentComments []*model.Comments
	commentsMap := make(map[int][]*model.Comments)
	comments, count := commentDao.GetCommentsByCommentID(page, limit, rootId)

	// 收集根评论
	for i := range comments {
		if comments[i].ParentId == 0 {
			parentComments = append(parentComments, comments[i])
		} else {
			commentsMap[comments[i].RootId] = append(commentsMap[comments[i].RootId], comments[i])
		}

	}
	// 如果根评论为空,说明是查询指定根评论下的子评论
	if len(parentComments) == 0 {
		return comments, count
	}
	for i := range parentComments {
		parentComments[i].ChildComments = commentsMap[parentComments[i].RootId]
	}

	return parentComments, count
}
