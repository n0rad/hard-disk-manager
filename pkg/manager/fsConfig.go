package manager

type FsConfig struct {
	PathConfig

	SearchConfigs bool
}

func (h *FsConfig) LoadFromDirIfExists(directory string) error {
	return loadFromDirIfExistsToStruct(directory, h)
}
