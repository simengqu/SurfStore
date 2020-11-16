package surfstore

import (
	"fmt"
	"log"
	"os"

	// "bufio"
	"io/ioutil"
	// "strings"
	"crypto/sha256"
	"encoding/hex"
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
	// metaStore := MetaStore{
	// 	FileMetaMap: make(map[string]FileMetaData),
	// }
	// metaStore := new(MetaStore)
	// blockStore := BlockStore{
	// 	BlockMap: make(map[string]Block),
	// }
	bl := new(bool)
	*bl = true
	// fmm := new(map[string]FileMetaData)	
	var fileMetaMap = map[string]FileMetaData{}
	client.GetFileInfoMap(bl, &fileMetaMap)
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
					Version:       -1,
					BlockHashList: []string{},
				}
				// client.UpdateFile(&fileMetaData, &fileMetaData.Version)
				for {
					_, err = fi.Read(block.BlockData)
					
					h := sha256.Sum256(block.BlockData)
					he := hex.EncodeToString(h[:])
					fileMetaData.BlockHashList = append(fileMetaData.BlockHashList, he)
					
					client.PutBlock(block, bl)
					fmt.Println(string(block.BlockData))

					if err != nil {
						fmt.Println("Error when reading...", err)
						break
					}
				}
				fmt.Println("fileMetaData.BlockHashList...", fileMetaData.BlockHashList)
				client.UpdateFile(&fileMetaData, &fileMetaData.Version)
				// for k := range client.BlockMap {
				// 	fileMetaData.BlockHashList = append(fileMetaData.BlockHashList, k)
				// }
				// fmt.Println(fileMetaData.BlockHashList)

				// fi.Write(block.BlockData)
				// blockStore.PutBlock(block, bl)
				// hashList := metaStore.FileMetaMap[f.Name()].BlockHashList
				// hashList := fileMetaData.BlockHashList
				// l := f.Name() + "," + "1" + "," + strings.Trim(fmt.Sprint(hashList), "[]")
				// indexF.WriteString(l+"\n")
				// fmt.Println(l)
				// metaStore.FileMetaMap = *new(map[string]FileMetaData)
				// metaStore.FileMetaMap[f.Name()] = fileMetaData

							
				fileMetaMap[f.Name()] = fileMetaData
				fmt.Println("fileMetaData.Version after calling...", fileMetaData.Version)
			}
		}
	}
	// var mp map[string]FileMetaData
	// mp = metaStore.FileMetaMap
	// fmt.Println(metaStore.FileMetaMap)
	// client.GetFileInfoMap(bl, &metaStore.FileMetaMap)
	// fmt.Println(metaStore.FileMetaMap)
	// fmt.Println(fmm)
	
	// fmt.Println(*bl)

	for key, _ := range fileMetaMap {
		_, err := os.Stat(path + "/" + key)
		// files not in dir, sync files
		if os.IsNotExist(err) {
			fmt.Println("File does not exist.")
			fi, err := os.Create(path + "/" + key)
			if err != nil {
				fmt.Println("Error when creating index.txt ...", err)
			}
			fmt.Println("Creating file...", key)
			block := Block{
				BlockData: make([]byte, client.BlockSize),
				BlockSize: client.BlockSize,
			}
			fileMetaData := fileMetaMap[key]
			blockHashList := fileMetaData.BlockHashList

			for _, hs := range blockHashList {
				fmt.Println("hs...", hs)
				client.GetBlock(hs, &block)
				fmt.Println("BlockData...", string(block.BlockData))
				fi.WriteString(string(block.BlockData))
			}
			
			// _, err = fi.Read(block.BlockData)
			// if err != nil {
			// 	fmt.Println("Error when reading...", err)
			// 	break
			// }
			
			fi.Close()
			
		}
	}

	client.GetFileInfoMap(bl, &fileMetaMap)
	fmt.Println(fileMetaMap)
	// for _, f := range files {
	// 	fmt.Println(f.Name())
	// 	if _, ok := fileMetaMap[f.Name()]; ok {
	// 		fmt.Println("File found in base...")
	// 	} else {
	// 		_, err := os.Create(path + "/" + f.Name())
	// 		if err != nil {
	// 			fmt.Println("Error when creating index.txt ...", err)
	// 		}
	// 		fi, err := os.Open(path + "/" + f.Name())
	// 		if err != nil {
	// 			fmt.Println("No such file...")
	// 		}
	// 		fmt.Println("Creating file...", f.Name())
	// 		block := Block{
	// 			BlockData: make([]byte, client.BlockSize),
	// 			BlockSize: client.BlockSize,
	// 		}
	// 		client.GetBlock("ss", &block)
	// 		_, err = fi.Read(block.BlockData)
	// 		if err != nil {
	// 			fmt.Println("Error when reading...", err)
	// 			break
	// 		}
	// 		fi.Write(block.BlockData)

	// 		fmt.Println(string(block.BlockData))
	// 	}
	// }

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
