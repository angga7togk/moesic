# Moesic

**Moesic** is a freely accessible and open collection of Japanese music and anime. The project is **open source**, with the goal of providing a simple, aesthetically pleasing music listening experience focused on Japanese and anime content.

## Requirements

* [FFmpeg](https://ffmpeg.org/) must be installed and available in your system PATH.

### Install FFmpeg

#### Linux (Debian/Ubuntu)

```bash
sudo apt update && sudo apt install ffmpeg
```

#### macOS (with Homebrew)

```bash
brew install ffmpeg
```

#### Windows

Download FFmpeg from the official website:

* [https://ffmpeg.org/download.html](https://ffmpeg.org/download.html)
* Or directly from: [https://www.gyan.dev/ffmpeg/builds/](https://www.gyan.dev/ffmpeg/builds/)

Extract the ZIP, then:

1. Copy the `bin` folder path (e.g., `C:\ffmpeg\bin`)
2. Add it to your **System Environment Variables > Path**
3. Open new terminal and run `ffmpeg -version` to verify

---

## Installation

### Linux / macOS

Run the following command to install Moesic:

```bash
curl -fsSL https://raw.githubusercontent.com/angga7togk/moesic/main/install.sh | bash
```

Alternatively, you can download and run the installation script manually:

```bash
wget https://raw.githubusercontent.com/angga7togk/moesic/main/install.sh
chmod +x install.sh
./install.sh
```

### Windows

#### Option 1: Using PowerShell

Run the following command in PowerShell:

```powershell
iwr -useb https://raw.githubusercontent.com/angga7togk/moesic/main/install.ps1 | iex
```

#### Option 2: Using Batch File

Download and run [`install.bat`](https://raw.githubusercontent.com/angga7togk/moesic/main/install.bat).

### Manual Installation

Please download the latest version of Moesic from the [Releases](https://github.com/angga7togk/moesic/releases) page and place the binary in your system path.

## Contributing

If you would like to contribute by adding your favorite music or playlist, please read the [Contributing Guide](data/CONTRIBUTING.md).

## Preview

![Moesic Logo](.github/img/preview.png)

```bash
 __  __  ___  ___ ___ ___ ___
|  \/  |/ _ \| __/ __|_ _/ __|
| |\/| | (_) | _|\__ \| | (__
|_|  |_|\___/|___|___/___\___|


⭐️ Star to support our work!
   https://github.com/angga7togk/moesic

Usage:
  moesic <options>              Moesic CLI

Options:
  --random, --play, -p          Play random flat moesic
  --random-playlist, -rp        Play random playlist
  --random-single, -rs          Play random single moesic
  --help, -h                    Command help
  --info, -i                    Moesic info
```
