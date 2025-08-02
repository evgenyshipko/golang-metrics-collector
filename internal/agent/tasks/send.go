package tasks

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/requests"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/setup"
	"github.com/evgenyshipko/golang-metrics-collector/internal/agent/types"
	"github.com/evgenyshipko/golang-metrics-collector/internal/common/logger"
)

func SendMetricsTask(ctx context.Context, cfg setup.AgentStartupValues, dataChan <-chan types.MetricMessage, errChan chan<- error) {
	ticker := time.NewTicker(cfg.ReportInterval)

	requester := requests.NewRequester(cfg)

	var wg sync.WaitGroup

	for w := 1; w <= cfg.RateLimit; w++ {
		wg.Add(1)
		go worker(ctx, w, dataChan, errChan, requester, ticker, &wg)
	}

	go func() {
		<-ctx.Done()
		logger.Instance.Debug("SendMetricsTask <-ctx.Done()")
		ticker.Stop() // останавливаем тикер
	}()

	wg.Wait() // Ждём завершения всех горутин
}

/*
Почему несколько проверок ctx.Done()?
1. Основная проверка в начале цикла
go
case <-ctx.Done():

	return

Зачем: Это "быстрый выход" - если контекст уже завершён, мы сразу выходим, не дожидаясь тикера или новых заданий.

2. Проверка при ожидании нового задания
go
case <-ctx.Done():

	return

Зачем: Когда мы ждём новое задание из канала jobs, операция чтения может блокироваться. Мы хотим иметь возможность прервать это ожидание, если пришёл сигнал завершения.

3. Проверка при отправке ошибки
go
select {
case errChan <- err:
case <-ctx.Done():

	    return
	}

Зачем: Отправка в канал errChan тоже может блокироваться (если канал заполнен). Мы не хотим зависнуть в этой операции, если программа завершается.

Глубокое объяснение необходимости:
Блокирующие операции:

Чтение из канала (<-jobs)

Отправка в канал (errChan <- err)

Ожидание тикера (<-ticker.C)
Все эти операции могут блокировать выполнение

Отзывчивость при завершении:

# Хотим реагировать на завершение максимально быстро

Нельзя допустить, чтобы воркер "завис" в блокирующей операции

Ресурсы и утечки:

# Если воркер не завершится вовремя, это может привести к утечкам

Особенно важно в долгоживущих приложениях
*/
func worker(ctx context.Context, id int, jobs <-chan types.MetricMessage, errChan chan<- error, requester *requests.Requester,
	ticker *time.Ticker, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Printf("Worker %d starting\n", id)
	for {
		select {
		case <-ctx.Done(): // Получен сигнал завершения
			return
		case <-ticker.C:
			select {
			case job, ok := <-jobs:
				if !ok { // Канал закрыт
					return
				}
				if job.Err != nil {
					logger.Instance.Warnw("Обработка ошибки", "error", job.Err)
					continue
				}

				err := requester.SendMetric(job.Data.Type, job.Data.Name, job.Data.Value)
				if err != nil {
					select {
					case errChan <- err:
					case <-ctx.Done():
						return
					}
				}
			case <-ctx.Done():
				return
			}
		}
	}
}
