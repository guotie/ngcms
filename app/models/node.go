package models

import (
	"github.com/robfig/revel"
	"strings"
	"sync"
	"time"
)

type Node struct {
	Id          uint32
	Name        string
	Parent      *Node
	Ancestors   []*Node
	Children    []*Node
	Silbings    []*Node
	Friends     []*Node
	Descp       string
	Extra_desc  string
	Protrait_xl string // 480 * 480
	Protrait_l  string // 240 * 240
	Protrait_m  string // 100 *100
	Protrait_s  string // 50 * 50
	Followers   int

	// private use
	pid       uint32
	friendstr string
}

var (
	nodelock    sync.RWMutex
	nodes       []*Node          = make([]*Node, 0)
	nodesbyID   map[uint32]*Node = make(map[uint32]*Node)
	nodesbyName map[string]*Node = make(map[string]*Node)
)

/*
Node 表结构如下
CREATE TABLE IF NOT EXISTS `nodes` (
	`id` INT(10) NOT NULL,
	`name` VARCHAR(50) NULL,
	`descp` VARCHAR(500) NULL,
	`extra_desc` VARCHAR(500) NULL,
	`pid` INT NULL,
	`friends` VARCHAR(300) NULL,
	`protrait_xl` VARCHAR(120) NULL,
	`protrait_l` VARCHAR(120) NULL,
	`protrait_m` VARCHAR(120) NULL,
	`protrait_s` VARCHAR(120) NULL
)
COLLATE='utf8_general_ci'
ENGINE=InnoDB;
*/

func init_nodes() {
	db := get_db()
	rows, err := db.Query(`SELECT id, name, descp, extra_desc, pid, friends, protrait_xl, protrait_l, protrait_m, protrait_s FROM nodes;`)
	if err != nil {
		revel.ERROR.Panicf("select table nodes failed: %s\n", err.Error())
	}
	for rows.Next() {
		node := new(Node)
		err := rows.Scan(&node.Id, &node.Name, &node.Descp, &node.Extra_desc, &node.pid,
			&node.friendstr, &node.Protrait_xl, &node.Protrait_l, &node.Protrait_m, &node.Protrait_s)
		if err != nil {
			revel.ERROR.Panicf("init nodes failed: %s\n", err.Error())
		}
		node.Friends = make([]*Node, 0)
		node.Children = make([]*Node, 0)
		node.Silbings = make([]*Node, 0)
		node.Ancestors = make([]*Node, 0)
		nodes = append(nodes, node)
	}

	if err = rows.Err(); err != nil {
		revel.ERROR.Panicf("init nodes failed: %s\n", err.Error())
	}

	build_nodes_relations()
}

func build_nodes_relations() {
	for _, node := range nodes {
		build_node_parent(node)
		build_node_children(node)
		build_node_siblings(node)
		build_node_friends(node)

		nodesbyID[node.Id] = node
		nodesbyName[node.Name] = node
	}
}

func build_node_parent(node *Node) {
	for _, p := range nodes {
		if p.Id == node.pid {
			node.Parent = p
			return
		}
	}
}

func build_node_children(node *Node) {
	for _, c := range nodes {
		if c.pid == node.Id {
			node.Children = append(node.Children, c)
		}
	}
}

func build_node_siblings(node *Node) {
	for _, s := range nodes {
		if s.pid == node.pid {
			node.Silbings = append(node.Silbings, s)
		}
	}
}

func build_node_friends(node *Node) {
	fs := strings.Split(node.friendstr, ",")
	if len(fs) == 0 {
		return
	}

	for _, f := range fs {
		fn := strings.TrimSpace(f)
		found := false
		for _, n := range nodes {
			if strings.EqualFold(n.Name, fn) {
				node.Friends = append(node.Friends, n)
				found = true
				break
			}
		}
		if !found {
			revel.WARN.Printf("Not found node(%s) friend-node by name: %s\n", node.Name, fn)
		}
	}
}

func Get_node_by_id(id int) *Node {
	nodelock.RLock()
	defer nodelock.RUnlock()

	node, ok := nodesbyID[uint32(id)]
	if !ok {
		return nil
	}
	return node
}

func Get_node_by_name(name string) *Node {
	nodelock.RLock()
	defer nodelock.RUnlock()

	node, ok := nodesbyName[name]
	if !ok {
		return nil
	}

	return node
}
