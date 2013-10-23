package models

import (
	"time"
)

type TopicInfo struct {
	Comments  int32 // 评论数
	CmtUsers  int32 // 评论人数
	Views     int32 // 浏览数
	CountUp   int32 // 正向评价
	CountDown int32 // 负向评价
	Favors    int32 // 收藏数/关注数
	Scores    int32 // 得分
}

type Topic struct {
	Id       int64
	Uid      int64  // Uid为0表示该topic是匿名发表
	Username string // 当Uid为0时，使用Username作为作者的名字显示
	Nid      int32  // Node ID
	Createtm time.Time
	createtm int64
	ClientIP string
	Title    string // 标题
	Content  string // 内容
	Length   uint32 // 内容长度
	Briefs   string // 简介
	Tags     string // 关键字
	Source   string // 来源
	Closed   string // 是否关闭
	Closeat  time.Time
	closeat  int64
	Appendix string // 补充内容
	HeadPic  string // 配图

	*TopicInfo

	*Node
	*User
}

func init_topic_table() {

}

func Get_topic() {

}

func Add_topic() {

}
func Reply_topic() {

}
func Countup_topic() {

}
func Countdown_topic() {

}
func Append_topic() {

}
func Close_topic() {

}

func Move_topic_node() {

}

func Favor_topic() {

}

func Score_topic() {

}
