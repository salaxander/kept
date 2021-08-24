# kept

kept is a CLI tool used to interact with Kubernetes [enhancements](https://github.com/kubernetes/enhancements). This is very much a beta project! In it's current state, it's primary use is to search for KEPs using a CLI interface. It also alows members of the Kubernetes release team to manage the lifecycle of KEPs.



# usage

### login

The first thing you'll want to do is login with GitHub. To do this, run:

```shell
$ kept login
```

This will display a code in your terminal and open a web browser window for you to enter the code in.



### list KEPs

To search for KEPs, run:

```shell
$ kept list
```

By default, this command will display all active KEPs. To filter them, you can use the supported flags.

| flag          | shorthand | function                                   |
| ------------- | --------- | ------------------------------------------ |
| `--milestone` | `-m`      | Filter KEPs by milestone (example  `1.23`) |
| `--sig`       | `-s`      | Filter KEPs by SIG (example `node`)        |
| `--stage`     | `-st`     | Filter KEPs by stage (example `beta`)      |
| `--tracked`   | `t`       | Use to display only tracked KEPs           |



### get KEP

to search an individual KEP using its issue number, run:

```shell
$ kept get {issue_number}
```

To immediately open the KEP in a brower, you can use the `--open` flag. Example:

```shell
$ kept get 2278 --open
```

