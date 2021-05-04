package main

import "strings"

type CharFilter interface {
	Filter([]string) []string
}

func newMappingCharFilter(mapper map[string]string) (*mappingCharFilter, error) {
	return &mappingCharFilter{
		mapper: mapper,
	}, nil
}

type mappingCharFilter struct {
	mapper map[string]string
}

func (f *mappingCharFilter) Filter(str []string) []string {
	result := make([]string, len(str))
	for i, s := range str {
		result[i] = s
		for k, v := range f.mapper {
			result[i] = strings.Replace(result[i], k, v, -1)
		}
	}
	return result
}
