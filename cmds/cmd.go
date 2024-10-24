package cmds

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jalpp/passdiy/config"
	hcp "github.com/jalpp/passdiy/hcpvault"
	opass "github.com/jalpp/passdiy/onepassword"
	cmd "github.com/jalpp/passdiy/password"
)

type CommandItem struct {
	title, desc string
	Subcmd      []CommandItem
}

var (
	pwpDesc             = fmt.Sprintf("Generate strong %d-word passphrase", config.PASSPHRASE_COUNT_NUM)
	saltDesc            = fmt.Sprintf("Generate password with extra %d-char salt on top", config.SALT_EXTRA_LENGTH)
	hashDesc            = "Generate hash value of a password with Argon2id or bcrypthash"
	argonhashDesc       = "Generate hash value of a password with Argon2id"
	bcrpthashDesc       = "Generate hash value of a password with bcrypt algorithm"
	hcpvaultstoreDesc   = "Store a new secret to Hashicorp Vault"
	hcpvaultconnectDesc = "Generate HCP API token and connect to Hashicorp Vault"
	hcpvaultlistDesc    = "List HCP Vault secrets log details"
	opassstoreDesc      = "Store a new secret to 1Password in password format"
	opasslistDesc       = "List 1Password Vault item names"
	mainpassDesc        = "Generate strong passwords from various algorithms"
	mainpinDesc         = "Generate strong pins from various algorithms"
	maintokenDesc       = "Generate strong token from various algorithms"
	hcpDesc             = "Manage Token/Password on Hashicorp Vault"
	opassDesc           = "Manage Token/Password on 1Password"
)

func (i CommandItem) Title() string       { return i.title }
func (i CommandItem) Description() string { return i.desc }
func (i CommandItem) FilterValue() string { return i.title }

func GetSingleCommandInfo(cmd string) string {
	return fmt.Sprintf("Generate a single %s", cmd)
}

func GetMulCommandInfo(cmd string) string {
	return fmt.Sprintf("Generate %d multiple %s all at once", config.MULTIPLE_VALUE_COUNT, cmd)
}

func GetHundCommandInfo(cmd string) string {
	return fmt.Sprintf("Generate a single %s by generating 100s and randomly picking one", cmd)
}

func GetTenKCommandInfo(cmd string) string {
	return fmt.Sprintf("Generate a single %s by generating 10000s and randomly picking one", cmd)
}

func CreateCommandItems() []list.Item {
	const config = config.LOTTERY_WHEEL_COUNT

	passItems := []CommandItem{
		{title: "pass", desc: GetSingleCommandInfo("strong password")},
		{title: "passmul", desc: GetMulCommandInfo("password")},
		{title: fmt.Sprintf("pass%d", config), desc: GetHundCommandInfo("password")},
		{title: fmt.Sprintf("pass%d", config*config), desc: GetTenKCommandInfo("password")},
	}

	pinItems := []CommandItem{
		{title: "pin", desc: GetSingleCommandInfo("pin")},
		{title: "pinmul", desc: GetMulCommandInfo("pin")},
		{title: fmt.Sprintf("pin%d", config), desc: GetHundCommandInfo("pin")},
		{title: fmt.Sprintf("pin%d", config*config), desc: GetTenKCommandInfo("pin")},
	}

	tokenItems := []CommandItem{
		{title: "token", desc: GetSingleCommandInfo("token")},
		{title: "tokenmul", desc: GetMulCommandInfo("token")},
		{title: fmt.Sprintf("token%d", config), desc: GetHundCommandInfo("token")},
		{title: fmt.Sprintf("token%d", config*config), desc: GetTenKCommandInfo("token")},
	}
	hashItems := []CommandItem{
		{title: "argonhash", desc: argonhashDesc},
		{title: "bcrypthash", desc: bcrpthashDesc},
	}

	hcpItems := []CommandItem{
		{title: "hcpvaultstore", desc: hcpvaultstoreDesc},
		{title: "hcpvaultconnect", desc: hcpvaultconnectDesc},
		{title: "hcpvaultlist", desc: hcpvaultlistDesc},
	}

	opassItems := []CommandItem{
		{title: "1passstore", desc: opassstoreDesc},
		{title: "1passlist", desc: opasslistDesc},
	}

	return []list.Item{
		CommandItem{title: "pass", desc: mainpassDesc, Subcmd: passItems},
		CommandItem{title: "pin", desc: mainpinDesc, Subcmd: pinItems},
		CommandItem{title: "token", desc: maintokenDesc, Subcmd: tokenItems},
		CommandItem{title: "salt", desc: saltDesc},
		CommandItem{title: "pwp", desc: pwpDesc},
		CommandItem{title: "hash", desc: hashDesc, Subcmd: hashItems},
		CommandItem{title: "hcpvault", desc: hcpDesc, Subcmd: hcpItems},
		CommandItem{title: "1pass", desc: opassDesc, Subcmd: opassItems},
	}
}

type CommandFinishedMsg struct {
	result string
}

func (i CommandFinishedMsg) Result() string { return i.result }

func ExecuteCommand(command, input string) tea.Cmd {
	return func() tea.Msg {
		result := HandleCommand(command, input)
		time.Sleep(1 * time.Second)
		return CommandFinishedMsg{result: result}
	}
}

func HandleCommand(input, userInput string) string {
	const config = config.LOTTERY_WHEEL_COUNT
	switch strings.TrimSpace(input) {
	case "pass":
		return cmd.GetStrongPassword()
	case "passmul":
		return cmd.GetMul("pass")
	case fmt.Sprintf("pass%d", config):
		return cmd.GetHundPick("pass")
	case fmt.Sprintf("pass%d", (config * config)):
		return cmd.GetTenKPick("pass")
	case "token":
		return cmd.GetAPIToken()
	case "tokenmul":
		return cmd.GetMul("token")
	case fmt.Sprintf("token%d", config):
		return cmd.GetHundPick("token")
	case fmt.Sprintf("token%d", config*config):
		return cmd.GetTenKPick("token")
	case "pin":
		return cmd.GetPin()
	case "pinmul":
		return cmd.GetMul("pin")
	case fmt.Sprintf("pin%d", config):
		return cmd.GetHundPick("pin")
	case fmt.Sprintf("pin%d", config*config):
		return cmd.GetTenKPick("pin")
	case "pwp":
		return cmd.GetPwp()
	case "salt":
		return cmd.AddSalt(userInput)
	case "argonhash":
		return cmd.HashFunc(userInput)
	case "bcrypthash":
		return cmd.BcryptHash(userInput)
	case "hash":
		return cmd.HashFunc(userInput)
	case "hcpvaultstore":
		parts := strings.SplitN(userInput, "=", 2)
		if len(parts) == 2 {
			name := parts[0]
			value := parts[1]
			return hcp.Create(name, value)
		}
		return "Invalid format. Use 'name=value'."
	case "hcpvaultconnect":
		return hcp.Connect()
	case "hcpvaultlist":
		var list string = hcp.List()
		if strings.Contains(list, "Unauthorized") {
			return "Please connect to Hashicorp vault via hcpvaultconnect"
		}
		return list
	case "1passstore":
		parts := strings.SplitN(userInput, "|", 3)
		if len(parts) == 3 {
			user := parts[0]
			pass := parts[1]
			url := parts[2]
			return opass.Create(user, pass, url)
		}
		return "Invalid format. use 'user|value|url'."
	case "1passlist":
		return opass.List()
	default:
		return fmt.Sprintf("Unknown command: %s", input)
	}
}

func CoverUp(pass string) string {
	return cmd.CoverUp(pass)
}
