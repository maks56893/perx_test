package cache

import (
	"sync"
	"time"
)

// Cache - структура кеша
type Cache struct {
	sync.RWMutex
	defaultExpiration time.Duration
	cleanupInterval   time.Duration
	items             map[string]item
}

type item struct {
	Value      interface{}
	Created    time.Time
	Expiration int64
}

// New - создать новый объект кеша
func New(defaultExpiration, cleanupInterval time.Duration) *Cache {

	// инициализируем карту(map) в паре ключ(string)/значение(Item)
	items := make(map[string]item)

	cache := Cache{
		items:             items,
		defaultExpiration: defaultExpiration,
		cleanupInterval:   cleanupInterval,
	}

	// Если интервал очистки больше 0, запускаем GC (удаление устаревших элементов)
	if cleanupInterval > 0 {
		cache.startGC() // данный метод рассматривается ниже
	}

	return &cache
}

// Set - положить значение с переданным ключом в кеш
func (c *Cache) Set(key string, value interface{}, duration time.Duration) {

	var expiration int64

	// Если продолжительность жизни равна 0 - используется значение по-умолчанию
	if duration == 0 {
		duration = c.defaultExpiration
	}

	// Устанавливаем время истечения кеша
	if duration > 0 {
		expiration = time.Now().Add(duration).UnixNano()
	}

	c.Lock()
	defer c.Unlock()

	c.items[key] = item{
		Value:      value,
		Expiration: expiration,
		Created:    time.Now(),
	}

	
}

// Get - получить объект из кеша по переданному ключу. Если объект не найден, вторым параметром возвращается false
func (c *Cache) Get(key string) (interface{}, bool) {

	c.RLock()

	defer c.RUnlock()

	item, found := c.items[key]

	// ключ не найден
	if !found {
		return nil, false
	}

	// Проверка на установку времени истечения, в противном случае он бессрочный
	if item.Expiration > 0 {

		// Если в момент запроса кеш устарел возвращаем nil
		if time.Now().UnixNano() > item.Expiration {
			return nil, false
		}

	}

	return item.Value, true
}

// Delete - удалить значение из кеша
func (c *Cache) Delete(key string) bool {

	c.Lock()

	defer c.Unlock()

	if _, found := c.items[key]; !found {
		return false
	}

	delete(c.items, key)

	return true
}

// Range - итерирование по объектам в кеше. Если переданная функция возвращаем false, то итерирование прекращается.
// При итерировании проверяется истекло ли время хранения объекта, если истекло, то объект удаляется
func (c *Cache) Range(f func(key string, value interface{}) bool) {
	for key, val := range c.items {
		if val.Expiration < time.Now().Unix() {
			c.Delete(key)
		}
        if !f(key, val.Value) {
            return
        }
	}
}


func (c *Cache) startGC() {
	go c.gc()
}

// gc - функция сборщик протухших значений
func (c *Cache) gc() {

	for {
		// ожидаем время установленное в cleanupInterval
		<-time.After(c.cleanupInterval)

		if c.items == nil {
			return
		}

		// Ищем элементы с истекшим временем жизни и удаляем из хранилища
		if keys := c.expiredKeys(); len(keys) != 0 {
			c.clearItems(keys)
		}

	}

}

// expiredKeys возвращает список "просроченных" ключей
func (c *Cache) expiredKeys() (keys []string) {

	c.RLock()

	defer c.RUnlock()

	for k, i := range c.items {
		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
			keys = append(keys, k)
		}
	}

	return
}

// clearItems удаляет ключи из переданного списка, в нашем случае "просроченные"
func (c *Cache) clearItems(keys []string) {

	c.Lock()

	defer c.Unlock()

	for _, k := range keys {
		delete(c.items, k)
	}
}
