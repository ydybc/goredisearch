package redisearch

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go-redisearch/goRedis"
	"strconv"
)

// Autocompleter implements a redisearch auto-completer API
type Autocompleter struct {
	IndexName string
	R         *redis.Client
	Ctx       context.Context
}

func (r Autocompleter) SetAutoCompleterIndex(name string) {
	r.IndexName = name
}

//// NewAutocompleter creates a new Autocompleter with the given pool and key name
//func NewAutocompleterFromPool(pool *redis.Pool, name string) *Autocompleter {
//	return &Autocompleter{IndexName: name, pool: pool}
//}
//
//// NewAutocompleter creates a new Autocompleter with the given host and key name
//func NewAutocompleter(addr, name string) *Autocompleter {
//	return &Autocompleter{
//		pool: redis.NewPool(func() (redis.Conn, error) {
//			return redis.Dial("tcp", addr)
//		}, maxConns),
//		IndexName: name,
//	}
//}

// Delete deletes the Autocompleter key for this AC
func (a *Autocompleter) Delete() error {
	_, err := a.R.Do(a.Ctx, "DEL", a.IndexName).Result()
	return err
}

// AddTerms pushes new term suggestions to the index
func (a *Autocompleter) AddTerms(terms ...Suggestion) error {
	var mErr MultiError
	conn := a.R.TxPipeline()
	for i, term := range terms {
		args := goRedis.Args{"FT.SUGADD", a.IndexName, term.Term, term.Score}
		if term.Incr {
			args = append(args, "INCR")
		}
		if term.Payload != "" {
			args = append(args, "PAYLOAD", term.Payload)
		}

		if err := conn.Do(a.Ctx, args...).Err(); err != nil {
			if mErr == nil {
				mErr = NewMultiError(len(terms))
			}
			mErr[i] = err

			return mErr
		}
	}
	cmders, err := conn.Exec(a.Ctx)
	if err != nil {
		return err
	}
	for k, cmder := range cmders {
		if cmder.Err() != nil {
			if mErr == nil {
				mErr = NewMultiError(len(mErr))
			}
			mErr[k] = cmder.Err()
		}
	}
	if mErr == nil {
		return nil
	}
	return mErr
}

// AddTerms pushes new term suggestions to the index
func (a *Autocompleter) DeleteTerms(terms ...Suggestion) error {
	var mErr MultiError
	conn := a.R.TxPipeline()
	for i, term := range terms {

		args := goRedis.Args{"FT.SUGDEL", a.IndexName, term.Term}
		if err := conn.Do(a.Ctx, args...).Err(); err != nil {
			if mErr == nil {
				mErr = NewMultiError(len(terms))
			}
			mErr[i] = err

			return mErr
		}
	}
	cmders, err := conn.Exec(a.Ctx)
	if err != nil {
		return err
	}
	for k, cmder := range cmders {
		if cmder.Err() != nil {
			if mErr == nil {
				mErr = NewMultiError(len(terms))
			}
			mErr[k] = err
		}
	}
	if mErr == nil {
		return nil
	}
	return mErr
}

// AddTerms pushes new term suggestions to the index
func (a *Autocompleter) Length() (len int64, err error) {
	len, err = a.R.Do(a.Ctx, "FT.SUGLEN", a.IndexName).Int64()
	return
}

// Suggest gets completion suggestions from the Autocompleter dictionary to the given prefix.
// If fuzzy is set, we also complete for prefixes that are in 1 Levenshten distance from the
// given prefix
//
// Deprecated: Please use SuggestOpts() instead
func (a *Autocompleter) Suggest(prefix string, num int, fuzzy bool) (ret []Suggestion, err error) {
	seropts := DefaultSuggestOptions
	seropts.Num = num
	seropts.Fuzzy = fuzzy
	args, inc := a.Serialize("FT.SUGGET", prefix, seropts)
	vals, err := goRedis.Strings(a.R.Do(a.Ctx, args...).Result())

	if err != nil {
		return nil, err
	}

	ret = ProcessSugGetVals(vals, inc, true, false)

	return
}

// SuggestOpts gets completion suggestions from the Autocompleter dictionary to the given prefix.
// SuggestOptions are passed allowing you specify if the returned values contain a payload, and scores.
// If SuggestOptions.Fuzzy is set, we also complete for prefixes that are in 1 Levenshtein distance from the
// given prefix
func (a *Autocompleter) SuggestOpts(prefix string, opts SuggestOptions) (ret []Suggestion, err error) {
	args, inc := a.Serialize("FT.SUGGET", prefix, opts)
	vals, err := goRedis.Strings(a.R.Do(a.Ctx, args...).Result())
	if err != nil {
		return nil, err
	}

	ret = ProcessSugGetVals(vals, inc, opts.WithScores, opts.WithPayloads)

	return
}

func (a *Autocompleter) Serialize(command, prefix string, opts SuggestOptions) (goRedis.Args, int) {
	inc := 1
	args := goRedis.Args{command, a.IndexName, prefix, "MAX", opts.Num}
	if opts.Fuzzy {
		args = append(args, "FUZZY")
	}
	if opts.WithScores {
		args = append(args, "WITHSCORES")
		inc++
	}
	if opts.WithPayloads {
		args = append(args, "WITHPAYLOADS")
		inc++
	}
	return args, inc
}

func ProcessSugGetVals(vals []string, inc int, WithScores, WithPayloads bool) (ret []Suggestion) {
	ret = make([]Suggestion, 0, len(vals)/inc)
	for i := 0; i < len(vals); i += inc {

		suggestion := Suggestion{Term: vals[i]}
		if WithScores {
			score, err := strconv.ParseFloat(vals[i+1], 64)
			if err != nil {
				continue
			}
			suggestion.Score = score
		}
		if WithPayloads {
			suggestion.Payload = vals[i+(inc-1)]
		}
		ret = append(ret, suggestion)

	}
	return
}
