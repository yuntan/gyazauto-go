# GyazAuto-go
Automatically upload captured screenshots to Gyazo. Useful to combine with other screenshot tools (ex. Windows standard screencapture functionary, Steam, etc).

Works on Windows and Linux. It may work on macOS but prebuilt binary is not provided (because I don't have mac).

## USAGE
1. Download executable.
2. Run it by double click or on command line.
3. Gyazo authorization page is opened. Grant access for this app.
4. Edit configuration file.
  - `%USERPROFILE%\AppData\Roaming\gyazauto\config.yaml` on Windows
  - `~/.config/gyazauto/config.yaml` on Linux
  
  ```yaml
  access_token: abcd0123... # do not change
  watch_dirs:
    # You can use ~
    # for Windows standard screenshot tool
    - ~\Pictures\Screenshots
    # for Steam game screenshots
    - C:\Program Files (x86)\Steam\userdata\0123456\789\remote\012345\screenshots
    # for Minecraft on Windows
    - ~\AppData\Roaming\.minecraft\screenshots
    # for Linux standard screenshot tools
    - ~/Picture/Screenshots/
    # for Minecraft on Linux
    - ~/.minecraft/screenshots
  ```
5. Restart app.
6. Register app to startup.

## TODO
Pullreq welcome

- add tags to screenshots
- write build instruction
