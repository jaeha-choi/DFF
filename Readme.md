# DFF -- Automatic rune/item page/spell
![mainImg](./image/DFF.png)

### D for Flash!
DFF is a program that fetches data from op.gg to set runes, item pages, and spells automatically.

### Installation
1. Go to [release page](https://github.com/jaeha-choi/DFF/releases) to download the latest DFF.
2. Extract downloaded `DFF_win_v0.x.zip`
3. Edit `config` file to update the League client directory.

### Execution
1. Double click the `DFF_win_v0.x.exe` to execute DFF
2. Play games as you normally would.

### Configuration (`config.json`) options

- `client_dir` : Game client directory, where League of Legends is installed.
- `enable_rune` : Enable automatic rune fetch.
- `enable_item` : Enable automatic item fetch.
- `enable_spell` : Enable automatic spell fetch.
- `interval` : Polling interval in seconds. Value from 1 ~ 5 can be set. 2 is reasonable.
- `d_flash` : If you are using left slot (D spell) for Flash, set as true.
    - Note: This option only works if Flash is a recommended spell.
- `debug` : Debugging option. Prints extra information when executed with a terminal.
- `language`: Language of rune page title. Only `en_US` and `ko_KR` show correctly on DFF. All languages show correctly in League of Legends client.

### Disclaimer
DFF was created under Riot Games' "Legal Jibber Jabber" policy using assets owned by Riot Games.  Riot Games does not endorse or sponsor this project.
