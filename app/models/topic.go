package models

import (
	"database/sql"
	"github.com/robfig/revel"
	"time"
)

type Appendix struct {
	Apptm   time.Time
	Content string
}

type TopicInfo struct {
	Cmts        int32 // 评论数
	CmtUsers    int32 // 评论人数
	Views       int32 // 浏览数
	CountUp     int32 // 正向评价
	CountDown   int32 // 负向评价
	Favors      int32 // 收藏数/关注数
	Scores      int32 // 得分
	LastReplied time.Time
	lastreplied int64
	Appendixs   []*Appendix // 补充内容
}

type Topic struct {
	Id       int64
	Uid      int64  // Uid为0表示该topic是匿名发表
	Username string // 当Uid为0时，使用Username作为作者的名字显示
	Nid      int32  // Node ID
	Ttype    int32  // topic type
	Tstate   int32  // topic state
	Createtm time.Time
	createtm int64
	ClientIP string
	Title    string // 标题
	Content  string // 内容
	Length   uint32 // 内容长度
	Briefs   string // 简介
	Tags     string // 关键字
	Source   string // 来源
	Closed   bool   // 是否关闭
	Closeat  time.Time
	closeat  int64
	HeadPic  string // 配图

	Node *Node
	User *User

	*TopicInfo
}

func init_topic_table() {
	var (
		cls string = `CREATE TABLE IF NOT EXISTS topic (
	id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
	uid INT(10) NOT NULL DEFAULT '0',
	username VARCHAR(50) NOT NULL,
	nid INT(10) NOT NULL DEFAULT '0',
	ttype INT(10) NOT NULL DEFAULT '0',
	tstate INT(10) NOT NULL DEFAULT '0',
	createtm INT(10) NOT NULL DEFAULT '0',
	clientip VARCHAR(50) NOT NULL DEFAULT '',
	title VARCHAR(50) NOT NULL DEFAULT '',
	content TEXT NULL,
	length INT(10) NOT NULL DEFAULT '0',
	briefs VARCHAR(200) NOT NULL DEFAULT '',
	tags VARCHAR(50) NOT NULL DEFAULT '',
	source VARCHAR(160) NOT NULL DEFAULT '',
	closed INT(10) NOT NULL DEFAULT '0',
	closeat INT(10) NOT NULL DEFAULT '0',
	headpic VARCHAR(160) NOT NULL DEFAULT '',
	PRIMARY KEY (id)
)
COLLATE='utf8_general_ci'
ENGINE=InnoDB;`
	)
	db := get_db()

	_, err := db.Exec(cls)
	if err != nil {
		revel.ERROR.Panicf("create table topic failed: %v\n", err)
	}
}

// 根据tid从数据库中查找topic。查找不到返回nil
func Get_topic(tid int64) *Topic {
	db := get_db()
	row := db.QueryRow(`SELECT uid, username, nid, ttype, tstate, createtm, clientip, title, content, length, briefs, tags, source, closed, closeat, headpic FROM topic WHERE id=?`, tid)
	t := new(Topic)
	err := row.Scan(&t.Uid, &t.Username, &t.Nid, &t.Ttype, &t.Tstate, &t.createtm, &t.ClientIP, &t.Title, &t.Content, &t.Length, &t.Briefs, &t.Tags, &t.Source, &t.Closed, &t.closeat, &t.HeadPic)

	if err == sql.ErrNoRows {
		revel.ERROR.Printf("Not found topic by id %v\n", tid)
		return nil
	} else if err != nil {
		revel.ERROR.Printf("Get topic (id=%v) failed: %v\n", tid, err)
		return nil
	}

	if t.Closed && t.closeat != 0 {
		t.Closeat = time.Unix(t.closeat, 0)
	}

	t.Createtm = time.Unix(t.createtm, 0)

	return t
}

// 新建topic
func Add_topic(u *User, username, title, content, clientip, briefs, tags, source, headpic string, nid, ttype, tstate int32) int64 {
	db := get_db()
	uid := int64(0)
	uname := username
	if u != nil {
		uid = u.Id
		uname = u.Username
	}
	createtm := time.Now().Unix()
	res, err := db.Exec(`INSERT INTO topic (uid, username, nid, ttype, tstate, createtm, clientip, title, content, length, briefs, tags, source, closed, headpic) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,false,?);`,
		uid, uname, nid, ttype, tstate, createtm, clientip, title, content, len(content), briefs, tags, source, headpic)
	if err != nil {
		revel.ERROR.Panicf("Add topic failed: %v\n", err)
		return -1
	}
	tid, err := res.LastInsertId()
	if err != nil {
		revel.ERROR.Panicf("Get inserted topic's LastInsertId failed: %v\n", err)
		return -1
	}
	return tid
}

// 查找Topic的User
func (t *Topic) GetUser() {
	if t.Uid == 0 {
		t.User = nil
		return
	}
}

// 填充Topic的Node
func (t *Topic) GetNode() {
	if t.Nid == 0 {
		t.Node = nil
		return
	}

	t.Node = Get_node_by_id(int(t.Nid))
}

func (t *Topic) Reply() {

}

func (t *Topic) Countup() {

}

func (t *Topic) Countdown() {

}

func (t *Topic) Append() {

}

func (t *Topic) Modify() {

}

func (t *Topic) Close() {

}

func (t *Topic) Move_node() {

}

func (t *Topic) Favor() {

}

func (t *Topic) Score() {

}
