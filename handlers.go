package main

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type App struct{ St *Store }

var tpl = template.Must(template.New("root").Funcs(template.FuncMap{
	"add": func(a, b int) int { return a + b },
	"sub": func(a, b int) int { return a - b },
	"f2":  func(f float64) string { return fmt.Sprintf("%.2f", f) },
	"f3":  func(f float64) string { return fmt.Sprintf("%.3f", f) },
}).Parse(templatesSrc))

func (a *App) render(w http.ResponseWriter, r *http.Request, page string, data map[string]any) {
	if data == nil {
		data = map[string]any{}
	}
	data["Page"] = page
	data["CSRF"] = getOrSetCSRF(w, r)
	if err := tpl.ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

// ---------- Resources ----------
func (a *App) ResourcesPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if !validateCSRF(r) {
			http.Error(w, "CSRF", 400)
			return
		}
		action, id, name := r.FormValue("action"), r.FormValue("id"), strings.TrimSpace(r.FormValue("name"))
		var err error
		switch action {
		case "create":
			if name == "" {
				err = errors.New("название обязательно")
			} else {
				_, err = a.St.ResourceCreate(name)
			}
		case "update":
			if name == "" {
				err = errors.New("название обязательно")
			} else {
				err = a.St.ResourceUpdate(id, name)
			}
		case "delete":
			err = a.deleteResource(id)
		}
		list, _ := a.St.ResourcesList()
		data := map[string]any{"Items": list}
		if err != nil {
			data["Error"] = err.Error()
			data["FormName"] = name
			data["FormID"] = id
		}
		a.render(w, r, "resources", data)
		return
	}
	list, _ := a.St.ResourcesList()
	a.render(w, r, "resources", map[string]any{"Items": list})
}

// п.7.1: нельзя удалить ресурс, если он в истории или в крафте
func (a *App) deleteResource(id string) error {
	hist, _ := a.St.HistoryList()
	for _, x := range hist {
		if x.ResourceID == id {
			return errInUse
		}
	}
	details, _ := a.St.DetailsList()
	for _, d := range details {
		for _, el := range d.Elements {
			if el.Kind == "resource" && el.ID == id {
				return errInUse
			}
		}
	}
	return a.St.ResourceDelete(id)
}

// ---------- Sources ----------
func (a *App) SourcesPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if !validateCSRF(r) {
			http.Error(w, "CSRF", 400)
			return
		}
		action, id, name := r.FormValue("action"), r.FormValue("id"), strings.TrimSpace(r.FormValue("name"))
		var err error
		switch action {
		case "create":
			if name == "" {
				err = errors.New("название обязательно")
			} else {
				_, err = a.St.SourceCreate(name)
			}
		case "update":
			if name == "" {
				err = errors.New("название обязательно")
			} else {
				err = a.St.SourceUpdate(id, name)
			}
		case "delete":
			err = a.St.SourceDelete(id)
		}
		list, _ := a.St.SourcesList()
		data := map[string]any{"Items": list}
		if err != nil {
			data["Error"] = err.Error()
			data["FormName"] = name
			data["FormID"] = id
		}
		a.render(w, r, "sources", data)
		return
	}
	list, _ := a.St.SourcesList()
	a.render(w, r, "sources", map[string]any{"Items": list})
}

// ---------- History ----------
const historyPageSize = 20

func (a *App) HistoryPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if !validateCSRF(r) {
			http.Error(w, "CSRF", 400)
			return
		}
		action, id := r.FormValue("action"), r.FormValue("id")
		var errText string
		switch action {
		case "create":
			e, msg := a.parseHistoryForm(r)
			if msg == "" {
				_, err := a.St.HistoryCreate(e)
				if err != nil {
					msg = err.Error()
				}
			}
			errText = msg
		case "update":
			e, msg := a.parseHistoryForm(r)
			if msg == "" {
				err := a.St.HistoryUpdate(id, e)
				if err != nil {
					msg = err.Error()
				}
			}
			errText = msg
		case "delete":
			_ = a.St.HistoryDelete(id)
		}
		if errText != "" {
			a.renderHistory(w, r, errText, r.Form)
			return
		}
		http.Redirect(w, r, "/history", http.StatusSeeOther)
		return
	}
	a.renderHistory(w, r, "", nil)
}

func (a *App) parseHistoryForm(r *http.Request) (HistoryEntry, string) {
	res := r.FormValue("resource_id")
	qtyS, costS, occ := r.FormValue("quantity"), r.FormValue("cost"), r.FormValue("occurred_at")
	if res == "" {
		return HistoryEntry{}, "выберите ресурс"
	}
	qty, err := strconv.ParseFloat(qtyS, 64)
	if err != nil || qty <= 0 {
		return HistoryEntry{}, "количество должно быть больше 0"
	}
	cost, err := strconv.ParseFloat(costS, 64)
	if err != nil || cost < 0 {
		return HistoryEntry{}, "стоимость должна быть >= 0"
	}
	if occ == "" {
		occ = now()
	}
	return HistoryEntry{ResourceID: res, Quantity: qty, Cost: cost, OccurredAt: occ}, ""
}

func (a *App) renderHistory(w http.ResponseWriter, r *http.Request, errText string, _ any) {
	resources, _ := a.St.ResourcesList()
	all, _ := a.St.HistoryList()

	fResource := r.URL.Query().Get("resource")
	fFrom := r.URL.Query().Get("date_from")
	fTo := r.URL.Query().Get("date_to")

	filtered := make([]HistoryEntry, 0, len(all))
	for _, x := range all {
		if fResource != "" && x.ResourceID != fResource {
			continue
		}
		if fFrom != "" && x.OccurredAt < fFrom+" 00:00:00" {
			continue
		}
		if fTo != "" && x.OccurredAt > fTo+" 23:59:59" {
			continue
		}
		filtered = append(filtered, x)
	}
	sort.Slice(filtered, func(i, j int) bool { return filtered[i].OccurredAt > filtered[j].OccurredAt })

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	pages := (len(filtered) + historyPageSize - 1) / historyPageSize
	if pages == 0 {
		pages = 1
	}
	if page > pages {
		page = pages
	}
	start := (page - 1) * historyPageSize
	end := start + historyPageSize
	if end > len(filtered) {
		end = len(filtered)
	}
	pageItems := filtered[start:end]

	resName := map[string]string{}
	for _, x := range resources {
		resName[x.ID] = x.Name
	}

	data := map[string]any{
		"Items":     pageItems,
		"Resources": resources,
		"ResName":   resName,
		"PageNum":   page,
		"Pages":     pages,
		"Total":     len(filtered),
		"FResource": fResource,
		"FDateFrom": fFrom,
		"FDateTo":   fTo,
		"Error":     errText,
	}
	a.render(w, r, "history", data)
}

// ---------- Summary ----------
func (a *App) SummaryPage(w http.ResponseWriter, r *http.Request) {
	resources, _ := a.St.ResourcesList()
	sel := r.URL.Query().Get("resource")
	if sel == "" && len(resources) > 0 {
		sel = resources[0].ID
	}

	var selName string
	for _, x := range resources {
		if x.ID == sel {
			selName = x.Name
			break
		}
	}

	all, _ := a.St.HistoryList()
	var totalQty, costSum float64
	var costCount int
	for _, x := range all {
		if x.ResourceID != sel {
			continue
		}
		totalQty += x.Quantity
		if x.Cost > 0 {
			costSum += x.Cost
			costCount++
		}
	}
	var avg float64
	if costCount > 0 {
		avg = costSum / float64(costCount)
	}
	total := avg * totalQty

	data := map[string]any{
		"Resources":    resources,
		"SelectedID":   sel,
		"SelectedName": selName,
		"TotalQty":     totalQty,
		"CostSum":      costSum,
		"CostCount":    costCount,
		"Avg":          avg,
		"Total":        total,
	}
	a.render(w, r, "summary", data)
}

// ---------- Details ----------
func (a *App) DetailsPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if !validateCSRF(r) {
			http.Error(w, "CSRF", 400)
			return
		}
		action := r.FormValue("action")
		switch action {
		case "create":
			d, msg := a.parseDetailForm(r)
			if msg == "" {
				_, err := a.St.DetailCreate(d)
				if err != nil {
					msg = err.Error()
				}
			}
			if msg != "" {
				a.renderDetails(w, r, msg)
				return
			}
		case "update":
			d, msg := a.parseDetailForm(r)
			if msg == "" {
				err := a.St.DetailUpdate(r.FormValue("id"), d)
				if err != nil {
					msg = err.Error()
				}
			}
			if msg != "" {
				a.renderDetails(w, r, msg)
				return
			}
		case "delete":
			if err := a.deleteDetail(r.FormValue("id")); err != nil {
				a.renderDetails(w, r, err.Error())
				return
			}
		}
		http.Redirect(w, r, "/details", http.StatusSeeOther)
		return
	}
	a.renderDetails(w, r, "")
}

func (a *App) deleteDetail(id string) error {
	details, _ := a.St.DetailsList()
	for _, d := range details {
		if d.ID == id {
			continue
		}
		for _, el := range d.Elements {
			if el.Kind == "detail" && el.ID == id {
				return errInUse
			}
		}
	}
	return a.St.DetailDelete(id)
}

// parseDetailForm — вся валидация п.6.
func (a *App) parseDetailForm(r *http.Request) (Detail, string) {
	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		return Detail{}, "название детали обязательно"
	}
	rarity := r.FormValue("rarity")
	if !rarityValid(rarity) {
		return Detail{}, "некорректная редкость"
	}

	isCraft := r.FormValue("is_craftable") == "on"

	count, _ := strconv.Atoi(r.FormValue("el_count"))
	if count < 0 {
		count = 0
	}

	var elements []CraftElement
	for i := 0; i < count; i++ {
		idx := strconv.Itoa(i)
		k := r.FormValue("el_kind_" + idx)
		id := r.FormValue("el_id_" + idx)
		qtyS := r.FormValue("el_qty_" + idx)

		if k == "" && id == "" && qtyS == "" {
			continue
		}
		if k != "detail" && k != "resource" {
			return Detail{}, "некорректный тип элемента"
		}
		if id == "" {
			return Detail{}, "элемент без выбранного объекта"
		}
		qty, err := strconv.ParseFloat(qtyS, 64)
		if err != nil || qty <= 0 {
			return Detail{}, "количество материала должно быть > 0"
		}
		elements = append(elements, CraftElement{Kind: k, ID: id, Qty: qty})
	}

	if isCraft && len(elements) == 0 {
		return Detail{}, "крафт включён, но нет элементов"
	}
	if !isCraft {
		elements = nil
	}

	saleS := strings.TrimSpace(r.FormValue("sale_price"))
	buyS := strings.TrimSpace(r.FormValue("buy_price"))
	var sale, buy *float64
	if saleS != "" || buyS != "" {
		if saleS == "" || buyS == "" {
			return Detail{}, "укажите обе цены или ни одной"
		}
		s, e1 := strconv.ParseFloat(saleS, 64)
		b, e2 := strconv.ParseFloat(buyS, 64)
		if e1 != nil || e2 != nil || s < 0 || b < 0 {
			return Detail{}, "цены должны быть числами >= 0"
		}
		sale, buy = &s, &b
	}
	if !isCraft && sale == nil {
		return Detail{}, "деталь не крафтится и без стоимости — недопустимо"
	}

	return Detail{Name: name, Rarity: rarity, IsCraftable: isCraft, Elements: elements, SalePrice: sale, BuyPrice: buy}, ""
}

// renderDetails восстанавливает форму из r.Form при ошибке (п.7.4)
func (a *App) renderDetails(w http.ResponseWriter, r *http.Request, errText string) {
	details, _ := a.St.DetailsList()
	resources, _ := a.St.ResourcesList()

	draft := struct {
		Name      string
		Rarity    string
		IsCraft   bool
		SalePrice string
		BuyPrice  string
		Rows      []map[string]string
	}{}
	if r.Method == http.MethodPost {
		draft.Name = r.FormValue("name")
		draft.Rarity = r.FormValue("rarity")
		draft.IsCraft = r.FormValue("is_craftable") == "on"
		draft.SalePrice = r.FormValue("sale_price")
		draft.BuyPrice = r.FormValue("buy_price")
		count, _ := strconv.Atoi(r.FormValue("el_count"))
		for i := 0; i < count; i++ {
			idx := strconv.Itoa(i)
			draft.Rows = append(draft.Rows, map[string]string{
				"Kind": r.FormValue("el_kind_" + idx),
				"ID":   r.FormValue("el_id_" + idx),
				"Qty":  r.FormValue("el_qty_" + idx),
			})
		}
	}

	resName := map[string]string{}
	for _, x := range resources {
		resName[x.ID] = x.Name
	}
	detName := map[string]string{}
	for _, x := range details {
		detName[x.ID] = x.Name
	}

	data := map[string]any{
		"Items":     details,
		"Resources": resources,
		"Rarities":  Rarities,
		"ResName":   resName,
		"DetName":   detName,
		"Draft":     draft,
		"Error":     errText,
	}
	a.render(w, r, "details", data)
}
