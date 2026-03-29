# Adding a New Command

Use `commands/workout.go` as the canonical template.

## Step 1: Create the command file

Create `commands/<name>.go` in package `commands`.

### Parent command

```go
func new<Name>Cmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "<name>",
        Short: "View <name> data",
    }
    cmd.AddCommand(new<Name>ListCmd())
    cmd.AddCommand(new<Name>GetCmd())
    return cmd
}
```

### List subcommand

```go
func new<Name>ListCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "list",
        Short: "List <name> records",
        RunE: func(cmd *cobra.Command, args []string) error {
            params := buildListParams(cmd)
            all, _ := cmd.InheritedFlags().GetBool("all")

            if all {
                records, err := state.svc.<Name>.ListAll(cmd.Context(), params)
                if err != nil {
                    return err
                }
                return state.formatter.Format(cmd.OutOrStdout(), records)
            }

            page, err := state.svc.<Name>.List(cmd.Context(), params)
            if err != nil {
                return err
            }
            if page.NextToken != "" {
                fmt.Fprintf(cmd.ErrOrStderr(),
                    "Note: more pages available. Use --all to fetch everything.\n")
            }
            return state.formatter.Format(cmd.OutOrStdout(), page)
        },
    }
    addListFlags(cmd)
    return cmd
}
```

### Get subcommand (numeric ID)

For integer IDs (like Cycle), validate with `strconv.Atoi`:

```go
func new<Name>GetCmd() *cobra.Command {
    return &cobra.Command{
        Use:   "get <id>",
        Short: "Get a <name> by ID",
        Args:  cobra.ExactArgs(1),
        RunE: func(cmd *cobra.Command, args []string) error {
            id, err := strconv.Atoi(args[0])
            if err != nil {
                return fmt.Errorf("invalid <name> ID: %s", args[0])
            }
            result, err := state.svc.<Name>.Get(cmd.Context(), id)
            if err != nil {
                return err
            }
            return state.formatter.Format(cmd.OutOrStdout(), result)
        },
    }
}
```

### Get subcommand (string ID)

For string IDs (like Workout, Sleep), pass directly:

```go
RunE: func(cmd *cobra.Command, args []string) error {
    result, err := state.svc.<Name>.Get(cmd.Context(), args[0])
    if err != nil {
        return err
    }
    return state.formatter.Format(cmd.OutOrStdout(), result)
},
```

## Step 2: Add service field to state

In `commands/root.go`, add a field to the `services` struct:

```go
type services struct {
    // ... existing fields ...
    <Name> api.<Name>Service
}
```

## Step 3: Wire the service in PersistentPreRunE

In `commands/root.go`, in the `PersistentPreRunE` block where services are initialized:

```go
state.svc = &services{
    // ... existing services ...
    <Name>: api.New<Name>Service(client),
}
```

## Step 4: Register the command

In `NewRootCmd()` in `commands/root.go`:

```go
root.AddCommand(new<Name>Cmd())
```

## Step 5: Auth exemption (if needed)

If the command should work without authentication, add its command path to `authExemptCommands`:

```go
var authExemptCommands = map[string]bool{
    // ... existing entries ...
    "whoop <name> <subcmd>": true,
}
```

## Checklist

- [ ] `commands/<name>.go` created with parent + subcommands
- [ ] Service field added to `services` struct in `root.go`
- [ ] Service wired in `PersistentPreRunE`
- [ ] Command registered via `root.AddCommand()`
- [ ] `addListFlags(cmd)` called for list commands
- [ ] Table formatter cases added (see `add-table-formatter.md`)
- [ ] Tests written (see `add-tests.md`)
