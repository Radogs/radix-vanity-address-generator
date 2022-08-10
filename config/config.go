package config

const (
	mainnetMode = "mainnet"
)

type Config struct {
	searchWords map[string]bool
	found       int
	total       int
	mode        string
}

func NewConfig() *Config {
	return &Config{
		searchWords: map[string]bool{},
		found:       0,
		total:       0,
		mode:        mainnetMode,
	}
}

func (c *Config) Match(word string) {
	for searchWord := range c.searchWords {
		if searchWord == word && !c.searchWords[searchWord] {
			c.searchWords[searchWord] = true
			c.found = c.found + 1
		}
	}
}

func (c *Config) SetWords(words []string) {
	for i, word := range words {
		c.searchWords[word] = false
		c.total = i
	}
}

func (c *Config) GetWords() (words []string) {
	for word, isFound := range c.searchWords {
		if !isFound {
			words = append(words, word)
		}
	}
	return
}

func (c *Config) Found() int {
	return c.found
}

func (c *Config) Total() int {
	return c.total
}
