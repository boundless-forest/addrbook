package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
)

type WorkSpaceCommand struct {
	fs *flag.FlagSet
}

func NewWorkSpaceCommand() *WorkSpaceCommand {
	wsc := WorkSpaceCommand{
		fs: flag.NewFlagSet("workspace", flag.ContinueOnError),
	}
	return &wsc
}
func (wsc *WorkSpaceCommand) Init(args []string) error {
	return wsc.fs.Parse(args)
}
func (wsc *WorkSpaceCommand) Name() string {
	return wsc.fs.Name()
}
func (wsc *WorkSpaceCommand) Run(db *DataBase) error {
	if len(os.Args) <= 2 {
		printUsage()
		return errors.New("workspace command requires subcommand")
	}

	subCmds := []Runner{
		NewWsNewCommand(),
		NewWsDelCommand(),
		NewWsListCommand(),
		NewWsOpenCommand(),
		NewWsSaveCommand(),
		NewWsUpdateCommand(),
		NewWsDeleteCommand(),
	}

	for _, cmd := range subCmds {
		if cmd.Name() == os.Args[2] {
			cmd.Init(os.Args[3:])
			return cmd.Run(db)
		}
	}

	return fmt.Errorf("unknown command %s", os.Args[2])
}

type HelpCommand struct {
	fs *flag.FlagSet
}

func NewHelpCommand() *HelpCommand {
	h := HelpCommand{
		fs: flag.NewFlagSet("help", flag.ContinueOnError),
	}
	return &h
}
func (h *HelpCommand) Init(args []string) error {
	return h.fs.Parse(args)
}
func (h *HelpCommand) Name() string {
	return h.fs.Name()
}
func (h *HelpCommand) Run(db *DataBase) error {
	printUsage()
	return nil
}

// workspace new ...

type WsNewCommand struct {
	fs *flag.FlagSet

	name string
}

func NewWsNewCommand() *WsNewCommand {
	new := WsNewCommand{
		fs: flag.NewFlagSet("new", flag.ContinueOnError),
	}
	new.fs.StringVar(&new.name, "name", "", "The name of the workspace")
	return &new
}
func (new *WsNewCommand) Name() string {
	return new.fs.Name()
}
func (new *WsNewCommand) Init(args []string) error {
	return new.fs.Parse(args)
}
func (new *WsNewCommand) Run(db *DataBase) error {
	if new.name == "" {
		printUsage()
		return errors.New("workspace name is required")
	}
	if err := db.CreateWorkSpace(new.name); err != nil {
		return err
	}
	if err := SaveToDB(db); err != nil {
		return errors.New("save workspaces error: " + err.Error())
	}

	fmt.Printf("workspace %s created successfully \n", new.name)
	return nil
}

// workspace del ...

type WsDelCommand struct {
	fs *flag.FlagSet

	name string
}

func NewWsDelCommand() *WsDelCommand {
	del := WsDelCommand{
		fs: flag.NewFlagSet("del", flag.ContinueOnError),
	}
	del.fs.StringVar(&del.name, "name", "", "The name of the workspace")
	return &del
}
func (del *WsDelCommand) Name() string {
	return del.fs.Name()
}
func (del *WsDelCommand) Init(args []string) error {
	return del.fs.Parse(args)
}
func (del *WsDelCommand) Run(db *DataBase) error {
	if del.name == "" {
		printUsage()
		return errors.New("workspace name is required")
	}

	if err := db.DeleteWorkSpace(del.name); err != nil {
		return err
	}
	if err := SaveToDB(db); err != nil {
		return errors.New("save workspaces error: " + err.Error())
	}

	fmt.Printf("workspace %s deleted successfully", del.name)
	return nil
}

// workspace list ...

type WsListCommand struct {
	fs *flag.FlagSet
}

func NewWsListCommand() *WsListCommand {
	list := WsListCommand{
		fs: flag.NewFlagSet("list", flag.ContinueOnError),
	}
	return &list
}
func (list *WsListCommand) Name() string {
	return list.fs.Name()
}
func (list *WsListCommand) Init(args []string) error {
	return list.fs.Parse(args)
}
func (list *WsListCommand) Run(db *DataBase) error {
	fmt.Println("workspaces:", db.ListWorkSpaces())
	return nil
}

// workspace open ...

type WsOpenCommand struct {
	fs *flag.FlagSet
}

func NewWsOpenCommand() *WsOpenCommand {
	open := WsOpenCommand{
		fs: flag.NewFlagSet("open", flag.ContinueOnError),
	}
	return &open
}
func (open *WsOpenCommand) Name() string {
	return open.fs.Name()
}
func (open *WsOpenCommand) Init(args []string) error {
	return open.fs.Parse(args)
}
func (open *WsOpenCommand) Run(db *DataBase) error {
	html, err := generateHtmlPage(db)
	if err != nil {
		return err
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(html))
	})

	http.ListenAndServe(":8080", nil)
	return nil
}

// workspace save ...

type WsSaveCommand struct {
	fs *flag.FlagSet

	workspace string
	contract  string
	address   string
	note      string
}

func NewWsSaveCommand() *WsSaveCommand {
	save := WsSaveCommand{
		fs: flag.NewFlagSet("save", flag.ContinueOnError),
	}
	save.fs.StringVar(&save.workspace, "workspace", "", "The name of the workspace")
	save.fs.StringVar(&save.contract, "contract", "", "The name of the contract")
	save.fs.StringVar(&save.address, "address", "", "The address of the contract")
	save.fs.StringVar(&save.note, "note", "", "The extra information of the contract")
	return &save
}
func (save *WsSaveCommand) Name() string {
	return save.fs.Name()
}
func (save *WsSaveCommand) Init(args []string) error {
	return save.fs.Parse(args)
}
func (save *WsSaveCommand) Run(db *DataBase) error {
	if err := db.Save(save.workspace, save.contract, save.address, save.note); err != nil {
		return err
	}
	if err := SaveToDB(db); err != nil {
		return errors.New("save workspaces error: " + err.Error())
	}

	fmt.Printf("The contract information [%s -> %s] has been saved successfully. \n", save.contract, save.address)
	return nil
}

// workspace update ...

type WsUpdateCommand struct {
	fs *flag.FlagSet

	workspace string
	contract  string
	address   string
	note      string
}

func NewWsUpdateCommand() *WsUpdateCommand {
	update := WsUpdateCommand{
		fs: flag.NewFlagSet("update", flag.ContinueOnError),
	}
	update.fs.StringVar(&update.workspace, "workspace", "", "The name of the workspace")
	update.fs.StringVar(&update.contract, "contract", "", "The name of the contract")
	update.fs.StringVar(&update.address, "address", "", "The address of the contract")
	update.fs.StringVar(&update.note, "note", "", "The extra information of the contract")
	return &update
}
func (update *WsUpdateCommand) Name() string {
	return update.fs.Name()
}
func (update *WsUpdateCommand) Init(args []string) error {
	return update.fs.Parse(args)
}
func (update *WsUpdateCommand) Run(db *DataBase) error {
	if err := db.Update(update.workspace, update.contract, update.address, update.note); err != nil {
		return err
	}
	if err := SaveToDB(db); err != nil {
		return errors.New("save workspaces error: " + err.Error())
	}
	fmt.Printf("The contract information [%s] has been updated successfully. \n", update.contract)
	return nil
}

// workspace delete ...

type WsDeleteCommand struct {
	fs *flag.FlagSet

	workspace string
	contract  string
}

func NewWsDeleteCommand() *WsDeleteCommand {
	delete := WsDeleteCommand{
		fs: flag.NewFlagSet("delete", flag.ContinueOnError),
	}
	delete.fs.StringVar(&delete.workspace, "workspace", "", "The name of the workspace")
	delete.fs.StringVar(&delete.contract, "contract", "", "The name of the contract")
	return &delete
}
func (delete *WsDeleteCommand) Name() string {
	return delete.fs.Name()
}
func (delete *WsDeleteCommand) Init(args []string) error {
	return delete.fs.Parse(args)
}
func (delete *WsDeleteCommand) Run(db *DataBase) error {
	if err := db.Delete(delete.workspace, delete.contract); err != nil {
		return err
	}
	if err := SaveToDB(db); err != nil {
		return errors.New("save workspaces error: " + err.Error())
	}
	fmt.Printf("The contract information [%s] has been deleted successfully.", delete.contract)
	return nil
}

type Runner interface {
	Init([]string) error
	Run(*DataBase) error
	Name() string
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  workspace new --name <name>,  Create a new workspace.")
	fmt.Println("  workspace del --name <name>,  Delete a new workspace.")
	fmt.Println("  workspace list,               List the workspaces managed by the current user.")
	fmt.Println("  workspace open,               Open the workspace with the default browser at http://127.0.0.1:8080.")
	fmt.Println("  workspace save   --workspace <workspace> --contract <name> --address <addr>, --note <something>, Save the contract address into the workspace.")
	fmt.Println("  workspace update --workspace <workspace> --contract <name> --address <addr>, --note <something>, Update the contract address in the specified workspace.")
	fmt.Println("  workspace delete --workspace <workspace> --contract <name> --address <addr>, --note <something>, Update the contract address in the specified workspace.")
}

func run(args []string) error {
	cmds := []Runner{
		NewWorkSpaceCommand(),
		NewHelpCommand(),
	}

	db := DataBase{}
	if err := LoadDB(&db); err != nil {
		return errors.New("load workspaces error: " + err.Error())
	}

	for _, cmd := range cmds {
		if cmd.Name() == args[0] {
			cmd.Init(os.Args[2:])
			return cmd.Run(&db)
		}
	}

	printUsage()
	return fmt.Errorf("unknown command %s", os.Args[1])
}
