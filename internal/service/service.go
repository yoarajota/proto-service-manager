package service

type Service struct {
	Name  string
	Start string
	Cwd   string
}

func FromConfig(cfg map[string]struct {
	Start string
	Cwd   string
}) map[string]Service {
	services := make(map[string]Service)

	for name, s := range cfg {
		services[name] = Service{
			Name:  name,
			Start: s.Start,
			Cwd:   s.Cwd,
		}
	}

	return services
}
