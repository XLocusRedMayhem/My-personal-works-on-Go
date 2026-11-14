package main

const templatesSrc = `
{{define "layout"}}
<!doctype html>
<html lang="ru">
<head>
<meta charset="utf-8">
<title>Учёт ресурсов</title>
<style>
* { box-sizing: border-box; }
body { font-family: system-ui, sans-serif; margin: 0; background: #f5f6f8; color: #222; }
.topbar { display:flex; align-items:center; gap:24px; padding:12px 24px; background:#1f2937; color:#fff; }
.topbar .brand { font-weight:600; }
.topbar nav a { color:#cbd5e1; text-decoration:none; margin-right:12px; }
.topbar nav a.active { color:#fff; font-weight:600; }
main { max-width: 1100px; margin: 24px auto; padding: 0 16px; }
h1 { margin-top: 0; }
.card { background:#fff; padding:16px; border-radius:10px; box-shadow:0 1px 4px rgba(0,0,0,.05); margin-bottom:16px; display:flex; flex-wrap:wrap; gap:12px; align-items:end; }
.card label { display:flex; flex-direction:column; font-size:13px; gap:4px; }
.card input, .card select { padding:6px 8px; border:1px solid #cbd5e1; border-radius:6px; }
button { padding:7px 12px; border:1px solid #2563eb; background:#2563eb; color:#fff; border-radius:6px; cursor:pointer; }
button.danger { background:#dc2626; border-color:#dc2626; }
button.secondary { background:#fff; color:#1f2937; border-color:#cbd5e1; }
.alert { background:#fee2e2; color:#991b1b; padding:10px 14px; border-radius:8px; margin-bottom:16px; }
.grid { width:100%; border-collapse:collapse; background:#fff; border-radius:10px; overflow:hidden; }
.grid th, .grid td { padding:8px 10px; border-bottom:1px solid #e5e7eb; text-align:left; }
.grid th { background:#f1f5f9; font-size:13px; }
.grid input[type=text], .grid input[type=number], .grid select { padding:5px 6px; border:1px solid #cbd5e1; border-radius:5px; width:100%; }
.actions { white-space:nowrap; }
.muted { color:#6b7280; font-size:12px; }
.kv { width:100%; border-collapse:collapse; }
.kv td { padding:8px 10px; border-bottom:1px solid #e5e7eb; }
.kv tr.total td { font-weight:700; font-size:18px; background:#f8fafc; }
.pager { display:flex; gap:16px; align-items:center; margin-top:12px; }
.pager a { color:#2563eb; text-decoration:none; }
.craft { border-top:1px dashed #cbd5e1; padding-top:12px; width:100%; }
.checkbox { flex-direction:row !important; align-items:center; gap:6px; }
.prices { display:flex; gap:12px; flex-wrap:wrap; width:100%; }
.element { display:flex; gap:8px; margin-bottom:8px; align-items:center; flex-wrap:wrap; }
.element select, .element input { padding:6px; border:1px solid #cbd5e1; border-radius:6px; }
</style>
</head>
<body>
<header class="topbar">
  <div class="brand">Учёт ресурсов</div>
  <nav>
    <a href="/resources" {{if eq .Page "resources"}}class="active"{{end}}>Ресурсы</a>
    <a href="/sources" {{if eq .Page "sources"}}class="active"{{end}}>Источники</a>
    <a href="/history" {{if eq .Page "history"}}class="active"{{end}}>История</a>
    <a href="/summary" {{if eq .Page "summary"}}class="active"{{end}}>Сумма</a>
    <a href="/details" {{if eq .Page "details"}}class="active"{{end}}>Детали</a>
  </nav>
</header>
<main>
  {{if .Error}}<div class="alert">{{.Error}}</div>{{end}}
  {{template "content" .}}
</main>
</body>
</html>
{{end}}

{{define "named_table"}}
<table class="grid">
  <thead><tr><th>Название</th><th>Создан</th><th>Изменён</th><th></th></tr></thead>
  <tbody>
  {{range .Items}}
    <tr>
      <form method="post" style="display:contents">
        <input type="hidden" name="csrf_token" value="{{$.CSRF}}">
        <input type="hidden" name="id" value="{{.ID}}">
        <td><input type="text" name="name" value="{{.Name}}" required></td>
        <td class="muted">{{.CreatedAt}}</td>
        <td class="muted">{{.UpdatedAt}}</td>
        <td class="actions">
          <button name="action" value="update">Сохранить</button>
          <button name="action" value="delete" class="danger" onclick="return confirm('Удалить?')">Удалить</button>
        </td>
      </form>
    </tr>
  {{else}}
    <tr><td colspan="4" class="muted">Пусто</td></tr>
  {{end}}
  </tbody>
</table>
{{end}}

{{define "content"}}
{{if eq .Page "resources"}}{{template "resources" .}}{{end}}
{{if eq .Page "sources"}}{{template "sources" .}}{{end}}
{{if eq .Page "history"}}{{template "history" .}}{{end}}
{{if eq .Page "summary"}}{{template "summary" .}}{{end}}
{{if eq .Page "details"}}{{template "details" .}}{{end}}
{{end}}

{{define "resources"}}
<h1>Ресурсы</h1>
<form method="post" class="card">
  <input type="hidden" name="csrf_token" value="{{.CSRF}}">
  <input type="hidden" name="action" value="create">
  <label>Название <input type="text" name="name" required></label>
  <button type="submit">Добавить</button>
</form>
{{template "named_table" .}}
{{end}}

{{define "sources"}}
<h1>Источники</h1>
<form method="post" class="card">
  <input type="hidden" name="csrf_token" value="{{.CSRF}}">
  <input type="hidden" name="action" value="create">
  <label>Название <input type="text" name="name" required></label>
  <button type="submit">Добавить</button>
</form>
{{template "named_table" .}}
{{end}}

{{define "history"}}
<h1>История добычи</h1>
<form method="post" class="card">
  <input type="hidden" name="csrf_token" value="{{.CSRF}}">
  <input type="hidden" name="action" value="create">
  <label>Ресурс
    <select name="resource_id" required>
      <option value="">—</option>
      {{range .Resources}}<option value="{{.ID}}">{{.Name}}</option>{{end}}
    </select>
  </label>
  <label>Количество <input type="number" step="any" min="0.0001" name="quantity" required></label>
  <label>Стоимость <input type="number" step="any" min="0" name="cost" value="0" required></label>
  <label>Дата/время <input type="text" name="occurred_at" placeholder="YYYY-MM-DD HH:MM:SS (пусто = сейчас)"></label>
  <button type="submit">Добавить</button>
</form>

<form method="get" class="card">
  <label>Ресурс
    <select name="resource">
      <option value="">Все</option>
      {{range .Resources}}<option value="{{.ID}}" {{if eq .ID $.FResource}}selected{{end}}>{{.Name}}</option>{{end}}
    </select>
  </label>
  <label>С <input type="date" name="date_from" value="{{.FDateFrom}}"></label>
  <label>По <input type="date" name="date_to" value="{{.FDateTo}}"></label>
  <button type="submit">Фильтр</button>
  <a href="/history" class="secondary" style="padding:7px 12px;border:1px solid #cbd5e1;border-radius:6px;text-decoration:none;color:#1f2937">Сбросить</a>
</form>

<table class="grid">
  <thead><tr><th>Дата</th><th>Ресурс</th><th>Кол-во</th><th>Стоимость</th><th></th></tr></thead>
  <tbody>
  {{range .Items}}
    <tr>
      <form method="post" style="display:contents">
        <input type="hidden" name="csrf_token" value="{{$.CSRF}}">
        <input type="hidden" name="id" value="{{.ID}}">
        <td><input type="text" name="occurred_at" value="{{.OccurredAt}}"></td>
        <td>
          <select name="resource_id">
            {{$rid := .ResourceID}}{{range $.Resources}}<option value="{{.ID}}" {{if eq .ID $rid}}selected{{end}}>{{.Name}}</option>{{end}}
          </select>
        </td>
        <td><input type="number" step="any" name="quantity" value="{{.Quantity}}"></td>
        <td><input type="number" step="any" name="cost" value="{{.Cost}}"></td>
        <td class="actions">
          <button name="action" value="update">Сохранить</button>
          <button name="action" value="delete" class="danger" onclick="return confirm('Удалить?')">Удалить</button>
        </td>
      </form>
    </tr>
  {{else}}
    <tr><td colspan="5" class="muted">Нет записей</td></tr>
  {{end}}
  </tbody>
</table>

<div class="pager">
  {{if gt .PageNum 1}}<a href="?page={{sub .PageNum 1}}&resource={{.FResource}}&date_from={{.FDateFrom}}&date_to={{.FDateTo}}">← Назад</a>{{end}}
  <span>Страница {{.PageNum}} из {{.Pages}} (всего {{.Total}})</span>
  {{if lt .PageNum .Pages}}<a href="?page={{add .PageNum 1}}&resource={{.FResource}}&date_from={{.FDateFrom}}&date_to={{.FDateTo}}">Вперёд →</a>{{end}}
</div>
{{end}}

{{define "summary"}}
<h1>Общая сумма</h1>
<form method="get" class="card">
  <label>Ресурс
    <select name="resource" onchange="this.form.submit()">
      {{$sel := .SelectedID}}{{range .Resources}}<option value="{{.ID}}" {{if eq .ID $sel}}selected{{end}}>{{.Name}}</option>{{end}}
    </select>
  </label>
  <button type="submit">Показать</button>
</form>
<div class="card" style="display:block">
  <h2 style="margin-top:0">{{.SelectedName}}</h2>
  <table class="kv">
    <tr><td>Суммарное количество (все записи)</td><td>{{f3 .TotalQty}}</td></tr>
    <tr><td>Записей со стоимостью &gt; 0</td><td>{{.CostCount}}</td></tr>
    <tr><td>Сумма стоимостей</td><td>{{f2 .CostSum}}</td></tr>
    <tr><td>Средняя стоимость</td><td>{{f2 .Avg}}</td></tr>
    <tr class="total"><td>Итог (средняя × суммарное кол-во)</td><td>{{f2 .Total}}</td></tr>
  </table>
</div>
{{end}}

{{define "details"}}
<h1>Детали</h1>

<form method="post" class="card" id="detailForm">
  <input type="hidden" name="csrf_token" value="{{.CSRF}}">
  <input type="hidden" name="action" value="create">
  <input type="hidden" name="el_count" value="{{len .Draft.Rows}}">

  <label>Название <input type="text" name="name" value="{{.Draft.Name}}" required></label>
  <label>Редкость
    <select name="rarity">
      {{$r := .Draft.Rarity}}{{range .Rarities}}<option value="{{.}}" {{if eq . $r}}selected{{end}}>{{.}}</option>{{end}}
    </select>
  </label>
  <label class="checkbox">
    <input type="checkbox" name="is_craftable" {{if .Draft.IsCraft}}checked{{end}}> Крафтится
  </label>
  <button type="submit" formaction="/details?craft_open=1" formmethod="get" class="secondary">Применить</button>

  {{if .Draft.IsCraft}}
  <div class="craft">
    <h3>Элементы крафта</h3>
    {{$rows := .Draft.Rows}}
    {{range $i, $row := $rows}}
      <div class="element">
        <input type="hidden" name="el_kind_{{$i}}" value="{{$row.Kind}}">
        <select disabled>
          <option value="resource" {{if eq $row.Kind "resource"}}selected{{end}}>Ресурс</option>
          <option value="detail" {{if eq $row.Kind "detail"}}selected{{end}}>Деталь</option>
        </select>
        <select name="el_id_{{$i}}">
          <option value="">—</option>
          {{if eq $row.Kind "detail"}}
            {{$id := $row.ID}}{{range $.Items}}<option value="{{.ID}}" {{if eq .ID $id}}selected{{end}}>{{.Name}}</option>{{end}}
          {{else}}
            {{$id := $row.ID}}{{range $.Resources}}<option value="{{.ID}}" {{if eq .ID $id}}selected{{end}}>{{.Name}}</option>{{end}}
          {{end}}
        </select>
        <input type="number" step="any" min="0.0001" name="el_qty_{{$i}}" value="{{$row.Qty}}" placeholder="Кол-во">
      </div>
    {{end}}
    <button type="submit" formaction="/details?add=1&kind=resource" formmethod="get" class="secondary">+ Ресурс</button>
    <button type="submit" formaction="/details?add=1&kind=detail" formmethod="get" class="secondary">+ Деталь</button>
  </div>
  {{end}}

  <div class="prices">
    <label>Цена продажи <input type="number" step="any" min="0" name="sale_price" value="{{.Draft.SalePrice}}"></label>
    <label>Цена покупки <input type="number" step="any" min="0" name="buy_price" value="{{.Draft.BuyPrice}}"></label>
  </div>

  <button type="submit">Добавить деталь</button>
</form>

<table class="grid">
  <thead><tr><th>Название</th><th>Редкость</th><th>Крафт</th><th>Цены</th><th></th></tr></thead>
  <tbody>
  {{range .Items}}
    <tr>
      <form method="post" style="display:contents">
        <input type="hidden" name="csrf_token" value="{{$.CSRF}}">
        <input type="hidden" name="id" value="{{.ID}}">
        <td>{{.Name}}</td>
        <td>{{.Rarity}}</td>
        <td>
          {{if .IsCraftable}}да{{else}}нет{{end}}
          <div class="muted">
            {{range .Elements}}
              {{if eq .Kind "detail"}}{{index $.DetName .ID}}{{else}}{{index $.ResName .ID}}{{end}} ×{{.Qty}}<br>
            {{end}}
          </div>
        </td>
        <td>
          {{if .SalePrice}}{{f2 .SalePrice}}{{else}}—{{end}} /
          {{if .BuyPrice}}{{f2 .BuyPrice}}{{else}}—{{end}}
        </td>
        <td class="actions">
          <button name="action" value="delete" class="danger" onclick="return confirm('Удалить?')">Удалить</button>
        </td>
      </form>
    </tr>
  {{else}}
    <tr><td colspan="5" class="muted">Нет деталей</td></tr>
  {{end}}
  </tbody>
</table>
{{end}}
`
