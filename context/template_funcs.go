package context

import "time"

func (c Context) ShowHeader() bool {
	return (c.Config.ShowTags && len(c.Tags) > 0) || (c.Config.ShowCategories && len(c.Categories) > 0) || (c.Config.ShowArchive && len(c.Articles) > 0)
}

// HumanizeDatetime returns a date or datetime depending of the datetime
// received.
// For example, if the datetime received doesn't have any hour/minutes, the
// hours/minutes part doesn't need to be shown.
func (c Context) HumanizeDatetime(datetime time.Time) string {
	if datetime.Hour()+datetime.Minute() == 0 {
		return datetime.Format("2006-01-02")
	}
	return datetime.Format("2006-01-02 15:04")
}

// ArrayOfPages is a dirty hack because we can not (or I don't know how) do a
// range from X to Y on the template
func (c Context) ArrayOfPages() (pages []int) {
	for i := 1; i < c.NumberOfPages()+1; i++ {
		pages = append(pages, i)
	}
	return pages
}
