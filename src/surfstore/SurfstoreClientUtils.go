package surfstore

import (
	"fmt"
	// "log"
	"os"

	// "bufio"
	"io/ioutil"
	"strings"
	"crypto/sha256"
	"encoding/hex"
	"bufio"
	"strconv"
)

/*
Implement the logic for a client syncing with the server here.
*/
func ClientSync(client RPCClient) {
	bl := new(bool) // bool to pass in rpc calls
	var tempmap = map[string]FileMetaData{}
	client.GetFileInfoMap(bl, &tempmap)
	// PrintMetaMap(tempmap)

	// read files in base
	path := client.BaseDir
	files, err := ioutil.ReadDir(path)
	if err != nil {
		// fmt.Println(err)
	}

	// create index.txt
	if _, err := os.Stat(path + "/index.txt"); err != nil {
		// fmt.Println("index.txt not found...")
		local_index, err := os.Create(path + "/index.txt")
		defer local_index.Close()
		if err != nil {
			// fmt.Println("Error when creating index.txt ...", err)
		}
	}

	// store local index in map
	var fileMetaMap_index = map[string]string{}
	var fileMetaMap_index_old = map[string]string{}
	f_index, err := os.Open(path + "/" + "index.txt")
	defer f_index.Close()
	if err != nil {
		// fmt.Println("Index.txt doesn't exist...")
	}
	rd := bufio.NewReader(f_index)
	for {
		line, err := rd.ReadString('\n')
		if err != nil {
			// fmt.Print("Error when reading index...", err)
			break
		}
		index_line := strings.Split(line, ",")
		fileMetaMap_index[index_line[0]] = strings.Trim(fmt.Sprint(line), "[]")
		fileMetaMap_index_old[index_line[0]] = strings.Trim(fmt.Sprint(line), "[]")
	}

	// create a map of FileMetaData(Filename, Version, BlockHashList[]), and obtain file info from the server
	var fileMetaMap = map[string]FileMetaData{}
	
	// scan the base dir for each file
	var file_name_base = map[string]int{}
	empty_base := true
	for _, f := range files {
		// fmt.Println("FileNames...", f.Name())
		
		if f.Name() != "index.txt" && f.Name() != ".DS_Store" { // ignore .DS_Store on mac
			empty_base = false
			file_name_base[f.Name()] = 0
			// open file
			fi, err := os.Open(path + "/" + f.Name())
			defer fi.Close()
			if err != nil {
				// fmt.Println("No such file...")
			}
			// fmt.Println("Writing in index.txt ...")
			block := Block{
				BlockData: make([]byte, client.BlockSize),
				BlockSize: client.BlockSize,
			}

			// in each file: read blocks, compute hash values
			blockHashList := []string{}
			fi_rd := bufio.NewReader(fi)
			blockList_base := []Block{}
			for {
				bl_byte_base := make([]byte, client.BlockSize)
				s_base, err := fi_rd.Read(bl_byte_base)
				if err != nil {
					// fmt.Println("Error when reading...", err)
					break
				}
				// fmt.Println("reading...", string(bl_byte_base[:s_base]))
				h := sha256.Sum256(bl_byte_base[:s_base])
				he := hex.EncodeToString(h[:])
				blockHashList = append(blockHashList, he)
				block.BlockData = bl_byte_base[:s_base]
				blockList_base = append(blockList_base, block)
				// fmt.Println("Block hash...", he)
			}

			// compare to local index
			f_index, err := os.Open(path + "/" + "index.txt")
			defer f_index.Close()
			if err != nil {
				// fmt.Println("Index.txt doesn't exist...")
			}
			rd := bufio.NewReader(f_index)
			new_file := true
			file_changed := false
			version_local := -1
			// read local index see if any new file or any file is changed
			for {
				line, err := rd.ReadString('\n')
				if err != nil {
					// fmt.Print("Error when reading index...", err)
					break
				}
				index_line := strings.Split(line, ",")
				// fmt.Println("file name in local index...", index_line[0], index_line)
				if f.Name() == index_line[0] {
					new_file = false
					hash_changed := strings.Trim(fmt.Sprint(blockHashList), "[]") + "\n"
					if hash_changed != index_line[2] {
						file_changed = true
						temp, _ := strconv.Atoi(index_line[1])
						version_local = temp
						// fmt.Println("File changed in local index...", f.Name())
						// fmt.Println("new hash...", hash_changed)
						// fmt.Println("old hash...", index_line[2])
					}
				}
			}

			client.GetFileInfoMap(bl, &fileMetaMap) // download an updated fileInfoMap
			// PrintMetaMap(fileMetaMap)
			if new_file {
				// fmt.Println("New file found in local index...", f.Name())
			}
			fileMetaData_base := FileMetaData{
				Filename:      f.Name(),
				Version:       1,
				BlockHashList: blockHashList,
			}
			// File in base is changed compared to local index
			if file_changed {
				// fmt.Println("File in base is changed compared to local index...", f.Name())
				if version_local != fileMetaMap[f.Name()].Version {
					// fmt.Println("Version in local and remote are diff...", version_local, fileMetaMap[f.Name()].Version)
				} else {
					// fmt.Println("Version in local and remote are same...", version_local)
					// sync local changes to cloud
					for _, bl_base := range blockList_base {
						// fmt.Println("bl_base...", bl_base)
						client.PutBlock(bl_base, bl)
					}
					version_local_update := new(int)
					*version_local_update = version_local + 1
					client.UpdateFile(&fileMetaData_base, version_local_update) // store in server
					fileMetaData_base.Version = *version_local_update
					client.GetFileInfoMap(bl, &fileMetaMap)
					// PrintMetaMap(fileMetaMap)
					l_base := f.Name() + "," + strconv.Itoa(fileMetaData_base.Version) + "," + strings.Trim(fmt.Sprint(blockHashList), "[]") + "\n"
					
					fileMetaMap_index[f.Name()] = l_base
					// fmt.Println("Writing to index...", l_base)
				}
			}


			// compare local index w/ remote index
			// if new file in local base dir that aren't in local index or remote index
			if _, ok := fileMetaMap[f.Name()]; ok { // in remote index
				// fmt.Println("File found in remote index...")
			} else { // file not in remote index
				// check if file in base is in local index, if in local index, check if file is changed
				// store blocks
				for _, bl_base := range blockList_base {
					// fmt.Println("bl_base...", bl_base)
					client.PutBlock(bl_base, bl)
				}
				
				version_update := new(int)
				*version_update = fileMetaData_base.Version
				client.UpdateFile(&fileMetaData_base, version_update) // store in server
				fileMetaData_base.Version = *version_update
				client.GetFileInfoMap(bl, &fileMetaMap)
				l_base := f.Name() + "," + strconv.Itoa(fileMetaData_base.Version) + "," + strings.Trim(fmt.Sprint(blockHashList), "[]") + "\n"
					
				fileMetaMap_index[f.Name()] = l_base
			}

			// file is deleted if its hashlist is tombstone in the map returned by GetFileInfoMap
			if v, ok := fileMetaMap[f.Name()]; ok {
				if strings.Trim(fmt.Sprint(v.BlockHashList), "[]") == "0" {
					fileMetaMap_index[f.Name()] = f.Name() + "," + strconv.Itoa(v.Version) + "," + "0" + "\n"
				}
			}
			
		}
	}

	// read local index see if  file is there in a client's index.txt but not in its base dir.
	for _, f := range files {
		file_name_base[f.Name()] = 0
	}
	// fmt.Println(fileMetaMap_index_old)
	for k, v := range fileMetaMap_index_old {
		// file is there in a client's index.txt but not in its base dir.
		if _, ok := file_name_base[k]; !ok {
			// fmt.Println("File is there in a client's index.txt but not in its base dir...", k)

			index_line := strings.Split(v, ",")
			tombVersion, _ := strconv.Atoi(index_line[1])
			metaData := FileMetaData {
				Filename: k,
				Version: 0,
				BlockHashList: []string{"0"},
			}
			metaData.Version = tombVersion
			if index_line[2] != "0\n" {
				client.UpdateFile(&metaData, &metaData.Version)
				tombVersion += 1
			}
			l_base := k + "," + strconv.Itoa(tombVersion) + "," + "0" + "\n"
			fileMetaMap_index[k] = l_base
			// fmt.Print("Deleted file...", fileMetaMap_index[k])
			
		}
	}

	if empty_base {
		// fmt.Println("Base is empty...")
	}
	// a file in remote index but not in local or base dir
	client.GetFileInfoMap(bl, &fileMetaMap)
	// PrintMetaMap(fileMetaMap)
	
	for k, v := range fileMetaMap {
		if _, ok := file_name_base[k]; ok {
			// fmt.Println("File in base...", k)
			hs_remote := strings.Trim(fmt.Sprint(v.BlockHashList), "[]")
			hs_index := strings.Split(fileMetaMap_index[k], ",")[2]
			// fmt.Println("hs_remote", hs_remote)
			// fmt.Println("hs_index", hs_index)
			if hs_remote + "\n" != hs_index {
				// download from server to base
				// fmt.Println("download from server to base")
				file_overwrite, err := os.Create(path + "/" + k)
				defer file_overwrite.Close()
				if err != nil {
					// fmt.Println("Error when open and overwriting file...", err)
				} else {
					// fmt.Println("Downloading file...", file_overwrite.Name())
				}
				line := v.BlockHashList // []string of hash values
				// fmt.Println("v block hash list...", v.BlockHashList)
				var block = new(Block)
				for _, hs := range line {
					client.GetBlock(hs, block)
					_, err := file_overwrite.Write(block.BlockData)
					if err != nil {
						// fmt.Println("Error when downloading file...", err)
					}
				}
				fileMetaMap_index[k] = k + "," + strconv.Itoa(v.Version) + "," + hs_remote + "\n"
			}
			
		} else {
			// fmt.Println("File not in base...", k)
			metaData := fileMetaMap[k]
			hs := strings.Trim(fmt.Sprint(metaData.BlockHashList), "[]")
			if hs != "0" { // file not deleted
				l_index := k + "," + strconv.Itoa(v.Version) + "," + strings.Trim(fmt.Sprint(v.BlockHashList), "[]") + "\n"
				// fmt.Println("Add to local index...", l_index)
				fileMetaMap_index[k] = l_index
				file_overwrite, err := os.Create(path + "/" + k)
				defer file_overwrite.Close()
				if err != nil {
					// fmt.Println("Error when open and overwriting file...", err)
				} else {
					// fmt.Println("Creating file...", file_overwrite.Name())
				}
				line := v.BlockHashList // []string of hash values
				// fmt.Println("v block hash list...", v.BlockHashList)
				var block = new(Block)
				// w := bufio.NewWriter(file_overwrite)
				for _, hs := range line {
					client.GetBlock(hs, block)
					// fmt.Println("block.BlockData...", block.BlockData)
					_, err := file_overwrite.Write(block.BlockData)
					if err != nil {
						// fmt.Println("Error when overwriting file...", err)
					}
				}
			} else {
				
			}
		}
		metaData := fileMetaMap[k]
		hs := strings.Trim(fmt.Sprint(metaData.BlockHashList), "[]")
		if hs == "0" {
			err := os.Remove(path + "/" + k)
			if err != nil {
				// fmt.Println("Fail to delete file...", err)
			}
		}
		
	}



	index_overwrite, err := os.OpenFile(path + "/index.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		// fmt.Println("Error when open and overwriting index...", err)
	}
	// fmt.Println("Overwriting index...")
	for _, v := range fileMetaMap_index {
		// fmt.Println("line in index...", k, v)
		_, err := index_overwrite.WriteString(v)
		if err != nil {
			// fmt.Println("Error when overwriting index...", err)
		}

	}

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
