package models

import (
	"database/sql"
	"fmt"
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
	Depth       int
	Descp       string
	Extra_desc  string
	Protrait_xl string // 480 * 480
	Protrait_l  string // 240 * 240
	Protrait_m  string // 100 *100
	Protrait_s  string // 50 * 50
	Followers   int
	Createtm    time.Time

	// private use
	pid       uint32
	friendstr string
	createtm  int64
}

var (
	nodelock    sync.RWMutex
	catalogs    []*Node          // node of depth = 1
	subcatas    []*Node          // nodes of depth = 2
	leafnodes   []*Node          //node of depth = 3
	nodes       []*Node          = make([]*Node, 0)
	nodesbyID   map[uint32]*Node = make(map[uint32]*Node)
	nodesbyName map[string]*Node = make(map[string]*Node)
)

/*
Node 表结构如下
CREATE TABLE IF NOT EXISTS `nodes` (
	`id` INT(10) NOT NULL,
	`name` VARCHAR(50) NOT NULL,
	`depth` INT NOT NULL DEFAULT '0',
	`descp` VARCHAR(500) NULL DEFAULT '',
	`extra_desc` VARCHAR(500) NULL DEFAULT '',
	`pid` INT NULL DEFAULT '0',
	`friends` VARCHAR(300) NULL DEFAULT '',
	`protrait_xl` VARCHAR(120) NULL,
	`protrait_l` VARCHAR(120) NULL,
	`protrait_m` VARCHAR(120) NULL,
	`protrait_s` VARCHAR(120) NULL,
	`followers` INT NULL default '0',
	`createtm` INT NULL
)
COLLATE='utf8_general_ci'
ENGINE=InnoDB;
*/

func init_nodes() {
	db := get_db()
	rows, err := db.Query(`SELECT id, name, depth, descp, extra_desc, pid, friends, protrait_xl, protrait_l, protrait_m, protrait_s, followers, createtm FROM nodes;`)
	if err != nil {
		revel.ERROR.Panicf("select table nodes failed: %s\n", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		node := new(Node)
		err := rows.Scan(&node.Id, &node.Name, &node.Depth, &node.Descp, &node.Extra_desc, &node.pid,
			&node.friendstr, &node.Protrait_xl, &node.Protrait_l, &node.Protrait_m, &node.Protrait_s,
			&node.Followers, &node.createtm)
		if err != nil {
			revel.ERROR.Panicf("init nodes failed: %s\n", err.Error())
		}
		node.Createtm = time.Unix(node.createtm, 0)
		nodes = append(nodes, node)
	}

	if err = rows.Err(); err != nil {
		revel.ERROR.Panicf("init nodes failed: %s\n", err.Error())
	}

	build_nodes_relations()

	catalogs = make([]*Node, 0)
	subcatas = make([]*Node, 0)
	leafnodes = make([]*Node, 0)
	for _, n := range nodes {
		if n.Depth == 1 {
			catalogs = append(catalogs, n)
		} else if n.Depth == 2 {
			subcatas = append(subcatas, n)
		} else if n.Depth == 3 {
			leafnodes = append(leafnodes, n)
		} else {
			panic(fmt.Sprintf("Node (%s %d) depth invalid: %d\n", n.Name, n.Id, n.Depth))
		}
	}
}

// 建立node之间的关系
func build_nodes_relations() {
	for _, node := range nodes {
		node.Friends = make([]*Node, 0)
		node.Children = make([]*Node, 0)
		node.Silbings = make([]*Node, 0)
		node.Ancestors = make([]*Node, 0)
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

// get node by node id
func Get_node_by_id(id int) *Node {
	nodelock.RLock()
	defer nodelock.RUnlock()

	node, ok := nodesbyID[uint32(id)]
	if !ok {
		return nil
	}
	return node
}

// get node by name
func Get_node_by_name(name string) *Node {
	nodelock.RLock()
	defer nodelock.RUnlock()

	node, ok := nodesbyName[name]
	if !ok {
		return nil
	}

	return node
}

// add new node
func Add_new_node(name, descp, extra_desc string, pid int, friends, protrait string) (node *Node, err error) {
	var (
		id    int64
		res   sql.Result
		depth int
	)

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("Cannot connect to database: %s\n", err.Error())
	}

	if name == "" {
		return nil, fmt.Errorf("Node's name should not be empty.\n")
	}

	if descp == "" {
		return nil, fmt.Errorf("Node's description should not be empty.\n")
	}

	if pid > 0 {
		parent := Get_node_by_id(pid)
		if parent == nil {
			return nil, fmt.Errorf("Parent node(id=%d) not exists.\n", pid)
		}
		if parent.Depth >= 2 {
			return nil, fmt.Errorf("Parent node's Depth is %d, cannot add child node.\n", parent.Depth)
		}
		depth = parent.Depth + 1
	} else {
		depth = 1
	}

	node = &Node{Name: name,
		Descp:       descp,
		Extra_desc:  extra_desc,
		Depth:       depth,
		pid:         uint32(pid),
		Protrait_xl: protrait,
		Protrait_l:  protrait,
		Protrait_m:  protrait,
		Protrait_s:  protrait,
		Createtm:    time.Now(),
		createtm:    time.Now().Unix(),
		Followers:   0,
		friendstr:   friends,
	}

	res, err = db.Exec(`INSERT INTO nodes (name, depth, descp, extra_desc, pid, createtm, friends, protrait_xl, protrait_l, protrait_m, protrait_s) VALUES(?,?,?,?,?,?,?,?,?);`,
		node.Name, node.Depth, node.Descp, node.Extra_desc, node.pid, node.createtm, node.friendstr, node.Protrait_xl, node.Protrait_l, node.Protrait_m, node.Protrait_s)
	if err != nil {
		return nil, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}
	node.Id = uint32(id)

	nodelock.Lock()
	nodes = append(nodes, node)
	build_nodes_relations()
	if node.Depth == 1 {
		catalogs = append(catalogs, node)
	} else if node.Depth == 2 {
		subcatas = append(subcatas, node)
	}
	nodelock.Unlock()

	return node, nil
}

// modify node's basic attr, such as desc, extra_desc, protrait, followers
func Modify_node(node *Node, pc bool) error {
	if err := db.Ping(); err != nil {
		return fmt.Errorf("Cannot connect to database: %s\n", err.Error())
	}
	cls := ""
	if pc {
		cls = fmt.Sprintf(`UPDATE nodes SET descp=%s, extra_desc=%s, followers=%s, protrait_xl=%s, protrait_l=%s, protrait_m=%s, protrait_s=%s where id=%d;`,
			node.Descp, node.Extra_desc, node.Followers, node.Protrait_xl, node.Protrait_l, node.Protrait_m, node.Protrait_s, node.Id)
	} else {
		cls = fmt.Sprintf(`UPDATE nodes SET name=%s, descp=%s, extra_desc=%s, followers=%s where id=%d;`,
			node.Descp, node.Extra_desc, node.Followers, node.Id)
	}
	_, err := db.Exec(cls)
	if err != nil {
		return fmt.Errorf("Modify node %s failed: %s\n", node.Name, err.Error())
	}

	return nil
}

func Modify_node_name(node *Node, nname string) {
	var (
		fs    []string
		found bool
	)

	nodelock.Lock()
	defer nodelock.Unlock()

	for _, nd := range nodes {
		if nd.friendstr == "" {
			continue
		}
		fs = strings.Split(nd.friendstr, ",")
		nfs := ""
		found = false
		for _, fn := range fs {
			name := strings.TrimSpace(fn)
			if strings.EqualFold(name, node.Name) {
				found = true
				nfs += nname
				nfs += ","
			} else {
				nfs += fn
				nfs += ","
			}
		}
		if found {
			_, err := db.Exec(`UPDATE nodes SET friends=? where id=?;`, nfs, nd.Id)
			if err != nil {
				revel.ERROR.Panicf("Modify node(%s) name to %s cause update its friend node(%s) failed: %s\n", node.Name, nname, nd.Name, err.Error())
				return
			} else {
				nd.friendstr = nfs
			}
		}
	}
	node.Name = nname
}

func Modify_node_pid(node *Node, npid int) {
	_, err := db.Exec(`UPDATE nodes SET pid=? where id=?;`, npid, node.Id)
	if err != nil {
		revel.ERROR.Panicf("Update node(%s) pid from %d to %d failed: %s\n", node.Name, node.Id, npid, err.Error())
		return
	}
	nodelock.Lock()
	defer nodelock.Unlock()
	node.Id = uint32(npid)
	build_nodes_relations()
}

func Del_node(node *Node) error {
	var id int
	err := db.QueryRow(`SELECT id from nodes where pid=?;`, node.Id).Scan(&id)
	if err == sql.ErrNoRows {
		db.Exec(`DELETE nodes where id=?`, node.Id)
	} else if err != nil {
		revel.ERROR.Panic(err)
		return nil
	} else {
		return fmt.Errorf("Cannot delete node %s because of it has child node.\n", node.Name)
	}

	nodelock.Lock()
	defer nodelock.Unlock()
	var i int
	for i, _ = range nodes {
		if i == int(node.Id) {
			break
		}
	}
	if i < len(nodes)-1 {
		nodes = append(nodes[:i], nodes[i+1:]...)
	} else {
		nodes = nodes[:len(nodes)-1]
	}
	build_nodes_relations()

	return nil
}

func Clear_all_nodes() {
	_, err := db.Exec(`truncate nodes;`)
	if err != nil {
		revel.ERROR.Panic(err)
		return
	}
	nodelock.Lock()
	defer nodelock.Unlock()

	nodes = make([]*Node, 0)
	nodesbyID = make(map[uint32]*Node)
	nodesbyName = make(map[string]*Node)
}
