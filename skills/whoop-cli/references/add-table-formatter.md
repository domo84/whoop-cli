# Adding Table Output for a New Type

Use `renderWorkouts` in `internal/output/table.go` (line 172) as the template.

## Step 1: Add type-switch cases

In `TableFormatter.Format()` in `internal/output/table.go`, add three cases inside the
`switch v := data.(type)` block:

```go
case []api.<Type>:
    return render<Types>(w, v)
case *api.<Type>:
    return render<Types>(w, []api.<Type>{*v})
case *api.PaginatedResponse[api.<Type>]:
    return render<Types>(w, v.Records)
```

These handle the three shapes of data that commands produce:
- `[]<Type>` from `ListAll`
- `*<Type>` from `Get`
- `*PaginatedResponse[<Type>]` from `List`

Place these cases before the `default` block.

## Step 2: Write the render function

```go
func render<Types>(w io.Writer, items []api.<Type>) error {
    table := noBorderTable(w, []string{"ID", "START", "FIELD1", "FIELD2", "STATE"})
    for _, item := range items {
        // Handle nil scores
        field1, field2 := "N/A", "N/A"
        if item.Score != nil {
            field1 = fmt.Sprintf("%.1f", item.Score.SomeFloat)
            field2 = strconv.Itoa(item.Score.SomeInt)
        }

        // Handle optional pointer fields
        optionalField := "N/A"
        if item.Score != nil && item.Score.OptionalField != nil {
            optionalField = fmt.Sprintf("%.1f%%", *item.Score.OptionalField)
        }

        table.Append([]string{
            strconv.Itoa(item.ID),                    // or item.ID for string IDs
            item.Start.Format("2006-01-02 15:04"),    // date format
            field1,
            field2,
            item.ScoreState,
        })
    }
    return table.Render()
}
```

## Conventions

- **Headers**: short uppercase strings (e.g. "ID", "START", "STRAIN", "AVG HR", "STATE")
- **Columns**: pick 5-7 most useful fields for the table view
- **Dates**: format as `"2006-01-02 15:04"` (Go reference time)
- **Nil scores**: always show `"N/A"` when `Score` is nil
- **Optional pointers** (`*float64`): check for nil before dereferencing
- **Duration**: use `formatDuration()` helper (already in table.go) for time differences
- **Booleans**: render as `"Yes"` / `"No"` (see sleep's nap field)

## Helper functions available

- `noBorderTable(w, headers)` -- creates a borderless ASCII table
- `formatDuration(d time.Duration)` -- formats as `"1h 30m"` or `"45m"`
