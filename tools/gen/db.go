package main

type Source struct {
	Name string `json:"name" yaml:"name"`
	URL  string `json:"url" yaml:"url"`
}

type Image struct {
	Width  int    `json:"width,omitempty" yaml:"width"`
	Source string `json:"src,omitempty" yaml:"src"`
}

type Record struct {
	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`

	Status      string `json:"status" yaml:"status"`
	LastVisited string `json:"last_visited" yaml:"last_visited"`

	Sources []Source `json:"sources,omitempty" yaml:"sources"`

	Image *Image `json:"image,omitempty" yaml:"image"`
}

type DB struct {
	Data []Record `json:"data" yaml:"data"`
}
