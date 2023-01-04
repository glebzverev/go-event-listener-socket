# Go Event-listener
Подключение осуществляется через соккет провайдера.
На гите была найдена ветка, где человек жаловался на то, что данныепо подписке перестают приходить через некоторое время. По какой-то причине данные перестают доходить. Интерфейсов ping-pong не нашлось, потому был создан таймер по истечение которого происходит переподклбчение к клиенту и соккету.
```go
if timer != nil {
	timer.Stop()
}
timer = time.NewTimer(time.Duration(time_in_seconds) * time.Second)
go func() {
	<-timer.C
	fmt.Println("Timer fired. Subscribe doesn't send any info ", time_in_seconds)
	sub.Unsubscribe()
	close(logs)
	AFTER_RECONNECT = true
	startBlockNumber = currentBlockNumber
	main()
}()
```

## Структура

-  В __recovery.go__ описана функция запроса эвентов за время переподключения
-  В __prepareEvent.go__ сырые функции выдающие наименование эвента и его TxHash
-  В __main.go__ описан цикл жизни программы, включающий подключение клиента, создание подписки и обработка ошибок  

## Пример выводимых данных в консоли

![alt text](https://github.com/glebzverev/go-event-listener-socket/blob/master/docs/console_example.png)
