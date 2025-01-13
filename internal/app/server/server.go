package server

type Server struct {
	articlesService ArticlesService
}

func NewServer(articlesService ArticlesService) Server {
	return Server{
		articlesService: articlesService,
	}
}
