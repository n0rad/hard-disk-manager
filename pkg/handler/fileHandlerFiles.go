package handler

//func init() {
//	DiskHandlerBuilders["crypto"] = handler{
//		func() Handler {
//			return &HandlerInfo{
//				CommonHandler: CommonHandler{
//					handlerName: "Info",
//				},
//			}
//		},
//	}
//}
//
//type HandlerInfo struct {
//	CommonHandler
//}
//
//
//
//func (b *BlockDeviceOLD) Index() (string, error) {
//	if b.Mountpoint == "" {
//		return "", errs.WithF(b.fields, "Cannot index, disk is not mounted")
//	}
//	// todo this should be a stream
//	output, err := b.server.Exec("sudo find " + b.Mountpoint + " -type f -printf '%A@ %s %P\n'")
//	if err != nil {
//		return "", errs.WithEF(err, b.fields, "Failed to find files in filesystem")
//	}
//	return string(output), nil
//}