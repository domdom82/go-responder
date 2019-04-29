package main

type Server interface {
	Run() error
	Stop() error
}
