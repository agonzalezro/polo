package site

// Len is needed to implement the sorting interface.
func (s Site) Len() int {
	return len(s.Articles)
}

// Less is a comparator to help us to sort the site by date DESC.
func (s Site) Less(i, j int) bool {
	return s.Articles[i].Date.After(s.Articles[j].Date)
}

// Swap is needed to implement the sorting interface.
func (s Site) Swap(i, j int) {
	s.Articles[i], s.Articles[j] = s.Articles[j], s.Articles[i]
}
