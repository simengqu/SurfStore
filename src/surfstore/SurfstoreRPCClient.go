package surfstore

import (
	"fmt"
	"net/rpc"
)

type RPCClient struct {
	ServerAddr string
	BaseDir    string
	BlockSize  int
}

func (surfClient *RPCClient) GetBlock(blockHash string, block *Block) error {
	// connect to the server
	conn, e := rpc.DialHTTP("tcp", surfClient.ServerAddr)
	if e != nil {
		return e
	}

	// perform the call
	e = conn.Call("Server.GetBlock", blockHash, block)
	if e != nil {
		conn.Close()
		return e
	}

	// close the connection
	return conn.Close()
}

func (surfClient *RPCClient) PutBlock(block Block, succ *bool) error {
	// panic("todo")
	// connect to the server
	conn, e := rpc.DialHTTP("tcp", surfClient.ServerAddr)
	if e != nil {
		return e
	}

	// perform the call
	e = conn.Call("Server.PutBlock", block, succ)
	if e != nil {
		conn.Close()
		return e
	}

	// close the connection
	return conn.Close()

	// bs := new(BlockStore)
	// bs.PutBlock(block, succ)
}

func (surfClient *RPCClient) HasBlocks(blockHashesIn []string, blockHashesOut *[]string) error {
	// panic("todo")
	// connect to the server
	conn, e := rpc.DialHTTP("tcp", surfClient.ServerAddr)
	if e != nil {
		return e
	}

	// perform the call
	e = conn.Call("Server.HasBlocks", blockHashesIn, blockHashesOut)
	if e != nil {
		conn.Close()
		return e
	}

	// close the connection
	return conn.Close()
}

func (surfClient *RPCClient) GetFileInfoMap(succ *bool, serverFileInfoMap *map[string]FileMetaData) error {
	// panic("todo")
	// connect to the server
	fmt.Println("In client.GetFileInfoMap>>>>>>>>>", *serverFileInfoMap)
	conn, e := rpc.DialHTTP("tcp", surfClient.ServerAddr)
	if e != nil {
		return e
	}

	// perform the call
	e = conn.Call("Server.GetFileInfoMap", succ, serverFileInfoMap)
	if e != nil {
		conn.Close()
		return e
	}
	fmt.Println("after client.GetFileInfoMap>>>>>>>>>", *serverFileInfoMap)
	// close the connection
	return conn.Close()
}

func (surfClient *RPCClient) UpdateFile(fileMetaData *FileMetaData, latestVersion *int) error {
	// panic("todo")
	// connect to the server
	fmt.Println("in client.UpdateFile latestVersion>>>>>>>>>", *latestVersion)
	conn, e := rpc.DialHTTP("tcp", surfClient.ServerAddr)
	if e != nil {
		return e
	}

	// perform the call
	e = conn.Call("Server.UpdateFile", fileMetaData, latestVersion)
	if e != nil {
		conn.Close()
		return e
	}

	// close the connection
	return conn.Close()
}

var _ Surfstore = new(RPCClient)

// Create an Surfstore RPC client
func NewSurfstoreRPCClient(hostPort, baseDir string, blockSize int) RPCClient {

	return RPCClient{
		ServerAddr: hostPort,
		BaseDir:    baseDir,
		BlockSize:  blockSize,
	}
}
