# MCSHGo
A **M**inecraft **S**erver **H**elper coded with **Go**. Based on stdin/out from server. Reforged from MCSH(my private repo, written in python)

## How to use it?
### Configuring

File structure of **MCSHGo** :

```
MCSHGO/
│  config.yml
│  MCSHGo
│
└─Scripts/
        serverName1.bat
        serverName2.bat
        ...
```
Fields in `config.yml`:
```
command_prefix: "#"
servers:
    serverName1:
        execOptions: ...
        execPath: ...
    serverName2:
        execOptions: ...
        execPath: ...
    ...
```
- `command_prefix`
    
    This is what MCSH will look for while you enter something into it.
    
    If any input contains it at the beginning, it will be considered as a **MCSH Command**
    
    check more in **MCSH Commands**
    
- `servers`
    It contains the information of all your server need to be managed with MCSHGo.
    For each server, the **key** should be a custom name for it(and there should be a script file in `Script/` folder with the same name of it to start the server.), and the **value** should have **key** `rootFolder` .
    
    - `execOptions`
        e.g. `-Xms4G -Xmm4G --nogui`.
    - `execPath`
        The path to the `.jar` file of your server.
        > - When doing bacnup jobs, MCSH will use the dir of this path to locate `world/` folder.
        > - Server will be using command `java -jar execOptions execPath` to start.

### IO for each server

- input:
    Write `xxx` to `serverName` with `serverName | xxx`.\
    Otherwise, MCSHGo will write the input to every server.
- output:
    It will present output from each server like this:
    ```
    YYYY-MM-DD HH:MM:SS [serverName/INFO]: ......
    YYYY-MM-DD HH:MM:SS [serverName/WARN]: ......
    ```

### MCSH Commands

A MCSH Commands is recognized with the prefix `#` for normal, you can configure it in `config.yml` .\

If you do `#xxx abc defgh ijk` \

`xxx abc defgh ijk` will be recognized as a **MCSH Command** \

> You can also use `serverName|#xxx` to run **MCSH Command** in specified server, or it will execute for every server.

- `backup [mode] [arg]`

	- enter ` `(empty) as `[mode]`

		> Show the backup list, Not developed yet.

	- enter ` make` as `[mode]`

		`arg` is `comment` , optional.\

		MCSH will copy your server's `serverRoot/world` to `Backups/` folder with a changed name in `servername - yyyy-mm-dd hh-mm-ss[ comment]` format

	- enter  `restore` as `[mode]`

		> stop the server, backup your server with comment `Restore2<backupName>` , Not developed yet.

- `run`
    - After you stoped some server, you can use this command to run it again.