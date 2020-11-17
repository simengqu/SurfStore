package surfstore

import (
	"fmt"
	"log"
	"os"

	// "bufio"
	"io/ioutil"
	"strings"
	"crypto/sha256"
	"encoding/hex"
	"bufio"
)

/*
Implement the logic for a client syncing with the server here.
*/
func ClientSync(client RPCClient) {
	// panic("todo")
	log.Println("In client sync")

	// read files in base
	path := client.BaseDir
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
	}

	// create index.txt
	if _, err := os.Stat(path + "/index.txt"); err != nil {
		fmt.Println("index.txt not found...")
		local_index, err := os.Create(path + "/index.txt")
		defer local_index.Close()
		if err != nil {
			fmt.Println("Error when creating index.txt ...", err)
		}
		// local_index.WriteString("test string")
	}

	bl := new(bool) // bool to pass in rpc calls
	// fmm := new(map[string]FileMetaData)	

	// create a map of FileMetaData(Filename, Version, BlockHashList[]), and obtain file info from the server
	var fileMetaMap = map[string]FileMetaData{}
	// client.GetFileInfoMap(bl, &fileMetaMap)
	// fmt.Println(metaStore.FileMetaMap)
	// client.GetFileInfoMap(bl, &metaStore.FileMetaMap)
	// fmt.Println(metaStore.FileMetaMap)

	// k = file name, v = FileMetaData
	// var fileMetaMap_base = map[string]FileMetaData{}
	// k = file name, v = string[file name, version, h1 h2 h3...]
	var fileMetaMap_index = map[string]string{}

	// scan the base dir for each file
	for _, f := range files {
		fmt.Println("FileNames...", f.Name())
		if f.Name() != "index.txt" && f.Name() != ".DS_Store" { // ignore .DS_Store on mac
			// open file
			fi, err := os.Open(path + "/" + f.Name())
			defer fi.Close()
			if err != nil {
				fmt.Println("No such file...")
			}
			fmt.Println("Writing in index.txt ...")
			block := Block{
				BlockData: make([]byte, client.BlockSize),
				BlockSize: client.BlockSize,
			}
			// f.Read(block.BlockData)
			
			// read blocks, compute hash values
			blockHashList := []string{}
			fi_rd := bufio.NewReader(fi)
			blockList_base := []Block{}
			for {
				// if err != nil { // end of file
				// 	fmt.Println("Error when reading...", err)
				// 	break
				// }
				
				// _, err = fi.Read(block.BlockData)
				bl_byte_base := make([]byte, client.BlockSize)
				s_base, err := fi_rd.Read(bl_byte_base)
				if err != nil {
					fmt.Println("Error when reading...", err)
					break
				}
				fmt.Println("reading...", string(bl_byte_base[:s_base]))
				h := sha256.Sum256(bl_byte_base[:s_base])
				he := hex.EncodeToString(h[:])
				blockHashList = append(blockHashList, he)
				block.BlockData = bl_byte_base[:s_base]
				blockList_base = append(blockList_base, block)
				fmt.Println("blockHashList:", he)
				fmt.Println("BlockData:", string(block.BlockData[:s_base]))
				fmt.Println("blockList_base:", blockList_base)
				// client.PutBlock(block, bl)
				// fmt.Println(string(block.BlockData))
			}

			fmt.Println("Hash list...", blockHashList)
			fmt.Println("blockList_base...", blockList_base)
			// compare to local index
			f_index, err := os.Open(path + "/" + "index.txt")
			defer f_index.Close()
			if err != nil {
				fmt.Println("Index.txt doesn't exist...")
			}
			rd := bufio.NewReader(f_index)
			new_file := true
			file_changed := false
			for {
				line, err := rd.ReadString('\n')
				if err != nil {
					fmt.Print("Error when reading index...", err)
					break
				}
				index_line := strings.Split(line, ",")
				fileMetaMap_index[index_line[0]] = strings.Trim(fmt.Sprint(line), "[]")
				fmt.Println("file name in local index...", index_line[0], index_line)
				if f.Name() == index_line[0] {
					new_file = false
					if strings.Join(blockHashList, "") != index_line[2] {
						file_changed = true
					}
				}
			}
			if new_file {
				fmt.Println("New file found in local index...", f.Name())
			}
			if file_changed {
				fmt.Println("File changed in local index...", f.Name())
			}

			// fileMetaData := FileMetaData{
			// 	Filename:      f.Name(),
			// 	Version:       -1,
			// 	BlockHashList: []string{},
			// }
			client.GetFileInfoMap(bl, &fileMetaMap) // download an updated FileInfoMap

			// compare local index w/ remote index
			// if new file in local base dir that aren't in local index or remote index
			if _, ok := fileMetaMap[f.Name()]; ok {
				fmt.Println("File found in remote index...")
			} else { // file not in remote index
				fmt.Println("File not found in remote index...")
				// check if file in base is in local index, if in local index, check if file is changed
				if new_file { // new file that is not in remote index or local index
					fmt.Println("File not found in local index...")
					// store blocks
					for _, bl_base := range blockList_base {
						fmt.Println("bl_base...", bl_base)
						client.PutBlock(bl_base, bl)
					}
					fileMetaData_base := FileMetaData{
						Filename:      f.Name(),
						Version:       -1,
						BlockHashList: blockHashList,
					}
					client.UpdateFile(&fileMetaData_base, &fileMetaData_base.Version) // store in server

					l_base := f.Name() + "," + "1" + "," + strings.Trim(fmt.Sprint(blockHashList), "[]")
					
					fileMetaMap_index[f.Name()] = l_base
					fmt.Println("Writing to index...", l_base)
				} else { // file in local index but not in remote index
					fmt.Println("File in local index but not in remote index...")
					if file_changed { // file in base is diff from file in index
						fmt.Println("File in local index is changed...")
						l_local := f.Name() + "," + "1" + "," + strings.Trim(fmt.Sprint(blockHashList), "[]") + "\n"
						fileMetaMap_index[f.Name()] = l_local
						fmt.Println("Change", f.Name(), "to", l_local)
					}
				}
			}
			// }
			// for k, v := range fileMetaMap {
			// 	if 
			// }


			// fmt.Println("fileMetaData.BlockHashList...", fileMetaData.BlockHashList)
			// client.UpdateFile(&fileMetaData, &fileMetaData.Version) // update or adding new files

			// fi.Write(block.BlockData)
			// blockStore.PutBlock(block, bl)
			// hashList := metaStore.FileMetaMap[f.Name()].BlockHashList
			// hashList := fileMetaData.BlockHashList
			// l := f.Name() + "," + "1" + "," + strings.Trim(fmt.Sprint(hashList), "[]")
			// indexF.WriteString(l+"\n")
			// fmt.Println(l)
			// metaStore.FileMetaMap = *new(map[string]FileMetaData)
			// metaStore.FileMetaMap[f.Name()] = fileMetaData

							
			// fileMetaMap[f.Name()] = fileMetaData
			// fmt.Println("fileMetaData.Version after calling...", fileMetaData.Version)
		}
	}
	index_overwrite, err := os.OpenFile(path + "/index.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	fmt.Println("Overwriting index...")
	for k, v := range fileMetaMap_index {
		fmt.Println("line in index...", k, v)
		_, err := index_overwrite.WriteString(v)
		if err != nil {
			fmt.Println("Error when overwriting index...", err)
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

	// for key, _ := range fileMetaMap {
	// 	_, err := os.Stat(path + "/" + key)
	// 	// files not in dir, sync files
	// 	if os.IsNotExist(err) {
	// 		fmt.Println("File does not exist.")
	// 		fi, err := os.Create(path + "/" + key)
	// 		if err != nil {
	// 			fmt.Println("Error when creating index.txt ...", err)
	// 		}
	// 		fmt.Println("Creating file...", key)
	// 		block := Block{
	// 			BlockData: make([]byte, client.BlockSize),
	// 			BlockSize: client.BlockSize,
	// 		}
	// 		fileMetaData := fileMetaMap[key]
	// 		blockHashList := fileMetaData.BlockHashList

	// 		for _, hs := range blockHashList {
	// 			fmt.Println("hs...", hs)
	// 			client.GetBlock(hs, &block)
	// 			fmt.Println("BlockData...", string(block.BlockData))
	// 			fi.WriteString(string(block.BlockData))
	// 		}
			
	// 		// _, err = fi.Read(block.BlockData)
	// 		// if err != nil {
	// 		// 	fmt.Println("Error when reading...", err)
	// 		// 	break
	// 		// }
			
	// 		fi.Close()
			
	// 	}
	// }

	// client.GetFileInfoMap(bl, &fileMetaMap)
	// fmt.Println(fileMetaMap)
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
