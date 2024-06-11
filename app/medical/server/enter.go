package server

import (
	"WlFrame-gin/app/medical/dao"
	"WlFrame-gin/app/medical/model"
	"WlFrame-gin/utils/response"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 结果
func AddResult(ctx *gin.Context) {
	res := &model.Result{}
	if err := ctx.ShouldBindJSON(res); err != nil {
		panic("绑定失败")
	}
	result := dao.InsertResult(res)
	if result.RowsAffected != 0 {
		response.ResponseDML(ctx, result.RowsAffected, result.Error)
	}
}

func GetResultList(ctx *gin.Context) {
	role, result := dao.SelectResultsList()
	response.ResponseDQL(ctx, role, result.RowsAffected, result.RowsAffected, result.Error)
}

func GetResultById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("id属性转换为int64类型失败，错误原因：%v", err))
	}
	res, result := dao.SelectResultById(id)
	response.ResponseDQL(ctx, res, result.RowsAffected, result.RowsAffected, result.Error)
}

func DropResult(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("id属性转换为int64类型失败，错误原因：%v", err))
	}
	result := dao.DeleteResult(id)
	response.ResponseDML(ctx, result.RowsAffected, result.Error)
}

// 社区
func AddCommunity(ctx *gin.Context) {
	res := &model.Community{}
	if err := ctx.ShouldBindJSON(res); err != nil {
		panic("绑定失败")
	}
	result := dao.InsertCommunity(res)
	if result.RowsAffected != 0 {
		response.ResponseDML(ctx, result.RowsAffected, result.Error)
	}
}

func GetCommunityList(ctx *gin.Context) {
	list, result := dao.SelectCommunityList()
	response.ResponseDQL(ctx, list, result.RowsAffected, result.RowsAffected, result.Error)
}

func DropCommunity(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("id属性转换为int64类型失败，错误原因：%v", err))
	}
	result := dao.DeleteCommunity(id)
	response.ResponseDML(ctx, result.RowsAffected, result.Error)
}

// 物品
func AddGoods(ctx *gin.Context) {
	goods := &model.Goods{}
	if err := ctx.ShouldBindJSON(goods); err != nil {
		panic(fmt.Sprintf("%v", err))
	}
	result := dao.InsertGoods(goods)
	if result.RowsAffected != 0 {
		response.ResponseDML(ctx, result.RowsAffected, result.Error)
	}
}

func GetGoodsList(ctx *gin.Context) {
	list, result := dao.SelectGoodsList()
	response.ResponseDQL(ctx, list, result.RowsAffected, result.RowsAffected, result.Error)
}

func DropGoods(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("id属性转换为int64类型失败，错误原因：%v", err))
	}
	result := dao.DeleteGoods(id)
	response.ResponseDML(ctx, result.RowsAffected, result.Error)
}

// 推送
func AddPush(ctx *gin.Context) {
	msg := &model.PushMsg{}
	if err := ctx.ShouldBindJSON(msg); err != nil {
		panic("绑定失败")
	}
	result := dao.InsertMsg(msg)
	if result.RowsAffected != 0 {
		response.ResponseDML(ctx, result.RowsAffected, result.Error)
	}
}

func GetPushList(ctx *gin.Context) {
	list, result := dao.SelectMsgList()
	response.ResponseDQL(ctx, list, result.RowsAffected, result.RowsAffected, result.Error)
}

func DropPush(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("id属性转换为int64类型失败，错误原因：%v", err))
	}
	result := dao.DeleteMsg(id)
	response.ResponseDML(ctx, result.RowsAffected, result.Error)
}

// 居民
func AddPeople(ctx *gin.Context) {
	people := &model.People{}
	if err := ctx.ShouldBindJSON(people); err != nil {
		panic("绑定失败")
	}
	result := dao.InsertPeople(people)
	if result.RowsAffected != 0 {
		response.ResponseDML(ctx, result.RowsAffected, result.Error)
	}
}

func GetPeopleList(ctx *gin.Context) {
	list, result := dao.SelectPeoplesList()
	response.ResponseDQL(ctx, list, result.RowsAffected, result.RowsAffected, result.Error)
}

func DropPeople(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("id属性转换为int64类型失败，错误原因：%v", err))
	}
	result := dao.DeletePeople(id)
	response.ResponseDML(ctx, result.RowsAffected, result.Error)
}

// 反馈
func AddFeedback(ctx *gin.Context) {
	feedback := &model.Feedback{}
	if err := ctx.ShouldBindJSON(feedback); err != nil {
		panic("绑定失败")
	}
	result := dao.InsertFeedback(feedback)
	if result.RowsAffected != 0 {
		response.ResponseDML(ctx, result.RowsAffected, result.Error)
	}
}

func GetFeedbackList(ctx *gin.Context) {
	list, result := dao.SelectFeedbacksList()
	response.ResponseDQL(ctx, list, result.RowsAffected, result.RowsAffected, result.Error)
}

func DropFeedback(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		panic(fmt.Sprintf("id属性转换为int64类型失败，错误原因：%v", err))
	}
	result := dao.DeleteFeedback(id)
	response.ResponseDML(ctx, result.RowsAffected, result.Error)
}
