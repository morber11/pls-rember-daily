# Pls Rember Daily

![image](https://i.redd.it/gxanhyonnlza1.png)

pls rember that wen u feel scare or frigten never forget ttimes wen u feeled happy

wen day is dark alway rember happy day

## Building

Requires Go 1.26.2 or newer with `cgo` enabled.

Output is inside the `dist\` folder

### Windows

Requires MinGW-w64 (`gcc` and `g++`) on your `PATH`.

```powershell
.\build.ps1
```

Optional flags:

- `-AllowMultipleRuns` plays the video every time instead of once per day
- `-SkipStartupDelay` skips the short delay before playing

### Linux

Requires `pkg-config` and the development files for GTK 4.8 or newer and X11.

```bash
./build-linux.sh
```
## Installing

Run the executable and pls rember

- **Windows:** copies itself to `%APPDATA%\PlsRemberDaily\`
- **Linux:** starts from wherever the binary is, move it somewhere permanent first e.g `~/.local/bin`.

## Uninstalling

Run the executable with `-uninstall`

- **Windows:**

  ```powershell
  .\PlsRemberDaily.exe -uninstall
  ```

Running from any location should automatically delete the exe in `%APPDATA%\PlsRemberDaily\` and startup entries

- **Linux:**

The binary itself is not removed, so delete it manually afterwards

  ```bash
  ./PlsRemberDaily -uninstall
  ```
