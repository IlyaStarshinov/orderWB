package cache

import (
	"log"
	"sync"

	"gorm.io/gorm"

	"github.com/IlyaStarshinov/orderWB/internal/model"
)

type Cache struct {
	mu     sync.RWMutex
	orders map[string]*model.Order
}

func NewCache() *Cache {
	return &Cache{
		orders: make(map[string]*model.Order),
	}
}

func (c *Cache) Set(order *model.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.orders[order.OrderUID] = order
}

func (c *Cache) Get(uid string) (*model.Order, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	order, exists := c.orders[uid]

	if exists {
		log.Printf("КЭШ HIT для заказа: %s", uid) // Добавлено
	} else {
		log.Printf("КЭШ MISS для заказа: %s", uid) // Добавлено
	}

	return order, exists
}
func (c *Cache) RestoreFromDB(db *gorm.DB) error {
	var orders []model.Order
	if err := db.Preload("Delivery").Preload("Payment").Preload("Items").Find(&orders).Error; err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	for i := range orders {
		c.orders[orders[i].OrderUID] = &orders[i]
	}
	return nil
}
