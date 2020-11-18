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
	*succ = false
	s.MetaStore.GetFileInfoMap(succ, serverFileInfoMap)
	return nil
}

func (s *Server) UpdateFile(fileMetaData *FileMetaData, latestVersion *int) error {
	s.MetaStore.UpdateFile(fileMetaData, latestVersion)
	return nil
}

func (s *Server) GetBlock(blockHash string, blockData *Block) error {
	s.BlockStore.GetBlock(blockHash, blockData)
	return nil
}

func (s *Server) PutBlock(blockData Block, succ *bool) error {
	s.BlockStore.PutBlock(blockData, succ)
	return nil
}

func (s *Server) HasBlocks(blockHashesIn []string, blockHashesOut *[]string) error {
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
	rpc.Register(&surfstoreServer)
	
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", hostAddr)
	if e != nil {
		log.Println("listen error:", e)
	}
	http.Serve(l, nil)
	return nil
}
