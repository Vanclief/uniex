package polygon

type optionApplyFunc func(*Polygon) error

type Option interface {
	applyOption(*Polygon) error
}

func (f optionApplyFunc) applyOption(p *Polygon) error {
	return f(p)
}

func WithAdjusted() Option {
	return optionApplyFunc(func(polygon *Polygon) error {
		polygon.adjusted = true
		return nil
	})
}

func WithTimespan(ts string) Option {
	return optionApplyFunc(func(polygon *Polygon) error {
		polygon.timespan = ts
		return nil
	})
}

func WithMultiplier(m string) Option {
	return optionApplyFunc(func(polygon *Polygon) error {
		polygon.multiplier = m
		return nil
	})
}

func WithLimit(l int) Option {
	return optionApplyFunc(func(polygon *Polygon) error {
		polygon.limit = l
		return nil
	})
}
