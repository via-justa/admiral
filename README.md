![Go Test](https://github.com/via-justa/admiral/workflows/Go/badge.svg)
![Language](https://img.shields.io/badge/Language-go-green)
[![Go Report Card](https://goreportcard.com/badge/github.com/via-justa/admiral)](https://goreportcard.com/report/github.com/via-justa/admiral)
[![Coverage Status](https://coveralls.io/repos/github/via-justa/admiral/badge.svg)](https://coveralls.io/github/via-justa/admiral)
[![license](https://img.shields.io/badge/license-CC-blue)](https://creativecommons.org/licenses/by-nc-sa/4.0/)
[![latest release](https://img.shields.io/badge/-latest_release-blue)](https://github.com/via-justa/admiral/releases/latest)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fvia-justa%2Fadmiral.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fvia-justa%2Fadmiral?ref=badge_shield)

# Admiral

Admiral - a lightweight, opinionated, command-line tool to manage [Ansible](https://www.ansible.com/) database based inventory and manage the hosts.

Managing hundreds or thousands of hosts on different infrastructure providers in a file-based inventory can be a complicated task that requires a lot of management overhead and can lead to misconfiguration, relationship loops between groups, and very long inventory files.

Admiral is meant to solve this issue by storing the inventory in a database and allow for better visibility of the relationships, making the inventory searchable and making the management of the inventory easier and error-prone.

As Prometheus is the most common monitoring tool these days and the one monitoring tool I favorite the most, I baked the option to export the inventory as  [file_sd_configs](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#file_sd_config) and provide the groups the host is part of as metric labels

## Admiral main features:
### Inventory feature
- Creation and Edit of hosts and groups in JSON structure using your favorite editor
- Creation of hosts in one command for use with CI/CD pipelines
- Bulk import hosts/groups / child-groups from JSON file
- Bulk enable/disable of hosts/monitoring in one command (none interactive)
- Command-line edit and delete of hosts, groups, and their relationships
- Create a new host/group from an existing one (copy) to save time and need for configuration
- Setting default common configurations for new hosts/groups
- MariaDB ssh proxy connection

### CLI features
- Full auto-completion of commands
- Realtime retrieval of hosts and groups for bash auto-completion
- Export of the inventory in ansible readable structure
- Export of the inventory in Prometheus static file structure
- Ansible ping command wrapper to validate ansible can communicate with the hosts
- SSH command with proxy-jump option to leverage the hostname auto-comple

Supported Operating systems
---------------------------

-   Windows 10
-   macOS
-   Linux - Tested on Debian x64 based operating systems (Debian,
    Ubuntu, Mint...) but should work on any x64 Linux distribution

For other Operating Systems you can build from source

Database backend support
------------
- MariaDB > 13 (recommended)
- SQLite3

Installation
------------

-   Download and extract the relevant version from the [release page](https://github.com/via-justa/admiral/releases) to a location in your `$PATH`
-   Add configuration file is detailed in the `Configuration File` section
-   Add bash-completion to your `.bashrc` or `.profile` (optional)
    ```shell
    . <(admiral completion)
    ```

Configuration File
-----------

The tool is expecting to find a `toml` configuration file with the database details in one of the following locations:
- /etc/admiral/config.toml
- ./config.toml
- $HOME/.admiral.toml

Example configuration file: [config-template.toml](./config-template.toml)

Usage
-----------

The tool comes with a full help menu that can be accessed with the flag `-h, --help`.
The compleat command documentation is also available [here](/docs/admiral.md)

Configuring the Database
-----------

A compatible `MariaDB > 13` scheme can be found [here](/fixtures/mariadb/01_scheme.sql).
A compatible `sqlite3` scheme can be found [here](/fixtures/dqlite/01_scheme.sql).

Use admiral for ssh connections
-----------

Admiral can be used to ssh to hosts with auto-completion of host names

SSH and Database proxy
-----------

If you need to connect to secured network where you have ssh gateway or the database server does not allow remote connection to the database, you can configure ssh proxy the tool will use to proxy the connection thru.
It is also possible to use `admiral ssh` to proxy SSH connections via the ssh proxy.

Using the prometheus `file_sd_configs` and labels to filter jobs
-----------

The easiest way to get the `file_sd_configs` generated and read by prometheus is by using a cron job or systemd timer.
```shell
*/1 * * * * "/usr/local/bin/admiral prometheus > /etc/prometheus/prometheus_file_sd.json.new && mv /etc/prometheus/prometheus_file_sd.json.new /etc/prometheus/prometheus_file_sd.json"
```
This [nginx-exporter](https://github.com/nginxinc/nginx-prometheus-exporter) job example will keep all hosts with direct group matching regex `web-.*` and from those drop host with direct group `web-proxy` using the relabel_configs mechanism.
```yaml
- job_name: 'nginx'
    file_sd_configs:
      - files:
        - "/etc/prometheus/prometheus_file_sd.json"
    relabel_configs:
      - source_labels: ['group']
        regex: 'web-.*'
        action: keep
      - source_labels: ['group']
        regex: 'web-proxy'
        action: drop
      - source_labels: [__address__]
        regex:  '(.*)'
        target_label: __address__
        replacement: '${1}:9113'
```

Usage examples
-----------

View all hosts
```
$ admiral view host
IP           |Hostname       |domain         |Enabled      |Monitored    |Direct Groups     |Inherited Groups
1.2.3.4      |host-2         |via-justa.com  |true         |true         |web               |app
1.2.3.5      |host-3         |via-justa.com  |true         |true         |web               |app
1.2.3.6      |host-4         |via-justa.com  |true         |true         |web               |app
1.2.3.7      |host-5         |via-justa.com  |true         |true         |web               |app
...
```

View specific host
```
$ admiral view host host-2
IP           |Hostname      |domain         |Enabled      |Monitored    |Direct Groups  |Inherited Groups
1.2.3.4      |host-2        |via-justa.com  |true         |true         |web            |app
```

View host as JSON (with vars)
```
$ admiral view host host-2
[
    {
        "ip": "1.2.3.4",
        "hostname": "host-2",
        "domain": "via-justa.com",
        "variables": {},
        "enable": true,
        "monitor": true,
        "direct_group": "web"
    }
]
```

Create New host
```
$ admiral create host host-2
<edit in editor>
IP           |Hostname       |domain         |Enabled      |Monitored    |Direct Groups |Inherited Groups
1.2.3.4      |host-2         |via-justa.com  |true         |true         |web           |app
Please confirm [y/n]: y
```

Create New host none-interactively (useful for CI/CD pipelines)
```
admiral create host host-2 --monitor=false --enable --ip 1.2.3.4 --group=web
IP           |Hostname       |domain         |Enabled      |Monitored    |Direct Groups |Inherited Groups
1.2.3.4      |host-2         |via-justa.com  |true         |true         |web           |
Please confirm [y/n]: y
```

Create a new host from an existing one
```
$ admiral copy host host-1 host-2
<edit in editor>
IP           |Hostname       |domain         |Enabled      |Monitored    |Direct Groups     |Inherited Groups
1.2.3.4      |host-2         |via-justa.com  |true         |true         |web               |
Please confirm [y/n]: y
```

Delete host
```
$ admiral delete host host-2
IP           |Hostname     |domain         |Enabled      |Monitored    |Direct Groups     |Inherited Groups
1.2.3.4      |host-2       |via-justa.com  |true         |true         |web               |
Please confirm [y/n]: y
```

Run Ansible ping to check Ansible connectivity
```
$ admiral ping host-2
host-2.via-justa.com | SUCCESS => {
    "ansible_facts": {
        "discovered_interpreter_python": "/usr/bin/python"
    },
    "changed": false,
    "ping": "pong"
}
```

Create child group relationship
```
$ admiral create child web app
Parent       |Parent ID    |Child         |Child ID
app          |2            |web           |3
Please confirm [y/n]: y
```

Issues and feature requests
-----------

I'm more than happy to reply to any issue of feature request via the github issue tracker.
When opening an issue, please provide the version you're using, Operating system and any useful information that can be used to investigate the issue.
When opening a feature request, please provide at least one detailed use case that the feature will solve.

License
-----------

<a rel="license" href="http://creativecommons.org/licenses/by-nc-sa/4.0/"><img alt="Creative Commons License" style="border-width:0" src="https://i.creativecommons.org/l/by-nc-sa/4.0/88x31.png" /></a><br />This work is licensed under a <a rel="license" href="http://creativecommons.org/licenses/by-nc-sa/4.0/">Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International License</a>.


## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fvia-justa%2Fadmiral.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fvia-justa%2Fadmiral?ref=badge_large)