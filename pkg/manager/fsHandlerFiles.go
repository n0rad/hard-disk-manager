package manager

func init() {
	fsHandlerBuilders["files"] = fsHandlerBuilder{
		new: func() FsHandler {
			return &FsHandlerFiles{}
		},
	}
}

type FsHandlerFiles struct {
	CommonFsHandler
}


func (h *FsHandlerFiles) Add() error {



	// todo index files
	//

	return nil
}


//
//
//
//func (b *BlockDeviceOLD) Index() (string, error) {
//	if b.Mountpoint == "" {
//		return "", errs.WithF(b.fields, "Cannot index, disk is not mounted")
//	}
//	// todo this should be a stream
//	output, err := b.server.Exec("find " + b.Mountpoint + " -type f -printf '%A@ %s %P\n'")
//	if err != nil {
//		return "", errs.WithEF(err, b.fields, "Failed to find files in filesystem")
//	}
//	return string(output), nil
//}