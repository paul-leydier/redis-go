package server

type storage struct {
	data map[string]string
}

func newStorage() *storage {
	return &storage{data: make(map[string]string)}
}

func (s *storage) Get(key string) string {
	return s.data[key]
}

func (s *storage) Set(key string, value string) {
	s.data[key] = value
}
