package surfstore

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	// "fmt"
)

type Server struct {
	BlockStore BlockStoreInterface
	MetaStore  MetaStoreInterface
}

func (s *Server) GetFileInfoMap(succ *bool, serverFileInfoMap *map[string]FileMetaData) error {
	// panic("todo")
	// fmt.Println("before server.GetFileInfoMap>>>>>>>>>", *serverFileInfoMap)
	// fmt.Println("before server.GetFileInfoMap>>>>>>>>>", *succ)
	*succ = false
	s.MetaStore.GetFileInfoMap(succ, serverFileInfoMap)
	// fmt.Println("after server.GetFileInfoMap>>>>>>>>>", *serverFileInfoMap)
	return nil
}

func (s *Server) UpdateFile(fileMetaData *FileMetaData, latestVersion *int) error {
	// panic("todo")
	// fmt.Println("in server.UpdateFile latestVersion>>>>>>>>>", *latestVersion)
	s.MetaStore.UpdateFile(fileMetaData, latestVersion)
	return nil
}

func (s *Server) GetBlock(blockHash string, blockData *Block) error {
	// panic("todo")
	s.BlockStore.GetBlock(blockHash, blockData)
	return nil
}

func (s *Server) PutBlock(blockData Block, succ *bool) error {
	// panic("todo")
	s.BlockStore.PutBlock(blockData, succ)
	return nil
}

func (s *Server) HasBlocks(blockHashesIn []string, blockHashesOut *[]string) error {
	// panic("todo")
	s.BlockStore.HasBlocks(blockHashesIn, blockHashesOut)
	return nil
}

// This line guarantees all method for surfstore are implemented
var _ Surfstore = new(Server)

func NewSurfstoreServer() Server {
	blockStore := BlockStore{BlockMap: map[string]Block{}}
	metaStore := MetaStore{FileMetaMap: map[string]FileMetaData{}}

	return Server{
		BlockStore: &blockStore,
		MetaStore:  &metaStore,
	}
}

func ServeSurfstoreServer(hostAddr string, surfstoreServer Server) error {
	// panic("todo")
	// surfstoreServer = NewSurfstoreServer()
	rpc.Register(&surfstoreServer)
	// rpc.Register(surfstoreServer.BlockStore)
	// rpc.Register(surfstoreServer.MetaStore)
	
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", hostAddr)
	if e != nil {
		log.Println("listen error:", e)
	}
	http.Serve(l, nil)
	// fmt.Println("Press enter key to end server")
	// fmt.Scanln()
	return nil
}
