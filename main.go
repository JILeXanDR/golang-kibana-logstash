package main

var logger = NewLogger("localhost:12201")

func main() {
	service1()
	service2()
	service3()
}

func service1() {
	logger.WithFields(map[string]interface{}{"module": "world"}).Debug("Hello!!!")
}
func service2() {
	logger.WithFields(map[string]interface{}{"module": "updater"}).Info("updating info...")
}
func service3() {
	logger.WithFields(map[string]interface{}{"module": "user"}).Error("failed to fetch user info")
}
