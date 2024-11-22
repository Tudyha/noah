package mux

import (
	"fmt"
	"sync"
)

// IDManager 管理ID的分配与回收
type iDManager struct {
	maxID   uint32     // 当前最大ID
	freeIDs []uint32   // 可用ID列表
	mu      sync.Mutex // 用于同步访问
}

// NewIDManager 创建一个新的IDManager实例
func newIDManager() *iDManager {
	return &iDManager{
		maxID:   0,
		freeIDs: make([]uint32, 0),
	}
}

// GetID 获取一个新的ID
func (m *iDManager) GetID() (uint32, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if len(m.freeIDs) > 0 {
		id := m.freeIDs[0]
		m.freeIDs = m.freeIDs[1:]
		return id, nil
	}

	m.maxID++
	return m.maxID, nil
}

// ReleaseID 释放一个ID，使其可以被再次使用
func (m *iDManager) ReleaseID(id uint32) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查ID是否已经存在或超出范围
	if id <= 0 || id > m.maxID {
		return fmt.Errorf("无效的ID: %d", id)
	}

	// 防止重复添加
	for _, v := range m.freeIDs {
		if v == id {
			return nil
		}
	}

	m.freeIDs = append(m.freeIDs, id)
	return nil
}
