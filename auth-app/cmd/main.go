package main

import "auth-app/internal/app/service"

func main() {
	s := service.NewService()
	err := s.Run()
	if err != nil {
		panic(err)
	}
}
