## admiral ssh

ssh to inventory host

### Synopsis

use host auto-complete for ssh, if `proxy` is set to true, proxy the connection / command through the server configured on `ssh-proxy` The domain will be appended to the hostname automatically

```
admiral ssh {hostname} [flags]
```

### Examples

```
admiral ssh host
admiral ssh host ls -l
```

### Options

```
  -h, --help   help for ssh
```

### SEE ALSO

* [admiral](admiral.md)	 - Admiral is a lightweight Ansible inventory database management tool

###### Auto generated by spf13/cobra on 3-Jan-2021
