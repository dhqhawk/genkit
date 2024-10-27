package safelist

import (
	"container/list"
	"sync"
)

type SafeList struct {
	list *list.List
	lock sync.RWMutex
}

func NewSafeList() *SafeList {
	return &SafeList{
		list: list.New(),
		lock: sync.RWMutex{},
	}
}

func (s *SafeList) Len() int {
	return s.list.Len()
}

func (s *SafeList) Front() *list.Element {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.list.Front()
}

func (s *SafeList) Back() *list.Element {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.list.Back()
}

func (s *SafeList) PopBack() any {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.list.Remove(s.list.Back())
}

func (s *SafeList) Remove(e *list.Element) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.list.Remove(e)
}

func (s *SafeList) PushFront(v any) *list.Element {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.list.PushFront(v)
}

func (s *SafeList) PushBack(v any) *list.Element {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.list.PushBack(v)
}

func (s *SafeList) InsertBefore(v any, mark *list.Element) *list.Element {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.list.InsertBefore(v, mark)
}

func (s *SafeList) InsertAfter(v any, mark *list.Element) *list.Element {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.list.InsertAfter(v, mark)
}

func (s *SafeList) MoveToFront(e *list.Element) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.list.MoveToFront(e)
}

func (s *SafeList) MoveToBack(e *list.Element) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.list.MoveToBack(e)
}

func (s *SafeList) MoveBefore(e, mark *list.Element) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.list.MoveBefore(e, mark)
}

func (s *SafeList) MoveAfter(e, mark *list.Element) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.list.MoveAfter(e, mark)
}

func (s *SafeList) PushBackList(other *list.List) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.list.PushBackList(other)
}

func (s *SafeList) PushFrontList(other *list.List) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.list.PushFrontList(other)
}

func (s *SafeList) PopBackALl() []any {
	s.lock.Lock()
	defer s.lock.Unlock()
	var interlist []any
	for e := s.list.Back(); e != nil; e = s.list.Back() {
		interlist = append(interlist, s.list.Remove(e))
	}
	return interlist
}

func (s *SafeList) PopFrontAll() []any {
	s.lock.Lock()
	defer s.lock.Unlock()
	var interlist []any
	for e := s.list.Front(); e != nil; e = s.list.Front() {
		interlist = append(interlist, s.list.Remove(e))
	}
	return interlist
}

func (s *SafeList) PushBackAll(val ...any) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, v := range val {
		s.list.PushBack(v)
	}
}

func (s *SafeList) PushFrontAll(val ...any) {
	s.lock.Lock()
	defer s.lock.Unlock()
	for _, v := range val {
		s.list.PushFront(v)
	}
}
