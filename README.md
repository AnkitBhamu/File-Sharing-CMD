# fsgoclient - File Sharing via Terminal

A command-line tool for sharing files between systems over the network using Go.

## Installation
```
run  bash ./install.sh
```
## Usage

```
fsgoclient -mode=<sender|receiver> [options]
```

### Modes

- `-mode=sender`  
  Start the tool in sender mode.

- `-mode=receiver`  
  Start the tool in receiver mode.

## Options

### Sender Mode

- `-rcvIp=<ip:port>`  
  IP address and port of the receiver (e.g., `localhost:8080`).

- `-sfdr=<path>`  
  Path to a file containing paths of files to be sent (one file path per line).

### Receiver Mode

- `-downloadDir=<path>`  
  Directory where received files will be saved.

- `-port=<port>`  
  Port on which the receiver will listen (default: `8080`).

## Examples

**Sender:**

```
fsgoclient -mode=sender -rcvIp=localhost:8080 -sfdr="C:\path\to\filestosend.txt"
```

**Receiver:**

```
fsgoclient -mode=receiver -port=8080 -downloadDir="C:\path\to\downloadDir"
```

## Notes

- Ensure both sender and receiver are on the same network or accessible via the given IP.
- The file passed to `-sfdr` should contain absolute paths to the files to be sent.
- Avoid using relative paths to prevent file resolution issues.
- Make sure the download directory exists before starting the receiver.
