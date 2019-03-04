package main

import "time"

var logger = NewLogger("localhost:12201")

func main() {
	logger.Debug("starting...")
	go func() {
		for ; ; {
			time.Sleep(1 * time.Second)
			service1()
		}
	}()
	go func() {
		for ; ; {
			time.Sleep(1 * time.Second)
			service2()
		}
	}()
	for ; ; {
		time.Sleep(1 * time.Second)
		service3()
	}
}

func service1() {
	logger.WithFields(map[string]interface{}{
		"service_name": "test_service 1",
	}).Debug("hey!")
}
func service2() {
	logger.WithFields(map[string]interface{}{
		"service_name": "test_service 2",
	}).Debug("how are you?")
}
func service3() {
	logger.WithFields(map[string]interface{}{
		"service_name": "test_service 3",
	}).Debug("hello!")
}
