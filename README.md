# Admiral

Admiral is a command line tool to manage ansible inventory. 

It can also expose the inventory to ansible as a full inventory structure. 

As monitoring is also important, the tool can also expose the inventory in Prometheus static file structure where all the host groups are set as host 'groups' label.

Supported Operating systems
---------------------------

-   Windows 10
-   MacOs
-   Linux - Tested on Debian x64 based operating systems (Debian,
    Ubuntu, Mint...) but should work on any x64 Linux distribution

For other Operating Systems you can build from source

Installation
------------

-   Download the relevant version from the [release page](https://github.com/via-justa/admiral/releases) to a location in your `$PATH`
-   Add configuration file is detailed in the `Configuration File` section

Configuration File
-----------

The tool is expecting to find a toml configuration file with the database details in one of the following locations:
- /etc/admiral/config.toml
- ./config.toml
- $HOME/.admiral.toml

Example configuration file:
```
[database]
user = "root"
password = "local"
host = "localhost:3306"
db = "ansible"
```

Usage
-----------
The tool comes with a full help menu that can be accessed with the flag `-h, --help`. 
The compleat command documentation is also available [here](https://github.com/via-justa/admiral/docs/admiral.md)

Configuring the Database
-----------
A compatible MariaDB > 13 scheme can be found [here](https://github.com/via-justa/admiral/blob/master/fixtures/scheme.sql).

Issues and feature requests
-----------

I'm more than happy to reply to any issue of feature request via the github issue tracker.
When opening an issue, please provide the version you're using, Operating system and any useful information that can be used to investigate the issue.
When opening a feature request, please provide at least one detailed use case that the feature will solve.

License
-----------

<a rel="license" href="http://creativecommons.org/licenses/by-nc-sa/4.0/"><img alt="Creative Commons License" style="border-width:0" src="https://i.creativecommons.org/l/by-nc-sa/4.0/88x31.png" /></a><br />This work is licensed under a <a rel="license" href="http://creativecommons.org/licenses/by-nc-sa/4.0/">Creative Commons Attribution-NonCommercial-ShareAlike 4.0 International License</a>.
