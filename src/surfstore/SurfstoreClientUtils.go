package surfstore

import (
	"fmt"
	"log"
	"os"

	// "bufio"
	"io/ioutil"
	// "strings"
)

/*
Implement the logic for a client syncing with the server here.
*/
func ClientSync(client RPCClient) {
	// panic("todo")
	log.Println("In client sync")
	path := client.BaseDir

	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}

	// if len(files) == 0 {
	// 	fmt.Println("Empty dir...")
	// } else {
	// 	for _, f := range files {
	// 		fmt.Println(f.Name())
	// 	}
	// }
	// create MetaStore
	metaStore := MetaStore{
		FileMetaMap: make(map[string]FileMetaData),
	}
	// metaStore := new(MetaStore)
	blockStore := BlockStore{
		BlockMap: make(map[string]Block),
	}
	bl := new(bool)
	*bl = true
	// fmm := new(map[string]FileMetaData)	
	var fmm = map[string]FileMetaData{}
	// fmt.Println(metaStore.FileMetaMap)
	// client.GetFileInfoMap(bl, &metaStore.FileMetaMap)
	// fmt.Println(metaStore.FileMetaMap)

	// create index.txt
	if _, err := os.Stat(path + "/index.txt"); err != nil {
		fmt.Println("index.txt not found...")
		// indexF, err := os.Create(path + "/index.txt")
		// if err != nil {
		// 	fmt.Println("Error when creating index.txt ...", err)
		// }
		for _, f := range files {
			if f.Name() != "index.txt" && f.Name() != ".DS_Store" {
				fi, err := os.Open(path + "/" + f.Name())
				if err != nil {
					fmt.Println("No such file...")
				}
				fmt.Println("Writing in index.txt ...")
				block := Block{
					BlockData: make([]byte, client.BlockSize),
					BlockSize: client.BlockSize,
				}
				// f.Read(block.BlockData)
				fileMetaData := FileMetaData{
					Filename:      f.Name(),
					Version:       1,
					BlockHashList: []string{},
				}
				for {
					_, err = fi.Read(block.BlockData)
					blockStore.PutBlock(block, bl)
					if err != nil {
						fmt.Println("Error when reading...", err)
						break
					}
					fmt.Println(string(block.BlockData))
				}

				for k := range blockStore.BlockMap {
					fileMetaData.BlockHashList = append(fileMetaData.BlockHashList, k)
				}
				fmt.Println(fileMetaData.BlockHashList)

				// fi.Write(block.BlockData)
				// blockStore.PutBlock(block, bl)
				// hashList := metaStore.FileMetaMap[f.Name()].BlockHashList
				// hashList := fileMetaData.BlockHashList
				// l := f.Name() + "," + "1" + "," + strings.Trim(fmt.Sprint(hashList), "[]")
				// indexF.WriteString(l+"\n")
				// fmt.Println(l)
				// metaStore.FileMetaMap = *new(map[string]FileMetaData)
				metaStore.FileMetaMap[f.Name()] = fileMetaData

							
				fmm[f.Name()] = fileMetaData
			}
		}
	}
	// var mp map[string]FileMetaData
	// mp = metaStore.FileMetaMap
	// fmt.Println(metaStore.FileMetaMap)
	// client.GetFileInfoMap(bl, &metaStore.FileMetaMap)
	// fmt.Println(metaStore.FileMetaMap)
	fmt.Println(fmm)
	client.GetFileInfoMap(bl, &fmm)
	fmt.Println(fmm)
	fmt.Println(*bl)
	for _, f := range files {
		fmt.Println(f.Name())
		// if strings.Contains(f.Name(), "index.txt") {
		// 	fmt.Println("index.txt found...")
		// 	hasIndex = true
		// }
	}

	// if files.Contains

	// index.txt
	// File1.dat,3,h0 h1 h2 h3

	// // read file
	// path := client.ServerAddr + client.BaseDir
	// fmt.Println("File path: ", path)
	// f, err := os.Open(path)
	// defer f.Close()
	// if err != nil {
	// 	fmt.Println("File not exist.")
	// }
	// buf := bufio.NewReader(f)
	// block := make([]byte, client.BlockSize)
	// // size, err := f.Read(block)
	// size, err := buf.Read(block)
	// if err != nil {
	// 	fmt.Println("Error when reading files.")
	// }
	// fmt.Println("Reading in", size, "bytes...")
}

/*
Helper function to print the contents of the metadata map.
*/
func PrintMetaMap(metaMap map[string]FileMetaData) {

	fmt.Println("--------BEGIN PRINT MAP--------")

	for _, filemeta := range metaMap {
		fmt.Println("\t", filemeta.Filename, filemeta.Version, filemeta.BlockHashList)
	}

	fmt.Println("---------END PRINT MAP--------")

}
