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
        rootFolder: ...
    serverName2:
        rootFolder: ...
    ...
```
- `command_prefix`
    
    This is what MCSH will look for while you enter something into it.
    
    If any input contains it at the beginning, it will be considered as a **MCSH Command**
    
    check more in **MCSH Commands**
    
- `servers`
    It contains the information of all your server need to be managed with MCSHGo.
    For each server, the **key** should be a custom name for it(and there should be a script file in `Script/` folder with the same name of it to start the server.), and the **value** should have **key** `rootFolder` .
    
    - `rootFolder`
        This is where your `.jar` file is.

MCSHGo will use the scriptFile with the same name of the server to start it up.\
The script file should contains "chdir" part and "run server" part.
> On linux, the sscript file should be a valid sh script.
> And it must have permission of executing. (`sudo chmod +x yourScriptName.sh`)
Example for Windows:
```bat
cd G:\_MC\_Server\_AzurCraft\1.16.2 flat && G: && java -jar -Xms6g -Xmx6g fabric-server-launch.jar --nogui
```
> `cd G:\_MC\_Server\_AzurCraft\1.16.2 flat && G:` is the "chdir" part, and the following things are "run server" part.
Example for Linux:
```sh
#!/bin/sh
cd /mnt/g/_MC/_Server/_AzurCraft/1.16.2\ flat && java -jar -Xms6g -Xmx6g fabric-server-launch.jar --nogui
```
> `cd /mnt/g/_MC/_Server/_AzurCraft/1.16.2\ flat` is the "chdir" part, and the following things are "run server" part.

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

A MCSH Commands is start with `#` for normal, you can configure it in `config.yml` .\

If you do `#xxx abc defgh ijk` \

`xxx abc defgh ijk` will be considered as a **MCSH Command** \

> You can also use `serverName|#xxx` to run **MCSH Command** in specified server, or it will execute for every server.

- `backup [mode] [arg]`

	- enter  ` `(empty) as `[mode]`

		> Show the backup list, Not developed yet.

	- enter  ` make` as `[mode]`

		`arg` is `comment` , optional.\

		MCSH will copy your server's `serverRoot/world` to `Backups/` folder with a changed name in `servername - yyyy-mm-dd hh-mm-ss[ comment]` format

	- enter  `restore` as `[mode]`

		> stop the server, backup your server with comment `Restore2<backupName>` , Not developed yet.

- `run`
    - After you stoped some server, you can use this command to run it again.