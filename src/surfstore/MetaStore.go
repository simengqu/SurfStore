package surfstore

// import "errors"
// import "fmt"

type MetaStore struct {
	FileMetaMap map[string]FileMetaData
}

func (m *MetaStore) GetFileInfoMap(_ignore *bool, serverFileInfoMap *map[string]FileMetaData) error {
	*serverFileInfoMap = m.FileMetaMap
	return nil
}

func (m *MetaStore) UpdateFile(fileMetaData *FileMetaData, latestVersion *int) (err error) {
	// file exist
	if _, ok := m.FileMetaMap[fileMetaData.Filename]; ok {
		fileMetaData.Version = m.FileMetaMap[fileMetaData.Filename].Version + 1
		m.FileMetaMap[fileMetaData.Filename] = *fileMetaData

		*latestVersion = fileMetaData.Version
		return nil
	} else { // file not exist
		m.FileMetaMap[fileMetaData.Filename] = *fileMetaData
		*latestVersion = 1
	}
	
	return nil
}

var _ MetaStoreInterface = new(MetaStore)
