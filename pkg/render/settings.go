package render

type Settings struct {
	FilterKeyPrefix string
	FilterGroup     string
	IncludeDisabled bool

	ShowBlankLines   bool
	ShowColors       bool
	ShowComments     bool
	ShowGroupBanners bool
	ShowPretty       bool

	Interpolate bool
}

func (rs *Settings) WithComments() bool {
	return rs.ShowPretty || rs.ShowComments
}

func (rs *Settings) WithGroupBanners() bool {
	return rs.ShowPretty || rs.ShowGroupBanners
}

func (rs *Settings) WithBlankLines() bool {
	return (rs.ShowBlankLines && rs.ShowComments)
}

func (rs *Settings) WithColors() bool {
	return rs.ShowColors
}
