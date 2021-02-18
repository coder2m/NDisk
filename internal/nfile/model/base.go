package model

func GetFileByHash(ht, hc string) (f *File, err error) {
	f = new(File)
	if err = f.SetHash(ht, hc); err == nil {
		err = MainDB().Find(f).Error
	}
	return
}

func GetFileById(id uint) (f *File, err error) {
	f = new(File)
	f.Model.ID = id
	err = MainDB().Find(f).Error
	return
}
