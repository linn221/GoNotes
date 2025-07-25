package views

type Page struct {
	PageTitle string
	Nav       string
	Timezone  string
	ResList   any
	Res       any
	Data      H
}

type ResourceData struct {
	Res  any
	Data H
}

type ResourceListData struct {
	ResList any
	Data    H
}

type H map[string]any
