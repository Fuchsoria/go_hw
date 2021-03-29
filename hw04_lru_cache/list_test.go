package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("Normal and Reverse", func(t *testing.T) {
		l := NewList()

		toRemove := l.PushFront(10)
		l.PushFront(20)
		l.PushBack("a")
		l.PushBack("b")

		l.Remove(toRemove)

		var prevItem *ListItem = l.Front()
		var nextItem *ListItem = l.Back()
		items := []interface{}{prevItem.Value}
		itemsReverse := []interface{}{nextItem.Value}

		for {
			if item := prevItem.Next; item != nil {
				prevItem = item

				items = append(items, item.Value)
			} else {
				break
			}
		}

		for {
			if item := nextItem.Prev; item != nil {
				nextItem = item

				itemsReverse = append(itemsReverse, item.Value)
			} else {
				break
			}
		}

		require.Equal(t, []interface{}{20, "a", "b"}, items)
		require.Equal(t, []interface{}{"b", "a", 20}, itemsReverse)
		require.Equal(t, 3, l.Len())
	})

	t.Run("Front and Back", func(t *testing.T) {
		l := NewList()

		l.PushFront(10)
		l.PushFront(20)
		l.PushBack("a")
		l.PushBack("b")

		require.Equal(t, 20, l.Front().Value)
		require.Equal(t, "b", l.Back().Value)
	})

	t.Run("Front and Back after remove", func(t *testing.T) {
		l := NewList()

		l.PushFront(10)
		toRemove := l.PushFront(20)
		l.PushBack("a")
		toRemoveSecond := l.PushBack("b")

		l.Remove(toRemove)
		l.Remove(toRemoveSecond)

		require.Equal(t, 10, l.Front().Value)
		require.Equal(t, "a", l.Back().Value)
	})

	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}
