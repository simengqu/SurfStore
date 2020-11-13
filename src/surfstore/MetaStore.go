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
	fmt.Println("before metaStore.GetFileInfoMap>>>>>>>>>", *_ignore)
	// *serverFileInfoMap = m.FileMetaMap
	*serverFileInfoMap = m.FileMetaMap
	// *serverFileInfoMap = nil
	fmt.Println("after metaStore.GetFileInfoMap>>>>>>>>>m", m.FileMetaMap)
	fmt.Println("before metaStore.GetFileInfoMap>>>>>>>>>s", *serverFileInfoMap)
	return nil
}

func (m *MetaStore) UpdateFile(fileMetaData *FileMetaData, latestVersion *int) (err error) {
	// panic("todo")

	if _, ok := m.FileMetaMap[fileMetaData.Filename]; ok {
		if fileMetaData.Version != *latestVersion+1 {
			return errors.New("The version is incorrect.")
		}


		fileMetaData.Version += 1
		m.FileMetaMap[fileMetaData.Filename] = *fileMetaData

		*latestVersion += 1
		return nil
	}

	return nil
}

var _ MetaStoreInterface = new(MetaStore)
