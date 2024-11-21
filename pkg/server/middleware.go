package server

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"music/pkg/log"

	"golang.org/x/time/rate"
)

var (
    // limiters будет хранить лимитеры для каждого IP-адреса
    limiters = make(map[string]*rate.Limiter)
    mu sync.Mutex
)

type writer struct {
    http.ResponseWriter
    statusCode int
}

// Middleware логирует запросы и время их выполнения
func Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        
        log.Logger.Info(fmt.Sprintf("Request %s %s", r.Method, r.URL.Path))
        
        wrappedWriter := &writer{w, http.StatusOK}
        next.ServeHTTP(wrappedWriter, r)

        log.Logger.Info(fmt.Sprintf("Completed %s %s with %d in %v", r.Method, r.URL.Path, wrappedWriter.statusCode, time.Since(start)))
    })
}

// getLimiter возвращает лимитер для заданного IP-адреса
func getLimiter(ip string) *rate.Limiter {
    mu.Lock()
    defer mu.Unlock()
    
    limiter, exists := limiters[ip]
    if !exists {
        limiter = rate.NewLimiter(1, 5) // 1 запрос в секунду с максимальным буфером в 5 запросов
        limiters[ip] = limiter
    }
    
    return limiter
}

// limitMiddleware ограничивает количество запросов с одного IP-адреса
func LimitMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := r.RemoteAddr
        limiter := getLimiter(ip)
        
        if !limiter.Allow() {
            http.Error(w, "Слишком много запросов", http.StatusTooManyRequests)
            return
        }
        
        next.ServeHTTP(w, r)
    })
}