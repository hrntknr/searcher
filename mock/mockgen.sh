#!/bin/bash
cd $(dirname $0)
mockgen -source ../charFilter.go -destination charFilter.go -package mock
mockgen -source ../db.go -destination db.go -package mock
mockgen -source ../sentenceSplitter.go -destination sentenceSplitter.go -package mock
mockgen -source ../service.go -destination service.go -package mock
mockgen -source ../tokenizer.go -destination tokenizer.go -package mock
mockgen -source ../wordFilter.go -destination wordFilter.go -package mock
