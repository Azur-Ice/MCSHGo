# MCSHGo
A **M**inecraft **S**erver **H**elper coded with **Go**. Based on stdin/out from server. Reforged from MCSH(my private repo, written in python)

## How to use it?
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
servers:
    serverName1:
        rootFolder: ...
    serverName2:
        rootFolder: ...
    ...
```
- `servers`
    It contains the information of all your server need to be managed with MCSHGo.
    For each server, the **key** should be a custom name for it(and there should be a script file in `Script/` folder with the same name of it to start the server.), and the **value** should have **key** `rootFolder` .
    - `rootFolder`
        This is where your `.jar` file is.

MCSHGo will use the scriptFile with the same name of the server to start it up.\
The script file should contains "chdir" part and "run server" part.
> On linux, the sscript file should be a valid sh script.
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

## IO for each server.
- input:
    Write `xxx` to `serverName` with `serverName | xxx`.\
    Otherwise, MCSHGo will write the input to every server.
- output:
    It will present output from each server like this:
    ```
    YYYY-MM-DD HH:MM:SS [<serverName>/INFO]: ......
    YYYY-MM-DD HH:MM:SS [<serverName>/WARN]: ......
    ```