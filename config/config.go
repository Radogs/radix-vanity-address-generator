package config

const (
	mainnetMode = "mainnet"
)

type Config struct {
	searchWords    map[string]bool
	found          int
	total          int
	mode           string
	isMnemonicMode bool
}

func NewConfig() *Config {
	return &Config{
		searchWords:    map[string]bool{},
		found:          0,
		total:          0,
		mode:           mainnetMode,
		isMnemonicMode: false,
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
		c.total = i + 1
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

func (c *Config) SetMnemonicMode(isMnemonicMode bool) {
	c.isMnemonicMode = isMnemonicMode
}

func (c *Config) GetMnemonicMode() bool {
	return c.isMnemonicMode
}

func (c *Config) Found() int {
	return c.found
}

func (c *Config) Total() int {
	return c.total
}
