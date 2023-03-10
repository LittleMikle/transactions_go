package main

func (s *server) routes() {
	s.echo.POST("/deposit", s.handleDeposit)
	s.echo.POST("/transfer", s.handleTransfer)
}
