# csgo-discord-rpc
A Discord RPC integration to the game CS:GO

## Usage

You can manually install, or use the [installer in the releases](https://github.com/marcos-travasso/csgo-discord-rpc/releases/latest)

### Manually installation
First you need to build the executable:
```
$ go build -ldflags -H=windowsgui
```

*Optionally: Move the executable into the Window's startup folder*

Then you move the *gamestate_integration_discordrpc.cfg* file into the [Counter-Strike cfg folder](https://developer.valvesoftware.com/wiki/Counter-Strike:_Global_Offensive_Game_State_Integration#Locating_CS:GO_Install_Directory)
