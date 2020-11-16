package surfstore

import "errors"
import "fmt"

type MetaStore struct {
	FileMetaMap map[string]FileMetaData
}

func (m *MetaStore) GetFileInfoMap(_ignore *bool, serverFileInfoMap *map[string]FileMetaData) error {
	// panic("todo")
	fmt.Println("before metaStore.GetFileInfoMap>>>>>>>>>m", m.FileMetaMap)
	fmt.Println("before metaStore.GetFileInfoMap>>>>>>>>>s", *serverFileInfoMap)

	*serverFileInfoMap = m.FileMetaMap
	

	fmt.Println("after metaStore.GetFileInfoMap>>>>>>>>>m", m.FileMetaMap)
	fmt.Println("before metaStore.GetFileInfoMap>>>>>>>>>s", *serverFileInfoMap)
	return nil
}

func (m *MetaStore) UpdateFile(fileMetaData *FileMetaData, latestVersion *int) (err error) {
	// panic("todo")
	fmt.Println("in metaStore.UpdateFile latestVersion>>>>>>>>>", *latestVersion)

	// file exist
	if _, ok := m.FileMetaMap[fileMetaData.Filename]; ok {
		if fileMetaData.Version != *latestVersion+1 {
			return errors.New("The version is incorrect.")
		}

		fmt.Println("Updating fiels...")
		fileMetaData.Version += 1
		m.FileMetaMap[fileMetaData.Filename] = *fileMetaData

		*latestVersion += 1
		return nil
	} else { // file not exist
		fmt.Println("Adding new files...")
		m.FileMetaMap[fileMetaData.Filename] = *fileMetaData
		*latestVersion = 1
		fmt.Println("After adding...")
	}

	return nil
}

var _ MetaStoreInterface = new(MetaStore)
