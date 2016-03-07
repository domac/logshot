package zookeeper

import (
	"errors"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"os"
	"path"
	"strings"
	"syscall"
	"time"
)

var (
	ErrNoChild      = errors.New("zk: children is nil")
	ErrNodeNotExist = errors.New("zk: node not exist")
)

//连接到zookeeper
func Connect(addr []string, timeout time.Duration) (*zk.Conn, error) {
	conn, _, err := zk.Connect(addr, timeout)
	if err != nil {
		log.Printf("zk.Connect(\"%v\", %d) error(%v)", addr, timeout, err)
		return nil, err
	}
	return conn, nil
}

func Create(conn *zk.Conn, fpath string) error {
	targetPath := ""
	for _, sourcePath := range strings.Split(fpath, "/")[1:] {
		targetPath = path.Join(targetPath, "/", sourcePath)
		fmt.Printf("create zookeeper path: \"%s\"", targetPath)
		_, err := conn.Create(targetPath, []byte(""), 0, zk.WorldACL(zk.PermAll))
		if err != nil {
			if err == zk.ErrNodeExists {
				fmt.Printf("zk.create(\"%s\") exists \n", targetPath)
			} else {
				fmt.Printf("zk.create(\"%s\") error(%v) \n", targetPath, err)
				return err
			}
		}
	}
	return nil
}

func process_quit() {
	if err := syscall.Kill(os.Getpid(), syscall.SIGQUIT); err != nil {
		fmt.Printf("syscall.Kill(%d, SIGQUIT) error(%v)", os.Getpid(), err)
	}
}

//创建节点,并返回watch
func GetNodesW(conn *zk.Conn, path string) ([]string, <-chan zk.Event, error) {
	nodes, stat, watch, err := conn.ChildrenW(path)
	if err != nil {
		if err == zk.ErrNoNode {
			return nil, nil, ErrNodeNotExist
		}
		fmt.Printf("zk.ChildrenW(\"%s\") error(%v)", path, err)
		return nil, nil, err
	}

	if stat == nil {
		return nil, nil, ErrNodeNotExist
	}
	if len(nodes) == 0 {
		return nil, nil, ErrNoChild
	}

	return nodes, watch, nil
}

func SetNodeData(conn *zk.Conn, path string, data []byte) error {
	_, err := conn.Set(path, data, -1)
	if err != nil {
		return err
	}
	return nil
}

func GetNodeData(conn *zk.Conn, path string) ([]byte, error) {
	data, _, _, err := conn.GetW(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//注册一个临时znode
func RegisterTempNode(conn *zk.Conn, fpath string, data []byte) error {
	targetPath, err := conn.Create(path.Join(fpath)+"/", data, zk.FlagEphemeral|zk.FlagSequence, zk.WorldACL(zk.PermAll))

	if err != nil {
		log.Printf("conn.Create(\"%s\", \"%s\", zk.FlagEphemeral|zk.FlagSequence) error(%v)", fpath, string(data), err)
		return err
	}
	log.Printf("create a zookeeper node: %s", targetPath)
	go func() {
		for {
			log.Printf("zk path: \"%s\" set a watch", targetPath)
			exist, _, watch, err := conn.ExistsW(targetPath)
			if err != nil {
				log.Printf("zk.ExistsW(\"%s\") error(%v)", targetPath, err)
				log.Printf("zk path: \"%s\" set watch failed, kill itself", targetPath)
				process_quit()
				return
			}
			if !exist {
				log.Printf("zk path: \"%s\" not exist, kill itself", targetPath)
				process_quit()
				return
			}
			event := <-watch
			log.Printf("zk path: \"%s\" receive a event %v", targetPath, event)
		}
	}()
	return nil
}
