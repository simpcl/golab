package strategy

const (
	BLOCK_SIZE                = 1024 * 1024
	OBJECT_MIN_INDEXES        = 32
	OBJECT_MAX_INDEXES        = 64
	DATAPARTGROUP_MIN_INDEXES = 256
	DATAPARTGROUP_MAX_INDEXES = 512
	DATAPART_MIN_INDEXES      = 1024
	DATAPART_MAX_INDEXES      = 2048
)

type DataRange struct {
	offset  int
	size    int
	indexes int
}

type DataBlock struct {
	DataRange
}

type DataPart struct {
	DataRange
	dumb   bool
	blocks []*DataBlock
}

func NewDataPart(offset int, dumb bool) *DataPart {
	return &DataPart{DataRange: DataRange{offset: offset, size: 0, indexes: 0}, dumb: dumb}
}

func (dp *DataPart) AddBlock(blockSize int) error {
	dp.size += blockSize
	dp.indexes += 1
	return nil
}

func (dp *DataPart) AddData(dataSize int) (int, error) {
	var err error
	var blockSize = 0
	for remainSize := dataSize; remainSize > 0; {
		if dp.dumb && dp.indexes >= OBJECT_MIN_INDEXES {
			dp.dumb = false
		} else if dp.indexes >= DATAPART_MIN_INDEXES {
			return remainSize, nil
		}
		if dataSize >= BLOCK_SIZE {
			blockSize = BLOCK_SIZE
		} else {
			blockSize = remainSize
		}
		err = dp.AddBlock(blockSize)
		if err != nil {
			return remainSize, err
		}
		remainSize -= blockSize
	}
	return 0, nil
}

func (dp *DataPart) GetBlocks() int {
	return dp.indexes
}

type DataPartGroup struct {
	DataRange
	dumb  bool
	parts []*DataPart
}

func NewDataPartGroup(offset int, dumb bool) *DataPartGroup {
	part := NewDataPart(offset, dumb)
	dps := []*DataPart{part}
	return &DataPartGroup{DataRange: DataRange{offset: offset, size: 0, indexes: 0}, dumb: dumb, parts: dps}
}

func (dpg *DataPartGroup) AddData(dataSize int) (int, error) {
	var err error
	remainSize := dataSize
	part := dpg.parts[len(dpg.parts)-1]
	for {
		remainSize, err = part.AddData(remainSize)
		if err != nil {
			return remainSize, err
		}
		if remainSize == 0 {
			return 0, nil
		}
		// create new datapart
		if part.dumb && dpg.indexes >= OBJECT_MIN_INDEXES {
			part.dumb = false
		} else if dpg.indexes >= DATAPART_MIN_INDEXES {
			return remainSize, nil
		}
		newpart := NewDataPart(part.offset+part.size, false)
		dpg.parts = append(dpg.parts, newpart)
		dpg.indexes += 1
		part = newpart
	}
	return 0, nil
}

type DataObject struct {
	size      int
	indexType int
	dpgs      []*DataPartGroup
}

func NewDataObject() *DataObject {
	dpgs := []*DataPartGroup{NewDataPartGroup(0, true)}
	dobj := &DataObject{size: 0, indexType: 0, dpgs: dpgs}

	return dobj
}

func (dobj *DataObject) AddData(dataSize int) error {
	var err error
	remainSize := dataSize
	dpg := dobj.dpgs[len(dobj.dpgs)-1]

	for {
		remainSize, err = dpg.AddData(remainSize)
		if err != nil {
			return err
		}
		if remainSize == 0 {
			break
		}
		newdpg := NewDataPartGroup(dpg.offset+dpg.size, false)
		dobj.dpgs = append(dobj.dpgs, newdpg)
		dpg = newdpg
	}

	dobj.size += dataSize
	dpg = dobj.dpgs[len(dobj.dpgs)-1]
	dobj.indexType = 0
	if dpg.dumb == false {
		dobj.indexType = 2
	} else if dpg.parts[len(dpg.parts)-1].dumb == false {
		dobj.indexType = 1
	}
	return nil
}
