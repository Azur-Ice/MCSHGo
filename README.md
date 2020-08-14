# MCSHGo
A minecraft server helper based on stdin/out from server. Reforged from MCSH(my private repo, wrote with python)

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

## `config.yml`
```
servers:
    serverName1:
        rootFolder: ...
        scriptFile: ...
    serverName2:
        rootFolder: ...
        scriptFile: ...
    ...
```
- `servers`
    It contains the information of all your server need to be managed with MCSHGo.
    For each server, the **key** should be a custom name for it, and the **value** should contains `rootFolder` and `scriptFile` .
    - `rootFolder`
        This is where your `.jar` file is.
    - `scriptFile`
        This is the name of the script file for starting the server.
        The script file should be located in `Scripts/` folder.
        > In the script, you should `cd` to your server root folder first, then start the server.

## IO for each server.
- input:
    Use `<serverName> | xxx` to write `xxx` to `serverName`.
    Otherwise, the input will be written to every server.
- output:
    It will present output from each server like this:
    ```
    YYYY-MM-DD HH:MM:SS [<serverName>/INFO]: ......
    YYYY-MM-DD HH:MM:SS [<serverName>/WARN]: ......
    ```