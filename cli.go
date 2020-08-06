package main

import (
	// "errors"
	"fmt"
	"github.com/gasiordev/go-cli"
	"os"
)

func getCLIVersionHandler(j *Jailguard) func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		fmt.Fprintf(c.GetStdout(), VERSION+"\n")
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIStateListHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}

		err := j.ListStateItems("")
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIStateRemoveHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}

		err := j.RemoveStateItem(c.Arg("item_type"), c.Arg("item_name"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIStateImportHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}

		err := j.ImportStateItem(c.Arg("item_type"), c.Arg("item_name"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}

		return 0
	}
	return fn
}

func (j *Jailguard) getCLIBaseDownloadHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		ow := false
		if c.Flag("overwrite") == "true" {
			ow = true
		}

		err := j.DownloadBase(c.Arg("release"), ow)
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIBaseRemoveHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.RemoveBase(c.Arg("release"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIBaseListHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		err := j.ListStateItems("bases")
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIJailListHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		err := j.ListStateItems("jails")
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIJailCreateHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		start := false
		if c.Flag("start") == "true" {
			start = true
		}
		err := j.CreateJail(c.Arg("file"), c.Flag("base"), start)
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIJailRemoveHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		stop := false
		if c.Flag("stop") == "true" {
			stop = true
		}
		err := j.RemoveJail(c.Arg("jail"), stop)
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIJailStopHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.StopJail(c.Arg("jail"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLIJailStartHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.StartJail(c.Arg("jail"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLINetifListHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		err := j.ListStateItems("netifs")
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLINetifCreateHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.CreateNetif(c.Arg("name"), c.Arg("ip_addr_begin"), c.Arg("ip_addr_end"), c.Arg("if_name"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLINetifDestroyHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.DestroyNetif(c.Arg("name"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLINetifAliasAddHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.AddNetifAlias(c.Arg("name"), c.Arg("ip_addr"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLINetifAliasDeleteHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.DeleteNetifAlias(c.Arg("name"), c.Arg("ip_addr"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func (j *Jailguard) getCLINetifAliasListHandler() func(*cli.CLI) int {
	fn := func(c *cli.CLI) int {
		if c.Flag("debug") == "true" {
			j.Debug = true
		}
		if c.Flag("quiet") == "true" {
			j.Quiet = true
		}

		err := j.ListNetifAliases(c.Arg("name"))
		if err != nil {
			j.Log(LOGERR, err.Error())
			return 2
		}
		return 0
	}
	return fn
}

func NewJailguardCLI(j *Jailguard) *cli.CLI {
	c := cli.NewCLI("Jailguard", "Create and manage jails in FreeBSD", "Mikolaj Gasior")

	_ = c.AddCmd("state_list", "Lists saved state items", j.getCLIStateListHandler())

	state_remove := c.AddCmd("state_remove", "Remove item from state", j.getCLIStateRemoveHandler())
	state_remove.AddArg("item_type", "TYPE", "", cli.TypeString|cli.Required)
	state_remove.AddArg("item_name", "NAME", "", cli.TypeString|cli.Required)
	state_remove.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	state_remove.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	state_import := c.AddCmd("state_import", "Import item to state", j.getCLIStateImportHandler())
	state_import.AddArg("item_type", "TYPE", "", cli.TypeString|cli.Required)
	state_import.AddArg("item_name", "NAME", "", cli.TypeAlphanumeric|cli.Required)
	state_import.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	state_import.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	base_download := c.AddCmd("base_download", "Downloads FreeBSD base", j.getCLIBaseDownloadHandler())
	// TODO: Change 'release' flag to TypeAlphanumeric once AllowHyphen gets implemented in go-cli
	base_download.AddArg("release", "RELEASE", "", cli.TypeString|cli.Required)
	base_download.AddFlag("overwrite", "w", "", "Overwrite if exists", cli.TypeBool)
	base_download.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	base_download.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	_ = c.AddCmd("base_list", "List FreeBSD bases", j.getCLIBaseListHandler())

	base_remove := c.AddCmd("base_remove", "Removes FreeBSD base", j.getCLIBaseRemoveHandler())
	base_remove.AddArg("release", "RELEASE", "", cli.TypeString|cli.Required)
	base_remove.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	base_remove.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	jail_create := c.AddCmd("jail_create", "Create jail source", j.getCLIJailCreateHandler())
	jail_create.AddArg("file", "FILE.JAIL", "", cli.TypePathFile|cli.MustExist|cli.Required)
	jail_create.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	jail_create.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)
	jail_create.AddFlag("base", "b", "", "Base to use", cli.TypeString)
	jail_create.AddFlag("start", "s", "", "Start jail after creating", cli.TypeBool)

	jail_remove := c.AddCmd("jail_remove", "Remove jail source", j.getCLIJailRemoveHandler())
	jail_remove.AddArg("jail", "JAIL", "", cli.TypeString|cli.Required)
	jail_remove.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	jail_remove.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)
	jail_remove.AddFlag("stop", "s", "", "Stop if running", cli.TypeBool)

	jail_stop := c.AddCmd("jail_stop", "Stop jail", j.getCLIJailStopHandler())
	jail_stop.AddArg("jail", "JAIL", "", cli.TypeString|cli.Required)
	jail_stop.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	jail_stop.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	jail_start := c.AddCmd("jail_start", "Start jail", j.getCLIJailStartHandler())
	jail_start.AddArg("jail", "JAIL", "", cli.TypeString|cli.Required)
	jail_start.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	jail_start.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	_ = c.AddCmd("jail_list", "List jails", j.getCLIJailListHandler())

	netif_create := c.AddCmd("netif_create", "Create network interface", j.getCLINetifCreateHandler())
	netif_create.AddArg("name", "NAME", "", cli.TypeString|cli.Required)
	netif_create.AddArg("ip_addr_begin", "IP_ADDR_BEGIN", "", cli.TypeString|cli.Required)
	netif_create.AddArg("ip_addr_end", "IP_ADDR_END", "", cli.TypeString|cli.Required)
	netif_create.AddArg("if_name", "INTERFACE_NAME", "", cli.TypeString)
	netif_create.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	netif_create.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	netif_destroy := c.AddCmd("netif_destroy", "Destroy network interface", j.getCLINetifDestroyHandler())
	netif_destroy.AddArg("name", "NAME", "", cli.TypeString|cli.Required)
	netif_destroy.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	netif_destroy.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	_ = c.AddCmd("netif_list", "List network interfaces", j.getCLINetifListHandler())

	netif_alias_add := c.AddCmd("netif_alias_add", "Add alias IP address to a network interface", j.getCLINetifAliasAddHandler())
	netif_alias_add.AddArg("name", "NAME", "", cli.TypeString|cli.Required)
	netif_alias_add.AddArg("ip_addr", "IP_ADDR", "", cli.TypeString|cli.Required)
	netif_alias_add.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	netif_alias_add.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	netif_alias_delete := c.AddCmd("netif_alias_delete", "Delete alias IP address from a network interface", j.getCLINetifAliasDeleteHandler())
	netif_alias_delete.AddArg("name", "NAME", "", cli.TypeString|cli.Required)
	netif_alias_delete.AddArg("ip_addr", "IP_ADDR", "", cli.TypeString|cli.Required)
	netif_alias_delete.AddFlag("quiet", "q", "", "Do not output anything", cli.TypeBool)
	netif_alias_delete.AddFlag("debug", "d", "", "Print more information", cli.TypeBool)

	netif_alias_list := c.AddCmd("netif_alias_list", "List network interface alias IP addresses", j.getCLINetifAliasListHandler())
	netif_alias_list.AddArg("name", "NAME", "", cli.TypeString|cli.Required)

	// pf_anchor_create : name
	// pf_anchor_remove : name

	// jail_portfwd_create : name host_port jail_port
	// jail_portfwd_remove : name [-a] host_port jail_port ---> -a removes all
	// jail_portfwd_list : name

	// jail_natpass_create : name gw_netif
	// jail_natpass_remove : name  ---> there can be only one
	// jail_natpass_list : name

	// whenever app is started, check for the configuration - if not then write about running
	// 'configure' or passing '--config'

	// whenever app is started, check if the configuration has correct values

	// config_list
	// config_get key            ------> /etc/jailguard.conf otherwise --config for all the commands
	// config_set key value
	// configure - before first use, tool has to be configured - all values can be passed interactively
	//   or as flags

	// when checking the jail conf: - if interface exists and if not then it has to be created
	//                              - if pf entries exist then pf has to be installed
	//                                                         and pf anchor created
	//                              - if networking then check the sysctl for allow_raw_sockets and if not present then display warning (which can be supressed by -q)

	// state_check for base and jail: checks if state is up-to-date

	_ = c.AddCmd("version", "Prints version", getCLIVersionHandler(j))

	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		os.Args = []string{"jailguard", "version"}
	}
	return c
}
