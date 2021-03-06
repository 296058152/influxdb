// Generated by tmpl
// https://github.com/benbjohnson/tmpl
//
// DO NOT EDIT!
// Source: array_cursor.gen.go.tmpl

package tsm1

import (
	"sort"

	"github.com/influxdata/influxdb/tsdb"
	"github.com/influxdata/influxdb/tsdb/cursors"
)

// Array Cursors

type floatArrayAscendingCursor struct {
	cache struct {
		values Values
		pos    int
	}

	tsm struct {
		buf       *tsdb.FloatArray
		values    *tsdb.FloatArray
		pos       int
		keyCursor *KeyCursor
	}

	end   int64
	res   *tsdb.FloatArray
	stats cursors.CursorStats
}

func newFloatArrayAscendingCursor() *floatArrayAscendingCursor {
	c := &floatArrayAscendingCursor{
		res: tsdb.NewFloatArrayLen(MaxPointsPerBlock),
	}
	c.tsm.buf = tsdb.NewFloatArrayLen(MaxPointsPerBlock)
	return c
}

func (c *floatArrayAscendingCursor) reset(seek, end int64, cacheValues Values, tsmKeyCursor *KeyCursor) {
	c.end = end
	c.cache.values = cacheValues
	c.cache.pos = sort.Search(len(c.cache.values), func(i int) bool {
		return c.cache.values[i].UnixNano() >= seek
	})

	c.tsm.keyCursor = tsmKeyCursor
	c.tsm.values = c.readArrayBlock()
	c.tsm.pos = sort.Search(c.tsm.values.Len(), func(i int) bool {
		return c.tsm.values.Timestamps[i] >= seek
	})
}

func (c *floatArrayAscendingCursor) Err() error { return nil }

// close closes the cursor and any dependent cursors.
func (c *floatArrayAscendingCursor) Close() {
	if c.tsm.keyCursor != nil {
		c.tsm.keyCursor.Close()
		c.tsm.keyCursor = nil
	}
	c.cache.values = nil
	c.tsm.values = nil
}

func (c *floatArrayAscendingCursor) Stats() cursors.CursorStats { return c.stats }

// Next returns the next key/value for the cursor.
func (c *floatArrayAscendingCursor) Next() *tsdb.FloatArray {
	pos := 0
	cvals := c.cache.values
	tvals := c.tsm.values

	c.res.Timestamps = c.res.Timestamps[:cap(c.res.Timestamps)]
	c.res.Values = c.res.Values[:cap(c.res.Values)]

	for pos < len(c.res.Timestamps) && c.tsm.pos < len(tvals.Timestamps) && c.cache.pos < len(cvals) {
		ckey := cvals[c.cache.pos].UnixNano()
		tkey := tvals.Timestamps[c.tsm.pos]
		if ckey == tkey {
			c.res.Timestamps[pos] = ckey
			c.res.Values[pos] = cvals[c.cache.pos].(FloatValue).RawValue()
			c.cache.pos++
			c.tsm.pos++
		} else if ckey < tkey {
			c.res.Timestamps[pos] = ckey
			c.res.Values[pos] = cvals[c.cache.pos].(FloatValue).RawValue()
			c.cache.pos++
		} else {
			c.res.Timestamps[pos] = tkey
			c.res.Values[pos] = tvals.Values[c.tsm.pos]
			c.tsm.pos++
		}

		pos++

		if c.tsm.pos >= len(tvals.Timestamps) {
			tvals = c.nextTSM()
		}
	}

	if pos < len(c.res.Timestamps) {
		if c.tsm.pos < len(tvals.Timestamps) {
			if pos == 0 {
				// optimization: all points served from TSM data
				copy(c.res.Timestamps, tvals.Timestamps)
				pos += copy(c.res.Values, tvals.Values)
				c.nextTSM()
			} else {
				// copy as much as we can
				n := copy(c.res.Timestamps[pos:], tvals.Timestamps[c.tsm.pos:])
				copy(c.res.Values[pos:], tvals.Values[c.tsm.pos:])
				pos += n
				c.tsm.pos += n
				if c.tsm.pos >= len(tvals.Timestamps) {
					c.nextTSM()
				}
			}
		}

		if c.cache.pos < len(cvals) {
			// TSM was exhausted
			for pos < len(c.res.Timestamps) && c.cache.pos < len(cvals) {
				c.res.Timestamps[pos] = cvals[c.cache.pos].UnixNano()
				c.res.Values[pos] = cvals[c.cache.pos].(FloatValue).RawValue()
				pos++
				c.cache.pos++
			}
		}
	}

	if pos > 0 && c.res.Timestamps[pos-1] >= c.end {
		pos -= 2
		for pos >= 0 && c.res.Timestamps[pos] >= c.end {
			pos--
		}
		pos++
	}

	c.res.Timestamps = c.res.Timestamps[:pos]
	c.res.Values = c.res.Values[:pos]

	c.stats.ScannedValues += len(c.res.Values)

	c.stats.ScannedBytes += len(c.res.Values) * 8

	return c.res
}

func (c *floatArrayAscendingCursor) nextTSM() *tsdb.FloatArray {
	c.tsm.keyCursor.Next()
	c.tsm.values = c.readArrayBlock()
	c.tsm.pos = 0
	return c.tsm.values
}

func (c *floatArrayAscendingCursor) readArrayBlock() *tsdb.FloatArray {
	values, _ := c.tsm.keyCursor.ReadFloatArrayBlock(c.tsm.buf)
	return values
}

type floatArrayDescendingCursor struct {
	cache struct {
		values Values
		pos    int
	}

	tsm struct {
		buf       *tsdb.FloatArray
		values    *tsdb.FloatArray
		pos       int
		keyCursor *KeyCursor
	}

	end   int64
	res   *tsdb.FloatArray
	stats cursors.CursorStats
}

func newFloatArrayDescendingCursor() *floatArrayDescendingCursor {
	c := &floatArrayDescendingCursor{
		res: tsdb.NewFloatArrayLen(MaxPointsPerBlock),
	}
	c.tsm.buf = tsdb.NewFloatArrayLen(MaxPointsPerBlock)
	return c
}

func (c *floatArrayDescendingCursor) reset(seek, end int64, cacheValues Values, tsmKeyCursor *KeyCursor) {
	c.end = end
	c.cache.values = cacheValues
	if len(c.cache.values) > 0 {
		c.cache.pos = sort.Search(len(c.cache.values), func(i int) bool {
			return c.cache.values[i].UnixNano() >= seek
		})
		if c.cache.pos == len(c.cache.values) {
			c.cache.pos--
		} else if c.cache.values[c.cache.pos].UnixNano() != seek {
			c.cache.pos--
		}
	} else {
		c.cache.pos = -1
	}

	c.tsm.keyCursor = tsmKeyCursor
	c.tsm.values = c.readArrayBlock()
	c.tsm.pos = sort.Search(c.tsm.values.Len(), func(i int) bool {
		return c.tsm.values.Timestamps[i] >= seek
	})
	if c.tsm.values.Len() > 0 {
		if c.tsm.pos == c.tsm.values.Len() {
			c.tsm.pos--
		} else if c.tsm.values.Timestamps[c.tsm.pos] != seek {
			c.tsm.pos--
		}
	} else {
		c.tsm.pos = -1
	}
}

func (c *floatArrayDescendingCursor) Err() error { return nil }

func (c *floatArrayDescendingCursor) Close() {
	if c.tsm.keyCursor != nil {
		c.tsm.keyCursor.Close()
		c.tsm.keyCursor = nil
	}
	c.cache.values = nil
	c.tsm.values = nil
}

func (c *floatArrayDescendingCursor) Stats() cursors.CursorStats { return c.stats }

func (c *floatArrayDescendingCursor) Next() *tsdb.FloatArray {
	pos := 0
	cvals := c.cache.values
	tvals := c.tsm.values

	c.res.Timestamps = c.res.Timestamps[:cap(c.res.Timestamps)]
	c.res.Values = c.res.Values[:cap(c.res.Values)]

	for pos < len(c.res.Timestamps) && c.tsm.pos >= 0 && c.cache.pos >= 0 {
		ckey := cvals[c.cache.pos].UnixNano()
		tkey := tvals.Timestamps[c.tsm.pos]
		if ckey == tkey {
			c.res.Timestamps[pos] = ckey
			c.res.Values[pos] = cvals[c.cache.pos].(FloatValue).RawValue()
			c.cache.pos--
			c.tsm.pos--
		} else if ckey > tkey {
			c.res.Timestamps[pos] = ckey
			c.res.Values[pos] = cvals[c.cache.pos].(FloatValue).RawValue()
			c.cache.pos--
		} else {
			c.res.Timestamps[pos] = tkey
			c.res.Values[pos] = tvals.Values[c.tsm.pos]
			c.tsm.pos--
		}

		pos++

		if c.tsm.pos < 0 {
			tvals = c.nextTSM()
		}
	}

	if pos < len(c.res.Timestamps) {
		// cache was exhausted
		if c.tsm.pos >= 0 {
			for pos < len(c.res.Timestamps) && c.tsm.pos >= 0 {
				c.res.Timestamps[pos] = tvals.Timestamps[c.tsm.pos]
				c.res.Values[pos] = tvals.Values[c.tsm.pos]
				pos++
				c.tsm.pos--
				if c.tsm.pos < 0 {
					tvals = c.nextTSM()
				}
			}
		}

		if c.cache.pos >= 0 {
			// TSM was exhausted
			for pos < len(c.res.Timestamps) && c.cache.pos >= 0 {
				c.res.Timestamps[pos] = cvals[c.cache.pos].UnixNano()
				c.res.Values[pos] = cvals[c.cache.pos].(FloatValue).RawValue()
				pos++
				c.cache.pos--
			}
		}
	}

	if pos > 0 && c.res.Timestamps[pos-1] <= c.end {
		pos -= 2
		for pos >= 0 && c.res.Timestamps[pos] <= c.end {
			pos--
		}
		pos++
	}

	c.res.Timestamps = c.res.Timestamps[:pos]
	c.res.Values = c.res.Values[:pos]

	return c.res
}

func (c *floatArrayDescendingCursor) nextTSM() *tsdb.FloatArray {
	c.tsm.keyCursor.Next()
	c.tsm.values = c.readArrayBlock()
	c.tsm.pos = len(c.tsm.values.Timestamps) - 1
	return c.tsm.values
}

func (c *floatArrayDescendingCursor) readArrayBlock() *tsdb.FloatArray {
	values, _ := c.tsm.keyCursor.ReadFloatArrayBlock(c.tsm.buf)

	c.stats.ScannedValues += len(values.Values)

	c.stats.ScannedBytes += len(values.Values) * 8

	return values
}

type integerArrayAscendingCursor struct {
	cache struct {
		values Values
		pos    int
	}

	tsm struct {
		buf       *tsdb.IntegerArray
		values    *tsdb.IntegerArray
		pos       int
		keyCursor *KeyCursor
	}

	end   int64
	res   *tsdb.IntegerArray
	stats cursors.CursorStats
}

func newIntegerArrayAscendingCursor() *integerArrayAscendingCursor {
	c := &integerArrayAscendingCursor{
		res: tsdb.NewIntegerArrayLen(MaxPointsPerBlock),
	}
	c.tsm.buf = tsdb.NewIntegerArrayLen(MaxPointsPerBlock)
	return c
}

func (c *integerArrayAscendingCursor) reset(seek, end int64, cacheValues Values, tsmKeyCursor *KeyCursor) {
	c.end = end
	c.cache.values = cacheValues
	c.cache.pos = sort.Search(len(c.cache.values), func(i int) bool {
		return c.cache.values[i].UnixNano() >= seek
	})

	c.tsm.keyCursor = tsmKeyCursor
	c.tsm.values = c.readArrayBlock()
	c.tsm.pos = sort.Search(c.tsm.values.Len(), func(i int) bool {
		return c.tsm.values.Timestamps[i] >= seek
	})
}

func (c *integerArrayAscendingCursor) Err() error { return nil }

// close closes the cursor and any dependent cursors.
func (c *integerArrayAscendingCursor) Close() {
	if c.tsm.keyCursor != nil {
		c.tsm.keyCursor.Close()
		c.tsm.keyCursor = nil
	}
	c.cache.values = nil
	c.tsm.values = nil
}

func (c *integerArrayAscendingCursor) Stats() cursors.CursorStats { return c.stats }

// Next returns the next key/value for the cursor.
func (c *integerArrayAscendingCursor) Next() *tsdb.IntegerArray {
	pos := 0
	cvals := c.cache.values
	tvals := c.tsm.values

	c.res.Timestamps = c.res.Timestamps[:cap(c.res.Timestamps)]
	c.res.Values = c.res.Values[:cap(c.res.Values)]

	for pos < len(c.res.Timestamps) && c.tsm.pos < len(tvals.Timestamps) && c.cache.pos < len(cvals) {
		ckey := cvals[c.cache.pos].UnixNano()
		tkey := tvals.Timestamps[c.tsm.pos]
		if ckey == tkey {
			c.res.Timestamps[pos] = ckey
			c.res.Values[pos] = cvals[c.cache.pos].(IntegerValue).RawValue()
			c.cache.pos++
			c.tsm.pos++
		} else if ckey < tkey {
			c.res.Timestamps[pos] = ckey
			c.res.Values[pos] = cvals[c.cache.pos].(IntegerValue).RawValue()
			c.cache.pos++
		} else {
			c.res.Timestamps[pos] = tkey
			c.res.Values[pos] = tvals.Values[c.tsm.pos]
			c.tsm.pos++
		}

		pos++

		if c.tsm.pos >= len(tvals.Timestamps) {
			tvals = c.nextTSM()
		}
	}

	if pos < len(c.res.Timestamps) {
		if c.tsm.pos < len(tvals.Timestamps) {
			if pos == 0 {
				// optimization: all points served from TSM data
				copy(c.res.Timestamps, tvals.Timestamps)
				pos += copy(c.res.Values, tvals.Values)
				c.nextTSM()
			} else {
				// copy as much as we can
				n := copy(c.res.Timestamps[pos:], tvals.Timestamps[c.tsm.pos:])
				copy(c.res.Values[pos:], tvals.Values[c.tsm.pos:])
				pos += n
				c.tsm.pos += n
				if c.tsm.pos >= len(tvals.Timestamps) {
					c.nextTSM()
				}
			}
		}

		if c.cache.pos < len(cvals) {
			// TSM was exhausted
			for pos < len(c.res.Timestamps) && c.cache.pos < len(cvals) {
				c.res.Timestamps[pos] = cvals[c.cache.pos].UnixNano()
				c.res.Values[pos] = cvals[c.cache.pos].(IntegerValue).RawValue()
				pos++
				c.cache.pos++
			}
		}
	}

	if pos > 0 && c.res.Timestamps[pos-1] >= c.end {
		pos -= 2
		for pos >= 0 && c.res.Timestamps[pos] >= c.end {
			pos--
		}
		pos++
	}

	c.res.Timestamps = c.res.Timestamps[:pos]
	c.res.Values = c.res.Values[:pos]

	c.stats.ScannedValues += len(c.res.Values)

	c.stats.ScannedBytes += len(c.res.Values) * 8

	return c.res
}

func (c *integerArrayAscendingCursor) nextTSM() *tsdb.IntegerArray {
	c.tsm.keyCursor.Next()
	c.tsm.values = c.readArrayBlock()
	c.tsm.pos = 0
	return c.tsm.values
}

func (c *integerArrayAscendingCursor) readArrayBlock() *tsdb.IntegerArray {
	values, _ := c.tsm.keyCursor.ReadIntegerArrayBlock(c.tsm.buf)
	return values
}

type integerArrayDescendingCursor struct {
	cache struct {
		values Values
		pos    int
	}

	tsm struct {
		buf       *tsdb.IntegerArray
		values    *tsdb.IntegerArray
		pos       int
		keyCursor *KeyCursor
	}

	end   int64
	res   *tsdb.IntegerArray
	stats cursors.CursorStats
}

func newIntegerArrayDescendingCursor() *integerArrayDescendingCursor {
	c := &integerArrayDescendingCursor{
		res: tsdb.NewIntegerArrayLen(MaxPointsPerBlock),
	}
	c.tsm.buf = tsdb.NewIntegerArrayLen(MaxPointsPerBlock)
	return c
}

func (c *integerArrayDescendingCursor) reset(seek, end int64, cacheValues Values, tsmKeyCursor *KeyCursor) {
	c.end = end
	c.cache.values = cacheValues
	if len(c.cache.values) > 0 {
		c.cache.pos = sort.Search(len(c.cache.values), func(i int) bool {
			return c.cache.values[i].UnixNano() >= seek
		})
		if c.cache.pos == len(c.cache.values) {
			c.cache.pos--
		} else if c.cache.values[c.cache.pos].UnixNano() != seek {
			c.cache.pos--
		}
	} else {
		c.cache.pos = -1
	}

	c.tsm.keyCursor = tsmKeyCursor
	c.tsm.values = c.readArrayBlock()
	c.tsm.pos = sort.Search(c.tsm.values.Len(), func(i int) bool {
		return c.tsm.values.Timestamps[i] >= seek
	})
	if c.tsm.values.Len() > 0 {
		if c.tsm.pos == c.tsm.values.Len() {
			c.tsm.pos--
		} else if c.tsm.values.Timestamps[c.tsm.pos] != seek {
			c.tsm.pos--
		}
	} else {
		c.tsm.pos = -1
	}
}

func (c *integerArrayDescendingCursor) Err() error { return nil }

func (c *integerArrayDescendingCursor) Close() {
	if c.tsm.keyCursor != nil {
		c.tsm.keyCursor.Close()
		c.tsm.keyCursor = nil
	}
	c.cache.values = nil
	c.tsm.values = nil
}

func (c *integerArrayDescendingCursor) Stats() cursors.CursorStats { return c.stats }

func (c *integerArrayDescendingCursor) Next() *tsdb.IntegerArray {
	pos := 0
	cvals := c.cache.values
	tvals := c.tsm.values

	c.res.Timestamps = c.res.Timestamps[:cap(c.res.Timestamps)]
	c.res.Values = c.res.Values[:cap(c.res.Values)]

	for pos < len(c.res.Timestamps) && c.tsm.pos >= 0 && c.cache.pos >= 0 {
		ckey := cvals[c.cache.pos].UnixNano()
		tkey := tvals.Timestamps[c.tsm.pos]
		if ckey == tkey {
			c.res.Timestamps[pos] = ckey
			c.res.Values[pos] = cvals[c.cache.pos].(IntegerValue).RawValue()
			c.cache.pos--
			c.tsm.pos--
		} else if ckey > tkey {
			c.res.Timestamps[pos] = ckey
			c.res.Values[pos] = cvals[c.cache.pos].(IntegerValue).RawValue()
			c.cache.pos--
		} else {
			c.res.Timestamps[pos] = tkey
			c.res.Values[pos] = tvals.Values[c.tsm.pos]
			c.tsm.pos--
		}

		pos++

		if c.tsm.pos < 0 {
			tvals = c.nextTSM()
		}
	}

	if pos < len(c.res.Timestamps) {
		// cache was exhausted
		if c.tsm.pos >= 0 {
			for pos < len(c.res.Timestamps) && c.tsm.pos >= 0 {
				c.res.Timestamps[pos] = tvals.Timestamps[c.tsm.pos]
				c.res.Values[pos] = tvals.Values[c.tsm.pos]
				pos++
				c.tsm.pos--
				if c.tsm.pos < 0 {
					tvals = c.nextTSM()
				}
			}
		}

		if c.cache.pos >= 0 {
			// TSM was exhausted
			for pos < len(c.res.Timestamps) && c.cache.pos >= 0 {
				c.res.Timestamps[pos] = cvals[c.cache.pos].UnixNano()
				c.res.Values[pos] = cvals[c.cache.pos].(IntegerValue).RawValue()
				pos++
				c.cache.pos--
			}
		}
	}

	if pos > 0 && c.res.Timestamps[pos-1] <= c.end {
		pos -= 2
		for pos >= 0 && c.res.Timestamps[pos] <= c.end {
			pos--
		}
		pos++
	}

	c.res.Timestamps = c.res.Timestamps[:pos]
	c.res.Values = c.res.Values[:pos]

	return c.res
}

func (c *integerArrayDescendingCursor) nextTSM() *tsdb.IntegerArray {
	c.tsm.keyCursor.Next()
	c.tsm.values = c.readArrayBlock()
	c.tsm.pos = len(c.tsm.values.Timestamps) - 1
	return c.tsm.values
}

func (c *integerArrayDescendingCursor) readArrayBlock() *tsdb.IntegerArray {
	values, _ := c.tsm.keyCursor.ReadIntegerArrayBlock(c.tsm.buf)

	c.stats.ScannedValues += len(values.Values)

	c.stats.ScannedBytes += len(values.Values) * 8

	return values
}

type unsignedArrayAscendingCursor struct {
	cache struct {
		values Values
		pos    int
	}

	tsm struct {
		buf       *tsdb.UnsignedArray
		values    *tsdb.UnsignedArray
		pos       int
		keyCursor *KeyCursor
	}

	end   int64
	res   *tsdb.UnsignedArray
	stats cursors.CursorStats
}

func newUnsignedArrayAscendingCursor() *unsignedArrayAscendingCursor {
	c := &unsignedArrayAscendingCursor{
		res: tsdb.NewUnsignedArrayLen(MaxPointsPerBlock),
	}
	c.tsm.buf = tsdb.NewUnsignedArrayLen(MaxPointsPerBlock)
	return c
}

func (c *unsignedArrayAscendingCursor) reset(seek, end int64, cacheValues Values, tsmKeyCursor *KeyCursor) {
	c.end = end
	c.cache.values = cacheValues
	c.cache.pos = sort.Search(len(c.cache.values), func(i int) bool {
		return c.cache.values[i].UnixNano() >= seek
	})

	c.tsm.keyCursor = tsmKeyCursor
	c.tsm.values = c.readArrayBlock()
	c.tsm.pos = sort.Search(c.tsm.values.Len(), func(i int) bool {
		return c.tsm.values.Timestamps[i] >= seek
	})
}

func (c *unsignedArrayAscendingCursor) Err() error { return nil }

// close closes the cursor and any dependent cursors.
func (c *unsignedArrayAscendingCursor) Close() {
	if c.tsm.keyCursor != nil {
		c.tsm.keyCursor.Close()
		c.tsm.keyCursor = nil
	}
	c.cache.values = nil
	c.tsm.values = nil
}

func (c *unsignedArrayAscendingCursor) Stats() cursors.CursorStats { return c.stats }

// Next returns the next key/value for the cursor.
func (c *unsignedArrayAscendingCursor) Next() *tsdb.UnsignedArray {
	pos := 0
	cvals := c.cache.values
	tvals := c.tsm.values

	c.res.Timestamps = c.res.Timestamps[:cap(c.res.Timestamps)]
	c.res.Values = c.res.Values[:cap(c.res.Values)]

	for pos < len(c.res.Timestamps) && c.tsm.pos < len(tvals.Timestamps) && c.cache.pos < len(cvals) {
		ckey := cvals[c.cache.pos].UnixNano()
		tkey := tvals.Timestamps[c.tsm.pos]
		if ckey == tkey {
			c.res.Timestamps[pos] = ckey
			c.res.Values[pos] = cvals[c.cache.pos].(UnsignedValue).RawValue()
			c.cache.pos++
			c.tsm.pos++
		} else if ckey < tkey {
			c.res.Timestamps[pos] = ckey
			c.res.Values[pos] = cvals[c.cache.pos].(UnsignedValue).RawValue()
			c.cache.pos++
		} else {
			c.res.Timestamps[pos] = tkey
			c.res.Values[pos] = tvals.Values[c.tsm.pos]
			c.tsm.pos++
		}

		pos++

		if c.tsm.pos >= len(tvals.Timestamps) {
			tvals = c.nextTSM()
		}
	}

	if pos < len(c.res.Timestamps) {
		if c.tsm.pos < len(tvals.Timestamps) {
			if pos == 0 {
				// optimization: all points served from TSM data
				copy(c.res.Timestamps, tvals.Timestamps)
				pos += copy(c.res.Values, tvals.Values)
				c.nextTSM()
			} else {
				// copy as much as we can
				n := copy(c.res.Timestamps[pos:], tvals.Timestamps[c.tsm.pos:])
				copy(c.res.Values[pos:], tvals.Values[c.tsm.pos:])
				pos += n
				c.tsm.pos += n
				if c.tsm.pos >= len(tvals.Timestamps) {
					c.nextTSM()
				}
			}
		}

		if c.cache.pos < len(cvals) {
			// TSM was exhausted
			for pos < len(c.res.Timestamps) && c.cache.pos < len(cvals) {
				c.res.Timestamps[pos] = cvals[c.cache.pos].UnixNano()
				c.res.Values[pos] = cvals[c.cache.pos].(UnsignedValue).RawValue()
				pos++
				c.cache.pos++
			}
		}
	}

	if pos > 0 && c.res.Timestamps[pos-1] >= c.end {
		pos -= 2
		for pos >= 0 && c.res.Timestamps[pos] >= c.end {
			pos--
		}
		pos++
	}

	c.res.Timestamps = c.res.Timestamps[:pos]
	c.res.Values = c.res.Values[:pos]

	c.stats.ScannedValues += len(c.res.Values)

	c.stats.ScannedBytes += len(c.res.Values) * 8

	return c.res
}

func (c *unsignedArrayAscendingCursor) nextTSM() *tsdb.UnsignedArray {
	c.tsm.keyCursor.Next()
	c.tsm.values = c.readArrayBlock()
	c.tsm.pos = 0
	return c.tsm.values
}

func (c *unsignedArrayAscendingCursor) readArrayBlock() *tsdb.UnsignedArray {
	values, _ := c.tsm.keyCursor.ReadUnsignedArrayBlock(c.tsm.buf)
	return values
}

type unsignedArrayDescendingCursor struct {
	cache struct {
		values Values
		pos    int
	}

	tsm struct {
		buf       *tsdb.UnsignedArray
		values    *tsdb.UnsignedArray
		pos       int
		keyCursor *KeyCursor
	}

	end   int64
	res   *tsdb.UnsignedArray
	stats cursors.CursorStats
}

func newUnsignedArrayDescendingCursor() *unsignedArrayDescendingCursor {
	c := &unsignedArrayDescendingCursor{
		res: tsdb.NewUnsignedArrayLen(MaxPointsPerBlock),
	}
	c.tsm.buf = tsdb.NewUnsignedArrayLen(MaxPointsPerBlock)
	return c
}

func (c *unsignedArrayDescendingCursor) reset(seek, end int64, cacheValues Values, tsmKeyCursor *KeyCursor) {
	c.end = end
	c.cache.values = cacheValues
	if len(c.cache.values) > 0 {
		c.cache.pos = sort.Search(len(c.cache.values), func(i int) bool {
			return c.cache.values[i].UnixNano() >= seek
		})
		if c.cache.pos == len(c.cache.values) {
			c.cache.pos--
		} else if c.cache.values[c.cache.pos].UnixNano() != seek {
			c.cache.pos--
		}
	} else {
		c.cache.pos = -1
	}

	c.tsm.keyCursor = tsmKeyCursor
	c.tsm.values = c.readArrayBlock()
	c.tsm.pos = sort.Search(c.tsm.values.Len(), func(i int) bool {
		return c.tsm.values.Timestamps[i] >= seek
	})
	if c.tsm.values.Len() > 0 {
		if c.tsm.pos == c.tsm.values.Len() {
			c.tsm.pos--
		} else if c.tsm.values.Timestamps[c.tsm.pos] != seek {
			c.tsm.pos--
		}
	} else {
		c.tsm.pos = -1
	}
}

func (c *unsignedArrayDescendingCursor) Err() error { return nil }

func (c *unsignedArrayDescendingCursor) Close() {
	if c.tsm.keyCursor != nil {
		c.tsm.keyCursor.Close()
		c.tsm.keyCursor = nil
	}
	c.cache.values = nil
	c.tsm.values = nil
}

func (c *unsignedArrayDescendingCursor) Stats() cursors.CursorStats { return c.stats }

func (c *unsignedArrayDescendingCursor) Next() *tsdb.UnsignedArray {
	pos := 0
	cvals := c.cache.values
	tvals := c.tsm.values

	c.res.Timestamps = c.res.Timestamps[:cap(c.res.Timestamps)]
	c.res.Values = c.res.Values[:cap(c.res.Values)]

	for pos < len(c.res.Timestamps) && c.tsm.pos >= 0 && c.cache.pos >= 0 {
		ckey := cvals[c.cache.pos].UnixNano()
		tkey := tvals.Timestamps[c.tsm.pos]
		if ckey == tkey {
			c.res.Timestamps[pos] = ckey
			c.res.Values[pos] = cvals[c.cache.pos].(UnsignedValue).RawValue()
			c.cache.pos--
			c.tsm.pos--
		} else if ckey > tkey {
			c.res.Timestamps[pos] = ckey
			c.res.Values[pos] = cvals[c.cache.pos].(UnsignedValue).RawValue()
			c.cache.pos--
		} else {
			c.res.Timestamps[pos] = tkey
			c.res.Values[pos] = tvals.Values[c.tsm.pos]
			c.tsm.pos--
		}

		pos++

		if c.tsm.pos < 0 {
			tvals = c.nextTSM()
		}
	}

	if pos < len(c.res.Timestamps) {
		// cache was exhausted
		if c.tsm.pos >= 0 {
			for pos < len(c.res.Timestamps) && c.tsm.pos >= 0 {
				c.res.Timestamps[pos] = tvals.Timestamps[c.tsm.pos]
				c.res.Values[pos] = tvals.Values[c.tsm.pos]
				pos++
				c.tsm.pos--
				if c.tsm.pos < 0 {
					tvals = c.nextTSM()
				}
			}
		}

		if c.cache.pos >= 0 {
			// TSM was exhausted
			for pos < len(c.res.Timestamps) && c.cache.pos >= 0 {
				c.res.Timestamps[pos] = cvals[c.cache.pos].UnixNano()
				c.res.Values[pos] = cvals[c.cache.pos].(UnsignedValue).RawValue()
				pos++
				c.cache.pos--
			}
		}
	}

	if pos > 0 && c.res.Timestamps[pos-1] <= c.end {
		pos -= 2
		for pos >= 0 && c.res.Timestamps[pos] <= c.end {
			pos--
		}
		pos++
	}

	c.res.Timestamps = c.res.Timestamps[:pos]
	c.res.Values = c.res.Values[:pos]

	return c.res
}

func (c *unsignedArrayDescendingCursor) nextTSM() *tsdb.UnsignedArray {
	c.tsm.keyCursor.Next()
	c.tsm.values = c.readArrayBlock()
	c.tsm.pos = len(c.tsm.values.Timestamps) - 1
	return c.tsm.values
}

func (c *unsignedArrayDescendingCursor) readArrayBlock() *tsdb.UnsignedArray {
	values, _ := c.tsm.keyCursor.ReadUnsignedArrayBlock(c.tsm.buf)

	c.stats.ScannedValues += len(values.Values)

	c.stats.ScannedBytes += len(values.Values) * 8

	return values
}

type stringArrayAscendingCursor struct {
	cache struct {
		values Values
		pos    int
	}

	tsm struct {
		buf       *tsdb.StringArray
		values    *tsdb.StringArray
		pos       int
		keyCursor *KeyCursor
	}

	end   int64
	res   *tsdb.StringArray
	stats cursors.CursorStats
}

func newStringArrayAscendingCursor() *stringArrayAscendingCursor {
	c := &stringArrayAscendingCursor{
		res: tsdb.NewStringArrayLen(MaxPointsPerBlock),
	}
	c.tsm.buf = tsdb.NewStringArrayLen(MaxPointsPerBlock)
	return c
}

func (c *stringArrayAscendingCursor) reset(seek, end int64, cacheValues Values, tsmKeyCursor *KeyCursor) {
	c.end = end
	c.cache.values = cacheValues
	c.cache.pos = sort.Search(len(c.cache.values), func(i int) bool {
		return c.cache.values[i].UnixNano() >= seek
	})

	c.tsm.keyCursor = tsmKeyCursor
	c.tsm.values = c.readArrayBlock()
	c.tsm.pos = sort.Search(c.tsm.values.Len(), func(i int) bool {
		return c.tsm.values.Timestamps[i] >= seek
	})
}

func (c *stringArrayAscendingCursor) Err() error { return nil }

// close closes the cursor and any dependent cursors.
func (c *stringArrayAscendingCursor) Close() {
	if c.tsm.keyCursor != nil {
		c.tsm.keyCursor.Close()
		c.tsm.keyCursor = nil
	}
	c.cache.values = nil
	c.tsm.values = nil
}

func (c *stringArrayAscendingCursor) Stats() cursors.CursorStats { return c.stats }

// Next returns the next key/value for the cursor.
func (c *stringArrayAscendingCursor) Next() *tsdb.StringArray {
	pos := 0
	cvals := c.cache.values
	tvals := c.tsm.values

	c.res.Timestamps = c.res.Timestamps[:cap(c.res.Timestamps)]
	c.res.Values = c.res.Values[:cap(c.res.Values)]

	for pos < len(c.res.Timestamps) && c.tsm.pos < len(tvals.Timestamps) && c.cache.pos < len(cvals) {
		ckey := cvals[c.cache.pos].UnixNano()
		tkey := tvals.Timestamps[c.tsm.pos]
		if ckey == tkey {
			c.res.Timestamps[pos] = ckey
			c.res.Values[pos] = cvals[c.cache.pos].(StringValue).RawValue()
			c.cache.pos++
			c.tsm.pos++
		} else if ckey < tkey {
			c.res.Timestamps[pos] = ckey
			c.res.Values[pos] = cvals[c.cache.pos].(StringValue).RawValue()
			c.cache.pos++
		} else {
			c.res.Timestamps[pos] = tkey
			c.res.Values[pos] = tvals.Values[c.tsm.pos]
			c.tsm.pos++
		}

		pos++

		if c.tsm.pos >= len(tvals.Timestamps) {
			tvals = c.nextTSM()
		}
	}

	if pos < len(c.res.Timestamps) {
		if c.tsm.pos < len(tvals.Timestamps) {
			if pos == 0 {
				// optimization: all points served from TSM data
				copy(c.res.Timestamps, tvals.Timestamps)
				pos += copy(c.res.Values, tvals.Values)
				c.nextTSM()
			} else {
				// copy as much as we can
				n := copy(c.res.Timestamps[pos:], tvals.Timestamps[c.tsm.pos:])
				copy(c.res.Values[pos:], tvals.Values[c.tsm.pos:])
				pos += n
				c.tsm.pos += n
				if c.tsm.pos >= len(tvals.Timestamps) {
					c.nextTSM()
				}
			}
		}

		if c.cache.pos < len(cvals) {
			// TSM was exhausted
			for pos < len(c.res.Timestamps) && c.cache.pos < len(cvals) {
				c.res.Timestamps[pos] = cvals[c.cache.pos].UnixNano()
				c.res.Values[pos] = cvals[c.cache.pos].(StringValue).RawValue()
				pos++
				c.cache.pos++
			}
		}
	}

	if pos > 0 && c.res.Timestamps[pos-1] >= c.end {
		pos -= 2
		for pos >= 0 && c.res.Timestamps[pos] >= c.end {
			pos--
		}
		pos++
	}

	c.res.Timestamps = c.res.Timestamps[:pos]
	c.res.Values = c.res.Values[:pos]

	c.stats.ScannedValues += len(c.res.Values)

	for _, v := range c.res.Values {
		c.stats.ScannedBytes += len(v)
	}

	return c.res
}

func (c *stringArrayAscendingCursor) nextTSM() *tsdb.StringArray {
	c.tsm.keyCursor.Next()
	c.tsm.values = c.readArrayBlock()
	c.tsm.pos = 0
	return c.tsm.values
}

func (c *stringArrayAscendingCursor) readArrayBlock() *tsdb.StringArray {
	values, _ := c.tsm.keyCursor.ReadStringArrayBlock(c.tsm.buf)
	return values
}

type stringArrayDescendingCursor struct {
	cache struct {
		values Values
		pos    int
	}

	tsm struct {
		buf       *tsdb.StringArray
		values    *tsdb.StringArray
		pos       int
		keyCursor *KeyCursor
	}

	end   int64
	res   *tsdb.StringArray
	stats cursors.CursorStats
}

func newStringArrayDescendingCursor() *stringArrayDescendingCursor {
	c := &stringArrayDescendingCursor{
		res: tsdb.NewStringArrayLen(MaxPointsPerBlock),
	}
	c.tsm.buf = tsdb.NewStringArrayLen(MaxPointsPerBlock)
	return c
}

func (c *stringArrayDescendingCursor) reset(seek, end int64, cacheValues Values, tsmKeyCursor *KeyCursor) {
	c.end = end
	c.cache.values = cacheValues
	if len(c.cache.values) > 0 {
		c.cache.pos = sort.Search(len(c.cache.values), func(i int) bool {
			return c.cache.values[i].UnixNano() >= seek
		})
		if c.cache.pos == len(c.cache.values) {
			c.cache.pos--
		} else if c.cache.values[c.cache.pos].UnixNano() != seek {
			c.cache.pos--
		}
	} else {
		c.cache.pos = -1
	}

	c.tsm.keyCursor = tsmKeyCursor
	c.tsm.values = c.readArrayBlock()
	c.tsm.pos = sort.Search(c.tsm.values.Len(), func(i int) bool {
		return c.tsm.values.Timestamps[i] >= seek
	})
	if c.tsm.values.Len() > 0 {
		if c.tsm.pos == c.tsm.values.Len() {
			c.tsm.pos--
		} else if c.tsm.values.Timestamps[c.tsm.pos] != seek {
			c.tsm.pos--
		}
	} else {
		c.tsm.pos = -1
	}
}

func (c *stringArrayDescendingCursor) Err() error { return nil }

func (c *stringArrayDescendingCursor) Close() {
	if c.tsm.keyCursor != nil {
		c.tsm.keyCursor.Close()
		c.tsm.keyCursor = nil
	}
	c.cache.values = nil
	c.tsm.values = nil
}

func (c *stringArrayDescendingCursor) Stats() cursors.CursorStats { return c.stats }

func (c *stringArrayDescendingCursor) Next() *tsdb.StringArray {
	pos := 0
	cvals := c.cache.values
	tvals := c.tsm.values

	c.res.Timestamps = c.res.Timestamps[:cap(c.res.Timestamps)]
	c.res.Values = c.res.Values[:cap(c.res.Values)]

	for pos < len(c.res.Timestamps) && c.tsm.pos >= 0 && c.cache.pos >= 0 {
		ckey := cvals[c.cache.pos].UnixNano()
		tkey := tvals.Timestamps[c.tsm.pos]
		if ckey == tkey {
			c.res.Timestamps[pos] = ckey
			c.res.Values[pos] = cvals[c.cache.pos].(StringValue).RawValue()
			c.cache.pos--
			c.tsm.pos--
		} else if ckey > tkey {
			c.res.Timestamps[pos] = ckey
			c.res.Values[pos] = cvals[c.cache.pos].(StringValue).RawValue()
			c.cache.pos--
		} else {
			c.res.Timestamps[pos] = tkey
			c.res.Values[pos] = tvals.Values[c.tsm.pos]
			c.tsm.pos--
		}

		pos++

		if c.tsm.pos < 0 {
			tvals = c.nextTSM()
		}
	}

	if pos < len(c.res.Timestamps) {
		// cache was exhausted
		if c.tsm.pos >= 0 {
			for pos < len(c.res.Timestamps) && c.tsm.pos >= 0 {
				c.res.Timestamps[pos] = tvals.Timestamps[c.tsm.pos]
				c.res.Values[pos] = tvals.Values[c.tsm.pos]
				pos++
				c.tsm.pos--
				if c.tsm.pos < 0 {
					tvals = c.nextTSM()
				}
			}
		}

		if c.cache.pos >= 0 {
			// TSM was exhausted
			for pos < len(c.res.Timestamps) && c.cache.pos >= 0 {
				c.res.Timestamps[pos] = cvals[c.cache.pos].UnixNano()
				c.res.Values[pos] = cvals[c.cache.pos].(StringValue).RawValue()
				pos++
				c.cache.pos--
			}
		}
	}

	if pos > 0 && c.res.Timestamps[pos-1] <= c.end {
		pos -= 2
		for pos >= 0 && c.res.Timestamps[pos] <= c.end {
			pos--
		}
		pos++
	}

	c.res.Timestamps = c.res.Timestamps[:pos]
	c.res.Values = c.res.Values[:pos]

	return c.res
}

func (c *stringArrayDescendingCursor) nextTSM() *tsdb.StringArray {
	c.tsm.keyCursor.Next()
	c.tsm.values = c.readArrayBlock()
	c.tsm.pos = len(c.tsm.values.Timestamps) - 1
	return c.tsm.values
}

func (c *stringArrayDescendingCursor) readArrayBlock() *tsdb.StringArray {
	values, _ := c.tsm.keyCursor.ReadStringArrayBlock(c.tsm.buf)

	c.stats.ScannedValues += len(values.Values)

	for _, v := range values.Values {
		c.stats.ScannedBytes += len(v)
	}

	return values
}

type booleanArrayAscendingCursor struct {
	cache struct {
		values Values
		pos    int
	}

	tsm struct {
		buf       *tsdb.BooleanArray
		values    *tsdb.BooleanArray
		pos       int
		keyCursor *KeyCursor
	}

	end   int64
	res   *tsdb.BooleanArray
	stats cursors.CursorStats
}

func newBooleanArrayAscendingCursor() *booleanArrayAscendingCursor {
	c := &booleanArrayAscendingCursor{
		res: tsdb.NewBooleanArrayLen(MaxPointsPerBlock),
	}
	c.tsm.buf = tsdb.NewBooleanArrayLen(MaxPointsPerBlock)
	return c
}

func (c *booleanArrayAscendingCursor) reset(seek, end int64, cacheValues Values, tsmKeyCursor *KeyCursor) {
	c.end = end
	c.cache.values = cacheValues
	c.cache.pos = sort.Search(len(c.cache.values), func(i int) bool {
		return c.cache.values[i].UnixNano() >= seek
	})

	c.tsm.keyCursor = tsmKeyCursor
	c.tsm.values = c.readArrayBlock()
	c.tsm.pos = sort.Search(c.tsm.values.Len(), func(i int) bool {
		return c.tsm.values.Timestamps[i] >= seek
	})
}

func (c *booleanArrayAscendingCursor) Err() error { return nil }

// close closes the cursor and any dependent cursors.
func (c *booleanArrayAscendingCursor) Close() {
	if c.tsm.keyCursor != nil {
		c.tsm.keyCursor.Close()
		c.tsm.keyCursor = nil
	}
	c.cache.values = nil
	c.tsm.values = nil
}

func (c *booleanArrayAscendingCursor) Stats() cursors.CursorStats { return c.stats }

// Next returns the next key/value for the cursor.
func (c *booleanArrayAscendingCursor) Next() *tsdb.BooleanArray {
	pos := 0
	cvals := c.cache.values
	tvals := c.tsm.values

	c.res.Timestamps = c.res.Timestamps[:cap(c.res.Timestamps)]
	c.res.Values = c.res.Values[:cap(c.res.Values)]

	for pos < len(c.res.Timestamps) && c.tsm.pos < len(tvals.Timestamps) && c.cache.pos < len(cvals) {
		ckey := cvals[c.cache.pos].UnixNano()
		tkey := tvals.Timestamps[c.tsm.pos]
		if ckey == tkey {
			c.res.Timestamps[pos] = ckey
			c.res.Values[pos] = cvals[c.cache.pos].(BooleanValue).RawValue()
			c.cache.pos++
			c.tsm.pos++
		} else if ckey < tkey {
			c.res.Timestamps[pos] = ckey
			c.res.Values[pos] = cvals[c.cache.pos].(BooleanValue).RawValue()
			c.cache.pos++
		} else {
			c.res.Timestamps[pos] = tkey
			c.res.Values[pos] = tvals.Values[c.tsm.pos]
			c.tsm.pos++
		}

		pos++

		if c.tsm.pos >= len(tvals.Timestamps) {
			tvals = c.nextTSM()
		}
	}

	if pos < len(c.res.Timestamps) {
		if c.tsm.pos < len(tvals.Timestamps) {
			if pos == 0 {
				// optimization: all points served from TSM data
				copy(c.res.Timestamps, tvals.Timestamps)
				pos += copy(c.res.Values, tvals.Values)
				c.nextTSM()
			} else {
				// copy as much as we can
				n := copy(c.res.Timestamps[pos:], tvals.Timestamps[c.tsm.pos:])
				copy(c.res.Values[pos:], tvals.Values[c.tsm.pos:])
				pos += n
				c.tsm.pos += n
				if c.tsm.pos >= len(tvals.Timestamps) {
					c.nextTSM()
				}
			}
		}

		if c.cache.pos < len(cvals) {
			// TSM was exhausted
			for pos < len(c.res.Timestamps) && c.cache.pos < len(cvals) {
				c.res.Timestamps[pos] = cvals[c.cache.pos].UnixNano()
				c.res.Values[pos] = cvals[c.cache.pos].(BooleanValue).RawValue()
				pos++
				c.cache.pos++
			}
		}
	}

	if pos > 0 && c.res.Timestamps[pos-1] >= c.end {
		pos -= 2
		for pos >= 0 && c.res.Timestamps[pos] >= c.end {
			pos--
		}
		pos++
	}

	c.res.Timestamps = c.res.Timestamps[:pos]
	c.res.Values = c.res.Values[:pos]

	c.stats.ScannedValues += len(c.res.Values)

	c.stats.ScannedBytes += len(c.res.Values) * 1

	return c.res
}

func (c *booleanArrayAscendingCursor) nextTSM() *tsdb.BooleanArray {
	c.tsm.keyCursor.Next()
	c.tsm.values = c.readArrayBlock()
	c.tsm.pos = 0
	return c.tsm.values
}

func (c *booleanArrayAscendingCursor) readArrayBlock() *tsdb.BooleanArray {
	values, _ := c.tsm.keyCursor.ReadBooleanArrayBlock(c.tsm.buf)
	return values
}

type booleanArrayDescendingCursor struct {
	cache struct {
		values Values
		pos    int
	}

	tsm struct {
		buf       *tsdb.BooleanArray
		values    *tsdb.BooleanArray
		pos       int
		keyCursor *KeyCursor
	}

	end   int64
	res   *tsdb.BooleanArray
	stats cursors.CursorStats
}

func newBooleanArrayDescendingCursor() *booleanArrayDescendingCursor {
	c := &booleanArrayDescendingCursor{
		res: tsdb.NewBooleanArrayLen(MaxPointsPerBlock),
	}
	c.tsm.buf = tsdb.NewBooleanArrayLen(MaxPointsPerBlock)
	return c
}

func (c *booleanArrayDescendingCursor) reset(seek, end int64, cacheValues Values, tsmKeyCursor *KeyCursor) {
	c.end = end
	c.cache.values = cacheValues
	if len(c.cache.values) > 0 {
		c.cache.pos = sort.Search(len(c.cache.values), func(i int) bool {
			return c.cache.values[i].UnixNano() >= seek
		})
		if c.cache.pos == len(c.cache.values) {
			c.cache.pos--
		} else if c.cache.values[c.cache.pos].UnixNano() != seek {
			c.cache.pos--
		}
	} else {
		c.cache.pos = -1
	}

	c.tsm.keyCursor = tsmKeyCursor
	c.tsm.values = c.readArrayBlock()
	c.tsm.pos = sort.Search(c.tsm.values.Len(), func(i int) bool {
		return c.tsm.values.Timestamps[i] >= seek
	})
	if c.tsm.values.Len() > 0 {
		if c.tsm.pos == c.tsm.values.Len() {
			c.tsm.pos--
		} else if c.tsm.values.Timestamps[c.tsm.pos] != seek {
			c.tsm.pos--
		}
	} else {
		c.tsm.pos = -1
	}
}

func (c *booleanArrayDescendingCursor) Err() error { return nil }

func (c *booleanArrayDescendingCursor) Close() {
	if c.tsm.keyCursor != nil {
		c.tsm.keyCursor.Close()
		c.tsm.keyCursor = nil
	}
	c.cache.values = nil
	c.tsm.values = nil
}

func (c *booleanArrayDescendingCursor) Stats() cursors.CursorStats { return c.stats }

func (c *booleanArrayDescendingCursor) Next() *tsdb.BooleanArray {
	pos := 0
	cvals := c.cache.values
	tvals := c.tsm.values

	c.res.Timestamps = c.res.Timestamps[:cap(c.res.Timestamps)]
	c.res.Values = c.res.Values[:cap(c.res.Values)]

	for pos < len(c.res.Timestamps) && c.tsm.pos >= 0 && c.cache.pos >= 0 {
		ckey := cvals[c.cache.pos].UnixNano()
		tkey := tvals.Timestamps[c.tsm.pos]
		if ckey == tkey {
			c.res.Timestamps[pos] = ckey
			c.res.Values[pos] = cvals[c.cache.pos].(BooleanValue).RawValue()
			c.cache.pos--
			c.tsm.pos--
		} else if ckey > tkey {
			c.res.Timestamps[pos] = ckey
			c.res.Values[pos] = cvals[c.cache.pos].(BooleanValue).RawValue()
			c.cache.pos--
		} else {
			c.res.Timestamps[pos] = tkey
			c.res.Values[pos] = tvals.Values[c.tsm.pos]
			c.tsm.pos--
		}

		pos++

		if c.tsm.pos < 0 {
			tvals = c.nextTSM()
		}
	}

	if pos < len(c.res.Timestamps) {
		// cache was exhausted
		if c.tsm.pos >= 0 {
			for pos < len(c.res.Timestamps) && c.tsm.pos >= 0 {
				c.res.Timestamps[pos] = tvals.Timestamps[c.tsm.pos]
				c.res.Values[pos] = tvals.Values[c.tsm.pos]
				pos++
				c.tsm.pos--
				if c.tsm.pos < 0 {
					tvals = c.nextTSM()
				}
			}
		}

		if c.cache.pos >= 0 {
			// TSM was exhausted
			for pos < len(c.res.Timestamps) && c.cache.pos >= 0 {
				c.res.Timestamps[pos] = cvals[c.cache.pos].UnixNano()
				c.res.Values[pos] = cvals[c.cache.pos].(BooleanValue).RawValue()
				pos++
				c.cache.pos--
			}
		}
	}

	if pos > 0 && c.res.Timestamps[pos-1] <= c.end {
		pos -= 2
		for pos >= 0 && c.res.Timestamps[pos] <= c.end {
			pos--
		}
		pos++
	}

	c.res.Timestamps = c.res.Timestamps[:pos]
	c.res.Values = c.res.Values[:pos]

	return c.res
}

func (c *booleanArrayDescendingCursor) nextTSM() *tsdb.BooleanArray {
	c.tsm.keyCursor.Next()
	c.tsm.values = c.readArrayBlock()
	c.tsm.pos = len(c.tsm.values.Timestamps) - 1
	return c.tsm.values
}

func (c *booleanArrayDescendingCursor) readArrayBlock() *tsdb.BooleanArray {
	values, _ := c.tsm.keyCursor.ReadBooleanArrayBlock(c.tsm.buf)

	c.stats.ScannedValues += len(values.Values)

	c.stats.ScannedBytes += len(values.Values) * 1

	return values
}
