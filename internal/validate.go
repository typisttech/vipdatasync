package internal

type Problemer interface {
	Text() string
	Len() int
}

type CheckerFunc func(DomainMap, []replacement) Problemer

type Validator struct {
	checkers []CheckerFunc
}

func NewValidator() Validator {
	return Validator{
		checkers: []CheckerFunc{
			func(dm DomainMap, _ []replacement) Problemer {
				return checkDuplicatedDestinations(dm)
			},
			func(_ DomainMap, rs []replacement) Problemer {
				return checkDuplicatedTos(rs)
			},
			func(_ DomainMap, rs []replacement) Problemer {
				return checkUnreplacedURLs(rs)
			},
			func(dm DomainMap, rs []replacement) Problemer {
				return checkUnusedDomainMapItems(dm, rs)
			},
		},
	}
}

func (v Validator) Validate(dm DomainMap, urls URLs) []Problemer {
	rs := dm.replace(urls)

	results := make([]Problemer, 0, len(v.checkers))

	for _, check := range v.checkers {
		r := check(dm, rs)
		results = append(results, r)
	}

	return results
}
